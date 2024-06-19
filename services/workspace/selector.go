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
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/workspace/servers"
	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/exp/maps"
)

func Workspace(token string) (workspace string, server string) {
	req, err := http.NewRequest("GET", services.WebServiceProfileUrl(), nil)
	if err != nil {
		fmt.Println(predefined.BuildError("Something wrong with GET request data:"), err)
		return "", ""
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(predefined.BuildError("Invalid token:"), err)
		return "", ""
	}

	if resp == nil {
		fmt.Println(predefined.BuildError("Bad request:"), http.StatusBadRequest)
		return "", ""
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(predefined.BuildError("Error:"), err)
		return "", ""
	}

	type Data struct {
		Identifier string            `json:"identifier"`
		Workspaces []string          `json:"workspaces"`
		Servers    map[string]string `json:"servers"`
	}

	// Unmarshal the JSON data into the struct
	var workspaceData Data
	allWorkspaces := map[int]string{}

	wsErr := json.Unmarshal([]byte(data), &workspaceData)
	if wsErr != nil {
		response.WrongResponseObserver(data)
		return "", ""
	}

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

	var selectedServer string = ""

	if len(workspaceData.Servers) > 1 {
		selectedServer = servers.Server(workspaceData.Servers)
	} else if len(workspaceData.Servers) == 1 {
		selectedServer = maps.Values(workspaceData.Servers)[0]
		fmt.Println(predefined.BuildAnsw("Your saved server: ", selectedServer))
	} else {
		fmt.Println(predefined.BuildWarning("You don't assigned to any server"))
	}

	return workspaceResult, selectedServer
}
