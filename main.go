package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"golang.design/x/clipboard"
	"os"
	"path/filepath"
	"strings"
)

type Configuration struct {
	AccountId     string `json:"accountId"`
	Username      string `json:"username"`
	Password      string `json:"apiKey"`
	APIUrl        string `json:"apiUrl"`
	ConfigVersion int    `json:"configVersion"`
}

const (
	// IssueId - BUG = "10004", FEATURE = "10003", TASK = "10002", STORY = "10001", EPIC = "10000"
	IssueId       = "10002"
	ConfigDir     = ".jirabuddy"
	ConfigFile    = "config.json"
	ConfigVersion = 1
)

func prompt(reader *bufio.Reader, label string) (string, error) {
	fmt.Printf("%s: ", label)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	if len(input) > 0 {
		return input, nil
	}

	return "", fmt.Errorf("Invalid input received for prompt %q", label)
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

func getConfigurationFromPrompt(reader *bufio.Reader) (Configuration, error) {
	username, err := prompt(reader, "Email Address")
	if err != nil {
		return Configuration{}, err
	}

	accountId, err := prompt(reader, "AccountId (your unique jira userId)")
	if err != nil {
		return Configuration{}, err
	}

	password, err := prompt(reader, "API Key")
	if err != nil {
		return Configuration{}, err
	}

	apiUrl, err := prompt(reader, "Jira URL")
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

func fetchConfiguration(reader *bufio.Reader, dir string, fileName string) (Configuration, error) {
	homeDir, err := os.UserHomeDir()
	configuration := Configuration{}
	if err != nil {
		return configuration, err
	}

	directoryPath := filepath.Join(homeDir, ConfigDir)

	os.Mkdir(directoryPath, os.ModePerm)
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

	os.Mkdir(dirPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return "", err
	}

	return dirPath, nil
}

func getConfiguration(reader *bufio.Reader) (Configuration, error) {
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

func main() {
	r := bufio.NewReader(os.Stdin)

	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	clipboardEnabled := false
	if err == nil {
		clipboardEnabled = true
	}

	configuration, err := getConfiguration(r)
	if err != nil {
		fmt.Printf("\nconfiguration error: %v\n", err)
		return
	}

	title, err := prompt(r, "Title")
	if err != nil {
		fmt.Printf("\nprompt error: %v\n", err)
		return
	}

	description, err := prompt(r, "Description")
	if err != nil {
		fmt.Printf("\nprompt error: %v\n", err)
		return
	}

	tp := jira.BasicAuthTransport{
		Username: configuration.Username,
		Password: configuration.Password,
	}

	client, err := jira.NewClient(tp.Client(), strings.TrimSpace(configuration.APIUrl))
	if err != nil {
		fmt.Printf("\nclient creation error: %v\n", err)
		return
	}

	p, _, err := client.Project.Get("THC")
	if err != nil {
		fmt.Printf("\nget project error error: %v\n", err)
		return
	}

	var components []*jira.Component
	components = append(components, &jira.Component{
		Name: "Query Language",
		ID:   "10232",
	})

	i := &jira.Issue{
		Fields: &jira.IssueFields{
			Summary:     title,
			Description: description,
			Type: jira.IssueType{
				ID: IssueId,
			},
			Reporter: &jira.User{
				AccountID: configuration.AccountId,
			},
			Project:    *p,
			Components: components,
			Assignee: &jira.User{
				AccountID: configuration.AccountId,
			},
		},
	}

	// create issue
	issue, _, err := client.Issue.Create(i)
	if err != nil {
		fmt.Printf("\nissue creation error: %v\n", err)
		return
	}

	// make "ready to work"
	_, err = client.Issue.DoTransition(issue.ID, "851")
	if err != nil {
		fmt.Printf("\ntransition error: %v\n", err)
		return
	}

	// make "in-progress"
	_, err = client.Issue.DoTransition(issue.ID, "711")
	if err != nil {
		fmt.Printf("\ntransition error: %v\n", err)
		return
	}

	fmt.Printf("%s/browse/%s\n", configuration.APIUrl, issue.Key)
	if clipboardEnabled {
		clipboard.Write(clipboard.FmtText, []byte(issue.Key))
	}
}
