/*
Copyright © 2024 Bridge Digital
*/
package config

import (
	"fmt"
	"os"
	"strings"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/envfile"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
)

func Execute(dumpPath string) {
	if len(strings.TrimSpace(dumpPath)) == 0 {
		fmt.Println(predefined.BuildError("The download DB dumps path cannot be empty"))
		return
	}

	if !envfile.IsEnvFileExist(false) {
		return
	}

	errDumpPath := createDumpPathDir(dumpPath)
	if errDumpPath != nil {
		fmt.Println(errDumpPath)
		return
	}

	savedConfig, err := envfile.ReadEnvFile()
	if err != nil {
		fmt.Println(predefined.BuildError("Error:"), err)
		return
	}

	savedConfig.DownloadDumpPath = dumpPath
	envfile.WriteEnvFile(savedConfig)

	fmt.Println(predefined.BuildOk("You have set the default path for database dumps"))
}

func createDumpPathDir(dumpPath string) error {
	if err := os.MkdirAll(dumpPath, os.ModePerm); err != nil {
		return fmt.Errorf(predefined.BuildError("can not get entered dump directory: %W"), err)
	}

	return nil
}
