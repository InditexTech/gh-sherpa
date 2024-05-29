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

type fakeCli struct {
	gh.Cli
	issue *ghIssue
}

var _ domain.GhCli = (*fakeCli)(nil)

func (f *fakeCli) GetRepo() (repo *domain.Repository, err error) {
	repo = &domain.Repository{
		Name:             "Repo 1",
		Owner:            "Owner 1",
		DefaultBranchRef: "main",
	}
	return
}

var errExecuteError = fmt.Errorf("execute error")

func (f *fakeCli) Execute(result any, _ []string) (err error) {
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

func TestGetIssue(t *testing.T) {
	newGhIssue := func(number int, labels []string) *ghIssue {
		labl := make([]Label, len(labels))
		for i, label := range labels {
			labl[i] = Label{Name: label}
		}

		return &ghIssue{
			Number: int64(number),
			Title:  "fake title",
			Body:   "fake body",
			Url:    "fake url",
			Labels: labl,
		}
	}

	newIssue := func(id string, typeLabel string, issueType issue_types.IssueType, labels []string) Issue {
		labl := make([]domain.Label, len(labels))
		for i, label := range labels {
			labl[i] = domain.Label{Name: label}
		}

		return Issue{
			id:        id,
			title:     "fake title",
			body:      "fake body",
			url:       "fake url",
			typeLabel: typeLabel,
			issueType: issueType,
			labels:    labl,
		}
	}

	tests := []struct {
		name           string
		identifier     string
		retrievedIssue *ghIssue
		want           Issue
		wantErr        bool
	}{
		{
			name:           "should return bug issue",
			identifier:     "1",
			retrievedIssue: newGhIssue(1, []string{"kind/bug"}),
			want:           newIssue("1", "kind/bug", issue_types.Bug, []string{"kind/bug"}),
		},
		{
			name:           "should return feature issue",
			identifier:     "1",
			retrievedIssue: newGhIssue(1, []string{"kind/enhancement", "other/label"}),
			want:           newIssue("1", "kind/enhancement", issue_types.Feature, []string{"kind/enhancement", "other/label"}),
		},
		{
			name:           "should return unknown issue if no label is present",
			identifier:     "1",
			retrievedIssue: newGhIssue(1, []string{}),
			want:           newIssue("1", "", issue_types.Unknown, []string{}),
		},
		{
			name:           "should return unknown issue if could not determine type",
			identifier:     "1",
			retrievedIssue: newGhIssue(1, []string{"random-label", "other-label"}),
			want:           newIssue("1", "", issue_types.Unknown, []string{"random-label", "other-label"}),
		},
		{
			name:    "should return error if could not execute",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			github := &Github{
				cli: &fakeCli{
					issue: tc.retrievedIssue,
				},
				cfg: Configuration{
					Github: config.Github{
						IssueLabels: config.GithubIssueLabels{
							issue_types.Bug:         {"kind/bug", "kind/bugfix"},
							issue_types.Feature:     {"kind/feature", "kind/enhancement"},
							issue_types.Refactoring: {},
						},
					},
				},
			}

			issue, err := github.GetIssue(tc.identifier)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, issue)
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
