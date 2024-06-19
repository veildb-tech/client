/*
Copyright © 2024 Bridge Digital
*/
package keypubfile

import (
	"fmt"
	"os"
	"strings"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
)

// Key file operations
func IsKeyFileExist(keyname string) bool {
	var result bool = true

	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		fmt.Printf(predefined.BuildError("Cannot get current APP directory: %W.\n"), errDir)
		return false
	}

	var ext string = services.PubKeyExt

	if strings.Contains(keyname, services.PubKeyExt) {
		ext = ""
	}

	_, err := os.Stat(configDir + "/" + keyname + ext)
	if err != nil {
		result = false
	}

	return result
}

func CreateKeyPubFile(keyname string) string {
	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		fmt.Printf(predefined.BuildError("Cannot get current APP directory: %W.\n"), errDir)
		return ""
	}

	var ext string = services.PubKeyExt

	if strings.Contains(keyname, ext) {
		ext = ""
	}

	keyFileName := keyname + ext

	file, err := os.Create(configDir + "/" + keyFileName)
	if err != nil {
		fmt.Println(predefined.BuildError("Cannot create key file:"), err)
		return ""
	}

	defer file.Close()

	return keyFileName
}

func WriteKeyPubFile(keyData string, keyFileName string) string {
	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		fmt.Printf(predefined.BuildError("Cannot get current APP directory: %W.\n"), errDir)
		return ""
	}

	data := []byte(keyData)

	var ext string = services.PubKeyExt

	if strings.Contains(keyFileName, services.PubKeyExt) {
		ext = ""
	}

	keyFileName = keyFileName + ext

	err := os.WriteFile(configDir+"/"+keyFileName, data, 0664)
	if err != nil {
		fmt.Println(predefined.BuildError("Cannot write key file:"), err)
	}

	return keyFileName
}

func ReadKeyPubFile(keyname string) string {
	var result string = ""

	if IsKeyFileExist(keyname) {
		configDir, errDir := services.CurrentAppDir()
		if errDir != nil {
			fmt.Printf(predefined.BuildError("Cannot get current APP directory: %W.\n"), errDir)
			return ""
		}

		var ext string = services.PubKeyExt

		if strings.Contains(keyname, services.PubKeyExt) {
			ext = ""
		}

		keyData, err := os.ReadFile(configDir + "/" + keyname + ext)
		if err != nil {
			fmt.Printf(predefined.BuildError("Cannot read the %s file: %W.\n"), keyname+ext, errDir)
			return ""
		}

		result = string(keyData)
	} else {
		fmt.Println(predefined.BuildWarning("Couldn't find key file. Ask the admin to give you a public key. Or create one if you have one using the save-key command."))
	}

	return result
}
