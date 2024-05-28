package github

import (
	"fmt"
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	"github.com/InditexTech/gh-sherpa/internal/gh"
	"github.com/stretchr/testify/assert"
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
	case *ghIssue:
		*result = ghIssue{
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
			want:      Issue{id: "1", title: "Issue 1", body: "Body 1", labels: []domain.Label{{Id: "1", Name: "Label 1"}}},
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

// TODO: MOVE THIS TEST WHERE IS NEEDED
func TestGetIssueType(t *testing.T) {

	createIssue := func(labelNames ...string) domain.Issue {
		labels := make([]domain.Label, len(labelNames))
		for i, labelName := range labelNames {
			labels[i] = domain.Label{Name: labelName}
		}
		return Issue{
			labels: labels,
			labelsConfig: config.GithubIssueLabels{
				issue_types.Bug:         {"bug", "bugfix"},
				issue_types.Feature:     {"feature", "enhancement"},
				issue_types.Refactoring: {},
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
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := tc.issue.IssueType()
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

//TODO: MOVE THIS TEST WHERE IS NEDED

func TestGetIssueTypeLabel(t *testing.T) {

	createIssue := func(labelNames ...string) domain.Issue {
		labels := make([]domain.Label, len(labelNames))
		for i, labelName := range labelNames {
			labels[i] = domain.Label{Name: labelName}
		}
		return Issue{
			labels: labels,
			labelsConfig: config.GithubIssueLabels{
				issue_types.Bug:         {"kind/bug", "kind/bugfix"},
				issue_types.Feature:     {"kind/feat"},
				issue_types.Refactoring: {},
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
			issue: createIssue("kind/feat"),
			want:  "kind/feat",
		},
		{
			name:  "Returns the same label in the issue if issue labels contains several mapped labels",
			issue: createIssue("kind/bugfix"),
			want:  "kind/bugfix",
		},
		{
			name:  "Get issue type label with multiple labels",
			issue: createIssue("not-a-label-kind", "kind/bugfix", "non-related-label"),
			want:  "kind/bugfix",
		},
		{
			name:  "Returns empty string if no label is present in the issue",
			issue: createIssue(),
			want:  "",
		},
		{
			name:  "Returns empty string if no kind label is present in the issue",
			issue: createIssue("not-a-label-kind", "non-related-label", "another-non-related-label"),
			want:  "",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := tc.issue.IssueTypeLabel()
			assert.Equal(t, tc.want, got)
		})
	}
}
