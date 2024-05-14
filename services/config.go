package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	//WebServiceUrl    string = "https://app.dbvisor.pro"
	WebServiceUrl    string = "https://db-manager.bridge2.digital"
	WebServiceApiUrl string = "api"
	AppName          string = "db-manager"
	ServiceToken     string = "SERVICE_TOKEN"
	EnvFileName      string = ".env.json"
)

type Config struct {
	ServiceToken string `json:"token"`
	Workspace    string `json:"workspace"`
	KeyFile      string `json:"key_file"`
}

func EnvFilePath() string {
	return EnvFileName
}

func ConfigData() []Config {
	data := []Config{
		{
			ServiceToken: "",
			Workspace:    "",
			KeyFile:      "",
		},
	}

	return data
}

func IsEnvFileExist() bool {
	var result bool = true

	c_dir, err_dir := currentAppDir()

	if err_dir != nil {
		fmt.Printf("Cannot get current APP directory: %W", err_dir)
		return false
	}

	_, err := os.ReadFile(c_dir + "/.env")
	if err != nil {
		fmt.Printf("Env file not found. Please run: %s install", AppName)
		result = false
	}

	return result
}

func CreateEnvFile(config []Config) {
	c_dir, err_dir := currentAppDir()

	if err_dir != nil {
		fmt.Printf("Cannot get current APP directory: %W", err_dir)
		return
	}

	file, err := os.Create(c_dir + "/" + EnvFileName)
	if err != nil {
		fmt.Println("Cannot create file:", err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		fmt.Println("Cannot write config data to file:", err)
		panic(err)
	}
}

func WriteEnvFile() {

}

func WebServiceAuthUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "login_check")
}

func WebServiceProfileUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "profile")
}

func currentAppDir() (string, error) {
	ex, err := os.Executable()

	dir := filepath.Dir(ex)
	if err != nil {
		//fmt.Errorf("Cannot get current APP directory: %W", err)
		return "", fmt.Errorf("can not get current app directory: %W", err)
	}

	return dir, err
}
