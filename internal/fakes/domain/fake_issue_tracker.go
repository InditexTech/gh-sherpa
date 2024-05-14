package domain

import (
	"fmt"
	"strings"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
)

type FakeIssueTracker struct {
	IssueTrackerType domain.IssueTrackerType
	Configurations   map[string]FakeIssueTrackerConfiguration
}

type FakeIssueTrackerConfiguration struct {
	IssueTrackerIdentifier string
	IssueType              issue_types.IssueType
	IssueTypeLabel         string
	Issue                  domain.Issue
}

var _ domain.IssueTracker = (*FakeIssueTracker)(nil)

func ErrIssueNotInitializedInMap(identifier string) error {
	return fmt.Errorf("issue %s not initialized in map", identifier)
}

func (f *FakeIssueTracker) GetIssue(issueId string) (issue domain.Issue, err error) {
	config, ok := f.Configurations[issueId]
	if !ok {
		return domain.Issue{}, ErrIssueNotInitializedInMap(issueId)
	}
	return config.Issue, nil
}

func (f *FakeIssueTracker) GetIssueType(issue domain.Issue) (issueType issue_types.IssueType) {
	return f.Configurations[issue.ID].IssueType
}

func (f *FakeIssueTracker) GetIssueTypeLabel(issue domain.Issue) string {
	return f.Configurations[issue.ID].IssueTypeLabel
}

func (f *FakeIssueTracker) IdentifyIssue(issueId string) bool {
	return f.Configurations[issueId].IssueTrackerIdentifier == string(f.IssueTrackerType)
}

func (f *FakeIssueTracker) FormatIssueId(issueId string) (formattedIssueId string) {
	issueTrackerType := f.GetIssueTrackerType()
	if issueTrackerType == domain.IssueTrackerTypeGithub {
		return fmt.Sprintf("GH-%s", issueId)
	}

	return issueId
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
