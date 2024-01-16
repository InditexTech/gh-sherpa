package github

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	"github.com/InditexTech/gh-sherpa/internal/gh"
)

var issuePattern = regexp.MustCompile(`^(?i:GH-)?(?P<issue_num>\d+)$`)

var ErrIssueNotFound = fmt.Errorf("the issue was not found")

type Github struct {
	cfg Configuration
	cli domain.GhCli
}

type Configuration struct {
	config.Github
}

type Issue struct {
	Number int64
	Title  string
	Body   string
	Labels []Label
	Url    string
}

type Label struct {
	Id          string
	Name        string
	Description string
	Color       string
}

var _ domain.IssueTracker = (*Github)(nil)

// New returns a new Github issue tracker with the given configuration
func New(cfg Configuration) (*Github, error) {

	return &Github{
		cfg: cfg,
		cli: &gh.Cli{},
	}, nil
}

func (g *Github) GetIssue(identifier string) (issue domain.Issue, err error) {
	command := []string{"issue", "view", identifier, "--json", "labels,number,title,body,url"}

	result := Issue{}

	err = g.cli.Execute(&result, command)
	if err != nil {
		if strings.Contains(err.Error(), "Could not resolve to an issue or pull request") {
			err = ErrIssueNotFound
		}
		return
	}

	labels := make([]domain.Label, len(result.Labels))

	for i, label := range result.Labels {
		labels[i] = domain.Label{
			Id:   label.Id,
			Name: label.Name,
		}
	}

	return domain.Issue{
		ID:           strconv.FormatInt(result.Number, 10),
		Title:        result.Title,
		Body:         result.Body,
		Url:          result.Url,
		Labels:       labels,
		IssueTracker: domain.IssueTrackerTypeGithub,
	}, nil
}

func (g *Github) GetIssueType(issue domain.Issue) (issueType issue_types.IssueType) {
	issueTypeLabel := g.GetIssueTypeLabel(issue)

	for issueType, cfgLabels := range g.cfg.Github.IssueLabels {
		if slices.Contains(cfgLabels, issueTypeLabel) {
			return issueType
		}
	}

	return issue_types.Unknown
}

func (g *Github) IdentifyIssue(identifier string) bool {
	return issuePattern.MatchString(identifier)
}

func (g *Github) CheckConfiguration() (err error) {
	// TODO: Check if configuration is valid
	return
}

func (g *Github) FormatIssueId(issueId string) (formattedIssueId string) {
	return fmt.Sprintf("GH-%s", issueId)
}

func (g *Github) ParseRawIssueId(identifier string) (issueId string) {
	match := issuePattern.FindStringSubmatch(identifier)

	if len(match) > 0 {
		return match[1]
	}

	return ""
}

func (g *Github) GetIssueTrackerType() domain.IssueTrackerType {
	return domain.IssueTrackerTypeGithub
}

// GetIssueTypeLabel returns the type label related to the issue or empty string if not found
func (g *Github) GetIssueTypeLabel(issue domain.Issue) string {
	for _, cfgLabels := range g.cfg.Github.IssueLabels {
		for _, label := range issue.Labels {
			if slices.Contains(cfgLabels, label.Name) {
				return label.Name
			}
		}
	}

	return ""
}
