package use_cases

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
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

// CreatePullRequestResult holds the outcome of a successful CreatePullRequest execution.
type CreatePullRequestResult struct {
	BranchName string `json:"branch"`
	PRURL      string `json:"pr_url"`
	Draft      bool   `json:"draft"`
}

// CreatePullRequestConfiguration contains the arguments for the CreatePullRequest use case
type CreatePullRequestConfiguration struct {
	IssueID             string
	BaseBranch          string
	FetchFromOrigin     bool
	DraftPR             bool
	IsInteractive       bool
	CloseIssue          bool
	TemplatePath        string
	BranchName          string   // --branch-name: bypass generation and use this name directly
	DryRun              bool     // --dry-run: print what would happen without executing
	OutputFormat        string   // --output: "" (default) or "json"
	PRTitle             string   // --pr-title: override the auto-generated PR title
	PRBody              string   // --pr-body: override the auto-generated PR body
	PRBodyFile          string   // --pr-body-file: read PR body from file
	NoUseExistingBranch bool     // --no-use-existing-branch: fail if a branch already exists (non-interactive)
	ExtraLabels         []string // --label: additional labels to apply to the PR
	Reviewers           []string // --reviewer: PR reviewers
	Assignees           []string // --assignee: PR assignees
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

func normalizeTemplatePath(git domain.GitProvider, path string) (string, error) {
	if path == "" {
		return "", nil
	}

	if !filepath.IsAbs(path) {
		repoRoot, err := git.GetRepositoryRoot()
		if err != nil {
			return "", fmt.Errorf("failed to determine repository root: %w", err)
		}
		path = filepath.Join(repoRoot, path)
	}

	return path, nil
}

func validateTemplateFile(templatePath string, git domain.GitProvider) error {
	if templatePath == "" {
		return nil // No template to validate
	}

	normalizedPath, err := normalizeTemplatePath(git, templatePath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(normalizedPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("template file does not exist: %s", templatePath)
		}

		return fmt.Errorf("error accessing template file %s: %w", templatePath, err)
	}

	return nil
}

// Execute executes the create pull request use case
func (cpr CreatePullRequest) Execute() (result CreatePullRequestResult, err error) {
	// Validate template if specified
	if err := validateTemplateFile(cpr.Cfg.TemplatePath, cpr.Git); err != nil {
		return result, err
	}

	isInteractive := cpr.Cfg.IsInteractive
	// --output json implies non-interactive: no stdin prompts must appear in JSON output.
	if cpr.Cfg.OutputFormat == "json" {
		isInteractive = false
	}
	// Normalize so helper methods that read cpr.Cfg.IsInteractive also see the effective value.
	cpr.Cfg.IsInteractive = isInteractive
	fromLocalBranch := cpr.Cfg.IssueID == ""

	repo, err := cpr.RepositoryProvider.GetRepository()
	if err != nil {
		logging.Debugf("error while getting repo: %s", err)
		return result, err
	}

	baseBranch := cpr.Cfg.BaseBranch
	if baseBranch == "" {
		baseBranch = repo.DefaultBranchRef
	}
	if err := cpr.fetchBranch(baseBranch); err != nil {
		return result, err
	}

	currentBranch, err := cpr.Git.GetCurrentBranch()
	if err != nil {
		return result, fmt.Errorf("could not get the current branch name because %s", err)
	}

	var issueID string
	if fromLocalBranch {
		issueID, err = cpr.extractIssueIdFromBranch(currentBranch)
		if err != nil {
			return result, err
		}
	} else {
		issueID = cpr.Cfg.IssueID
	}

	issue, err := cpr.IssueTrackerProvider.GetIssue(issueID)
	if err != nil {
		return result, err
	}

	var branchExists bool
	if fromLocalBranch {
		branchExists = true
	} else {
		formattedIssueID := issue.FormatID()
		currentBranch, branchExists = cpr.Git.FindBranch(fmt.Sprintf("/%s-", formattedIssueID))

		if branchExists {
			if !isInteractive {
				if cpr.Cfg.NoUseExistingBranch {
					return result, fmt.Errorf("the branch %s already exists", logging.PaintWarning(currentBranch))
				}
				// Default non-interactive behavior: reuse the existing branch silently
				logging.Debugf("Reusing existing branch %s", currentBranch)
			} else {
				logging.PrintWarn(
					fmt.Sprintf("there is already a local branch named %s for this issue", logging.PaintInfo(currentBranch)))
			}
		}
	}

	branchConfirmed := true
	if branchExists && isInteractive {
		branchConfirmed, err = cpr.UserInteractionProvider.AskUserForConfirmation(
			"Do you want to use this branch to create the pull request", true)
		if err != nil {
			return result, err
		}

		if fromLocalBranch && !branchConfirmed {
			// Clean exit since user do not want to continue
			return result, nil
		}
	}

	// Fetch branch to get the latest changes. We don't check the error because
	// the branch may not exist in the remote repository.
	_ = cpr.fetchBranch(currentBranch)

	if branchExists && branchConfirmed {
		if err := cpr.Git.CheckoutBranch(currentBranch); err != nil {
			return result, fmt.Errorf("could not switch to the branch because %w", err)
		}
	} else {
		currentBranch, err = cpr.BranchProvider.GetBranchName(issue, *repo)
		if err != nil {
			return result, err
		}

		// After stablishing the branch, fetch it to get the latest changes.
		// We don't check the error because the branch may not exist in the
		// remote repository.
		_ = cpr.fetchBranch(currentBranch)

		cancel, err := cpr.createBranch(currentBranch, baseBranch)
		if err != nil {
			return result, err
		}
		if cancel {
			return result, nil
		}
	}

	hasPendingCommits, err := cpr.hasPendingCommits(currentBranch)
	if err != nil {
		return result, err
	}

	if hasPendingCommits {
		if cpr.Cfg.OutputFormat != "json" {
			logging.PrintWarn("the branch contains commits that have not been pushed yet")
		}

		if isInteractive {
			confirmed, err := cpr.UserInteractionProvider.AskUserForConfirmation(
				"Do you want to continue pushing all pending commits in this branch and create the pull request", true)
			if err != nil {
				return result, err
			}
			if !confirmed {
				return result, nil
			}
		}
	}

	// Early exit for dry-run: report what would happen without touching the remote.
	if cpr.Cfg.DryRun {
		result.BranchName = currentBranch
		result.Draft = cpr.Cfg.DraftPR
		if cpr.Cfg.OutputFormat == "json" {
			jsonBytes, jsonErr := json.Marshal(result)
			if jsonErr != nil {
				return result, fmt.Errorf("failed to serialize dry-run result: %w", jsonErr)
			}
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Printf("[dry-run] Would create PR from %s to %s\n",
				logging.PaintInfo(currentBranch), logging.PaintInfo(baseBranch))
		}
		return result, nil
	}

	if !hasPendingCommits && !cpr.Git.RemoteBranchExists(currentBranch) {
		if err = cpr.createEmptyCommit(); err != nil {
			return result, fmt.Errorf("could not do the empty commit because %s", err)
		}
	}

	if err := cpr.pushChanges(currentBranch); err != nil {
		return result, err
	}

	pr, err := cpr.PullRequestProvider.GetPullRequestForBranch(currentBranch)
	if err != nil {
		return result, fmt.Errorf("error while getting pull request for branch: %w", err)
	}

	if pr != nil && !pr.Closed {
		return result, fmt.Errorf("a pull request %s for this branch already exists", pr.Url)
	}

	prURL, err := cpr.createPullRequestFromIssue(issue, baseBranch, currentBranch)
	if err != nil {
		return result, err
	}

	result.BranchName = currentBranch
	result.PRURL = prURL
	result.Draft = cpr.Cfg.DraftPR

	if cpr.Cfg.OutputFormat == "json" {
		jsonBytes, jsonErr := json.Marshal(result)
		if jsonErr != nil {
			return result, fmt.Errorf("failed to serialize result: %w", jsonErr)
		}
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("\nThe pull request %s have been created!\nYou are now working on the branch %s\n",
			logging.PaintInfo(prURL), logging.PaintInfo(currentBranch))
	}

	return result, nil
}

func (cpr *CreatePullRequest) extractIssueIdFromBranch(currentBranch string) (string, error) {
	branchNameInfo := branches.ParseBranchName(currentBranch)
	if branchNameInfo == nil || branchNameInfo.IssueId == "" {
		return "", fmt.Errorf("could not find an issue identifier in the current branch named %s", logging.PaintWarning(currentBranch))
	}

	if cpr.Cfg.OutputFormat != "json" {
		logging.PrintInfo(fmt.Sprintf("The current branch named %s is available to create a pull request", logging.PaintWarning(currentBranch)))
	}

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
	// --pr-title and --pr-body fully override auto-generation
	if cpr.Cfg.PRTitle != "" {
		title = cpr.Cfg.PRTitle
		body = cpr.Cfg.PRBody
		return
	}

	// --pr-body-file overrides the body
	if cpr.Cfg.PRBodyFile != "" {
		content, readErr := os.ReadFile(cpr.Cfg.PRBodyFile)
		if readErr != nil {
			err = fmt.Errorf("failed to read --pr-body-file: %w", readErr)
			return
		}
		body = string(content)
		// Still auto-generate the title from the issue
		switch issue.TrackerType() {
		case domain.IssueTrackerTypeGithub:
			title = issue.Title()
		case domain.IssueTrackerTypeJira:
			title = fmt.Sprintf("[%s] %s", issue.ID(), issue.Title())
		default:
			err = fmt.Errorf("issue tracker %s is not supported", issue.TrackerType())
		}
		return
	}

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

	// --pr-body overrides only the body (when --pr-title is not set)
	if cpr.Cfg.PRBody != "" {
		body = cpr.Cfg.PRBody
		return
	}

	// If template path is provided, read and append it after the issue reference
	// (we already validated that the file exists in validateTemplateFile)
	if cpr.Cfg.TemplatePath != "" {
		var templateContent []byte

		// Get normalized template path (absolute path)
		templatePath, err := normalizeTemplatePath(cpr.Git, cpr.Cfg.TemplatePath)
		if err != nil {
			return "", "", err
		}

		// Read template file (we know it exists because we checked in validateTemplateFile)
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

	if cpr.Cfg.OutputFormat != "json" {
		fmt.Printf("\nA new pull request is going to be created from %s to %s branch\n",
			logging.PaintInfo(branch), logging.PaintInfo(baseBranch))
	}

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

func (cpr *CreatePullRequest) createPullRequestFromIssue(issue domain.Issue, baseBranch string, headBranch string) (prURL string, err error) {
	title, body, err := cpr.getPullRequestTitleAndBody(issue)
	if err != nil {
		return "", err
	}

	labels := []string{}
	typeLabel := issue.TypeLabel()
	if typeLabel != "" {
		labels = append(labels, typeLabel)
	}
	labels = append(labels, cpr.Cfg.ExtraLabels...)

	prURL, err = cpr.PullRequestProvider.CreatePullRequest(title, body, baseBranch, headBranch, cpr.Cfg.DraftPR, labels, cpr.Cfg.Reviewers, cpr.Cfg.Assignees)
	if err != nil {
		return "", fmt.Errorf("could not create the pull request because %s", err)
	}

	return prURL, nil
}
