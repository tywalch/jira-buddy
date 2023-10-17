package jira

type IssueStatus struct {
	Name string
	Id   string
}

type ValidIssueStatuses = [3]IssueStatus

func GetIssueStatuses() ValidIssueStatuses {
	return ValidIssueStatuses{
		{Name: "Open", Id: ""},
		{Name: "Ready To Work", Id: "851"},
		{Name: "In Progress", Id: "711"},
	}
}

func IsValidIssueStatusId(transitions ValidIssueStatuses, statusId string) bool {
	for _, t := range transitions {
		if t.Id == statusId {
			return true
		}
	}
	return false
}
