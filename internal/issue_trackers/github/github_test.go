package github

import (
	"errors"
	"fmt"
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	"github.com/InditexTech/gh-sherpa/internal/gh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type fakeCli struct {
	gh.Cli
	issue *ghIssue
	err   error
}

func (f *fakeCli) setError() {
	f.err = errors.New("error")
}

func (f *fakeCli) setIssue(number int) {
	f.issue = &ghIssue{
		Number: int64(number),
		Title:  "Issue Title",
		Body:   "Issue Body",
		Url:    "https://github.example.com/owner/repo/issues/1",
		Labels: []Label{
			{
				Id:          1,
				Name:        "kind/feature",
				Description: "feature kind label",
				Color:       "#fff",
			},
		},
	}
}

func (f *fakeCli) resetLabels() {
	f.issue.Labels = []Label{}
}

func (f *fakeCli) addIssueTypeLabel(issueType issue_types.IssueType) {
	f.issue.Labels = append(f.issue.Labels, Label{
		Id:          len(issueType),
		Name:        fmt.Sprintf("kind/%s", issueType),
		Description: fmt.Sprintf("%s kind label", issueType),
	})
}

var _ githubCli = (*fakeCli)(nil)

func (f *fakeCli) GetRepository() (repo *domain.Repository, err error) {
	repo = &domain.Repository{
		Name:             "Repo 1",
		Owner:            "Owner 1",
		DefaultBranchRef: "main",
	}
	return
}

var errExecuteError = fmt.Errorf("execute error")

func (f *fakeCli) Execute(result any, _ []string) (err error) {
	if f.err != nil {
		return f.err
	}
	if f.issue == nil {
		return errExecuteError
	}
	switch result := result.(type) {
	case *ghIssue:
		*result = *f.issue
	default:
		panic("unexpected type")
	}

	return
}

type GithubTestSuite struct {
	suite.Suite
	github         *Github
	fakeCli        *fakeCli
	defaultIssueID string
	expectedIssue  *Issue
	newGhCli       func() githubCli
}

func TestGithubSuite(t *testing.T) {
	suite.Run(t, new(GithubTestSuite))
}

func (s *GithubTestSuite) SetupSuite() {
	s.defaultIssueID = "1"
	s.newGhCli = newGhCli
}

func (s *GithubTestSuite) TearDownSuite() {
	newGhCli = s.newGhCli
}

func (s *GithubTestSuite) SetupSubTest() {
	s.fakeCli = &fakeCli{}
	s.fakeCli.setIssue(1)

	newGhCli = func() githubCli {
		return s.fakeCli
	}

	cfg := Configuration{
		Github: config.Github{
			IssueLabels: config.GithubIssueLabels{
				issue_types.Bug:         {"kind/bug", "kind/bugfix"},
				issue_types.Feature:     {"kind/feature", "kind/enhancement"},
				issue_types.Refactoring: {},
			},
		},
	}

	g, err := New(cfg)
	s.Require().NoError(err)

	s.github = g

	s.expectedIssue = &Issue{
		id:        s.defaultIssueID,
		title:     "Issue Title",
		body:      "Issue Body",
		url:       "https://github.example.com/owner/repo/issues/1",
		typeLabel: "kind/feature",
		issueType: issue_types.Feature,
		labels: []domain.Label{
			{
				Id:          "1",
				Name:        "kind/feature",
				Description: "feature kind label",
				Color:       "#fff",
			},
		},
	}
}

func (s *GithubTestSuite) TestGetIssue() {
	s.Run("should return error if could not execute", func() {
		s.fakeCli.setError()

		issue, err := s.github.GetIssue(s.defaultIssueID)

		s.Error(err)
		s.Nil(issue)
	})

	s.Run("should return bug issue", func() {
		s.fakeCli.resetLabels()
		s.fakeCli.addIssueTypeLabel(issue_types.Bug)

		issue, err := s.github.GetIssue(s.defaultIssueID)

		s.NoError(err)
		s.Require().NotNil(issue)
		s.Equal(issue_types.Bug, issue.Type())
		s.Equal("kind/bug", issue.TypeLabel())
	})

	s.Run("should return unknown issue if no label is present", func() {
		s.fakeCli.resetLabels()

		issue, err := s.github.GetIssue(s.defaultIssueID)

		s.NoError(err)
		s.Require().NotNil(issue)
		s.Equal(issue_types.Unknown, issue.Type())
	})

	s.Run("should return unknown issue if could not determine label type", func() {
		s.fakeCli.resetLabels()
		s.fakeCli.addIssueTypeLabel("random-label")

		issue, err := s.github.GetIssue(s.defaultIssueID)

		s.NoError(err)
		s.Require().NotNil(issue)
		s.Equal(issue_types.Unknown, issue.Type())
	})

	s.Run("should return issue", func() {
		issue, err := s.github.GetIssue(s.defaultIssueID)

		s.NoError(err)
		s.Require().NotNil(s.expectedIssue)
		s.Equal(*s.expectedIssue, issue)
	})

	s.Run("should return error if given issue id is a pull request number", func() {
		s.fakeCli.issue.PullRequest = &ghPullRequest{}

		issueId := "99"
		issue, err := s.github.GetIssue(issueId)

		s.ErrorContains(err, ErrIdIsPullRequestNumber(issueId).Error())
		s.Nil(issue)
	})

}

func Test_CheckConfiguration(t *testing.T) {
	type fields struct {
		Cli githubCli
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Github{
				cli: tt.fields.Cli,
			}
			tt.wantErr(t, g.CheckConfiguration(), "CheckConfiguration()")
		})
	}
}

func Test_IdentifyIssue(t *testing.T) {
	type fields struct {
		Cli githubCli
	}
	type args struct {
		identifier string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "IdentifyIssue",
			args:   args{identifier: "1"},
			fields: fields{Cli: &fakeCli{}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Github{
				cli: tt.fields.Cli,
			}
			assert.Equalf(t, tt.want, g.IdentifyIssue(tt.args.identifier), "IdentifyIssue(%v)", tt.args.identifier)
		})
	}
}

func TestGithub_FormatIssueId(t *testing.T) {
	type args struct {
		issue domain.Issue
	}
	tests := []struct {
		name        string
		args        args
		wantIssueId string
	}{
		{
			name:        "FormatIssueId",
			args:        args{issue: Issue{id: "1"}},
			wantIssueId: "GH-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantIssueId, tt.args.issue.FormatID(), "FormatIssueId(%v)", tt.args.issue.ID())
		})
	}
}
