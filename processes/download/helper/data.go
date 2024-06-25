/*
Copyright © 2024 Bridge Digital
*/
package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/envfile"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/request"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/response"
	"github.com/AlecAivazis/survey/v2"
)

const (
	StatusReady              string = "ready"
	StatusReadyWithError     string = "ready_with_error"
	DBStatus                 string = "enabled"
	DefaultDBDumpsFolderName        = "DBVisor DB Dumps"
)

func GetServerData(dbUid string, token string) map[string]string {
	data, err := request.CreateGetRequest(services.WebServiceDatabaseListUrl()+"/"+dbUid, &token)
	if err != nil {
		fmt.Println(predefined.BuildError("Error:"), err)
		return nil
	}

	type Data struct {
		Server string `json:"server"`
		Name   string `json:"name"`
	}

	var (
		dbData Data
	)

	dbServerData := map[string]string{}

	dbErr := json.Unmarshal([]byte(data), &dbData)
	if dbErr != nil {
		err := response.WrongResponseObserver(data)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	if len(dbData.Server) > 0 {
		dbServerId := dbData.Server
		_, dbServerId, _ = strings.Cut(dbServerId, "/api/servers/")
		dbServerData["serverId"] = dbServerId
		dbServerData["dbName"] = dbData.Name
	}

	return dbServerData
}

func GetServerPubKey(servers map[string]envfile.Server, serverId string) string {
	if len(servers) == 0 {
		fmt.Println(predefined.BuildError("There are no saved server data records"))
		return ""
	}

	var result = ""

	for _, s := range servers {
		if s.ServerId == serverId {
			result = s.KeyFile
			break
		}
	}

	return result
}

func GetDbUidByDump(dumpUid string, token string) string {
	data, err := request.CreateGetRequest(services.WebServiceDatabaseDumpUrl()+"/"+dumpUid, &token)
	if err != nil {
		fmt.Println(predefined.BuildError("Error:"), err)
		return ""
	}

	type Db struct {
		Uid string `json:"uid"`
	}

	type Data struct {
		Db Db `json:"db"`
	}

	var (
		dumpData Data
		dbUid    string = ""
	)

	dbErr := json.Unmarshal([]byte(data), &dumpData)
	if dbErr != nil {
		err := response.WrongResponseObserver(data)
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}

	if len(dumpData.Db.Uid) > 0 {
		dbUid = dumpData.Db.Uid
	}

	return dbUid
}

func DefaultDumpPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf(predefined.BuildError("can not get current HOME user directory: %W"), err)
	}

	var appDir string = homeDir + "/" + DefaultDBDumpsFolderName

	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		return "", fmt.Errorf(predefined.BuildError("can not get default DB dumps directory: %W"), err)
	}

	return appDir, err
}

func GetDumpUid(dbUid string, token string, selectedWorkspace string) (string, error) {
	var (
		//Uncomment if you need to load workspaces and not use them from the env file.
		//selectedWorkspace string = workspace.Workspace(token)
		requestUrl string = services.WebServiceDatabaseDumpUrl() + "?db.uid=" + dbUid + "&workspace=" + selectedWorkspace +
			"&status[]=" + StatusReady + "&status[]=" + StatusReadyWithError
	)

	data, err := request.CreateGetRequest(requestUrl, &token)
	if err != nil {
		return "", fmt.Errorf(predefined.BuildError("Error:"), err)
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
		err := response.WrongResponseObserver(data)
		if err != nil {
			return "", err
		}
	}

	readyDumps := map[string]string{}

	if len(dumps) > 0 {
		for _, dump := range dumps {
			dumpName := dump.Uuid + "[" + dump.Date + "]"
			allDumps = append(allDumps, dumpName)
			readyDumps[dumpName] = dump.Uuid
		}

		sort.Strings(allDumps)

		var selectedDb int

		prompt := &survey.Select{
			Message: "Please select dump to process with:",
			Options: allDumps,
		}

		survey.AskOne(prompt, &selectedDb)

		return readyDumps[allDumps[selectedDb]], nil
	} else {
		return "", fmt.Errorf(predefined.BuildWarning("Not found active dumps for selected DB"))
	}
}

func GetDbUid(token string) string {
	data, err := request.CreateGetRequest(services.WebServiceDatabaseListUrl(), &token)
	if err != nil {
		fmt.Println(predefined.BuildError("Error:"), err)
		return ""
	}

	type Data struct {
		Name   string `json:"name"`
		Uid    string `json:"uid"`
		Status string `json:"status"`
	}

	var (
		dbData        []Data
		allDbDataName []string
	)

	dbErr := json.Unmarshal([]byte(data), &dbData)
	if dbErr != nil {
		err := response.WrongResponseObserver(data)
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}

	activeDBs := map[string]string{}

	if len(dbData) > 0 {
		for _, db := range dbData {
			if db.Status == DBStatus {
				allDbDataName = append(allDbDataName, db.Name)
				activeDBs[db.Name] = db.Uid
			}
		}

		if len(allDbDataName) == 0 {
			fmt.Println(predefined.BuildWarning("Not found active databases"))
			return ""
		}

		sort.Strings(allDbDataName)

		var selectedDb int

		prompt := &survey.Select{
			Message: "Please select database to process with:",
			Options: allDbDataName,
		}

		survey.AskOne(prompt, &selectedDb)

		return activeDBs[allDbDataName[selectedDb]]
	} else {
		fmt.Println(predefined.BuildWarning("Not found active databases"))
	}

	return ""
}
