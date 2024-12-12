package domain

import (
	"errors"
	"fmt"
	"slices"

	"github.com/InditexTech/gh-sherpa/internal/domain"
)

type FakePullRequestProvider struct {
	PullRequests           map[string]*domain.PullRequest
	PullRequestsWithErrors []string
	CreatedPRs            []CreatedPR
}

type CreatedPR struct {
	Title       string
	Body        string
	BaseBranch  string
	HeadBranch  string
	Draft       bool
	Labels      []string
}

var _ domain.PullRequestProvider = (*FakePullRequestProvider)(nil)

func NewFakePullRequestProvider() *FakePullRequestProvider {
	return &FakePullRequestProvider{
		PullRequests:           map[string]*domain.PullRequest{},
		PullRequestsWithErrors: []string{},
		CreatedPRs:            []CreatedPR{},
	}
}

func (f *FakePullRequestProvider) AddPullRequest(branchName string, pr domain.PullRequest) {
	f.PullRequests[branchName] = &pr
}

func (f *FakePullRequestProvider) HasPullRequestForBranch(branch string) bool {
	pr := f.PullRequests[branch]

	return pr != nil
}

func (f *FakePullRequestProvider) GetPullRequestForBranch(branch string) (pullRequest *domain.PullRequest, err error) {
	pr := f.PullRequests[branch]

	return pr, nil
}

func ErrPrAlreadyExists(branch string) error {
	return fmt.Errorf("pr already exists for branch %s", branch)
}

var ErrPullRequestWithError = errors.New("pull request with error")

func (f *FakePullRequestProvider) CreatePullRequest(title string, body string, baseBranch string, headBranch string, draft bool, labels []string) (prUrl string, err error) {
	if slices.Contains(f.PullRequestsWithErrors, headBranch) {
		return "", ErrPullRequestWithError
	}

	f.CreatedPRs = append(f.CreatedPRs, CreatedPR{
		Title:       title,
		Body:        body,
		BaseBranch:  baseBranch,
		HeadBranch:  headBranch,
		Draft:       draft,
		Labels:      labels,
	})

	pr := &domain.PullRequest{
		Title:       title,
		HeadRefName: headBranch,
		BaseRefName: baseBranch,
		State:       "OPEN",
		Closed:      false,
		Url:         fmt.Sprintf("https://github.com/owner/repo/pull/%d", len(f.PullRequests)+1),
	}

	f.PullRequests[headBranch] = pr

	return pr.Url, nil
}
