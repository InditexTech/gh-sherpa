package use_cases

import (
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/logging"
)

type CreateBranchConfiguration struct {
	IssueID         string
	BaseBranch      string
	FetchFromOrigin bool
	IsInteractive   bool
}

type CreateBranch struct {
	Cfg                     CreateBranchConfiguration
	Git                     domain.GitProvider
	RepositoryProvider      domain.RepositoryProvider
	IssueTrackerProvider    domain.IssueTrackerProvider
	UserInteractionProvider domain.UserInteractionProvider
	BranchProvider          domain.BranchProvider
}

// Execute executes the create branch use case
func (cb CreateBranch) Execute() (err error) {
	if cb.Cfg.IssueID == "" {
		return fmt.Errorf("sherpa needs an valid issue identifier")
	}

	repo, err := cb.RepositoryProvider.GetRepository()
	if err != nil {
		return err
	}

	baseBranch := cb.Cfg.BaseBranch
	if baseBranch == "" {
		logging.Debugf("Base branch not set, using default branch, %s", repo.DefaultBranchRef)
		baseBranch = repo.DefaultBranchRef
	}

	issueTrackerProvider, err := cb.IssueTrackerProvider.GetIssueTracker(cb.Cfg.IssueID)
	if err != nil {
		return err
	}

	branchName, err := cb.BranchProvider.GetBranchName(issueTrackerProvider, cb.Cfg.IssueID, *repo)
	if err != nil {
		return err
	}

	fmt.Printf("\nA new local branch named %s is going to be created\n", logging.PaintInfo(branchName))
	if cb.Cfg.IsInteractive {
		confirmed, err := cb.UserInteractionProvider.AskUserForConfirmation("Do you want to continue?", true)
		if err != nil {
			return err
		}
		if !confirmed {
			return nil
		}
	}

	return cb.checkoutBranch(branchName, baseBranch, !cb.Cfg.FetchFromOrigin)
}

func (cb CreateBranch) checkoutBranch(branchName string, baseBranch string, fetch bool) error {
	if cb.Git.BranchExists(branchName) {
		return fmt.Errorf("a local branch with the name %s already exists", branchName)
	}

	if fetch {
		if err := cb.Git.FetchBranchFromOrigin(baseBranch); err != nil {
			return fmt.Errorf("error while fetching the branch %s: %s", baseBranch, err)
		}
	}

	if err := cb.Git.CheckoutNewBranchFromOrigin(branchName, baseBranch); err != nil {
		return err
	}

	fmt.Printf("A local branch named %s has been created!\n", logging.PaintInfo(branchName))

	return nil
}
