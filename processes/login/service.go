/*
Copyright © 2024 Bridge Digital
*/
package login

import (
	"fmt"
	"strings"

	saveKey "github.com/dbvisor-pro/client/processes/savekey"
	"github.com/dbvisor-pro/client/services"
	"github.com/dbvisor-pro/client/services/envfile"
	"github.com/dbvisor-pro/client/services/predefined"
	workspacePac "github.com/dbvisor-pro/client/services/workspace"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func Execute(cmd *cobra.Command) string {
	// Check if service URL is configured, if not prompt for it
	if !isServiceUrlConfigured() {
		serviceUrl := promptForServiceUrl()
		if len(serviceUrl) == 0 {
			return ""
		}
		saveServiceUrl(serviceUrl)
	}

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

// isServiceUrlConfigured checks if service URL is already saved in config
func isServiceUrlConfigured() bool {
	if !envfile.IsEnvFileExist(true) {
		return false
	}

	config, err := envfile.ReadEnvFile()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(config.ServiceUrl)) > 0
}

// promptForServiceUrl prompts user to enter service URL
func promptForServiceUrl() string {
	var serviceUrl string

	qUrl := &survey.Question{
		Prompt: &survey.Input{
			Message: "Service URL:",
			Default: services.DefaultWebServiceUrl,
			Help:    "Enter the VeilDB service URL (e.g., https://app.veildb.com)",
		},
		Validate: func(val interface{}) error {
			if str, _ := val.(string); len(strings.TrimSpace(str)) == 0 {
				return fmt.Errorf(predefined.BuildError("the Service URL cannot be empty"))
			}
			return nil
		},
	}

	survey.AskOne(qUrl.Prompt, &serviceUrl, survey.WithValidator(qUrl.Validate))
	return strings.TrimSpace(serviceUrl)
}

// saveServiceUrl saves the service URL to config file
func saveServiceUrl(serviceUrl string) {
	// Reset cache so new URL is used
	services.ResetConfigCache()

	if !envfile.IsEnvFileExist(true) {
		// Create a minimal config with just the service URL
		config := envfile.Config{
			ServiceUrl: serviceUrl,
			Data:       make(map[string]envfile.Workspace),
		}
		envfile.CreateEnvFile(config)
		return
	}

	config, err := envfile.ReadEnvFile()
	if err != nil {
		fmt.Println(predefined.BuildError("Error reading config:"), err)
		return
	}

	config.ServiceUrl = serviceUrl
	envfile.WriteEnvFile(config)
}
