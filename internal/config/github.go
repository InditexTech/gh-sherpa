package config

import "github.com/InditexTech/gh-sherpa/internal/domain/issue_types"

type Github struct {
	IssueLabels GithubIssueLabels `mapstructure:"issue_labels"`
}

type GithubIssueLabels map[issue_types.IssueType][]string
