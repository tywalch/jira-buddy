package configuration

import (
	"encoding/json"
	"fmt"
	"github.com/tywalch/jira-buddy/app/prompt"
	"os"
	"path/filepath"
	"strings"
)

const (
	ConfigDir     = ".jirabuddy"
	ConfigFile    = "config.json"
	ConfigVersion = 1
)

type Configuration struct {
	AccountId     string `json:"accountId"`
	Username      string `json:"username"`
	Password      string `json:"apiKey"`
	APIUrl        string `json:"apiUrl"`
	ConfigVersion int    `json:"configVersion"`
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getConfigurationFromFile(path string) (Configuration, error) {
	var configuration Configuration
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("\nread file error: %v\n", err)
		return configuration, err
	}

	err = json.Unmarshal(file, &configuration)

	if err != nil {
		return configuration, err
	}

	return configuration, nil
}

func getConfigurationFromPrompt(reader *prompt.PromptReader) (Configuration, error) {
	username, err := reader.GetString("Email Address")
	if err != nil {
		return Configuration{}, err
	}

	accountId, err := reader.GetString("AccountId (your unique jira userId)")
	if err != nil {
		return Configuration{}, err
	}

	password, err := reader.GetString("API Key")
	if err != nil {
		return Configuration{}, err
	}

	apiUrl, err := reader.GetString("Jira URL")
	if err != nil {
		return Configuration{}, err
	}

	configuration := Configuration{
		Username:      strings.TrimSpace(username),
		AccountId:     strings.TrimSpace(accountId),
		Password:      strings.TrimSpace(password),
		APIUrl:        strings.TrimSpace(apiUrl),
		ConfigVersion: ConfigVersion,
	}

	return configuration, nil
}

func writeConfiguration(path string, configuration Configuration) error {
	outputJson, err := json.MarshalIndent(configuration, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, outputJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func fetchConfiguration(reader *prompt.PromptReader, dir string, fileName string) (Configuration, error) {
	homeDir, err := os.UserHomeDir()
	configuration := Configuration{}
	if err != nil {
		return configuration, err
	}

	directoryPath := filepath.Join(homeDir, ConfigDir)

	err = os.Mkdir(directoryPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("\ndirectory error: %v\n", err)
		return configuration, err
	}

	configurationPath := filepath.Join(dir, fileName)

	fileExists, err := fileExists(configurationPath)
	if err != nil {
		return configuration, err
	}

	if fileExists {
		configuration, err = getConfigurationFromFile(configurationPath)
		if err != nil {
			return configuration, err
		}
	} else {
		fmt.Print("No configuration found. Please enter your Jira credentials.\n")
		configuration, err = getConfigurationFromPrompt(reader)
		if err != nil {
			return configuration, err
		}
	}

	if configuration.ConfigVersion != ConfigVersion {
		fmt.Print("\nConfiguration update required\n")
		configuration, err = getConfigurationFromPrompt(reader)
		if err != nil {
			return configuration, err
		}
	}

	err = writeConfiguration(configurationPath, configuration)
	if err != nil {
		return configuration, err
	}

	return configuration, nil
}

func getConfigDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dirPath := filepath.Join(homeDir, ConfigDir)

	err = os.Mkdir(dirPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return "", err
	}

	return dirPath, nil
}

func GetConfiguration(reader *prompt.PromptReader) (Configuration, error) {
	dirPath, err := getConfigDirectory()
	configuration := Configuration{}
	if err != nil {
		return configuration, err
	}

	configuration, err = fetchConfiguration(reader, dirPath, ConfigFile)
	if err != nil {
		return configuration, err
	}

	return configuration, nil
}
