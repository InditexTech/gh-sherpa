// SPDX-FileCopyrightText: 2023 INDITEX S.A
//
// SPDX-License-Identifier: Apache-2.0

package config

import "github.com/InditexTech/gh-sherpa/internal/domain/issue_types"

type Github struct {
	IssueLabels GithubIssueLabels `mapstructure:"issue_labels" validate:"required,validIssueTypeKeys,uniqueMapValues"`
}

type GithubIssueLabels map[issue_types.IssueType][]string
