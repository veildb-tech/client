/*
Copyright © 2024 Bridge Digital
*/
package request

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
)

func CreatePostRequest(data []byte, url string, token *string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf(predefined.BuildError("something wrong with POST request data: %w"), err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	var errMsg string = ""

	if token != nil {
		req.Header.Set("Authorization", "Bearer "+*token)
		errMsg = predefined.BuildWarning("Your token has expired. Use the [login] command to update it")
	} else {
		errMsg = predefined.BuildError("Invalid credentials")
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(predefined.BuildError("invalid credentials: %w"), err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf(predefined.BuildError("bad status: %s. %s"), resp.Status, errMsg)
		} else {
			return nil, fmt.Errorf(predefined.BuildError("bad status: %s"), resp.Status)
		}
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(predefined.BuildError("error: %w"), err)
	}

	return result, nil
}

func CreateGetRequest(url string, token *string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf(predefined.BuildError("something wrong with GET request data: %w"), err)
	}

	req.Header.Set("Accept", "application/json")

	if token != nil {
		req.Header.Set("Authorization", "Bearer "+*token)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(predefined.BuildError("invalid token: %w"), err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf(predefined.BuildError("bad status: %s. Your token has expired. Use the [login] command to update it"), resp.Status)
		} else {
			return nil, fmt.Errorf(predefined.BuildError("bad status: %s"), resp.Status)
		}
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(predefined.BuildError("error: %w"), err)
	}

	return result, nil
}
