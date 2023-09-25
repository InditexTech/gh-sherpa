package branches

import (
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	"github.com/stretchr/testify/assert"
)

func TestParseIssueContext(t *testing.T) {
	tests := []struct {
		name  string
		given string
		want  string
	}{
		{
			name:  "issue context with special chars",
			given: "  /hello world//test/.test..test@{test\\\\test.lock  ",
			want:  "hello-world-test-test-testtesttest",
		},
		{
			name:  "issue context with special chars",
			given: "Begoña Caçadora_renombrádo' Él !cÓncepto $de \"bloqueo\" en cc%, úÍ",
			want:  "begona-cazadora_renombrado-el-concepto-de-bloqueo-en-cc-ui",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := parseIssueContext(tt.given)

			assert.Equal(t, tt.want, context)
		})
	}
}

func TestFormatBranchName(t *testing.T) {
	repositoryName := "InditexTech/gh-sherpa"

	type args struct {
		repository           string
		branchType           string
		issueId              string
		issueContext         string
		branchPrefixOverride map[issue_types.IssueType]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Does format branch name",
			args: args{repository: repositoryName, branchType: "feature", issueId: "GH-1", issueContext: "my-title"},
			want: "feature/GH-1-my-title",
		},
		{
			name: "Does format branch name with override",
			args: args{
				repository:           repositoryName,
				branchType:           "feature",
				issueId:              "GH-1",
				issueContext:         "my-title",
				branchPrefixOverride: map[issue_types.IssueType]string{issue_types.Feature: "feat"},
			},
			want: "feat/GH-1-my-title",
		},
		{
			name: "Does format long branch name",
			args: args{repository: repositoryName, branchType: "feature", issueId: "GH-1", issueContext: "my-title-is-too-long-and-it-should-be-truncated"},
			want: "feature/GH-1-my-title-is-too-long-and-it-s",
		},
		{
			name: "Does format long branch name with override",
			args: args{
				repository:           repositoryName,
				branchType:           "feature",
				issueId:              "GH-1",
				issueContext:         "my-title-is-too-long-and-it-should-be-truncated",
				branchPrefixOverride: map[issue_types.IssueType]string{issue_types.Feature: "feat"},
			},
			want: "feat/GH-1-my-title-is-too-long-and-it-shou",
		},
		{
			name: "Does format branch with empty override",
			args: args{
				repository:           repositoryName,
				branchType:           "refactoring",
				issueId:              "GH-1",
				issueContext:         "refactor-issue",
				branchPrefixOverride: map[issue_types.IssueType]string{issue_types.Refactoring: ""},
			},
			want: "refactoring/GH-1-refactor-issue",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			b := BranchProvider{
				cfg: Configuration{
					Branches: config.Branches{
						Prefixes: tt.args.branchPrefixOverride,
					},
				},
			}
			branchName := b.formatBranchName(tt.args.repository, tt.args.branchType, tt.args.issueId, tt.args.issueContext)

			assert.Equal(t, tt.want, branchName)
		})
	}
}

func TestParseBranchName(t *testing.T) {
	for _, tc := range []struct {
		name       string
		branchName string
		want       *BranchNameInfo
	}{
		{
			branchName: "feature/GH-1-my-title",
			want:       &BranchNameInfo{BranchType: "feature", IssueId: "GH-1", IssueContext: "my-title"},
		},
		{
			branchName: "bugfix/PROJECTKEY-1-my-title",
			want:       &BranchNameInfo{BranchType: "bugfix", IssueId: "PROJECTKEY-1", IssueContext: "my-title"},
		},
		{
			branchName: "feature/GH-1-my-title-is-too-long-and-it-should-not-matter",
			want:       &BranchNameInfo{BranchType: "feature", IssueId: "GH-1", IssueContext: "my-title-is-too-long-and-it-should-not-matter"},
		},
		{
			branchName: "randomprefix/A_PROJECT_KEY-99-issue-tittle-here",
			want:       &BranchNameInfo{BranchType: "randomprefix", IssueId: "A_PROJECT_KEY-99", IssueContext: "issue-tittle-here"},
		},
	} {
		tc := tc
		t.Run(tc.branchName, func(t *testing.T) {
			branchInfo := ParseBranchName(tc.branchName)

			assert.Equal(t, tc.want, branchInfo)
		})
	}
}
