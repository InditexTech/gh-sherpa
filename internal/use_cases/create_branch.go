package use_cases

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/logging"
)

// CreateBranchResult holds the outcome of a successful CreateBranch execution.
type CreateBranchResult struct {
	BranchName   string `json:"branch"`
	WorktreePath string `json:"worktree_path,omitempty"`
}

type CreateBranchConfiguration struct {
	IssueID         string
	BaseBranch      string
	FetchFromOrigin bool
	IsInteractive   bool
	BranchName      string // --branch-name: bypass generation and use this name directly
	DryRun          bool   // --dry-run: print what would happen without executing
	OutputFormat    string // --output: "" (default) or "json"
	UseWorktree     bool
	WorktreePath    string
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
	if cb.Cfg.UseWorktree || cb.Cfg.WorktreePath != "" {
		result.WorktreePath, err = cb.resolveWorktreePath(branchName)
		if err != nil {
			return result, err
		}
	}

	if cb.Cfg.DryRun {
		if cb.Cfg.OutputFormat == "json" {
			jsonBytes, jsonErr := json.Marshal(result)
			if jsonErr != nil {
				return result, fmt.Errorf("failed to serialize dry-run result: %w", jsonErr)
			}
			fmt.Println(string(jsonBytes))
		} else if result.WorktreePath != "" {
			logging.Info(fmt.Sprintf("[dry-run] Would create worktree for branch %s from %s at %s",
				branchName, baseBranch, result.WorktreePath))
		} else {
			fmt.Printf("[dry-run] Would create branch: %s from %s\n", logging.PaintInfo(branchName), logging.PaintInfo(baseBranch))
		}
		return result, nil
	}

	if cb.Cfg.OutputFormat != "json" && result.WorktreePath != "" {
		logging.Info(fmt.Sprintf("\nA new worktree for local branch %s is going to be created at %s",
			branchName, result.WorktreePath))
	} else if cb.Cfg.OutputFormat != "json" {
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

	if err := cb.createBranch(branchName, baseBranch, result.WorktreePath); err != nil {
		return result, err
	}

	if cb.Cfg.OutputFormat == "json" {
		jsonBytes, jsonErr := json.Marshal(result)
		if jsonErr != nil {
			return result, fmt.Errorf("failed to serialize result: %w", jsonErr)
		}
		fmt.Println(string(jsonBytes))
	} else if result.WorktreePath != "" {
		logging.Info(fmt.Sprintf("A worktree for local branch %s has been created at %s!",
			branchName, result.WorktreePath))
	} else {
		fmt.Printf("A local branch named %s has been created!\n", logging.PaintInfo(branchName))
	}

	return result, nil
}

func (cb CreateBranch) createBranch(branchName string, baseBranch string, worktreePath string) error {
	if cb.Git.BranchExists(branchName) {
		return fmt.Errorf("a local branch with the name %s already exists", branchName)
	}

	if cb.Cfg.FetchFromOrigin {
		if err := cb.Git.FetchBranchFromOrigin(baseBranch); err != nil {
			return fmt.Errorf("error while fetching the branch %s: %s", baseBranch, err)
		}
	}

	if worktreePath != "" {
		return cb.Git.CreateWorktree(worktreePath, branchName, baseBranch)
	}

	if err := cb.Git.CheckoutNewBranchFromOrigin(branchName, baseBranch); err != nil {
		return err
	}

	return nil
}

func (cb CreateBranch) resolveWorktreePath(branchName string) (string, error) {
	worktreePath := cb.Cfg.WorktreePath
	if worktreePath == "" {
		repositoryRoot, err := cb.Git.GetRepositoryRoot()
		if err != nil {
			return "", err
		}
		worktreePath = filepath.Join(repositoryRoot, "..", "worktrees", filepath.FromSlash(branchName))
		return filepath.Clean(worktreePath), nil
	}

	absolutePath, err := filepath.Abs(worktreePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve worktree path: %s", err)
	}

	return filepath.Clean(absolutePath), nil
}
