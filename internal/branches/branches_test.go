package branches

import (
	"testing"

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
			context := ParseIssueContext(tt.given)

			assert.Equal(t, tt.want, context)
		})
	}
}

func TestFormatBranchName(t *testing.T) {
	type args struct {
		repository   string
		branchType   string
		issueId      string
		issueContext string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Does format branch name",
			args: args{repository: "InditexTech/gh-sherpa", branchType: "feature", issueId: "GH-1", issueContext: "my-title"},
			want: "feature/GH-1-my-title",
		},
		{
			name: "Does format branch name",
			args: args{repository: "InditexTech/gh-sherpa", branchType: "feature", issueId: "GH-1", issueContext: "my-title-is-too-long-and-it-should-be-truncated"},
			want: "feature/GH-1-my-title-is-too-long-and-it-shoul",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			branchName := FormatBranchName(tt.args.repository, tt.args.branchType, tt.args.issueId, tt.args.issueContext)

			assert.Equal(t, tt.want, branchName)
		})
	}
}
