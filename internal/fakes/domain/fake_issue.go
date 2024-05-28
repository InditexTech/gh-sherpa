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

// Body implements domain.Issue.
func (f *FakeIssue) Body() string {
	return "fake body"
}

// FormatID implements domain.Issue.
func (f *FakeIssue) FormatID() string {
	if f.issueTrackerType == domain.IssueTrackerTypeGithub {
		return fmt.Sprintf("GH-%s", f.id)
	}

	return f.id
}

// ID implements domain.Issue.
func (f *FakeIssue) ID() string {
	return f.id
}

// IssueTrackerType implements domain.Issue.
func (f *FakeIssue) IssueTrackerType() domain.IssueTrackerType {
	return f.issueTrackerType
}

// IssueType implements domain.Issue.
func (f *FakeIssue) IssueType() issue_types.IssueType {
	return f.issueType
}

// IssueTypeLabel implements domain.Issue.
func (f *FakeIssue) IssueTypeLabel() string {
	return fmt.Sprintf("kind/%s", f.issueType)
}

// Title implements domain.Issue.
func (f *FakeIssue) Title() string {
	return "fake title"
}

// URL implements domain.Issue.
func (f *FakeIssue) URL() string {
	return "fake url"
}
