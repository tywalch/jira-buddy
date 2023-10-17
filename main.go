package main

import (
	"fmt"
	"github.com/tywalch/jira-buddy/app"
	"github.com/tywalch/jira-buddy/app/configuration"
	"github.com/tywalch/jira-buddy/app/jira"

	"github.com/tywalch/jira-buddy/app/prompt"
	clippy "golang.design/x/clipboard"
)

func promptTicketInput(reader *prompt.PromptReader, config configuration.Configuration) (app.TicketInput, error) {
	title, err := reader.GetString("Title")
	if err != nil {
		return app.TicketInput{}, err
	}

	description, err := reader.GetString("Description")
	if err != nil {
		return app.TicketInput{}, err
	}

	status, err := promptIssueStatus(reader)
	if err != nil {
		return app.TicketInput{}, err
	}

	issueType, err := promptIssueType(reader)
	if err != nil {
		return app.TicketInput{}, err
	}

	return app.TicketInput{
		Title:       title,
		Description: description,
		Status:      status,
		AccountId:   config.AccountId,
		Type:        issueType,
	}, nil
}

func getIssueStatusPickerOptions(statuses jira.ValidIssueStatuses) []prompt.StringPickerOption {
	options := make([]prompt.StringPickerOption, len(statuses))
	for i, t := range statuses {
		options[i] = prompt.StringPickerOption{
			Name:  t.Name,
			Value: t.Id,
		}
	}

	return options
}

func promptIssueStatus(reader *prompt.PromptReader) (string, error) {
	issueStatuses := jira.GetIssueStatuses()
	options := getIssueStatusPickerOptions(issueStatuses)
	status, err := reader.PickString("Ticket Status", options)
	if err != nil {
		return "", err
	}

	if !jira.IsValidIssueStatusId(issueStatuses, status) {
		return "", fmt.Errorf("Invalid issue status provided: '%s'", status)
	}

	return status, nil
}

func getIssueTypePickerOptions(issueTypes jira.ValidIssueTypes) []prompt.StringPickerOption {
	options := make([]prompt.StringPickerOption, len(issueTypes))
	for i, t := range issueTypes {
		options[i] = prompt.StringPickerOption{
			Name:  t.Name,
			Value: t.Id,
		}
	}

	return options
}

func promptIssueType(reader *prompt.PromptReader) (string, error) {
	issueTypes := jira.GetIssueTypes()
	options := getIssueTypePickerOptions(issueTypes)
	issueType, err := reader.PickString("Ticket Type", options)
	if err != nil {
		return "", err
	}

	if !jira.IsValidIssueTypeId(issueTypes, issueType) {
		return "", fmt.Errorf("Invalid issue type provided: '%s'", issueType)
	}

	return issueType, nil
}

func NewClipboard() func(value string) {
	err := clippy.Init()
	clipboardEnabled := false
	if err == nil {
		clipboardEnabled = true
	}

	return func(value string) {
		if clipboardEnabled {
			clippy.Write(clippy.FmtText, []byte(value))
		}
	}
}

func main() {
	clipboard := NewClipboard()

	reader := app.NewReader()

	config, err := configuration.GetConfiguration(reader)
	if err != nil {
		fmt.Printf("\nError retrieving configuration: %v\n", err)
		return
	}

	jiraBuddy, err := app.NewJiraBuddy(config)
	if err != nil {
		fmt.Printf("\nError initializing client: %v\n", err)
		return
	}

	ticketInput, err := promptTicketInput(reader, config)
	if err != nil {
		fmt.Printf("\nError processing ticket input: %v\n", err)
		return
	}

	ticket, err := jiraBuddy.CreateNewTicket(ticketInput)
	if err != nil {
		fmt.Printf("\nError creating new ticket: %v\n", err)
		return
	}

	fmt.Println(ticket.Url)
	clipboard(ticket.Key)
}
