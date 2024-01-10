// SPDX-FileCopyrightText: 2023 INDITEX S.A
//
// SPDX-License-Identifier: Apache-2.0

package gh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCli_CreatePullRequest(t *testing.T) {
	type args struct {
		title      string
		body       string
		baseBranch string
		headBranch string
		draft      bool
		labels     []string
	}
	tests := []struct {
		name        string
		args        args
		wantPrURL   string
		wantErr     bool
		executeArgs []string
	}{
		{
			name:        "CreatePR",
			args:        args{title: "title", body: "body", baseBranch: "develop", headBranch: "asd", draft: false, labels: []string{}},
			wantPrURL:   "https://github.com/InditexTech/gh-sherpa/pulls/1",
			wantErr:     false,
			executeArgs: []string{"pr", "create", "-B", "develop", "-H", "asd", "-t", "title", "-b", "body"},
		},
		{
			name:        "CreatePR draft and default",
			args:        args{title: "", body: "", baseBranch: "develop", headBranch: "asd", draft: true, labels: []string{}},
			wantPrURL:   "https://github.com/InditexTech/gh-sherpa/pulls/1",
			wantErr:     false,
			executeArgs: []string{"pr", "create", "-B", "develop", "-H", "asd", "-d", "-f"},
		},
		{
			name:        "CreatePR with labels",
			args:        args{title: "", body: "", baseBranch: "develop", headBranch: "asd", draft: true, labels: []string{"label1", "label2"}},
			wantPrURL:   "https://github.com/InditexTech/gh-sherpa/pulls/1",
			wantErr:     false,
			executeArgs: []string{"pr", "create", "-B", "develop", "-H", "asd", "-d", "-f", "-l", "label1", "-l", "label2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{}

			var executeArgs []string
			ExecuteStringResult = func(args []string) (result string, err error) {
				executeArgs = args
				return "https://github.com/InditexTech/gh-sherpa/pulls/1\n", nil
			}

			gotPrURL, err := c.CreatePullRequest(tt.args.title, tt.args.body, tt.args.baseBranch, tt.args.headBranch, tt.args.draft, tt.args.labels)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.executeArgs, executeArgs)
			assert.Equal(t, tt.wantPrURL, gotPrURL)
		})
	}
}
