package config

import "github.com/InditexTech/gh-sherpa/internal/domain/issue_types"

type Github struct {
	IssueLabels      GithubIssueLabels `mapstructure:"issue_labels" validate:"required,validIssueTypeKeys,uniqueMapValues"`
	ForkOrganization string            `mapstructure:"fork_organization"`
}

type GithubIssueLabels map[issue_types.IssueType][]string
