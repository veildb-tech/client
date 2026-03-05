/*
Copyright © 2024 Bridge Digital
*/
package savekey

import (
	"fmt"
	"strings"

	"github.com/dbvisor-pro/client/services"
	"github.com/dbvisor-pro/client/services/envfile"
	"github.com/dbvisor-pro/client/services/keypubfile"
	"github.com/dbvisor-pro/client/services/predefined"
	"github.com/dbvisor-pro/client/services/workspace"
	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/exp/maps"
)

func Execute(isNew bool, keyName string) string {
	if !isNew {
		reCreate()
		return ""
	}

	var (
		options                 = []string{"Yes", "No"}
		selectedOption, keyData string
	)

	if len(keyName) == 0 {
		return ""
	}

	if !keypubfile.IsKeyFileExist(keyName) {
		keypubfile.CreateKeyPubFile(keyName)
	} else {
		prompt := &survey.Select{
			Message: "Key file is already exists. Do you want to override existing file?",
			Options: options,
		}

		survey.AskOne(prompt, &selectedOption)

		if selectedOption == "No" {
			return keyName + services.PubKeyExt
		}
	}

	qKeyData := &survey.Question{
		Prompt: &survey.Multiline{Message: "Enter public key:"},
		Validate: func(val interface{}) error {
			if str, _ := val.(string); len(strings.TrimSpace(str)) == 0 {
				return fmt.Errorf(predefined.BuildError("the key cannot be empty"))
			}
			return nil
		},
	}

	survey.AskOne(qKeyData.Prompt, &keyData, survey.WithValidator(qKeyData.Validate))

	return keypubfile.WriteKeyPubFile(keyData, keyName)
}

// reCreate is called by the save-key command to add or update server keys.
func reCreate() {
	savedConfig, err := envfile.ReadEnvFile()
	if err != nil {
		fmt.Println(predefined.BuildError("Error:"), err)
		return
	}

	currentWorkspace := savedConfig.CurrentWorkspace
	fmt.Println(predefined.BuildAnsw("Your workspace: ", currentWorkspace))

	savedServers := savedConfig.Data[currentWorkspace].Servers
	savedServersKeys := maps.Keys(savedServers)

	// Decide action
	const actionAdd = "Add new server key"
	const actionUpdate = "Update existing server key"

	actionOptions := []string{actionAdd}
	if len(savedServersKeys) > 0 {
		actionOptions = append(actionOptions, actionUpdate)
	}

	var selectedAction int
	if len(actionOptions) > 1 {
		prompt := &survey.Select{
			Message: "What do you want to do?",
			Options: actionOptions,
		}
		survey.AskOne(prompt, &selectedAction)
	}

	if actionOptions[selectedAction] == actionAdd {
		addNewServerKey(savedConfig, currentWorkspace)
	} else {
		updateExistingServerKey(savedConfig, currentWorkspace, savedServers, savedServersKeys)
	}
}

func addNewServerKey(savedConfig envfile.Config, currentWorkspace string) {
	// Fetch available servers from API using saved token
	profileData := workspace.GetProfileData(savedConfig.ServiceToken)

	if len(profileData.Servers) == 0 {
		fmt.Println(predefined.BuildError("No servers found in your profile"))
		return
	}

	// Build list of server names for selection
	serverIds := maps.Keys(profileData.Servers)
	serverNames := make([]string, len(serverIds))
	for i, id := range serverIds {
		serverNames[i] = profileData.Servers[id]
	}

	var selectedServerIndex int
	if len(serverNames) > 1 {
		prompt := &survey.Select{
			Message: "Select server to add key for:",
			Options: serverNames,
		}
		survey.AskOne(prompt, &selectedServerIndex)
	} else {
		fmt.Println(predefined.BuildAnsw("Server: ", serverNames[0]))
	}

	selectedServerName := serverNames[selectedServerIndex]
	selectedServerId := serverIds[selectedServerIndex]

	keyName := currentWorkspace + "_" + selectedServerName
	saveKeyForServer(savedConfig, currentWorkspace, selectedServerName, selectedServerId, keyName)
}

func updateExistingServerKey(savedConfig envfile.Config, currentWorkspace string, savedServers map[string]envfile.Server, savedServersKeys []string) {
	var selectedServerIndex int

	if len(savedServersKeys) > 1 {
		prompt := &survey.Select{
			Message: "Select one of your saved servers:",
			Options: savedServersKeys,
		}
		survey.AskOne(prompt, &selectedServerIndex)
	} else {
		fmt.Println(predefined.BuildAnsw("Your saved server: ", savedServersKeys[0]))
	}

	currentServerName := savedServersKeys[selectedServerIndex]
	savedKeyName := savedServers[currentServerName].KeyFile
	selectedServerId := savedServers[currentServerName].ServerId

	var keyName string

	if !keypubfile.IsKeyFileExist(savedKeyName) {
		keyName = currentWorkspace + "_" + currentServerName
		keyName = keypubfile.CreateKeyPubFile(keyName)
		if len(keyName) == 0 {
			fmt.Println(predefined.BuildError("Something is wrong with getting the server key name."))
			return
		}
	} else {
		var selectedOption string
		prompt := &survey.Select{
			Message: "Key file already exists. Do you want to override it?",
			Options: []string{"Yes", "No"},
		}
		survey.AskOne(prompt, &selectedOption)

		if selectedOption == "No" {
			return
		}
		keyName = savedKeyName
	}

	saveKeyForServer(savedConfig, currentWorkspace, currentServerName, selectedServerId, keyName)
}

func saveKeyForServer(savedConfig envfile.Config, currentWorkspace, serverName, serverId, keyName string) {
	var keyData string

	qKeyData := &survey.Question{
		Prompt: &survey.Multiline{Message: "Enter public key:"},
		Validate: func(val interface{}) error {
			if str, _ := val.(string); len(strings.TrimSpace(str)) == 0 {
				return fmt.Errorf(predefined.BuildError("the key cannot be empty"))
			}
			return nil
		},
	}

	survey.AskOne(qKeyData.Prompt, &keyData, survey.WithValidator(qKeyData.Validate))

	if len(keyData) == 0 {
		return
	}

	keyName = keypubfile.WriteKeyPubFile(keyData, keyName)

	if len(keyName) > 0 {
		configData := map[string]string{
			"token":     savedConfig.ServiceToken,
			"workspace": currentWorkspace,
			"keyName":   keyName,
			"server":    serverName,
			"serverId":  serverId,
		}

		envfile.WriteEnvFile(envfile.ConfigData(configData))
		fmt.Println(predefined.BuildOk("The public key has been saved successfully"))
	}
}
