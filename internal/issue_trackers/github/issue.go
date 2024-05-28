package github

import (
	"fmt"
	"slices"

	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
)

type Issue struct {
	id           string
	title        string
	body         string
	url          string
	labels       []domain.Label
	labelsConfig config.GithubIssueLabels
}

var _ domain.Issue = (*Issue)(nil)

func (i Issue) FormatID() string {
	return fmt.Sprintf("GH-%s", i.id)
}

func (i Issue) ID() string {
	return i.id
}

func (i Issue) Title() string {
	return i.title
}

func (i Issue) Body() string {
	return i.body
}

func (i Issue) URL() string {
	return i.url
}

func (i Issue) TypeLabel() string {
	for _, cfgLabels := range i.labelsConfig {
		for _, label := range i.labels {
			if slices.Contains(cfgLabels, label.Name) {
				return label.Name
			}
		}
	}

	return ""
}

func (i Issue) TrackerType() domain.IssueTrackerType {
	return domain.IssueTrackerTypeGithub
}

func (i Issue) Type() issue_types.IssueType {
	issueTypeLabel := i.TypeLabel()

	for issueType, cfgLabels := range i.labelsConfig {
		if slices.Contains(cfgLabels, issueTypeLabel) {
			return issueType
		}
	}

	return issue_types.Unknown
}
