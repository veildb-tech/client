//go:build dev
// +build dev

/*
Copyright © 2024 Bridge Digital
*/
package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dbvisor-pro/client/services/predefined"
)

const (
	DefaultWebServiceUrl string = "https://db-manager.bridge2.digital"
	WebServiceApiUrl     string = "api"
	AppName              string = "veildb"
	EnvFileName          string = ".env.json"
	PubKeyExt            string = ".pem"
)

// configCache stores the service URL to avoid reading config file repeatedly
var configCache struct {
	loaded     bool
	serviceUrl string
}

// GetServiceUrl returns the service URL from config or default
func GetServiceUrl() string {
	if configCache.loaded {
		return configCache.serviceUrl
	}

	configDir, err := CurrentAppDir()
	if err != nil {
		configCache.loaded = true
		configCache.serviceUrl = DefaultWebServiceUrl
		return DefaultWebServiceUrl
	}

	file, err := os.ReadFile(configDir + "/" + EnvFileName)
	if err != nil {
		configCache.loaded = true
		configCache.serviceUrl = DefaultWebServiceUrl
		return DefaultWebServiceUrl
	}

	var config struct {
		ServiceUrl string `json:"service_url"`
	}

	if err := json.Unmarshal(file, &config); err != nil || config.ServiceUrl == "" {
		configCache.loaded = true
		configCache.serviceUrl = DefaultWebServiceUrl
		return DefaultWebServiceUrl
	}

	configCache.loaded = true
	configCache.serviceUrl = config.ServiceUrl
	return config.ServiceUrl
}

// ResetConfigCache resets the config cache (call after updating config)
func ResetConfigCache() {
	configCache.loaded = false
	configCache.serviceUrl = ""
}

// API login_check url
func WebServiceAuthUrl() string {
	return fmt.Sprintf("%s/%s/%s", GetServiceUrl(), WebServiceApiUrl, "login_check")
}

// API profile url
func WebServiceProfileUrl() string {
	return fmt.Sprintf("%s/%s/%s", GetServiceUrl(), WebServiceApiUrl, "profile")
}

// API database list url
func WebServiceDatabaseListUrl() string {
	return fmt.Sprintf("%s/%s/%s", GetServiceUrl(), WebServiceApiUrl, "databases")
}

// API database dump url
func WebServiceDatabaseDumpUrl() string {
	return fmt.Sprintf("%s/%s/%s", GetServiceUrl(), WebServiceApiUrl, "database_dumps")
}

// API database download link url
func WebServiceDownLoadLinkUrl() string {
	return fmt.Sprintf("%s/%s/%s", GetServiceUrl(), WebServiceApiUrl, "get_download_link")
}

func CurrentAppDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf(predefined.BuildError("can not get current HOME user directory: %w"), err)
	}

	var appDir string = homeDir + "/.veildb"

	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		return "", fmt.Errorf(predefined.BuildError("can not get current app directory: %w"), err)
	}

	return appDir, err
}
