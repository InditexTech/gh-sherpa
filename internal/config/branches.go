package config

import "github.com/InditexTech/gh-sherpa/internal/domain/issue_types"

type Branches struct {
	Prefixes map[issue_types.IssueType]string
}
