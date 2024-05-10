package domain

import (
	"fmt"
	"strings"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
)

type FakeIssueTracker struct {
	IssueTrackerType       domain.IssueTrackerType
	IssueTrackerIdentifier map[string]string
	IssueTypes             map[string]issue_types.IssueType
	IssueTypesLabel        map[string]string
	Issues                 map[string]domain.Issue
}

var _ domain.IssueTracker = (*FakeIssueTracker)(nil)

func NewFakeGitHubIssueTracker() *FakeIssueTracker {
	// TODO: IMPLEMENT THIS METHOD
	return &FakeIssueTracker{}
}

func NewFakeJiraIssueTracker() *FakeIssueTracker {
	// TODO: IMPLEMENTE THIS METHOD
	return &FakeIssueTracker{}
}

func ErrIssueNotInitializedInMap(identifier string) error {
	return fmt.Errorf("issue %s not initialized in map", identifier)
}

func (f *FakeIssueTracker) GetIssue(identifier string) (issue domain.Issue, err error) {
	issue, ok := f.Issues[identifier]
	if !ok {
		return domain.Issue{}, ErrIssueNotInitializedInMap(identifier)
	}
	return issue, nil
}

func (f *FakeIssueTracker) GetIssueType(issue domain.Issue) (issueType issue_types.IssueType) {
	return f.IssueTypes[issue.ID]
}

func (f *FakeIssueTracker) GetIssueTypeLabel(issue domain.Issue) string {
	return f.IssueTypesLabel[issue.ID]
}

func (f *FakeIssueTracker) IdentifyIssue(identifier string) bool {
	return f.IssueTrackerIdentifier[identifier] == string(f.IssueTrackerType)
}

func (f *FakeIssueTracker) FormatIssueId(issueId string) (formattedIssueId string) {
	issue, _ := f.GetIssue(issueId)
	issueType := f.GetIssueType(issue)
	issueTrackerType := f.GetIssueTrackerType()
	issueTracker := "GH"
	if issueTrackerType == domain.IssueTrackerTypeJira {
		issueTracker = "PROJECTKEY"
	}
	return fmt.Sprintf("%s/%s-%s-generated-from-issue-tracker", issueType, issueTracker, issueId)
}

func (f *FakeIssueTracker) ParseRawIssueId(identifier string) (issueId string) {
	prefix := ""

	switch f.IssueTrackerType {
	case domain.IssueTrackerTypeGithub:
		prefix = "GH-"
	case domain.IssueTrackerTypeJira:
		prefix = "PROJECTKEY-"
	}

	if strings.HasPrefix(identifier, prefix) {
		return strings.ReplaceAll(identifier, prefix, "")
	}

	return ""
}

func (f *FakeIssueTracker) GetIssueTrackerType() domain.IssueTrackerType {
	return f.IssueTrackerType
}
