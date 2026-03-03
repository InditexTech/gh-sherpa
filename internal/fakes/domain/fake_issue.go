package domain

import (
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
)

type FakeIssue struct {
	id               string
	title            string
	body             string
	url              string
	issueType        issue_types.IssueType
	issueTrackerType domain.IssueTrackerType
	typeLabel        string
	labels           []domain.Label
}

var _ domain.Issue = (*FakeIssue)(nil)

func (f *FakeIssue) SetTitle(title string) {
	f.title = title
}

func (f *FakeIssue) SetType(issueType issue_types.IssueType) {
	f.issueType = issueType
}

func (f *FakeIssue) SetTypeLabel(label string) {
	f.typeLabel = label
}

func (f *FakeIssue) AddLabel(label domain.Label) {
	f.labels = append(f.labels, label)
}

func NewFakeIssue(id string, issueType issue_types.IssueType, issueTrackerType domain.IssueTrackerType) *FakeIssue {
	return &FakeIssue{
		id:               id,
		title:            "fake title",
		body:             "fake body",
		url:              "fake url",
		issueType:        issueType,
		issueTrackerType: issueTrackerType,
		typeLabel:        fmt.Sprintf("kind/%s", issueType),
	}
}

func (f *FakeIssue) Body() string {
	return f.body
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
	return f.typeLabel
}

func (f *FakeIssue) Title() string {
	return f.title
}

func (f *FakeIssue) URL() string {
	return f.url
}

func (f *FakeIssue) HasLabel(labelName string) bool {
	for _, label := range f.labels {
		if label.Name == labelName {
			return true
		}
	}
	return false
}
