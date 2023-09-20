package github

import (
	"fmt"
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	"github.com/InditexTech/gh-sherpa/internal/gh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockedCli struct {
	gh.Cli
}

var _ domain.GhCli = (*mockedCli)(nil)

func (m *mockedCli) GetRepo() (repo *domain.Repository, err error) {
	repo = &domain.Repository{
		Name:             "Repo 1",
		Owner:            "Owner 1",
		DefaultBranchRef: "main",
	}
	return
}

func (m *mockedCli) Execute(result any, _ []string) (err error) {
	switch result := result.(type) {
	case *Issue:
		*result = Issue{
			Number: 1,
			Title:  "Issue 1",
			Body:   "Body 1",
			Labels: []Label{{Id: "1", Name: "Label 1"}},
		}
	default:
		panic("unexpected type")
	}

	return
}

type mockedCliWithErr struct {
	mockedCli
}

func (m *mockedCliWithErr) Execute(result any, _ []string) (err error) {
	return fmt.Errorf("error")
}

func Test_GetIssue(t *testing.T) {
	type args struct {
		identifier string
	}
	tests := []struct {
		name      string
		args      args
		want      domain.Issue
		wantErr   bool
		mockedCli domain.GhCli
	}{
		{
			name:      "GetIssue",
			args:      args{identifier: "1"},
			want:      domain.Issue{ID: "1", Title: "Issue 1", Body: "Body 1", Labels: []domain.Label{{Id: "1", Name: "Label 1"}}, IssueTracker: domain.IssueTrackerTypeGithub},
			mockedCli: &mockedCli{},
			wantErr:   false,
		},
		{
			name:      "GetIssue with error",
			args:      args{identifier: "1"},
			mockedCli: &mockedCliWithErr{},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			github := &Github{cli: tt.mockedCli}
			issue, err := github.GetIssue(tt.args.identifier)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, issue)
		})
	}
}

func Test_CheckConfiguration(t *testing.T) {
	type fields struct {
		Cli domain.GhCli
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

func TestGetIssueType(t *testing.T) {

	createIssue := func(labelNames ...string) domain.Issue {
		labels := make([]domain.Label, len(labelNames))
		for i, labelName := range labelNames {
			labels[i] = domain.Label{Name: labelName}
		}
		return domain.Issue{Labels: labels}
	}

	cfg := Configuration{
		Github: config.Github{
			IssueLabels: map[issue_types.IssueType][]string{
				issue_types.Bug:         {"bug", "bugfix"},
				issue_types.Feature:     {"feature", "enhancement"},
				issue_types.Refactoring: {},
			},
		},
	}

	g, err := New(cfg)
	require.NoError(t, err)

	for _, tc := range []struct {
		name  string
		issue domain.Issue
		want  issue_types.IssueType
	}{
		{
			name:  "GetIssueType bug",
			issue: createIssue("bug"),
			want:  issue_types.Bug,
		},
		{
			name:  "GetIssueType bugfix",
			issue: createIssue("bugfix"),
			want:  issue_types.Bug,
		},
		{
			name:  "GetIssueType feature",
			issue: createIssue("feature"),
			want:  issue_types.Feature,
		},
		{
			name:  "GetIssueType enhancement",
			issue: createIssue("enhancement"),
			want:  issue_types.Feature,
		},
		{
			name:  "GetIssueType unknown refactoring",
			issue: createIssue("refactoring"),
			want:  issue_types.Unknown,
		},
		{
			name:  "GetIssueType unknown",
			issue: createIssue("unknown"),
			want:  issue_types.Unknown,
		},
		{
			name:  "GetIssueType bug with several labels",
			issue: createIssue("feature", "bug", "refactoring"),
			want:  issue_types.Bug,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := g.GetIssueType(tc.issue)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_IdentifyIssue(t *testing.T) {
	type fields struct {
		Cli domain.GhCli
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
			fields: fields{Cli: &mockedCli{}},
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
	type fields struct {
		Cli domain.GhCli
	}
	type args struct {
		issue domain.Issue
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantIssueId string
	}{
		{
			name:        "FormatIssueId",
			args:        args{issue: domain.Issue{ID: "1"}},
			wantIssueId: "GH-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Github{
				cli: tt.fields.Cli,
			}
			assert.Equalf(t, tt.wantIssueId, g.FormatIssueId(tt.args.issue.ID), "FormatIssueId(%v)", tt.args.issue)
		})
	}
}
