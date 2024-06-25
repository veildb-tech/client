/*
Copyright © 2024 Bridge Digital
*/
package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
)

const DestinationPath = "/usr/local/bin/"

func Execute() {
	_, err := os.Stat(DestinationPath + services.AppName)
	if err == nil {
		fmt.Println(predefined.BuildOk("The application has been installed successfully"))
		return
	}

	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		fmt.Printf(predefined.BuildError("Cannot get current APP directory: %W.\n"), errDir)
		return
	}

	var sourcePath string = configDir + "/bin/" + services.AppName

	_, errApp := os.Stat(sourcePath)
	if errApp != nil {
		fmt.Println(predefined.BuildError("Executable file missing"))
		return
	}

	/* if envfile.IsEnvFileExist(true) {
		fmt.Println(predefined.BuildOk("Application has already installed"))
		return
	} else { */
	//createLink()
	testw(sourcePath)
	//}
}

func createLink() {
	configDir, errDir := services.CurrentAppDir()
	if errDir != nil {
		fmt.Printf(predefined.BuildError("Cannot get current APP directory: %W.\n"), errDir)
		return
	}

	command := fmt.Sprintf("export PATH=\"%s/bin:$PATH\" \n", configDir)
	bashProfileCandidates := []string{".bashrc", ".bash_profile"}
	homeDir := os.Getenv("HOME")

	for _, bashProfileCandidate := range bashProfileCandidates {
		candidateFilePath := filepath.Join(homeDir, bashProfileCandidate)

		if _, err := os.Stat(candidateFilePath); err == nil {
			file, err := os.OpenFile(candidateFilePath, os.O_APPEND|os.O_WRONLY, 0644)

			if err != nil {
				fmt.Printf(predefined.BuildError("Error opening file %s: %v"), candidateFilePath, err)
				continue
			}

			defer file.Close()

			if _, err := fmt.Fprintln(file, command); err != nil {
				fmt.Printf(predefined.BuildError("Error writing to file %s: %v"), candidateFilePath, err)
				return
			}
		}
	}
}

func testw(sourcePath string) {
	// Define the source and destination paths
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		fmt.Println("HOME environment variable is not set.")
		return
	}
	//source := configDir + "/bin/" + services.AppName
	//destination := "/usr/local/bin/db-manager"

	// Create the command to execute
	cmd := exec.Command("sudo", "ln", "-s", sourcePath, DestinationPath)

	// Set the environment variables for the command
	cmd.Env = os.Environ()

	// Run the command and capture any errors
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		return
	}

	// Print the output of the command
	//fmt.Printf("Command executed successfully: %s\n", string(output))
}
