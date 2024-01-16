package jira

import (
	"reflect"
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	gojira "github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		expected := domain.Issue{
			ID:    "ISSUE-1",
			Title: "Summary",
			Body:  "Description",
			Type: domain.IssueType{
				Id:          "1",
				Name:        "Bug",
				Description: "Fixes a bug",
			},
			Url:          "https://jira.example.com/jira/browse/ISSUE-1",
			IssueTracker: domain.IssueTrackerTypeJira,
		}

		s.Truef(reflect.DeepEqual(expected, result), "expected: %v, got: %v", expected, result)
	})
}

func TestGetIssueType(t *testing.T) {

	createIssue := func(issueTypeId string) domain.Issue {
		return domain.Issue{Type: domain.IssueType{Id: issueTypeId}}
	}

	cfg := Configuration{
		Jira: config.Jira{
			IssueTypes: config.JiraIssueTypes{
				issue_types.Bug:         {"1"},
				issue_types.Feature:     {"3", "5"},
				issue_types.Improvement: {},
			},
		},
	}

	j, err := New(cfg)
	require.NoError(t, err)

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
			got := j.GetIssueType(tc.issue)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetIssueTypeLabel(t *testing.T) {
	createIssue := func(issueTypeId string) domain.Issue {
		return domain.Issue{Type: domain.IssueType{Id: issueTypeId}}
	}

	cfg := Configuration{
		Jira: config.Jira{
			IssueTypes: config.JiraIssueTypes{
				issue_types.Bug:         {"1"},
				issue_types.Feature:     {"3", "5"},
				issue_types.Improvement: {},
			},
		},
		IssueTypeLabels: map[issue_types.IssueType][]string{
			issue_types.Bug:     {"kind/bug", "kind/bugfix"},
			issue_types.Feature: {"kind/feat"},
		},
	}

	j, err := New(cfg)
	require.NoError(t, err)

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
			got := j.GetIssueTypeLabel(tc.issue)
			assert.Equal(t, tc.want, got)
		})
	}
}
