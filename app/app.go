package app

import (
	"bufio"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"golang.design/x/clipboard"
	"os"
	"strings"
)

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

func main() {
	r := bufio.NewReader(os.Stdin)

	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	clipboardEnabled := false
	if err == nil {
		clipboardEnabled = true
	}

	configuration, err := GetConfiguration(r)
	ticketDetail, err := getTicketDetail(r)
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
			Summary:     ticketDetail.Title,
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
