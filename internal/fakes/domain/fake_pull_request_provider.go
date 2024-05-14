package domain

import (
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/domain"
)

type FakePullRequestProvider struct {
	PullRequests map[string]*domain.PullRequest
}

var _ domain.PullRequestProvider = (*FakePullRequestProvider)(nil)

func NewFakePullRequestProvider() *FakePullRequestProvider {
	return &FakePullRequestProvider{
		PullRequests: map[string]*domain.PullRequest{
			"feature/GH-3-pull-request-sample": {
				Title:       "GH-3-pull-request-sample",
				Number:      3,
				State:       "OPEN",
				Closed:      false,
				Url:         "https://github.com/inditextech/gh-sherpa-test-repo/pulls/3",
				HeadRefName: "feature/GH-3-pull-request-sample",
				BaseRefName: "main",
				Labels: []domain.Label{
					{
						Id:   "kind/feature",
						Name: "kind/feature",
					},
				},
			},
			"feature/PROJECTKEY-3-pull-request-sample": {
				Title:       "PROJECTKEY-4-pull-request-sample",
				Number:      3,
				State:       "OPEN",
				Closed:      false,
				Url:         "https://github.com/inditextech/gh-sherpa-test-repo/pulls/3",
				HeadRefName: "feature/PROJECTKEY-3-pull-request-sample",
				BaseRefName: "main",
				Labels: []domain.Label{
					{
						Id:   "kind/feature",
						Name: "kind/feature",
					},
				},
			},
			"feature/GH-1-sample-issue":         nil,
			"feature/PROJECTKEY-1-sample-issue": nil,
			"feature/GH-3-local-branch":         nil,
			"feature/PROJECTKEY-3-local-branch": nil,
		},
	}
}

func (f *FakePullRequestProvider) HasPullRequest(branch string) bool {
	pr := f.PullRequests[branch]

	return pr != nil
}

func (f *FakePullRequestProvider) GetPullRequestForBranch(branch string) (pullRequest *domain.PullRequest, err error) {
	pr, ok := f.PullRequests[branch]
	if !ok {
		return nil, nil
	}

	return pr, nil
}

func ErrPrAlreadyExists(branch string) error {
	return fmt.Errorf("pr already exists for branch %s", branch)
}

func ErrBranchNotDefined(branch string) error {
	return fmt.Errorf("branch %s is not defined in the pull request map", branch)
}

func (f *FakePullRequestProvider) CreatePullRequest(title string, body string, baseBranch string, headBranch string, draft bool, labels []string) (prUrl string, err error) {
	pr, ok := f.PullRequests[headBranch]
	if pr != nil {
		return "", ErrPrAlreadyExists(headBranch)
	}
	if !ok {
		return "", ErrBranchNotDefined(headBranch)
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
