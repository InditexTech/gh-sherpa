package use_cases

import (
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	"github.com/InditexTech/gh-sherpa/internal/interactive"
	"github.com/InditexTech/gh-sherpa/internal/logging"
)

type CreateBranchArgs struct {
	IssueValue       string
	BaseValue        string
	NoFetchValue     bool
	UseDefaultValues bool
}

type CreateBranch struct {
	BranchPrefixOverride    map[issue_types.IssueType]string
	Git                     domain.GitProvider
	GhCli                   domain.GhCli
	IssueTrackerProvider    domain.IssueTrackerProvider
	UserInteractionProvider domain.UserInteractionProvider
	BranchProvider          domain.BranchProvider
}

// Execute executes the create branch use case
func (cb CreateBranch) Execute(args CreateBranchArgs) (err error) {
	if args.IssueValue == "" {
		return fmt.Errorf("sherpa needs an valid issue identifier")
	}

	repo, err := cb.GhCli.GetRepo()
	if err != nil {
		return err
	}

	baseBranch := args.BaseValue
	if baseBranch == "" {
		logging.Debugf("Base branch not set, using default branch, %s", repo.DefaultBranchRef)
		baseBranch = repo.DefaultBranchRef
	}

	issueTrackerProvider, err := cb.IssueTrackerProvider.GetIssueTracker(args.IssueValue)
	if err != nil {
		return err
	}

	branchName, err := cb.BranchProvider.GetBranchName(issueTrackerProvider, args.IssueValue, *repo, args.UseDefaultValues)
	if err != nil {
		return err
	}

	if err := cb.confirmBranchName(branchName, args.UseDefaultValues); err != nil {
		return err
	}

	return cb.checkoutBranch(branchName, baseBranch, !args.NoFetchValue)
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

func (cb CreateBranch) confirmBranchName(branchName string, useDefaultValues bool) (err error) {
	if !useDefaultValues {
		fmt.Println()
		fmt.Printf("A new local branch named %s is going to be created", logging.PaintInfo(branchName))
		fmt.Println()

		confirmed, err := cb.UserInteractionProvider.AskUserForConfirmation("Do you want to continue?", true)
		if err != nil {
			return err
		}

		if !confirmed {
			return interactive.ErrOpCanceled
		}
	}

	return nil
}
