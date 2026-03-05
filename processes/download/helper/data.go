/*
Copyright © 2024 Bridge Digital
*/
package helper

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dbvisor-pro/client/services"
	"github.com/dbvisor-pro/client/services/envfile"
	"github.com/dbvisor-pro/client/services/predefined"
	"github.com/dbvisor-pro/client/services/request"
	"github.com/dbvisor-pro/client/services/response"
	"github.com/AlecAivazis/survey/v2"
)

const (
	StatusReady          string = "ready"
	StatusReadyWithError string = "ready_with_error"
	DBStatus             string = "enabled"
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

func GetLatestDumpUid(dbUid string, token string, selectedWorkspace string) (string, error) {
	var (
		requestUrl string = services.WebServiceDatabaseDumpUrl() + "?db.uid=" + dbUid + "&workspace=" + selectedWorkspace +
			"&status[]=" + StatusReady + "&status[]=" + StatusReadyWithError
	)

	data, err := request.CreateGetRequest(requestUrl, &token)
	if err != nil {
		return "", fmt.Errorf(predefined.BuildError("Error:"), err)
	}

	type Data struct {
		Uuid string `json:"uuid"`
	}

	var dumps []Data

	dbErr := json.Unmarshal([]byte(data), &dumps)
	if dbErr != nil {
		err := response.WrongResponseObserver(data)
		if err != nil {
			return "", err
		}
	}

	if len(dumps) == 0 {
		return "", fmt.Errorf(predefined.BuildWarning("Not found active dumps for selected DB"))
	}

	return dumps[0].Uuid, nil
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
		Uuid      string `json:"uuid"`
		Filename  string `json:"filename"`
		CreatedAt string `json:"created_at"`
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
			label := formatDumpLabel(dump.CreatedAt, dump.Filename)
			allDumps = append(allDumps, label)
			readyDumps[label] = dump.Uuid
		}

		sort.Sort(sort.Reverse(sort.StringSlice(allDumps)))

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

func formatDumpLabel(createdAt, filename string) string {
	t, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05-07:00", createdAt)
	}

	var dateStr string
	if err == nil {
		dateStr = t.Format("2006-01-02 15:04")
	} else {
		dateStr = createdAt
	}

	if filename != "" {
		return dateStr + " (" + filename + ")"
	}
	return dateStr
}
