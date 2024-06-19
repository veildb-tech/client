/*
Copyright © 2024 Bridge Digital
*/
package response

import (
	"encoding/json"
	"fmt"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
)

type WrongData struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func WrongResponseObserver(data []byte) {
	var wrongData WrongData

	dbErr := json.Unmarshal([]byte(data), &wrongData)
	if dbErr != nil {
		fmt.Println(predefined.BuildError("Error decoding from json:"), dbErr)
		return
	}

	if wrongData.Code == 401 && wrongData.Msg == "Invalid JWT Token" {
		fmt.Println(predefined.BuildWarning("Your token has expired. Use the [login] command to update it"))
		return
	} else {
		if wrongData.Code > 0 && len(wrongData.Msg) > 0 {
			fmt.Printf(predefined.BuildError("Code: %d. Message: %s \n"), wrongData.Code, wrongData.Msg)
		} else {
			fmt.Println(predefined.BuildError("Something wrong!"))
		}
		return
	}
}
