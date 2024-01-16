package domain

import "github.com/InditexTech/gh-sherpa/internal/domain/issue_types"

type IssueTrackerType string

func (i IssueTrackerType) String() string {
	return string(i)
}

const (
	IssueTrackerTypeGithub IssueTrackerType = "github"
	IssueTrackerTypeJira   IssueTrackerType = "jira"
)

type IssueTracker interface {
	// GetIssue returns the issue information
	GetIssue(identifier string) (issue Issue, err error)
	// GetIssueType returns the issue type
	GetIssueType(issue Issue) (issueType issue_types.IssueType)
	// GetIssueTypeLabel returns the associated label for the issue type
	GetIssueTypeLabel(issue Issue) string
	// IdentifyIssue checks if the identifier is a valid issue identifier
	IdentifyIssue(identifier string) bool
	// FormatIssueId formats the issue identifier
	FormatIssueId(issueId string) (formattedIssueId string)
	// ParseRawIssueId parses the raw issue identifier and returns the issue id
	ParseRawIssueId(identifier string) (issueId string)
	// GetIssueTrackerType returns the issue tracker type
	GetIssueTrackerType() IssueTrackerType
}

type IssueTrackerProvider interface {
	GetIssueTracker(identifier string) (issueTracker IssueTracker, err error)
	ParseIssueId(identifier string) (issueId string)
}
