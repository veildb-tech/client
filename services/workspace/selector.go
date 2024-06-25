/*
Copyright © 2024 Bridge Digital
*/
package workspace

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/response"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/token"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/workspace/servers"
	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/exp/maps"
)

type Data struct {
	Identifier string            `json:"identifier"`
	Workspaces []string          `json:"workspaces"`
	Servers    map[string]string `json:"servers"`
}

func Workspace(credentials map[string]string) map[string]string {
	tokenRow := token.JwtToken(credentials)
	if len(tokenRow) == 0 {
		fmt.Println(predefined.BuildError("Something wrong. Token is empty"))
		return nil
	}

	workspaceData := getProfileData(tokenRow)

	allWorkspaces := map[int]string{}
	var workspaceResult string = ""

	if len(workspaceData.Workspaces) > 1 {
		for k, workspace := range workspaceData.Workspaces {
			allWorkspaces[k] = workspace
		}

		var selectedWorkspace int

		prompt := &survey.Select{
			Message: "Select workspace:",
			Options: workspaceData.Workspaces,
		}

		survey.AskOne(prompt, &selectedWorkspace)

		workspaceResult = allWorkspaces[selectedWorkspace]
	} else if len(workspaceData.Workspaces) == 1 {
		workspaceResult = workspaceData.Workspaces[0]
		fmt.Println(predefined.BuildAnsw("Your saved workspace: ", workspaceResult))
	} else {
		fmt.Println(predefined.BuildWarning("You don't assigned to any workspace"))
	}

	credentials["workspace"] = workspaceResult

	tokenRow = token.JwtToken(credentials)
	if len(tokenRow) == 0 {
		fmt.Println(predefined.BuildError("Something wrong. Token is empty"))
		return nil
	}

	var (
		options                                              = []string{"Yes", "No"}
		selectedOption, selectedServerName, selectedServerId string
		configData                                           = map[string]string{}
	)

	prompt := &survey.Select{
		Message: "Do you want to save a new server key or update an existing one?",
		Options: options,
	}

	survey.AskOne(prompt, &selectedOption)

	if selectedOption == "No" {
		configData = map[string]string{
			"token":     tokenRow,
			"workspace": workspaceResult,
			"server":    selectedServerName,
			"serverId":  selectedServerId,
		}

		return configData
	}

	workspaceData = getProfileData(tokenRow)

	if len(workspaceData.Servers) > 1 {
		selectedServerId, selectedServerName = servers.Server(workspaceData.Servers)
	} else if len(workspaceData.Servers) == 1 {
		selectedServerId = maps.Keys(workspaceData.Servers)[0]
		selectedServerName = workspaceData.Servers[selectedServerId]

		fmt.Println(predefined.BuildAnsw("Your saved server: ", selectedServerName))
	} else {
		fmt.Println(predefined.BuildWarning("You don't assigned to any server"))
	}

	configData = map[string]string{
		"token":     tokenRow,
		"workspace": workspaceResult,
		"server":    selectedServerName,
		"serverId":  selectedServerId,
	}

	return configData
}

func getProfileData(token string) Data {
	var workspaceData Data

	req, err := http.NewRequest("GET", services.WebServiceProfileUrl(), nil)
	if err != nil {
		fmt.Println(predefined.BuildError("Something wrong with GET request data:"), err)
		return workspaceData
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(predefined.BuildError("Invalid token:"), err)
		return workspaceData
	}

	if resp == nil {
		fmt.Println(predefined.BuildError("Bad request:"), http.StatusBadRequest)
		return workspaceData
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(predefined.BuildError("Error:"), err)
		return workspaceData
	}

	// Unmarshal the JSON data into the struct
	wsErr := json.Unmarshal([]byte(data), &workspaceData)
	if wsErr != nil {
		err := response.WrongResponseObserver(data)
		if err != nil {
			fmt.Println(err)
		}
	}

	return workspaceData
}
