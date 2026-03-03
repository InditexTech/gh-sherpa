package domain

import "github.com/InditexTech/gh-sherpa/internal/domain/issue_types"

type Issue interface {
	FormatID() string
	ID() string
	Title() string
	Body() string
	URL() string
	TypeLabel() string
	TrackerType() IssueTrackerType
	Type() issue_types.IssueType
	HasLabel(labelName string) bool
}
