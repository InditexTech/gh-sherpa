package use_cases

import (
	"encoding/json"
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/logging"
)

// CreateBranchResult holds the outcome of a successful CreateBranch execution.
type CreateBranchResult struct {
	BranchName string `json:"branch"`
}

type CreateBranchConfiguration struct {
	IssueID          string
	BaseBranch       string
	FetchFromOrigin  bool
	IsInteractive    bool
	BranchName       string // --branch-name: bypass generation and use this name directly
	DryRun           bool   // --dry-run: print what would happen without executing
	OutputFormat     string // --output: "" (default) or "json"
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
func (cb CreateBranch) Execute() (result CreateBranchResult, err error) {
	if cb.Cfg.IssueID == "" {
		return result, fmt.Errorf("sherpa needs an valid issue identifier")
	}

	repo, err := cb.RepositoryProvider.GetRepository()
	if err != nil {
		return result, err
	}

	baseBranch := cb.Cfg.BaseBranch
	if baseBranch == "" {
		logging.Debugf("Base branch not set, using default branch, %s", repo.DefaultBranchRef)
		baseBranch = repo.DefaultBranchRef
	}

	var branchName string
	if cb.Cfg.BranchName != "" {
		branchName = cb.Cfg.BranchName
	} else {
		issue, err := cb.IssueTrackerProvider.GetIssue(cb.Cfg.IssueID)
		if err != nil {
			return result, err
		}

		branchName, err = cb.BranchProvider.GetBranchName(issue, *repo)
		if err != nil {
			return result, err
		}
	}

	result.BranchName = branchName

	if cb.Cfg.DryRun {
		if cb.Cfg.OutputFormat == "json" {
			jsonBytes, jsonErr := json.Marshal(result)
			if jsonErr != nil {
				return result, fmt.Errorf("failed to serialize dry-run result: %w", jsonErr)
			}
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Printf("[dry-run] Would create branch: %s from %s\n", logging.PaintInfo(branchName), logging.PaintInfo(baseBranch))
		}
		return result, nil
	}

	if cb.Cfg.OutputFormat != "json" {
		fmt.Printf("\nA new local branch named %s is going to be created\n", logging.PaintInfo(branchName))
	}
	if cb.Cfg.IsInteractive {
		confirmed, err := cb.UserInteractionProvider.AskUserForConfirmation("Do you want to continue?", true)
		if err != nil {
			return result, err
		}
		if !confirmed {
			return result, nil
		}
	}

	if err := cb.checkoutBranch(branchName, baseBranch, !cb.Cfg.FetchFromOrigin); err != nil {
		return result, err
	}

	if cb.Cfg.OutputFormat == "json" {
		jsonBytes, jsonErr := json.Marshal(result)
		if jsonErr != nil {
			return result, fmt.Errorf("failed to serialize result: %w", jsonErr)
		}
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("A local branch named %s has been created!\n", logging.PaintInfo(branchName))
	}

	return result, nil
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

	return nil
}
