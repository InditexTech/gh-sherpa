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

var ErrIdIsPullRequestNumber = func(identifier string) error {
	return fmt.Errorf("given identifier %s is a Pull Request number", identifier)
}

type githubCli interface {
	Execute(result interface{}, command []string) error
	domain.RepositoryProvider
}

type Github struct {
	cfg Configuration
	cli githubCli
}

type Configuration struct {
	config.Github
}

type ghIssue struct {
	Number      int64
	Title       string
	Body        string
	Labels      []Label
	Url         string
	PullRequest *ghPullRequest `json:"pull_request"`
}

func (i ghIssue) isPullRequest() bool {
	return i.PullRequest != nil
}

type ghPullRequest map[string]any

type Label struct {
	Id          int
	Name        string
	Description string
	Color       string
}

var newGhCli = func() githubCli {
	return &gh.Cli{}
}

// New returns a new Github issue tracker with the given configuration
func New(cfg Configuration) (*Github, error) {

	return &Github{
		cfg: cfg,
		cli: newGhCli(),
	}, nil
}

func (g *Github) GetIssue(identifier string) (issue domain.Issue, err error) {
	repo, err := g.cli.GetRepository()
	if err != nil {
		return nil, err
	}

	apiPath := fmt.Sprintf("/repos/%s/issues/%s", repo.NameWithOwner, identifier)
	command := []string{"api", apiPath}

	result := ghIssue{}

	err = g.cli.Execute(&result, command)
	if err != nil {
		if strings.Contains(err.Error(), "Could not resolve to an issue or pull request") {
			err = ErrIssueNotFound
		}
		return
	}

	if result.isPullRequest() {
		return nil, ErrIdIsPullRequestNumber(identifier)
	}

	labels := make([]domain.Label, len(result.Labels))

	for i, label := range result.Labels {
		labels[i] = domain.Label{
			Id:          fmt.Sprintf("%d", label.Id),
			Name:        label.Name,
			Description: label.Description,
			Color:       label.Color,
		}
	}

	issueTypeLabel := g.getIssueTypeLabel(labels)

	return Issue{
		id:        strconv.FormatInt(result.Number, 10),
		title:     result.Title,
		body:      result.Body,
		url:       result.Url,
		labels:    labels,
		typeLabel: issueTypeLabel,
		issueType: g.getIssueType(issueTypeLabel),
	}, nil

}

func (g *Github) getIssueType(issueTypeLabel string) issue_types.IssueType {
	for issueType, cfgLabels := range g.cfg.IssueLabels {
		if slices.Contains(cfgLabels, issueTypeLabel) {
			return issueType
		}
	}

	return issue_types.Unknown
}

func (g *Github) getIssueTypeLabel(labels []domain.Label) string {
	for _, cfgLabels := range g.cfg.IssueLabels {
		for _, label := range labels {
			if slices.Contains(cfgLabels, label.Name) {
				return label.Name
			}
		}
	}

	return ""
}

func (g *Github) IdentifyIssue(identifier string) bool {
	return issuePattern.MatchString(identifier)
}

func (g *Github) CheckConfiguration() (err error) {
	// TODO: Check if configuration is valid
	return
}

func (g *Github) ParseRawIssueId(identifier string) (issueId string) {
	match := issuePattern.FindStringSubmatch(identifier)

	if len(match) > 0 {
		return match[1]
	}

	return ""
}
