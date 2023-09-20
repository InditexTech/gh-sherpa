package jira

import (
	"reflect"
	"testing"

	gojira "github.com/andygrunwald/go-jira"
	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
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
