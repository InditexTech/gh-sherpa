package jira

import (
	"reflect"
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	gojira "github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type JiraTestSuite struct {
	suite.Suite
	jira *Jira
}

func TestJiraTestSuite(t *testing.T) {
	suite.Run(t, new(JiraTestSuite))
}

func (s *JiraTestSuite) SetupSuite() {
	cfg := Configuration{
		Jira: config.Jira{
			Auth: config.JiraAuth{
				Host: "https://jira.example.com/jira",
			},
		},
	}

	j, err := New(cfg)
	s.Require().NoError(err)

	s.jira = j
}

func (s *JiraTestSuite) TestGojiraIssueToDomainIssue() {
	s.Run("Should convert a gojira issue to a domain issue", func() {
		issue := gojira.Issue{
			Key: "ISSUE-1",
			Fields: &gojira.IssueFields{
				Summary:     "Summary",
				Description: "Description",
				Type: gojira.IssueType{
					ID:          "1",
					Name:        "Bug",
					Description: "Fixes a bug",
				},
			},
		}

		result := s.jira.goJiraIssueToIssue(issue)

		expected := Issue{
			id:    "ISSUE-1",
			title: "Summary",
			body:  "Description",
			issueType: IssueType{
				Id:          "1",
				Name:        "Bug",
				Description: "Fixes a bug",
			},
			url: "https://jira.example.com/jira/browse/ISSUE-1",
		}

		s.Truef(reflect.DeepEqual(expected, result), "expected: %v, got: %v", expected, result)
	})
}

func TestGetIssueType(t *testing.T) {

	createIssue := func(issueTypeId string) domain.Issue {
		return Issue{
			issueType: IssueType{
				Id: issueTypeId,
			},
			issueTypesConfig: config.JiraIssueTypes{
				issue_types.Bug:         {"1"},
				issue_types.Feature:     {"3", "5"},
				issue_types.Improvement: {},
			},
		}
	}

	for _, tc := range []struct {
		name  string
		issue domain.Issue
		want  issue_types.IssueType
	}{
		{
			name:  "GetIssueType bug",
			issue: createIssue("1"),
			want:  issue_types.Bug,
		},
		{
			name:  "GetIssueType feature",
			issue: createIssue("3"),
			want:  issue_types.Feature,
		},
		{
			name:  "GetIssueType unknown",
			issue: createIssue("-1"),
			want:  issue_types.Unknown,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := tc.issue.Type()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetIssueTypeLabel(t *testing.T) {
	createIssue := func(issueTypeId string) domain.Issue {
		return Issue{
			issueType: IssueType{
				Id: issueTypeId,
			},
			issueTypesConfig: config.JiraIssueTypes{
				issue_types.Bug:         {"1"},
				issue_types.Feature:     {"3", "5"},
				issue_types.Improvement: {},
			},
			labelsConfig: map[issue_types.IssueType][]string{
				issue_types.Bug:     {"kind/bug", "kind/bugfix"},
				issue_types.Feature: {"kind/feat"},
			},
		}
	}

	for _, tc := range []struct {
		name  string
		issue domain.Issue
		want  string
	}{
		{
			name:  "Get issue type label with single mapped label",
			issue: createIssue("5"),
			want:  "kind/feat",
		},
		{
			name:  "Returns first issue label if multiple labels are mapped to the same issue type",
			issue: createIssue("1"),
			want:  "kind/bug",
		},
		{
			name:  "Returns empty string if no kind is present in the issue",
			issue: createIssue("-1"),
			want:  "",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := tc.issue.TypeLabel()
			assert.Equal(t, tc.want, got)
		})
	}
}
