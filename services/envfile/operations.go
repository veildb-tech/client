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
	KeyFile  string `json:"key_file"`
	ServerId string `json:"server_id"`
}

type Workspace struct {
	Servers map[string]Server `json:"servers"`
}

type Config struct {
	DownloadDumpPath string               `json:"dump_path"`
	ServiceToken     string               `json:"token"`
	CurrentWorkspace string               `json:"current_workspace"`
	Data             map[string]Workspace `json:"data"`
}

type ConfigDataService interface {
	ConfigData()
}

func ConfigData(userData map[string]string) Config {
	var workspace string = userData["workspace"]
	var server string = userData["server"]

	currentServer := Server{
		KeyFile:  userData["keyName"],
		ServerId: userData["serverId"],
	}

	serversData := map[string]Server{
		server: currentServer,
	}

	data := Workspace{
		Servers: serversData,
	}

	var configData = map[string]Workspace{}
	configData[workspace] = data

	config := Config{
		ServiceToken:     userData["token"],
		CurrentWorkspace: workspace,
		Data:             configData,
	}

	return config
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

func CreateEnvFile(config Config) {
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

func WriteEnvFile(config Config) {
	workspace := config.CurrentWorkspace

	if !IsEnvFileExist(true) {
		CreateEnvFile(config)
	}

	configFromFile, err := ReadEnvFile()
	if err != nil {
		fmt.Printf(predefined.BuildError("Error: %s"), err)
		return
	}

	serverKeys := maps.Keys(config.Data[workspace].Servers)
	server := serverKeys[0]

	if len(server) > 0 {
		workspaceDataFromFile, v := configFromFile.Data[workspace]
		if v {
			serverDataFromFile, vs := workspaceDataFromFile.Servers[server]
			if vs {
				serverDataFromFile.KeyFile = config.Data[workspace].Servers[server].KeyFile
				workspaceDataFromFile.Servers[server] = serverDataFromFile
			} else {
				workspaceDataFromFile.Servers[server] = Server{
					KeyFile:  config.Data[workspace].Servers[server].KeyFile,
					ServerId: config.Data[workspace].Servers[server].ServerId,
				}
			}
			configFromFile.Data[workspace] = workspaceDataFromFile
		} else {
			configFromFile.Data[workspace] = config.Data[workspace]
		}
	}

	if len(config.DownloadDumpPath) > 0 {
		configFromFile.DownloadDumpPath = config.DownloadDumpPath
	}

	configFromFile.ServiceToken = config.ServiceToken
	configFromFile.CurrentWorkspace = config.CurrentWorkspace

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

func ReadEnvFile() (Config, error) {
	config := Config{}

	if !IsEnvFileExist(true) {
		return config, fmt.Errorf(predefined.BuildError("env file not found. Please use the [login] command to update it"))
	}

	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		return config, fmt.Errorf(predefined.BuildError("cannot get current APP directory: %W"), errDir)
	}

	file, err := os.ReadFile(configDir + "/" + services.EnvFileName)
	if err != nil {
		return config, fmt.Errorf(predefined.BuildError("env file is not readable: %W"), errDir)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, fmt.Errorf(predefined.BuildError("the settings record is not readable: %W"), errDir)
	}

	return config, nil
}
