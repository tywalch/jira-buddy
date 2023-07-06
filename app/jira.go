package app

import (
	"bufio"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"strings"
)

const (
	IssueId = "10002"
)

type JiraClient struct {
	*jira.Client
	Configuration Configuration
}

func CreateClient(configuration Configuration) (JiraClient, error) {
	tp := jira.BasicAuthTransport{
		Username: configuration.Username,
		Password: configuration.Password,
	}

	client, err := jira.NewClient(tp.Client(), strings.TrimSpace(configuration.APIUrl))
	if err != nil {
		fmt.Printf("\nclient creation error: %v\n", err)
		return JiraClient{client, configuration}, nil
	}

	return JiraClient{client, configuration}, nil
}

type JiraBuddy interface {
	CreateIssue(title string, description string) (*jira.Issue, error)
}

type TicketDetail struct {
	Title       string
	Description string
	IssueId     string
}

func getTicketDetail(r *bufio.Reader) (TicketDetail, error) {
	title, err := prompt(r, "Title")
	if err != nil {
		return TicketDetail{}, err
	}

	description, err := prompt(r, "Description")
	if err != nil {
		return TicketDetail{}, err
	}

	return TicketDetail{
		Title:       title,
		Description: description,
	}, nil
}

func (client JiraClient) CreateIssue(title string, description string) (*jira.Issue, error) {
	p, _, err := client.Project.Get("THC")
	if err != nil {
		fmt.Printf("\nget project error error: %v\n", err)
		return &jira.Issue{}, err
	}

	var components []*jira.Component
	components = append(components, &jira.Component{
		Name: "Query Language",
		ID:   "10232",
	})

	i := &jira.Issue{
		Fields: &jira.IssueFields{
			Summary:     title, /**/
			Description: description,
			Type: jira.IssueType{
				ID: IssueId,
			},
			Reporter: &jira.User{
				AccountID: client.Configuration.AccountId,
			},
			Project:    *p,
			Components: components,
			Assignee: &jira.User{
				AccountID: client.Configuration.AccountId,
			},
		},
	}

	// create issue
	issue, _, err := client.Issue.Create(i)
	if err != nil {
		return i, nil
	}

	// make "ready to work"
	_, err = client.Issue.DoTransition(issue.ID, "851")
	if err != nil {
		return i, nil
	}

	// make "in-progress"
	_, err = client.Issue.DoTransition(issue.ID, "711")
	if err != nil {
		return i, nil
	}

	return issue, nil
}
