package github

import (
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
)

type Issue struct {
	id        string
	title     string
	body      string
	url       string
	typeLabel string
	issueType issue_types.IssueType
	labels    []domain.Label
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
	return i.typeLabel
}

func (i Issue) TrackerType() domain.IssueTrackerType {
	return domain.IssueTrackerTypeGithub
}

func (i Issue) Type() issue_types.IssueType {
	return i.issueType
}

func (i Issue) HasLabel(labelName string) bool {
	for _, label := range i.labels {
		if label.Name == labelName {
			return true
		}
	}
	return false
}
