package jira

import (
	"errors"
	"github.com/andygrunwald/go-jira"
)

type IssueInput struct {
	Title       string `default:""`
	Description string `default:""`
	AccountId   string `default:""`
	Type        string `default:""`
	Status      string `default:""`
}

func CreateNewIssue(client *JiraClient, input IssueInput) (*jira.Issue, error) {
	p, _, err := client.Project.Get("THC")
	if err != nil {
		return &jira.Issue{}, err
	}

	var components []*jira.Component
	components = append(components, &jira.Component{
		Name: "Query Language",
		ID:   "10232",
	})

	i := &jira.Issue{
		Fields: &jira.IssueFields{
			Summary:     input.Title,
			Description: input.Description,
			Type: jira.IssueType{
				ID: input.Type,
			},
			Reporter: &jira.User{
				AccountID: input.AccountId,
			},
			Project:    *p,
			Components: components,
			Assignee: &jira.User{
				AccountID: input.AccountId,
			},
		},
	}

	// create issue
	issue, _, err := client.Issue.Create(i)
	if err != nil {
		return issue, err
	}

	// transition issue
	err = setIssueStatus(client, issue.ID, input.Status)
	if err != nil {
		return issue, err
	}

	return issue, nil
}

func setIssueStatus(client *JiraClient, issueId string, statusId string) error {
	transitions := GetIssueStatuses()
	isValid := IsValidIssueStatusId(transitions, statusId)
	if !isValid {
		return errors.New("invalid transition name")
	}

	for _, t := range transitions {
		if len(t.Id) > 0 {
			_, err := client.Issue.DoTransition(issueId, t.Id)
			if err != nil {
				return err
			}
		}

		if t.Id == statusId {
			return nil
		}
	}

	return nil
}
