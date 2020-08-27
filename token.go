package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not detect your home directory - %s", err)
	}

	return path.Join(home, ".protopkg"), nil
}

func setToken(token string) error {
	config, err := configPath()
	if err != nil {
		return err
	}

	contents := fmt.Sprintf("token:%s", token)
	return ioutil.WriteFile(config, []byte(contents), 0644)
}

func fallbackToken() string {
	if token := os.Getenv("GITHUB_PACKAGE_PULL_TOKEN"); token != "" {
		return token
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token
	}

	return ""
}

func getToken() string {
	config, err := configPath()
	if err != nil {
		return fallbackToken()
	}

	contents, err := ioutil.ReadFile(config)
	if err != nil {
		return fallbackToken()
	}

	if string(contents) != "" {
		parts := strings.Split(string(contents), ":")
		if len(parts) == 2 {
			return parts[1]
		}
	}

	return fallbackToken()
}
