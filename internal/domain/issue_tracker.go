package domain

type IssueTrackerType string

func (i IssueTrackerType) String() string {
	return string(i)
}

const (
	IssueTrackerTypeGithub IssueTrackerType = "github"
	IssueTrackerTypeJira   IssueTrackerType = "jira"
)

type IssueTrackerProvider interface {
	GetIssue(identifier string) (issue Issue, err error)
	ParseIssueId(identifier string) (issueId string)
}
