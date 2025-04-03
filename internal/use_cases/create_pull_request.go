package use_cases

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/InditexTech/gh-sherpa/internal/branches"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/logging"
)

// ErrRemoteBranchAlreadyExists is returned when the remote branch already exists
func ErrRemoteBranchAlreadyExists(branchName string) error {
	return fmt.Errorf("there is already a remote branch named %s for this issue. Please checkout that branch", branchName)
}

// ErrFetchBranch is returned when the branch could not be fetched
func ErrFetchBranch(branch string, err error) error {
	return fmt.Errorf("could not fetch the branch %s: %w", branch, err)
}

// ErrPushChanges is returned when the remote branch could not be created
func ErrPushChanges(branch string, err error) error {
	return fmt.Errorf("could not push to remote branch %s: %w", branch, err)
}

// CreatePullRequestConfiguration contains the arguments for the CreatePullRequest use case
type CreatePullRequestConfiguration struct {
	IssueID         string
	BaseBranch      string
	FetchFromOrigin bool
	DraftPR         bool
	IsInteractive   bool
	CloseIssue      bool
	TemplatePath    string
}

type CreatePullRequest struct {
	Cfg                     CreatePullRequestConfiguration
	Git                     domain.GitProvider
	RepositoryProvider      domain.RepositoryProvider
	IssueTrackerProvider    domain.IssueTrackerProvider
	UserInteractionProvider domain.UserInteractionProvider
	PullRequestProvider     domain.PullRequestProvider
	BranchProvider          domain.BranchProvider
}

// Execute executes the create pull request use case
func (cpr CreatePullRequest) Execute() error {
	isInteractive := cpr.Cfg.IsInteractive
	fromLocalBranch := cpr.Cfg.IssueID == ""

	repo, err := cpr.RepositoryProvider.GetRepository()
	if err != nil {
		logging.Debugf("error while getting repo: %s", err)
		return err
	}

	baseBranch := cpr.Cfg.BaseBranch
	if baseBranch == "" {
		baseBranch = repo.DefaultBranchRef
	}
	if err := cpr.fetchBranch(baseBranch); err != nil {
		return err
	}

	currentBranch, err := cpr.Git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("could not get the current branch name because %s", err)
	}

	var issueID string
	if fromLocalBranch {
		issueID, err = cpr.extractIssueIdFromBranch(currentBranch)
		if err != nil {
			return err
		}
	} else {
		issueID = cpr.Cfg.IssueID
	}

	issue, err := cpr.IssueTrackerProvider.GetIssue(issueID)
	if err != nil {
		return err
	}

	var branchExists bool
	if fromLocalBranch {
		branchExists = true
	} else {
		formattedIssueID := issue.FormatID()
		currentBranch, branchExists = cpr.Git.FindBranch(fmt.Sprintf("/%s-", formattedIssueID))

		if branchExists {
			if !isInteractive {
				return fmt.Errorf("the branch %s already exists", logging.PaintWarning(currentBranch))
			}

			logging.PrintWarn(
				fmt.Sprintf("there is already a local branch named %s for this issue", logging.PaintInfo(currentBranch)))
		}
	}

	branchConfirmed := true
	if branchExists && isInteractive {
		branchConfirmed, err = cpr.UserInteractionProvider.AskUserForConfirmation(
			"Do you want to use this branch to create the pull request", true)
		if err != nil {
			return err
		}

		if fromLocalBranch && !branchConfirmed {
			// Clean exit since user do not want to continue
			return nil
		}
	}

	// Fetch branch to get the latest changes. We don't check the error because
	// the branch may not exist in the remote repository.
	_ = cpr.fetchBranch(currentBranch)

	if branchExists && branchConfirmed {
		if err := cpr.Git.CheckoutBranch(currentBranch); err != nil {
			return fmt.Errorf("could not switch to the branch because %w", err)
		}
	} else {
		currentBranch, err = cpr.BranchProvider.GetBranchName(issue, *repo)
		if err != nil {
			return err
		}

		// After stablishing the branch, fetch it to get the latest changes.
		// We don't check the error because the branch may not exist in the
		// remote repository.
		_ = cpr.fetchBranch(currentBranch)

		cancel, err := cpr.createBranch(currentBranch, baseBranch)
		if err != nil {
			return err
		}
		if cancel {
			return nil
		}
	}

	hasPendingCommits, err := cpr.hasPendingCommits(currentBranch)
	if err != nil {
		return err
	}

	if hasPendingCommits {
		logging.PrintWarn("the branch contains commits that have not been pushed yet")

		confirmed, err := cpr.UserInteractionProvider.AskUserForConfirmation(
			"Do you want to continue pushing all pending commits in this branch and create the pull request", true)
		if err != nil {
			return err
		}
		if !confirmed {
			return nil
		}
	} else if !cpr.Git.RemoteBranchExists(currentBranch) {
		if err = cpr.createEmptyCommit(); err != nil {
			return fmt.Errorf("could not do the empty commit because %s", err)
		}
	}

	if err := cpr.pushChanges(currentBranch); err != nil {
		return err
	}

	pr, err := cpr.PullRequestProvider.GetPullRequestForBranch(currentBranch)
	if err != nil {
		return fmt.Errorf("error while getting pull request for branch: %w", err)
	}

	if pr != nil && !pr.Closed {
		return fmt.Errorf("a pull request %s for this branch already exists", pr.Url)
	}

	if err := cpr.createPullRequestFromIssue(issue, baseBranch, currentBranch); err != nil {
		return err
	}

	return nil
}

