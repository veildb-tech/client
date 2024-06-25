/*
Copyright © 2024 Bridge Digital
*/
package login

import (
	"fmt"
	"strings"

	saveKey "gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/processes/savekey"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/envfile"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
	workspacePac "gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/workspace"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func Execute(cmd *cobra.Command) string {
	credentials := map[string]string{
		"username": "",
		"password": "",
	}

	var username, password, keyFileName string

	qUsername := &survey.Question{
		Prompt: &survey.Input{Message: "Username:"},
		Validate: func(val interface{}) error {
			if str, _ := val.(string); len(strings.TrimSpace(str)) == 0 {
				return fmt.Errorf(predefined.BuildError("the Username cannot be empty"))
			}
			return nil
		},
	}

	survey.AskOne(qUsername.Prompt, &username, survey.WithValidator(qUsername.Validate))
	if len(username) == 0 {
		return ""
	}

	credentials["username"] = username

	qPwd := &survey.Question{
		Prompt: &survey.Password{Message: "Password:"},
		Validate: func(val interface{}) error {
			if str, _ := val.(string); len(strings.TrimSpace(str)) == 0 {
				return fmt.Errorf(predefined.BuildError("the Password cannot be empty"))
			}
			return nil
		},
	}

	survey.AskOne(qPwd.Prompt, &password, survey.WithValidator(qPwd.Validate))
	if len(string(password)) == 0 {
		return ""
	}

	credentials["password"] = string(password)

	configData := workspacePac.Workspace(credentials)
	if len(configData["workspace"]) == 0 {
		fmt.Println(predefined.BuildError("Something wrong. Workspace is empty"))
		return ""
	}

	if !envfile.IsEnvFileExist(true) {
		configData["keyName"] = ""
		envfile.CreateEnvFile(envfile.ConfigData(configData))
	}

	if len(configData["server"]) > 0 && len(configData["serverId"]) > 0 {
		keyFileName = saveKey.Execute(true, configData["workspace"]+"_"+configData["server"])

		if len(keyFileName) > 0 {
			configData["keyName"] = keyFileName
		}
	}

	envfile.WriteEnvFile(envfile.ConfigData(configData))

	return predefined.BuildOk("You logged in successfully")
}
