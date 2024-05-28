package domain

import (
	"errors"
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/domain"
)

type FakeBranchProvider struct{}

var _ domain.BranchProvider = (*FakeBranchProvider)(nil)

func NewFakeBranchProvider() FakeBranchProvider {
	return FakeBranchProvider{}
}

func (f *FakeBranchProvider) GetBranchName(issue domain.Issue, repo domain.Repository) (branchName string, err error) {
	switch issue.TrackerType() {
	case domain.IssueTrackerTypeGithub:
		return fmt.Sprintf("feature/GH-%s-generated-branch-name", issue.ID()), nil
	case domain.IssueTrackerTypeJira:
		return fmt.Sprintf("feature/%s-generated-branch-name", issue.ID()), nil
	default:

	}

	return "", errors.New("unknown issue tracker type")
}