func (cpr *CreatePullRequest) extractIssueIdFromBranch(currentBranch string) (string, error) {
	branchNameInfo := branches.ParseBranchName(currentBranch)
	if branchNameInfo == nil || branchNameInfo.IssueId == "" {
		return "", fmt.Errorf("could not find an issue identifier in the current branch named %s", logging.PaintWarning(currentBranch))
	}

	logging.PrintInfo(fmt.Sprintf("The current branch named %s is available to create a pull request", logging.PaintWarning(currentBranch)))

	return cpr.IssueTrackerProvider.ParseIssueId(branchNameInfo.IssueId), nil
}

func (cpr *CreatePullRequest) createEmptyCommit() error {
	return cpr.Git.CommitEmpty("chore: initial commit")
}

func (cpr *CreatePullRequest) pushChanges(branchName string) (err error) {
	// 19. PUSH CHANGES
	err = cpr.Git.PushBranch(branchName)
	if err != nil {
		return ErrPushChanges(branchName, err)
	}

	return
}

func (cpr *CreatePullRequest) getPullRequestTitleAndBody(issue domain.Issue) (title string, body string, err error) {
	// Determine base title and issue reference body based on tracker type
	switch issue.TrackerType() {
	case domain.IssueTrackerTypeGithub:
		title = issue.Title()

		keyword := "Related to"
		if cpr.Cfg.CloseIssue {
			keyword = "Closes"
		}
		body = fmt.Sprintf("%s #%s", keyword, issue.ID())

	case domain.IssueTrackerTypeJira:
		title = fmt.Sprintf("[%s] %s", issue.ID(), issue.Title())
		body = fmt.Sprintf("Relates to [%s](%s)", issue.ID(), issue.URL())

	default:
		return "", "", fmt.Errorf("issue tracker %s is not supported", issue.TrackerType())
	}

	// If template path is provided, read and append it after the issue reference
	if cpr.Cfg.TemplatePath != "" {
		var templateContent []byte

		// If not an absolute path, resolve from repository root
		templatePath := cpr.Cfg.TemplatePath
		if !filepath.IsAbs(templatePath) {
			repoRoot, err := cpr.Git.GetRepositoryRoot()
			if err != nil {
				return "", "", fmt.Errorf("failed to determine repository root: %w", err)
			}
			templatePath = filepath.Join(repoRoot, templatePath)
		}

		// Check if template file exists
		if _, statErr := os.Stat(templatePath); os.IsNotExist(statErr) {
			return "", "", fmt.Errorf("template file does not exist: %s", cpr.Cfg.TemplatePath)
		}

		// Read template file
		templateContent, err = os.ReadFile(templatePath)
		if err != nil {
			return "", "", fmt.Errorf("failed to read template file: %w", err)
		}

		// Reference to issue first, then template content
		body = body + "\n\n" + string(templateContent)
	}

	return
}

func (cpr *CreatePullRequest) hasPendingCommits(currentBranch string) (bool, error) {
	commitsToPush, err := cpr.Git.GetCommitsToPush(currentBranch)
	if err != nil {
		return false, err
	}

	return len(commitsToPush) > 0, nil
}

func (cpr *CreatePullRequest) createNewLocalBranch(currentBranch string, baseBranch string) error {
	// Create the new branch from the base branch
	if err := cpr.Git.CheckoutNewBranchFromOrigin(currentBranch, baseBranch); err != nil {
		return fmt.Errorf("could not create the local branch because %s", err)
	}

	return nil
}

func (cpr *CreatePullRequest) fetchBranch(branch string) error {
	if cpr.Cfg.FetchFromOrigin {
		if err := cpr.Git.FetchBranchFromOrigin(branch); err != nil {
			return ErrFetchBranch(branch, err)
		}
	}

	return nil
}

func (cpr *CreatePullRequest) createBranch(branch string, baseBranch string) (cancel bool, err error) {
	if cpr.Git.RemoteBranchExists(branch) {
		err = ErrRemoteBranchAlreadyExists(branch)
		return
	}

	fmt.Printf("\nA new pull request is going to be created from %s to %s branch\n",
		logging.PaintInfo(branch), logging.PaintInfo(baseBranch))

	if cpr.Cfg.IsInteractive {
		var confirmed bool
		confirmed, err = cpr.UserInteractionProvider.AskUserForConfirmation("Do you want to continue?", true)
		if err != nil {
			return
		}
		if !confirmed {
			cancel = true
			return
		}
	}

	if err = cpr.createNewLocalBranch(branch, baseBranch); err != nil {
		return
	}

	return
}

func (cpr *CreatePullRequest) createPullRequestFromIssue(issue domain.Issue, baseBranch string, headBranch string) error {
	title, body, err := cpr.getPullRequestTitleAndBody(issue)
	if err != nil {
		return err
	}

	labels := []string{}
	typeLabel := issue.TypeLabel()
	if typeLabel != "" {
		labels = append(labels, typeLabel)
	}

	prURL, err := cpr.PullRequestProvider.CreatePullRequest(title, body, baseBranch, headBranch, cpr.Cfg.DraftPR, labels)
	if err != nil {
		return fmt.Errorf("could not create the pull request because %s", err)
	}

	fmt.Printf("\nThe pull request %s have been created!\nYou are now working on the branch %s\n",
		logging.PaintInfo(prURL), logging.PaintInfo(headBranch))

	return nil
}
