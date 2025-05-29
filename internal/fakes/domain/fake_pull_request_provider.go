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
}

var _ domain.PullRequestProvider = (*FakePullRequestProvider)(nil)

func NewFakePullRequestProvider() *FakePullRequestProvider {
	return &FakePullRequestProvider{
		PullRequests:           map[string]*domain.PullRequest{},
		PullRequestsWithErrors: []string{},
	}
}

func (f *FakePullRequestProvider) AddPullRequest(branchName string, pr domain.PullRequest) {
	f.PullRequests[branchName] = &pr
}

// GetLastCreatedPR returns the last pull request created for a branch
func (f *FakePullRequestProvider) GetPullRequestForBranchBody(branch string) (body string, exists bool) {
	pr := f.PullRequests[branch]
	if pr == nil {
		return "", false
	}
	return pr.Body, true
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

	pr := f.PullRequests[headBranch]
	if pr != nil && !pr.Closed {
		return "", ErrPrAlreadyExists(headBranch)
	}

	prLabels := make([]domain.Label, len(labels))
	for i, label := range labels {
		prLabels[i] = domain.Label{
			Id:   label,
			Name: label,
		}
	}
	pr = &domain.PullRequest{
		Title:       title,
		Number:      5,
		State:       "OPEN",
		Closed:      false,
		Url:         "https://github.com/inditextech/gh-sherpa-test-repo/pulls/5",
		HeadRefName: headBranch,
		BaseRefName: baseBranch,
		Labels:      prLabels,
		Body:        body,
	}

	f.PullRequests[headBranch] = pr

	return pr.Url, nil
}
