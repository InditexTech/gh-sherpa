package jira

import (
	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
)

type Issue struct {
	id               string
	title            string
	body             string
	url              string
	issueType        IssueType
	issueTypesConfig config.JiraIssueTypes
	labelsConfig     map[issue_types.IssueType][]string
}

var _ domain.Issue = (*Issue)(nil)

type IssueType struct {
	Id          string
	Name        string
	Description string
}

// Body implements domain.Issue.
func (i Issue) Body() string {
	return i.body
}

// FormatID implements domain.Issue.
func (i Issue) FormatID() string {
	return i.id
}

// ID implements domain.Issue.
func (i Issue) ID() string {
	return i.id
}

// Title implements domain.Issue.
func (i Issue) Title() string {
	return i.title
}

// TrackerType implements domain.Issue.
func (i Issue) TrackerType() domain.IssueTrackerType {
	return domain.IssueTrackerTypeJira
}

// Type implements domain.Issue.
func (i Issue) Type() issue_types.IssueType {
	for issueType, ids := range i.issueTypesConfig {
		for _, id := range ids {
			if id == i.issueType.Id {
				return issueType
			}
		}
	}

	return issue_types.Unknown
}

// TypeLabel implements domain.Issue.
func (i Issue) TypeLabel() string {
	issueType := i.Type()

	for mappedIssueType, labels := range i.labelsConfig {
		if issueType == mappedIssueType && len(labels) > 0 {
			return labels[0]
		}
	}

	return ""
}

// URL implements domain.Issue.
func (i Issue) URL() string {
	return i.url
}
