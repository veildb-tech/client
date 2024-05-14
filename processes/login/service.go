package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

func Execute(cmd *cobra.Command) string {

	credentials := map[string]string{
		"username": "",
		"password": "",
	}

	var token, username, workspace string

	//reader := bufio.NewReader(os.Stdin)

USERNAME:
	fmt.Println("Username: ")
	fmt.Scanln(&username)
	//username, _ := reader.ReadString('\n')

	if len(strings.TrimSpace(username)) == 0 {
		fmt.Println("The Username cannot be empty")
		goto USERNAME
	} else {
		credentials["username"] = username
	}

PASSWORD:
	fmt.Println("Password: ")

	password, _ := gopass.GetPasswdMasked()

	if len(strings.TrimSpace(string(password))) == 0 {
		fmt.Println("The Password cannot be empty")
		goto PASSWORD
	} else {
		credentials["password"] = string(password)
	}

	token = jwtToken(credentials)
	workspace = getWorkspace(token)

	fmt.Println(workspace)

	if len(token) == 0 || len(workspace) == 0 {
		return ""
	}

	if !services.IsEnvFileExist() {
		services.CreateEnvFile(services.ConfigData())
	} else {
		//services.WriteEnvFile(services.ConfigData())
	}

	return "You logged in successfully"
}

// Get token from server
func jwtToken(credentials map[string]string) string {
	enc_creds, err := json.Marshal(credentials)
	if err != nil {
		fmt.Println("Error encoding to json:", err)
		return ""
	}

	req, err := http.NewRequest("POST", services.WebServiceAuthUrl(), bytes.NewBuffer(enc_creds))
	if err != nil {
		fmt.Println("Something wrong with POST request data:", err)
		return ""
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Invalid credentials:", err)
		return ""
	}

	if resp == nil {
		return fmt.Sprint(http.StatusBadRequest)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	var c_data map[string]string

	c_err := json.Unmarshal(data, &c_data)
	if c_err != nil {
		fmt.Println("Error decoding from json:", c_err)
		return ""
	}

	if len(c_data["token"]) > 0 {
		return c_data["token"]
	}

	return ""
}

func getWorkspace(token string) string {
	req, err := http.NewRequest("GET", services.WebServiceProfileUrl(), nil)
	if err != nil {
		fmt.Println("Something wrong with POST request data:", err)
		return ""
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Invalid token:", err)
		return ""
	}

	if resp == nil {
		return fmt.Sprint(http.StatusBadRequest)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	//fmt.Print("returned: ", string(data))
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	type Data struct {
		Identifier string   `json:"identifier"`
		Workspaces []string `json:"workspaces"`
	}

	// Unmarshal the JSON data into the struct
	var w_data Data
	all_workspaces := map[int]string{}

	w_err := json.Unmarshal([]byte(data), &w_data)
	if w_err != nil {
		fmt.Println("Error decoding from json:", w_err)
		return ""
	}

	if len(w_data.Workspaces) > 0 {
	WORKSPACE:
		fmt.Println("Select workspace: ")

		for k, workspace := range w_data.Workspaces {
			fmt.Println(k+1, ": ", workspace)
			all_workspaces[k+1] = workspace
		}

		var selected_workspace int

		fmt.Scanln(&selected_workspace)

		if selected_workspace > 0 {
			return all_workspaces[selected_workspace]
		} else {
			fmt.Println("You must select workspace")
			goto WORKSPACE
		}
	} else {
		fmt.Println("You don't assigned to any workspace")

	}

	return ""
}
