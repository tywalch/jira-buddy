package jira

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/tywalch/jira-buddy/app/configuration"
	"strings"
)

type JiraClient struct {
	*jira.Client
}

func NewJiraClient(configuration configuration.Configuration) (*JiraClient, error) {
	tp := jira.BasicAuthTransport{
		Username: configuration.Username,
		Password: configuration.Password,
	}

	client, err := jira.NewClient(tp.Client(), strings.TrimSpace(configuration.APIUrl))
	if err != nil {
		fmt.Printf("\nclient creation error: %v\n", err)
		return &JiraClient{client}, err
	}

	return &JiraClient{client}, nil
}
