package domain

import (
	"fmt"
	"strings"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
)

type FakeIssueTracker struct {
	IssueTrackerType domain.IssueTrackerType
	Configurations   map[string]fakeIssueTrackerConfiguration
	// IssueTrackerIdentifier map[string]string
	// IssueTypes             map[string]issue_types.IssueType
	// IssueTypesLabel        map[string]string
	// Issues                 map[string]domain.Issue
}

type fakeIssueTrackerConfiguration struct {
	IssueTrackerIdentifier string
	IssueType              issue_types.IssueType
	IssueTypeLabel         string
	Issue                  domain.Issue
}

var _ domain.IssueTracker = (*FakeIssueTracker)(nil)

func newFakeIssueTracker() *FakeIssueTracker {
	return &FakeIssueTracker{
		Configurations: map[string]fakeIssueTrackerConfiguration{
			"1": {
				IssueTrackerIdentifier: domain.IssueTrackerTypeGithub.String(),
				IssueType:              issue_types.Feature,
				IssueTypeLabel:         "kind/feature",
				Issue: domain.Issue{
					ID:           "1",
					IssueTracker: domain.IssueTrackerTypeGithub,
					Title:        "Sample issue",
					Body:         "Sample issue body",
					Labels: []domain.Label{
						{
							Id:   "kind/feature",
							Name: "kind/feature",
						},
					},
					Url: "https://github.com/InditexTech/gh-sherpa-repo-test/issues/1",
				},
			},
			"PROJECTKEY-1": {
				IssueTrackerIdentifier: domain.IssueTrackerTypeJira.String(),
				IssueType:              issue_types.Feature,
				IssueTypeLabel:         "kind/feature",
				Issue: domain.Issue{
					ID:           "PROJECTKEY-1",
					IssueTracker: domain.IssueTrackerTypeJira,
					Title:        "Sample issue",
					Body:         "Sample issue body",
					Labels: []domain.Label{
						{
							Id:   "kind/feature",
							Name: "kind/feature",
						},
					},
					Url: "https://sample.jira.com/PROJECTKEY-1",
				},
			},
			"3": {
				IssueTrackerIdentifier: domain.IssueTrackerTypeGithub.String(),
				IssueType:              issue_types.Feature,
				IssueTypeLabel:         "kind/documentation",
				Issue: domain.Issue{
					ID:           "3",
					IssueTracker: domain.IssueTrackerTypeGithub,
					Title:        "Sample documentation issue",
					Body:         "Sample documentation issue body",
					Labels: []domain.Label{
						{
							Id:   "kind/documentation",
							Name: "kind/documentation",
						},
					},
					Url: "https://github.com/InditexTech/gh-sherpa-repo-test/issues/3",
				},
			},
			"PROJECTKEY-3": {
				IssueTrackerIdentifier: domain.IssueTrackerTypeJira.String(),
				IssueType:              issue_types.Feature,
				IssueTypeLabel:         "kind/documentation",
				Issue: domain.Issue{
					ID:           "PROJECTKEY-3",
					IssueTracker: domain.IssueTrackerTypeJira,
					Title:        "Sample documentation issue",
					Body:         "Sample documentation issue body",
					Labels: []domain.Label{
						{
							Id:   "kind/documentation",
							Name: "kind/documentation",
						},
					},
					Url: "https://sample.jira.com/PROJECTKEY-3",
				},
			},
			"6": {
				IssueTrackerIdentifier: domain.IssueTrackerTypeGithub.String(),
				IssueType:              issue_types.Feature,
				IssueTypeLabel:         "kind/refactoring",
				Issue: domain.Issue{
					ID:           "6",
					IssueTracker: domain.IssueTrackerTypeJira,
					Title:        "Sample refactoring issue",
					Body:         "Sample refactoring issue body",
					Labels: []domain.Label{
						{
							Id:   "kind/refactoring",
							Name: "kind/refactoring",
						},
					},
					Url: "https://github.com/InditexTech/gh-sherpa-repo-test/issues/6",
				},
			},
			"PROJECTKEY-6": {
				IssueTrackerIdentifier: domain.IssueTrackerTypeJira.String(),
				IssueType:              issue_types.Feature,
				IssueTypeLabel:         "kind/refactoring",
				Issue: domain.Issue{
					ID:           "PROJECTKEY-6",
					IssueTracker: domain.IssueTrackerTypeJira,
					Title:        "Sample refactoring issue",
					Body:         "Sample refactoring issue body",
					Labels: []domain.Label{
						{
							Id:   "kind/refactoring",
							Name: "kind/refactoring",
						},
					},
					Url: "https://sample.jira.com/PROJECTKEY-6",
				},
			},
		},
	}
}

func NewFakeGitHubIssueTracker() *FakeIssueTracker {
	issueTracker := newFakeIssueTracker()
	issueTracker.IssueTrackerType = domain.IssueTrackerTypeGithub
	return issueTracker
}

func NewFakeJiraIssueTracker() *FakeIssueTracker {
	issueTracker := newFakeIssueTracker()
	issueTracker.IssueTrackerType = domain.IssueTrackerTypeJira
	return issueTracker
}

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
