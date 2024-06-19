/*
Copyright © 2024 Bridge Digital
*/
package login

import (
	"encoding/json"
	"fmt"
	"strings"

	saveKey "gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/processes/savekey"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/envfile"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/request"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/response"
	workspacePac "gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/workspace"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func Execute(cmd *cobra.Command) string {

	credentials := map[string]string{
		"username": "",
		"password": "",
	}

	var token, username, password, workspace, keyFileName, server string

	qUsername := &survey.Question{
		Name:   "Username",
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
		Name:   "Password",
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

	token = jwtToken(credentials)
	if len(token) == 0 {
		return ""
	}

	workspace, server = workspacePac.Workspace(token)
	if len(workspace) == 0 {
		return ""
	}

	configData := map[string]string{
		"token":     token,
		"workspace": workspace,
		"server":    server,
	}

	if !envfile.IsEnvFileExist(false) {
		configData["keyName"] = ""
		envfile.CreateEnvFile(envfile.ConfigData(configData))
	}

	keyFileName = saveKey.Execute(true, workspace+"_"+server)

	if len(keyFileName) > 0 {
		configData["keyName"] = keyFileName
	}

	envfile.WriteEnvFile(envfile.ConfigData(configData))

	return predefined.BuildOk("You logged in successfully")
}

// Get token from server
func jwtToken(credentials map[string]string) string {
	credsInJson, err := json.Marshal(credentials)
	if err != nil {
		fmt.Println(predefined.BuildError("Error encoding to json:"), err)
		return ""
	}

	data, err := request.CreatePostRequest(credsInJson, services.WebServiceAuthUrl(), nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var configData map[string]string

	configErr := json.Unmarshal(data, &configData)
	if configErr != nil {
		response.WrongResponseObserver(data)
		return ""
	}

	if len(configData["token"]) > 0 {
		return configData["token"]
	}

	return ""
}
