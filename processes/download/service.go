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
	"sort"
	"strings"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/encrypter"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/envfile"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/request"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/response"
	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/exp/maps"
)

const (
	DefaultDumpDBName    string = "backup"
	DefaultDumpDBExt     string = ".sql"
	StatusReady          string = "ready"
	StatusReadyWithError string = "ready_with_error"
)

func Execute(dbUid, dumpUid string) {
	savedWorkspaces, err := envfile.ReadEnvFile()
	if err != nil {
		fmt.Println(predefined.BuildError("Error:"), err)
		return
	}

	var (
		selectedWorkspaceIndex, selectedServerIndex int
		savedWorkspacesKeys, savedServersKeys       []string
	)

	savedWorkspacesKeys = maps.Keys(savedWorkspaces)

	if len(savedWorkspacesKeys) > 1 {
		prompt := &survey.Select{
			Message: "Select one of your saved workspaces:",
			Options: savedWorkspacesKeys,
		}

		survey.AskOne(prompt, &selectedWorkspaceIndex)
	} else {
		selectedWorkspaceIndex = 0
		fmt.Println(predefined.BuildAnsw("Your saved workspaces: ", savedWorkspacesKeys[selectedWorkspaceIndex]))
	}

	selectedToken := savedWorkspaces[savedWorkspacesKeys[selectedWorkspaceIndex]].ServiceToken

	if dbUid == "" {
		dbUid = getDbUid(selectedToken)
	}

	if dbUid == "" {
		return
	}

	if dumpUid == "" {
		dumpUid = getDumpUid(dbUid, selectedToken, savedWorkspacesKeys[selectedWorkspaceIndex])
	}

	if dumpUid == "" {
		return
	}

	dumpDbData := map[string]string{
		"dbuuid":   dbUid,
		"dumpuuid": dumpUid,
	}

	savedServers := savedWorkspaces[savedWorkspacesKeys[selectedWorkspaceIndex]].Servers
	savedServersKeys = maps.Keys(savedServers)

	if len(savedServersKeys) > 1 {
		promptS := &survey.Select{
			Message: "Select one of your saved servers:",
			Options: savedServersKeys,
		}

		survey.AskOne(promptS, &selectedServerIndex)
	} else {
		selectedServerIndex = 0
		fmt.Println(predefined.BuildAnsw("Your saved server: ", savedServersKeys[selectedServerIndex]))
	}

	selectedKeyPubName := savedServers[savedServersKeys[selectedServerIndex]].KeyFile

	encryptedData := encrypter.EncryptData(dumpDbData, selectedKeyPubName)
	if encryptedData == nil {
		return
	}

	download(dumpDbData, encryptedData, selectedToken)
}

func getDbUid(token string) string {
	data, err := request.CreateGetRequest(services.WebServiceDatabaseListUrl(), &token)
	if err != nil {
		fmt.Println(predefined.BuildError("Error:"), err)
		return ""
	}

	type Data struct {
		Name string `json:"name"`
		Uid  string `json:"uid"`
	}

	var (
		dbData        []Data
		allDbDataName []string
	)

	dbErr := json.Unmarshal([]byte(data), &dbData)
	if dbErr != nil {
		response.WrongResponseObserver(data)
		return ""
	}

	if len(dbData) > 0 {
		for _, uid := range dbData {
			allDbDataName = append(allDbDataName, uid.Name)
		}

		sort.Strings(allDbDataName)

		var selectedDb int

		prompt := &survey.Select{
			Message: "Please select database to process with:",
			Options: allDbDataName,
		}

		survey.AskOne(prompt, &selectedDb)

		return dbData[selectedDb].Uid
	} else {
		fmt.Println(predefined.BuildWarning("Not found active databases"))
	}

	return ""
}

func getDumpUid(dbUid string, token string, selectedWorkspace string) string {
	var (
		//Uncomment if you need to load workspaces and not use them from the env file.
		//selectedWorkspace string = workspace.Workspace(token)
		requestUrl string = services.WebServiceDatabaseDumpUrl() + "?db.uid=" + dbUid + "&workspace=" + selectedWorkspace +
			"&status[]=" + StatusReady + "&status[]=" + StatusReadyWithError
	)

	data, err := request.CreateGetRequest(requestUrl, &token)
	if err != nil {
		fmt.Println(predefined.BuildError("Error:"), err)
		return ""
	}

	type Data struct {
		Uuid string `json:"uuid"`
		Date string `json:"updated_at"`
	}

	var (
		dumps    []Data
		allDumps []string
	)

	dbErr := json.Unmarshal([]byte(data), &dumps)
	if dbErr != nil {
		response.WrongResponseObserver(data)
		return ""
	}

	if len(dumps) > 0 {
		for _, uid := range dumps {
			allDumps = append(allDumps, uid.Uuid+"["+uid.Date+"]")
		}

		sort.Strings(allDumps)

		var selectedDb int

		prompt := &survey.Select{
			Message: "Please select dump to process with:",
			Options: allDumps,
		}

		survey.AskOne(prompt, &selectedDb)

		return dumps[selectedDb].Uuid
	} else {
		fmt.Println(predefined.BuildWarning("Not found active dumps for selected DB"))
	}

	return ""
}

func download(dumpDbData map[string]string, encryptedData []byte, token string) {
	var saveDumpPath, saveDumpName string

	prompt := &survey.Input{
		Message: "Specify path to save dump:",
		Help:    "By default, the save directory is the location directory of the console application",
	}

	survey.AskOne(prompt, &saveDumpPath)

	prompt = &survey.Input{
		Message: "Specify filename:",
		Help:    "Default DB name is " + DefaultDumpDBName + DefaultDumpDBExt,
		Default: DefaultDumpDBName,
	}

	survey.AskOne(prompt, &saveDumpName)

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
		response.WrongResponseObserver(data)
		return
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

		if len(strings.TrimSpace(saveDumpPath)) == 0 {
			saveDumpPath = configDir + "/"
		} else {
			saveDumpPath = strings.TrimRight(saveDumpPath, "/")
			saveDumpPath += "/"
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
