/*
Copyright © 2024 Bridge Digital
*/
package download

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/processes/download/helper"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/encrypter"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/envfile"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/request"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/response"
)

const (
	DefaultDumpDBName string = "backup"
	DefaultDumpDBExt  string = ".sql"
)

func Execute(dbUid, dumpUid string) {
	savedConfigData, err := envfile.ReadEnvFile()
	if err != nil {
		fmt.Println(predefined.BuildError("Error:"), err)
		return
	}

	var (
		selectedToken, selectedKeyPubName, selectedWorkspace, serverId, dbName string
	)

	selectedWorkspace = savedConfigData.CurrentWorkspace
	selectedToken = savedConfigData.ServiceToken

	selectedDataByWorkspace, ok := savedConfigData.Data[selectedWorkspace]
	if !ok {
		fmt.Println(predefined.BuildError("There is no record of a saved workspace"))
		return
	}

	if dbUid == "" {
		dbUid = helper.GetDbUid(selectedToken)
	}

	if len(dbUid) > 0 {
		serverData := helper.GetServerData(dbUid, selectedToken)
		if len(serverData) == 0 {
			fmt.Println(predefined.BuildError("Failed to retrieve selected database data"))
			return
		}

		serverId = serverData["serverId"]
		dbName = serverData["dbName"]

		selectedKeyPubName = helper.GetServerPubKey(selectedDataByWorkspace.Servers, serverId)

		if len(selectedKeyPubName) == 0 {
			fmt.Println(predefined.BuildError("There is no saved public key for the server. Use the [login] command to create it"))
			return
		}
	}

	if dbUid == "" {
		fmt.Println(predefined.BuildError("Failed to get DB UID"))
		return
	}

	if dumpUid == "" {
		dumpUid, err = helper.GetDumpUid(dbUid, selectedToken, selectedWorkspace)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		if len(dumpUid) > 0 && dbUid == "" {
			dbUid = helper.GetDbUidByDump(dumpUid, selectedToken)
			if len(dumpUid) == 0 {
				fmt.Println(predefined.BuildError("Failed to get DB UID"))
				return
			}

			serverData := helper.GetServerData(dbUid, selectedToken)
			if serverData == nil {
				fmt.Println(predefined.BuildError("Failed to retrieve selected database data"))
				return
			}

			serverId = serverData["serverId"]
			dbName = serverData["dbName"]

			selectedKeyPubName = helper.GetServerPubKey(selectedDataByWorkspace.Servers, serverId)
			if len(selectedKeyPubName) == 0 {
				fmt.Println(predefined.BuildError("There is no saved public key for the server. Use the [login] command to create it"))
				return
			}
		}
	}

	if dumpUid == "" {
		fmt.Println(predefined.BuildError("Failed to get dump UID"))
		return
	}

	var defaultDumpPath string = ""

	if len(savedConfigData.DownloadDumpPath) > 0 {
		defaultDumpPath = savedConfigData.DownloadDumpPath
	}

	currentTime := time.Now()
	dumpDbData := map[string]string{
		"dbuuid":   dbUid,
		"dumpuuid": dumpUid,
		"dumpname": DefaultDumpDBName + "_" + selectedWorkspace + "_" + dbName + "_" + currentTime.Format("2000-01-01 00:00:00"),
		"dumppath": defaultDumpPath,
	}

	encryptedData := encrypter.EncryptData(dumpDbData, selectedKeyPubName)
	if encryptedData == nil {
		return
	}

	download(dumpDbData, encryptedData, selectedToken)
}

func download(dumpDbData map[string]string, encryptedData []byte, token string) {
	var saveDumpName string = dumpDbData["dumpname"]
	if len(saveDumpName) == 0 {
		return
	}

	fmt.Println(predefined.BuildOk("Downloading..."))

	var requestUrl string = services.WebServiceDownLoadLinkUrl() + "/" + dumpDbData["dbuuid"] + "/" + dumpDbData["dumpuuid"]

	data, err := request.CreateGetRequest(requestUrl, &token)
	if err != nil {
		fmt.Println(predefined.BuildError("Error:"), err)
		return
	}

	var (
		dumLinkData map[string]string
		link        string = ""
		ok          bool
	)

	configErr := json.Unmarshal(data, &dumLinkData)
	if configErr != nil {
		err := response.WrongResponseObserver(data)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if len(dumLinkData) > 0 {
		link, ok = dumLinkData["link"]
		if !ok {
			link = ""
		}
	}

	if len(link) > 0 {
		configDir, errDir := services.CurrentAppDir()
		if errDir != nil {
			fmt.Printf(predefined.BuildError("Cannot get current APP directory: %W.\n"), errDir)
			return
		}

		var saveDumpPath string = ""

		if len(dumpDbData["dumpname"]) > 0 {
			saveDumpPath = dumpDbData["dumpname"]
			saveDumpPath = strings.TrimRight(saveDumpPath, "/")
			saveDumpPath += "/"
		} else {
			saveDumpPath, err = helper.DefaultDumpPath()
			if err != nil {
				fmt.Println(predefined.BuildError("Error:"), err)
				return
			}

			if len(strings.TrimSpace(saveDumpPath)) == 0 {
				saveDumpPath = configDir + "/"
			} else {
				saveDumpPath = strings.TrimRight(saveDumpPath, "/")
				saveDumpPath += "/"
			}
		}

		saveDumpName = strings.TrimSuffix(saveDumpName, DefaultDumpDBExt)

		var fullFilePath string = saveDumpPath + saveDumpName + DefaultDumpDBExt

		downloadFile(link, encryptedData, fullFilePath)
	} else {
		fmt.Println(predefined.BuildWarning("The download dump URL is empty."))
		return
	}
}

func downloadFile(link string, encryptedData []byte, fullFilePath string) {
	req, err := http.NewRequest("POST", link, bytes.NewBuffer(encryptedData))
	if err != nil {
		fmt.Println(predefined.BuildError("Something wrong with POST request data:"), err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(predefined.BuildError("Invalid credentials:"), err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			fmt.Printf(predefined.BuildError("bad status: %s. Your token has expired. Use the login command to update it \n"), resp.Status)
			return
		} else {
			fmt.Printf(predefined.BuildError("Bad status: %s \n"), resp.Status)
			return
		}
	}

	defer resp.Body.Close()

	outFile, err := os.Create(fullFilePath)
	if err != nil {
		fmt.Printf(predefined.BuildError("Error creating file: %v"), err)
		return
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		fmt.Printf(predefined.BuildError("Error copying response body to file: %v"), err)
		return
	}
}
