package jira

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	gojira "github.com/andygrunwald/go-jira"
)

var issuePattern = regexp.MustCompile(`^(?P<issue_key>\w+)-(?P<issue_num>\d+)$`)

type Jira struct {
	cfg    Configuration
	client gojiraClient
}

type gojiraClient interface {
	getIssue(issueID string) (*gojira.Issue, *gojira.Response, error)
}

type Configuration struct {
	config.Jira
	IssueTypeLabels map[issue_types.IssueType][]string
}

// New returns a new Jira issue tracker with the given configuration
func New(cfg Configuration) (jira *Jira, err error) {

	jira = &Jira{cfg: cfg}

	bearerClient, err := createBearerClient(cfg.Auth.Token, cfg.Auth.Host, cfg.Auth.SkipTLSVerify)
	if err != nil {
		return nil, fmt.Errorf("could not create a Jira client: %s", err)
	}

	jira.client = bearerClient

	return
}

// var _ domain.IssueTracker = (*Jira)(nil)

func (j *Jira) GetIssue(identifier string) (issue domain.Issue, err error) {
	issueGot, res, err := j.client.getIssue(identifier)

	if err != nil {
		if res == nil {
			err = fmt.Errorf("could not get response from host '%s'. Check your jira configuration", j.cfg.Auth.Host)
			return
		}

		switch res.StatusCode {
		case http.StatusUnauthorized:
			err = errors.New("your PAT is invalid or revoked")
		case http.StatusForbidden:
			err = errors.New("you do not have permission to get this issue")
		case http.StatusNotFound:
			err = errors.New("the issue was not found")
		default:
			err = fmt.Errorf("could not get issue: %s", err)
		}

		return
	}

	issue = j.goJiraIssueToIssue(*issueGot)

	return
}

// func (j *Jira) GetIssueType(issue domain.Issue) (issueType issue_types.IssueType) {
// 	for issueType, ids := range j.cfg.Jira.IssueTypes {
// 		for _, id := range ids {
// 			if id == issue.Type.Id {
// 				return issueType
// 			}
// 		}
// 	}

// 	return issue_types.Unknown
// }

func (j *Jira) IdentifyIssue(identifier string) bool {
	return issuePattern.MatchString(identifier)
}

func (j *Jira) CheckConfiguration() (err error) {
	// TODO: Check if configuration is valid
	return
}

// func (j *Jira) FormatIssueId(issueId string) (formattedIssueId string) {
// 	return issueId
// }

// func (j *Jira) GetIssueTrackerType() domain.IssueTrackerType {
// 	return domain.IssueTrackerTypeJira
// }

func (j *Jira) ParseRawIssueId(identifier string) (issueId string) {
	return identifier
}

func (j *Jira) goJiraIssueToIssue(issue gojira.Issue) domain.Issue {

	// issueTypeLabel := j.getIssueTypeLabel(issue.Fields.Type)

	issueType := j.getIssueType(issue.Fields.Type.ID)

	return Issue{
		id:    issue.Key,
		title: issue.Fields.Summary,
		body:  issue.Fields.Description,
		url:   fmt.Sprintf("%s/browse/%s", j.cfg.Auth.Host, issue.Key),
		jiraIssueType: JiraIssueType{
			Id:          issue.Fields.Type.ID,
			Name:        issue.Fields.Type.Name,
			Description: issue.Fields.Type.Description,
		},
		issueType: issueType,
		typeLabel: j.getIssueTypeLabel(issueType),
	}
}

func (j *Jira) getIssueType(issueID string) issue_types.IssueType {
	for issueType, ids := range j.cfg.IssueTypes {
		for _, id := range ids {
			if id == issueID {
				return issueType
			}
		}
	}

	return issue_types.Unknown
}

func (j *Jira) getIssueTypeLabel(issueType issue_types.IssueType) string {
	for mappedIssueType, labels := range j.cfg.IssueTypeLabels {
		if issueType == mappedIssueType && len(labels) > 0 {
			return labels[0]
		}
	}

	return ""
}

// func (j *Jira) GetIssueTypeLabel(issue domain.Issue) string {
// 	issueType := j.GetIssueType(issue)

// 	for mappedIssueType, labels := range j.cfg.IssueTypeLabels {
// 		if issueType == mappedIssueType && len(labels) > 0 {
// 			return labels[0]
// 		}
// 	}

// 	return ""
// }
