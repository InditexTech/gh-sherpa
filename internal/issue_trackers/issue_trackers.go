package issue_trackers

import (
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/issue_trackers/github"
	"github.com/InditexTech/gh-sherpa/internal/issue_trackers/jira"
	"github.com/InditexTech/gh-sherpa/internal/logging"
)

type Configuration struct {
	Jira   jira.Configuration
	Github github.Configuration
}
type Provider struct {
	cfg    Configuration
	github github.Github
	jira   jira.Jira
}

var _ domain.IssueTrackerProvider = (*Provider)(nil)

// New returns a new issue tracker provider
func New(cfg Configuration) (*Provider, error) {
	g, err := github.New(cfg.Github)
	if err != nil {
		return nil, err
	}

	j, err := jira.New(cfg.Jira)
	if err != nil {
		return nil, err
	}

	return &Provider{
		cfg:    cfg,
		github: *g,
		jira:   *j,
	}, nil
}

// NewConfiguration returns a new configuration from the given global configuration
func NewFromConfiguration(globalConfig config.Configuration) (*Provider, error) {
	return New(Configuration{
		Jira: jira.Configuration{
			Jira:            globalConfig.Jira,
			IssueTypeLabels: globalConfig.Github.IssueLabels,
		},
		Github: github.Configuration{
			Github: globalConfig.Github,
		},
	})
}

func (p Provider) GetIssue(identifier string) (domain.Issue, error) {
	if p.github.IdentifyIssue(identifier) {
		logging.Debugf("Issue %s identified as a Github issue", identifier)
		return p.github.GetIssue(identifier)
	}

	if p.jira.IdentifyIssue(identifier) {
		logging.Debugf("Issue %s identified as a Jira issue", identifier)
		return p.jira.GetIssue(identifier)
	}

	return nil, fmt.Errorf("could not identify issue %s", identifier)
}

func (p Provider) ParseIssueId(identifier string) (issueId string) {
	if p.github.IdentifyIssue(identifier) {
		logging.Debugf("Issue %s identified as a Github issue", identifier)
		return p.github.ParseRawIssueId(identifier)
	}

	if p.jira.IdentifyIssue(identifier) {
		logging.Debugf("Issue %s identified as a Jira issue", identifier)
		return p.jira.ParseRawIssueId(identifier)
	}

	return
}
