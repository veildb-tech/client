//go:build dev
// +build dev

/*
Copyright © 2024 Bridge Digital
*/
package services

import (
	"fmt"
	"os"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
)

const (
	WebServiceUrl    string = "https://db-manager.bridge2.digital"
	WebServiceApiUrl string = "api"
	AppName          string = "dbvisor"
	EnvFileName      string = ".env.json"
	PubKeyExt        string = ".pem"
)

// API login_check url
func WebServiceAuthUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "login_check")
}

// API profile url
func WebServiceProfileUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "profile")
}

// API database list url
func WebServiceDatabaseListUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "databases")
}

// API database dump url
func WebServiceDatabaseDumpUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "database_dumps")
}

// API database download link url
func WebServiceDownLoadLinkUrl() string {
	return fmt.Sprintf("%s/%s/%s", WebServiceUrl, WebServiceApiUrl, "get_download_link")
}

func CurrentAppDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf(predefined.BuildError("can not get current HOME user directory: %W"), err)
	}

	var appDir string = homeDir + "/.dbvisor"

	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		return "", fmt.Errorf(predefined.BuildError("can not get current app directory: %W"), err)
	}

	return appDir, err
}
