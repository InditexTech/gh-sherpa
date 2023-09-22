package use_cases

import (
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
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

	var branch string
	canceled, err := cb.askUserForNewBranchName(&branch, issueTrackerProvider, args.IssueValue, *repo, args.UseDefaultValues)
	if err != nil {
		return err
	}

	if canceled {
		return nil
	}

	return checkoutBranch(branch, baseBranch, !args.NoFetchValue, cb.Git)
}

func checkoutBranch(branchName string, baseBranch string, fetch bool, git domain.GitProvider) error {
	if git.BranchExists(branchName) {
		return fmt.Errorf("a local branch with the name %s already exists", branchName)
	}

	if fetch {
		if err := git.FetchBranchFromOrigin(baseBranch); err != nil {
			return fmt.Errorf("error while fetching the branch %s: %s", baseBranch, err)
		}
	}

	if err := git.CheckoutNewBranchFromOrigin(branchName, baseBranch); err != nil {
		return err
	}

	fmt.Printf("A local branch named %s has been created!\n", logging.PaintInfo(branchName))

	return nil
}

func (cb CreateBranch) askUserForNewBranchName(branchName *string, issueTracker domain.IssueTracker, issueID string, repo domain.Repository, useDefaultValues bool) (cancelled bool, err error) {

	provs := providers{
		UserInteraction: cb.UserInteractionProvider,
	}
	if err := askBranchName(cb.BranchPrefixOverride, branchName, issueTracker, issueID, repo, useDefaultValues, provs); err != nil {
		return false, err
	}

	if !useDefaultValues {
		fmt.Println()
		fmt.Printf("A new local branch named %s is going to be created", logging.PaintInfo(*branchName))
		fmt.Println()

		confirmation, err := cb.UserInteractionProvider.AskUserForConfirmation("Do you want to continue?", true)
		if err != nil {
			return false, err
		}

		if !confirmation {
			return true, nil
		}
	}

	return false, nil
}
