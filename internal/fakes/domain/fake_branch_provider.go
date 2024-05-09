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

func (f *FakeBranchProvider) GetBranchName(issueTracker domain.IssueTracker, issueIdentifier string, repo domain.Repository) (branchName string, err error) {
	switch issueTracker.GetIssueTrackerType() {
	case domain.IssueTrackerTypeGithub:
		return fmt.Sprintf("feature/GH-%s-generated-branch-name", issueIdentifier), nil
	case domain.IssueTrackerTypeJira:
		return fmt.Sprintf("feature/%s-generated-branch-name", issueIdentifier), nil
	default:

	}

	return "", errors.New("unknown issue tracker type")
}
