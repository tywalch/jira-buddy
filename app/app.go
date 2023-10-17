package app

import (
	"bufio"
	"fmt"
	"github.com/tywalch/jira-buddy/app/configuration"
	"github.com/tywalch/jira-buddy/app/jira"
	"github.com/tywalch/jira-buddy/app/prompt"
	"os"
)

type JiraBuddy struct {
	Client        *jira.JiraClient
	Configuration configuration.Configuration
}

func NewJiraBuddy(config configuration.Configuration) (JiraBuddy, error) {
	client, err := jira.NewJiraClient(config)
	if err != nil {
		return JiraBuddy{&jira.JiraClient{}, config}, err
	}

	return JiraBuddy{client, config}, nil
}

type TicketInput struct {
	Title       string `default:""`
	Description string `default:""`
	AccountId   string `default:""`
	Type        string `default:""`
	Status      string `default:""`
}

type JiraBuddyHelper interface {
	CreateNewTicket(ticketInput TicketInput) (Ticket, error)
}

type Ticket struct {
	Key string
	Url string
}

func (jiraBuddy JiraBuddy) CreateNewTicket(ticketInput TicketInput) (Ticket, error) {
	issue, err := jira.CreateNewIssue(jiraBuddy.Client, jira.IssueInput{
		Title:       ticketInput.Title,
		Description: ticketInput.Description,
		Status:      ticketInput.Status,
		AccountId:   ticketInput.AccountId,
		Type:        ticketInput.Type,
	})

	if err != nil {
		return Ticket{}, err
	}

	return Ticket{
		Key: issue.Key,
		Url: fmt.Sprintf("%s/browse/%s\n", jiraBuddy.Configuration.APIUrl, issue.Key),
	}, nil
}

func NewReader() *prompt.PromptReader {
	return &prompt.PromptReader{Reader: bufio.NewReader(os.Stdin)}
}
