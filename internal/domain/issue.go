package domain

import "github.com/InditexTech/gh-sherpa/internal/domain/issue_types"

// type Issue struct {
// 	ID           string
// 	Title        string
// 	Body         string
// 	Url          string
// 	IssueTracker IssueTrackerType
// 	// Used in GitHub
// 	Labels []Label
// 	// Used in Jira
// 	Type IssueType
// }

// type IssueType struct {
// 	Id          string
// 	Name        string
// 	Description string
// }

type Issue interface {
	FormatID() string
	ID() string
	Title() string
	Body() string
	URL() string
	IssueTypeLabel() string
	IssueTrackerType() IssueTrackerType
	IssueType() issue_types.IssueType
}
