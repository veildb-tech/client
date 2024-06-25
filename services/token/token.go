/*
Copyright © 2024 Bridge Digital
*/
package token

import (
	"encoding/json"
	"fmt"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/request"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/response"
)

// Get token from server
func JwtToken(credentials map[string]string) string {
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
		err := response.WrongResponseObserver(data)
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}

	if len(configData["token"]) > 0 {
		return configData["token"]
	}

	return ""
}
