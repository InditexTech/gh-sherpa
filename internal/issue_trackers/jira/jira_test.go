package jira

import (
	"errors"
	"net/http"
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	gojira "github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/suite"
)

type fakeClient struct {
	issue    *gojira.Issue
	response *gojira.Response
	err      error
}

func (f *fakeClient) setError() {
	f.err = errors.New("error")
}

func (f *fakeClient) setIssue(key string) {
	f.issue = &gojira.Issue{
		Key: key,
		Fields: &gojira.IssueFields{
			Summary:     "Issue Summary",
			Description: "Issue Description",
			Type: gojira.IssueType{
				ID:          "3",
				Name:        "Jira Issue Name",
				Description: "Jira Issue Description",
			},
		},
	}
}

func (f *fakeClient) changeIssueType(id string) {
	if f.issue != nil {
		f.issue.Fields.Type.ID = id
	}
}

func (f *fakeClient) setResponse(statusCode int) {
	f.response = &gojira.Response{
		Response: &http.Response{
			StatusCode: statusCode,
		},
	}
}

func (f *fakeClient) getIssue(identifier string) (*gojira.Issue, *gojira.Response, error) {
	return f.issue, f.response, f.err
}

type JiraTestSuite struct {
	suite.Suite
	jira               *Jira
	createBearerClient func(token string, host string, skipTLSVerify bool) (gojiraClient, error)
	fakeClient         *fakeClient
	defaultKey         string
	expectedIssue      *Issue
}

func TestJiraTestSuite(t *testing.T) {
	suite.Run(t, new(JiraTestSuite))
}

func (s *JiraTestSuite) SetupSuite() {
	s.defaultKey = "PROJECTKEY-1"
	s.createBearerClient = createBearerClient
}

func (s *JiraTestSuite) TearDownSuite() {
	createBearerClient = s.createBearerClient
}

func (s *JiraTestSuite) SetupSubTest() {
	s.fakeClient = &fakeClient{}
	s.fakeClient.setIssue(s.defaultKey)
	s.fakeClient.setResponse(http.StatusOK)

	createBearerClient = func(token, host string, skipTLSVerify bool) (gojiraClient, error) {
		return s.fakeClient, nil
	}

	cfg := Configuration{
		Jira: config.Jira{
			Auth: config.JiraAuth{
				Host: "https://jira.example.com/jira",
			},
			IssueTypes: config.JiraIssueTypes{
				issue_types.Bug:         {"1"},
				issue_types.Feature:     {"3", "5"},
				issue_types.Refactoring: {},
			},
		},
		IssueTypeLabels: map[issue_types.IssueType][]string{
			issue_types.Bug:         {"kind/bug", "kind/bugfix"},
			issue_types.Feature:     {"kind/feature"},
			issue_types.Refactoring: {},
		},
	}

	j, err := New(cfg)
	s.Require().NoError(err)

	s.jira = j

	s.expectedIssue = &Issue{
		id:    s.defaultKey,
		title: "Issue Summary",
		body:  "Issue Description",
		url:   "https://jira.example.com/jira/browse/PROJECTKEY-1",
		jiraIssueType: JiraIssueType{
			Id:          "3",
			Name:        "Jira Issue Name",
			Description: "Jira Issue Description",
		},
		typeLabel: "kind/feature",
		issueType: issue_types.Feature,
	}
}

func (s *JiraTestSuite) TestGetIssue() {
	s.Run("should return error if could not execute", func() {
		s.fakeClient.setError()

		issue, err := s.jira.GetIssue(s.defaultKey)

		s.Error(err)
		s.Nil(issue)
	})

	s.Run("should return bug issue", func() {
		s.fakeClient.changeIssueType("1")

		issue, err := s.jira.GetIssue(s.defaultKey)

		s.NoError(err)
		s.Require().NotNil(issue)
		s.Equal(issue_types.Bug, issue.Type())
		s.Equal("kind/bug", issue.TypeLabel())
	})

	s.Run("should return unknown issue", func() {

		s.fakeClient.changeIssueType("99")

		issue, err := s.jira.GetIssue(s.defaultKey)

		s.NoError(err)
		s.Require().NotNil(issue)
		s.Equal(issue_types.Unknown, issue.Type())
	})

	s.Run("should return issue", func() {
		issue, err := s.jira.GetIssue(s.defaultKey)

		s.NoError(err)
		s.Require().NotNil(s.expectedIssue)
		s.Equal(*s.expectedIssue, issue)
	})
}
