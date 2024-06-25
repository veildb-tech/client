/*
Copyright © 2024 Bridge Digital
*/
package servers

import (
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/util"
	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/exp/maps"
)

func Server(servers map[string]string) (serverId string, serverName string) {
	var (
		selectedServerIndex int
		serversValues       []string
	)

	serversValues = maps.Values(servers)

	prompt := &survey.Select{
		Message: "Select server:",
		Options: serversValues,
	}

	survey.AskOne(prompt, &selectedServerIndex)

	selectedServerId := util.MapKeyByValue(servers, serversValues[selectedServerIndex])

	return selectedServerId, serversValues[selectedServerIndex]
}
