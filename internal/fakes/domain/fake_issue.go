package domain

import (
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
)

type FakeIssue struct {
	id               string
	issueType        issue_types.IssueType
	issueTrackerType domain.IssueTrackerType
}

var _ domain.Issue = (*FakeIssue)(nil)

func NewFakeIssue(id string, issueType issue_types.IssueType, issueTrackerType domain.IssueTrackerType) *FakeIssue {
	return &FakeIssue{
		id:               id,
		issueType:        issueType,
		issueTrackerType: issueTrackerType,
	}
}

func (f *FakeIssue) Body() string {
	return "fake body"
}

func (f *FakeIssue) FormatID() string {
	if f.issueTrackerType == domain.IssueTrackerTypeGithub {
		return fmt.Sprintf("GH-%s", f.id)
	}

	return f.id
}

func (f *FakeIssue) ID() string {
	return f.id
}

func (f *FakeIssue) TrackerType() domain.IssueTrackerType {
	return f.issueTrackerType
}

func (f *FakeIssue) Type() issue_types.IssueType {
	return f.issueType
}

func (f *FakeIssue) TypeLabel() string {
	return fmt.Sprintf("kind/%s", f.issueType)
}

func (f *FakeIssue) Title() string {
	return "fake title"
}

func (f *FakeIssue) URL() string {
	return "fake url"
}
