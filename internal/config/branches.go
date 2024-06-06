package config

import "github.com/InditexTech/gh-sherpa/internal/domain/issue_types"

type Branches struct {
	Prefixes  BranchesPrefixes `validate:"validIssueTypeKeys"`
	MaxLength int              `mapstructure:"max_length" validate:"gte=0"`
}

type BranchesPrefixes map[issue_types.IssueType]string
