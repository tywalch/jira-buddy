package jira

type IssueType struct {
	Name string
	Id   string
}

type ValidIssueTypes = [5]IssueType

func GetIssueTypes() ValidIssueTypes {
	return ValidIssueTypes{
		{Name: "Bug", Id: "10004"},
		{Name: "Feature", Id: "10003"},
		{Name: "Task", Id: "10002"},
		{Name: "Story", Id: "10001"},
		{Name: "Epic", Id: "10000"},
	}
}

func IsValidIssueTypeId(issueTypes ValidIssueTypes, issueTypeId string) bool {
	for _, s := range issueTypes {
		if s.Id == issueTypeId {
			return true
		}
	}

	return false
}
