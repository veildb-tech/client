/*
Copyright © 2024 Bridge Digital
*/
package envfile

import (
	"encoding/json"
	"fmt"
	"os"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
	"golang.org/x/exp/maps"
)

type Server struct {
	KeyFile string `json:"key_file"`
}

type Workspace struct {
	ServiceToken string            `json:"token"`
	Servers      map[string]Server `json:"servers"`
}

type ConfigDataService interface {
	ConfigData()
}

func ConfigData(userData map[string]string) map[string]Workspace {
	var workspace string = userData["workspace"]
	var server string = userData["server"]

	currentServer := Server{
		KeyFile: userData["keyName"],
	}

	serversData := map[string]Server{
		server: currentServer,
	}

	data := Workspace{
		ServiceToken: userData["token"],
		Servers:      serversData,
	}

	var configData = map[string]Workspace{}
	configData[workspace] = data

	return configData
}

func IsEnvFileExist(msgSupress bool) bool {
	var result bool = true

	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		fmt.Printf(predefined.BuildError("Cannot get current APP directory: %W.\n"), errDir)
		return false
	}

	_, err := os.Stat(configDir + "/" + services.EnvFileName)
	if err != nil {
		if !msgSupress {
			fmt.Printf(predefined.BuildWarning("Env file not found. Please run: %s login.\n"), services.AppName)
		}
		result = false
	}

	return result
}

func CreateEnvFile(config map[string]Workspace) {
	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		fmt.Printf(predefined.BuildError("Cannot get current APP directory: %W.\n"), errDir)
		return
	}

	file, err := os.Create(configDir + "/" + services.EnvFileName)
	if err != nil {
		fmt.Println(predefined.BuildError("Cannot create file:"), err)
		return
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		fmt.Println(predefined.BuildError("Cannot write config data to file:"), err)
		return
	}
}

func WriteEnvFile(config map[string]Workspace) {
	workspaceKeys := maps.Keys(config)
	workspace := workspaceKeys[0]

	if !IsEnvFileExist(true) {
		CreateEnvFile(config)
	}

	configFromFile, err := ReadEnvFile()
	if err != nil {
		fmt.Printf(predefined.BuildError("Error: %s"), err)
		return
	}

	serverKeys := maps.Keys(config[workspace].Servers)
	server := serverKeys[0]

	workspaceDataFromFile, v := configFromFile[workspace]
	if v {
		serverDataFromFile, vs := workspaceDataFromFile.Servers[server]
		if vs {
			serverDataFromFile.KeyFile = config[workspace].Servers[server].KeyFile
			workspaceDataFromFile.Servers[server] = serverDataFromFile
		} else {
			workspaceDataFromFile.Servers[server] = Server{KeyFile: config[workspace].Servers[server].KeyFile}
		}

		workspaceDataFromFile.ServiceToken = config[workspace].ServiceToken

		configFromFile[workspace] = workspaceDataFromFile
	} else {
		configFromFile[workspace] = config[workspace]
	}

	data, errData := json.Marshal(configFromFile)
	if errData != nil {
		fmt.Println(predefined.BuildError("Cannot encode config data: "), err)
		return
	}

	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		fmt.Printf(predefined.BuildError("Cannot get current APP directory: %W.\n"), errDir)
		return
	}

	err = os.WriteFile(configDir+"/"+services.EnvFileName, data, 0644)
	if err != nil {
		fmt.Println(predefined.BuildError("Cannot write to env file:"), err)
		return
	}
}

func ReadEnvFile() (map[string]Workspace, error) {
	if !IsEnvFileExist(true) {
		return nil, fmt.Errorf(predefined.BuildError("env file not found. Please use the [login] command to update it"))
	}

	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		return nil, fmt.Errorf(predefined.BuildError("cannot get current APP directory: %W"), errDir)
	}

	file, err := os.ReadFile(configDir + "/" + services.EnvFileName)
	if err != nil {
		return nil, fmt.Errorf(predefined.BuildError("env file is not readable: %W"), errDir)
	}

	config := make(map[string]Workspace)

	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf(predefined.BuildError("the settings record is not readable: %W"), errDir)
	}

	return config, nil
}
