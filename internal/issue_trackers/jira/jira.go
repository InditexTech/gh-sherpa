package jira

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	gojira "github.com/andygrunwald/go-jira"
	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
)

var issuePattern = regexp.MustCompile(`^(?P<issue_key>\w+)-(?P<issue_num>\d+)$`)

type Jira struct {
	cfg    Configuration
	client JiraClient
}

type JiraClient struct {
	gojira.Client
}

type Configuration struct {
	config.Jira
}

// Validate validates the configuration
func (c Configuration) Validate() (err error) {
	//TODO: Validate configuration
	return nil
}

// New returns a new Jira issue tracker with the given configuration
func New(cfg Configuration) (jira *Jira, err error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	jira = &Jira{cfg: cfg}

	gojiraClient, err := createBearerClient(cfg.Auth.Token, cfg.Auth.Host)
	if err != nil {
		return nil, fmt.Errorf("could not create a Jira client: %s", err)
	}

	jira.client = *gojiraClient

	return
}

var _ domain.IssueTracker = (*Jira)(nil)

func (j *Jira) GetIssue(identifier string) (issue domain.Issue, err error) {
	issueGot, res, err := j.client.Issue.Get(identifier, &gojira.GetQueryOptions{Fields: "issuetype,summary"})

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

func (j *Jira) GetIssueType(issue domain.Issue) (issueType issue_types.IssueType, err error) {
	for issueType, ids := range j.cfg.Jira.IssueTypes {
		for _, id := range ids {
			if id == issue.Type.Id {
				return issueType, nil
			}
		}
	}

	return issue_types.Unknown, nil
}

func (j *Jira) IdentifyIssue(identifier string) bool {
	return issuePattern.MatchString(identifier)
}

func (j *Jira) CheckConfiguration() (err error) {

	return
}

func (j *Jira) FormatIssueId(issueId string) (formattedIssueId string) {
	return issueId
}

func (j *Jira) GetIssueTrackerType() domain.IssueTrackerType {
	return domain.IssueTrackerTypeJira
}

func createBearerClient(token string, host string) (client *JiraClient, err error) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	tp := gojira.BearerAuthTransport{
		Token:     token,
		Transport: customTransport,
	}

	gojiraClient, err := gojira.NewClient(tp.Client(), host)

	if err != nil {
		return
	}

	client = &JiraClient{*gojiraClient}

	return
}

func (j *Jira) ParseRawIssueId(identifier string) (issueId string) {
	return identifier
}

func (j *Jira) goJiraIssueToIssue(issue gojira.Issue) domain.Issue {
	return domain.Issue{
		ID:    issue.Key,
		Title: issue.Fields.Summary,
		Body:  issue.Fields.Description,
		Url:   fmt.Sprintf("%s/browse/%s", j.cfg.Auth.Host, issue.Key),
		Type: domain.IssueType{
			Id:          issue.Fields.Type.ID,
			Name:        issue.Fields.Type.Name,
			Description: issue.Fields.Type.Description,
		},
		IssueTracker: domain.IssueTrackerTypeJira,
	}
}
