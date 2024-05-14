package domain

import (
	"errors"
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/domain"
)

type FakePullRequestProvider struct {
	PullRequests           map[string]*domain.PullRequest
	PullRequestsWithErrors map[string]error
}

var _ domain.PullRequestProvider = (*FakePullRequestProvider)(nil)

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
	prError, ok := f.PullRequestsWithErrors[headBranch]
	if ok {
		if prError == nil {
			return "", ErrPullRequestWithError
		}

		return "", prError
	}

	pr := f.PullRequests[headBranch]
	if pr != nil {
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
	}

	f.PullRequests[headBranch] = pr

	return pr.Url, nil
}
