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
	return i.issueType
}

// TypeLabel implements domain.Issue.
func (i Issue) TypeLabel() string {
	return i.typeLabel
}

// URL implements domain.Issue.
func (i Issue) URL() string {
	return i.url
}
