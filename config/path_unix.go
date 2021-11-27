//go:build linux

package config

import (
	"errors"
	"os"
)

func getConfigPath() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dirname += "/.config/sshelter"

	err = os.MkdirAll(dirname, 0700)
	if err != nil {
		return "", err
	}

	dirname += "/sshelter.yml"

	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		return "", errors.New("config file not found")
	}

	return dirname, err
}

func GetSshConfigFile() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dirname += "/.ssh/config"

	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		return "", errors.New("ssh config file not found")
	}

	return dirname, err
}
