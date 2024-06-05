package domain

import (
	"errors"
	"strings"

	"github.com/InditexTech/gh-sherpa/internal/domain"
)

type FakeIssueTrackerProvider struct {
	Issues []domain.Issue
}

var _ domain.IssueTrackerProvider = (*FakeIssueTrackerProvider)(nil)

func NewFakeIssueTrackerProvider() *FakeIssueTrackerProvider {
	return &FakeIssueTrackerProvider{
		Issues: []domain.Issue{},
	}
}

var ErrNoIssue = errors.New("no issue")

func (f *FakeIssueTrackerProvider) AddIssue(issue domain.Issue) {
	f.Issues = append(f.Issues, issue)
}

func (f *FakeIssueTrackerProvider) GetIssue(identifier string) (issue domain.Issue, err error) {
	for _, i := range f.Issues {
		if i.ID() == identifier {
			return i, nil
		}
	}

	return nil, ErrNoIssue
}

func (f *FakeIssueTrackerProvider) ParseIssueId(identifier string) (issueId string) {
	return strings.TrimPrefix(identifier, "GH-")
}
