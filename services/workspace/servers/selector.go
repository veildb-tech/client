/*
Copyright © 2024 Bridge Digital
*/
package servers

import (
	"github.com/AlecAivazis/survey/v2"
)

func Server(servers map[string]string) string {
	allServers := []string{}

	for _, server := range servers {
		allServers = append(allServers, server)
	}

	var selectedServer int

	prompt := &survey.Select{
		Message: "Select server:",
		Options: allServers,
	}

	survey.AskOne(prompt, &selectedServer)

	return allServers[selectedServer]
}
