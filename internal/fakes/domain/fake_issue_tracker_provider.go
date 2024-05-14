package domain

import (
	"errors"
	"strconv"
	"strings"

	"github.com/InditexTech/gh-sherpa/internal/domain"
)

type FakeIssueTrackerProvider struct {
	IssueTracker domain.IssueTracker
}

var _ domain.IssueTrackerProvider = (*FakeIssueTrackerProvider)(nil)

func (f *FakeIssueTrackerProvider) GetIssueTracker(identifier string) (issueTracker domain.IssueTracker, err error) {
	_, err = strconv.Atoi(identifier)
	if strings.HasPrefix(identifier, "GH-") || strings.HasPrefix(identifier, "PROJECTKEY-") || err == nil {
		return f.IssueTracker, nil
	}

	return nil, errors.New("issue tracker not found")
}

func (f *FakeIssueTrackerProvider) ParseIssueId(identifier string) (issueId string) {
	return strings.TrimPrefix(identifier, "GH-")
}
