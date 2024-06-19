/*
Copyright © 2024 Bridge Digital
*/
package savekey

import (
	"fmt"
	"strings"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/envfile"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/keypubfile"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/exp/maps"
)

func Execute(isNew bool, keyName string) string {
	if !isNew {
		reCreate()
		return ""
	}

	var (
		options        = []string{"Yes", "No"}
		selectedOption string
	)

	prompt := &survey.Select{
		Message: "Want to create a public Pem key for downloading?",
		Options: options,
	}

	survey.AskOne(prompt, &selectedOption)

	if selectedOption == "No" {
		return ""
	}

	var keyData string

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
		Name:   "Key Data",
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

// Function for regenerating a key
func reCreate() {
	savedWorkspaces, err := envfile.ReadEnvFile()
	if err != nil {
		fmt.Println(predefined.BuildError("Error:"), err)
		return
	}

	var (
		selectedWorkspaceIndex, selectedServerIndex int
		savedWorkspacesKeys, savedServersKeys       []string
		options                                     = []string{"Yes", "No"}
		selectedOption, currentWorkspace            string
	)

	savedWorkspacesKeys = maps.Keys(savedWorkspaces)

	if len(savedWorkspacesKeys) > 1 {
		promptW := &survey.Select{
			Message: "Select one of your saved workspaces:",
			Options: savedWorkspacesKeys,
		}

		survey.AskOne(promptW, &selectedWorkspaceIndex)
	} else {
		selectedWorkspaceIndex = 0
		fmt.Println(predefined.BuildAnsw("Your saved workspaces: ", savedWorkspacesKeys[selectedWorkspaceIndex]))
	}

	currentWorkspace = savedWorkspacesKeys[selectedWorkspaceIndex]

	savedConfigData := savedWorkspaces[currentWorkspace]

	savedServers := savedWorkspaces[currentWorkspace].Servers
	savedServersKeys = maps.Keys(savedServers)

	if len(savedServersKeys) > 1 {
		promptS := &survey.Select{
			Message: "Select one of your saved servers:",
			Options: savedServersKeys,
		}

		survey.AskOne(promptS, &selectedServerIndex)
	} else {
		selectedServerIndex = 0
		fmt.Println(predefined.BuildAnsw("Your saved server: ", savedServersKeys[selectedServerIndex]))
	}

	savedKeyName := savedServers[savedServersKeys[selectedServerIndex]].KeyFile

keyNameAsk:

	var keyName, keyData, currentServerName string

	currentServerName = savedServersKeys[selectedServerIndex]

	if !keypubfile.IsKeyFileExist(savedKeyName) {
		keyName = currentWorkspace + "_" + currentServerName

		if len(keyName) > 0 {
			keyName = keypubfile.CreateKeyPubFile(keyName)
		} else {
			fmt.Println(predefined.BuildError("Something is wrong with getting the server key name."))
			return
		}
	} else {
		prompt := &survey.Select{
			Message: "Key file is already exists. Do you want to override existing file?",
			Options: options,
		}

		survey.AskOne(prompt, &selectedOption)

		if selectedOption == "No" {
			goto keyNameAsk
		}

		keyName = savedKeyName
	}

	qKeyData := &survey.Question{
		Name:   "Key Data",
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

	if envfile.IsEnvFileExist(false) {
		if len(keyName) > 0 {
			configData := map[string]string{
				"token":     savedConfigData.ServiceToken,
				"workspace": savedWorkspacesKeys[selectedWorkspaceIndex],
				"keyName":   keyName,
				"server":    currentServerName,
			}

			envfile.WriteEnvFile(envfile.ConfigData(configData))

			fmt.Println(predefined.BuildOk("The public key has been saved successfully"))
		}
	}
}
