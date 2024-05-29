package jira

import (
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
)

type Issue struct {
	id            string
	title         string
	body          string
	url           string
	jiraIssueType JiraIssueType
	typeLabel     string
	issueType     issue_types.IssueType
}

var _ domain.Issue = (*Issue)(nil)

type JiraIssueType struct {
	Id          string
	Name        string
	Description string
}

func (i Issue) Body() string {
	return i.body
}

func (i Issue) FormatID() string {
	return i.id
}

func (i Issue) ID() string {
	return i.id
}

func (i Issue) Title() string {
	return i.title
}

func (i Issue) TrackerType() domain.IssueTrackerType {
	return domain.IssueTrackerTypeJira
}

func (i Issue) Type() issue_types.IssueType {
	return i.issueType
}

func (i Issue) TypeLabel() string {
	return i.typeLabel
}

func (i Issue) URL() string {
	return i.url
}
