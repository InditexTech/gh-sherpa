package use_cases

import (
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/branches"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/logging"
)

// ErrRemoteBranchAlreadyExists is returned when the remote branch already exists
func ErrRemoteBranchAlreadyExists(branchName string) error {
	return fmt.Errorf("there is already a remote branch named %s for this issue. Please checkout that branch", branchName)
}

// CreatePullRequestConfiguration contains the arguments for the CreatePullRequest use case
type CreatePullRequestConfiguration struct {
	IssueID         string
	BaseBranch      string
	FetchFromOrigin bool
	DraftPR         bool
	IsInteractive   bool
	CloseIssue      bool
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

	repo, err := cpr.RepositoryProvider.GetRepository()
	if err != nil {
		logging.Debugf("error while getting repo: %s", err)
		return err
	}

	baseBranch := cpr.Cfg.BaseBranch
	if baseBranch == "" {
		baseBranch = repo.DefaultBranchRef
	}

	currentBranch, err := cpr.Git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("could not get the current branch name because %s", err)
	}

	isIssueIdProvided := cpr.Cfg.IssueID != ""

	// Extract issue
	issue, err := cpr.extractIssue(currentBranch)
	if err != nil {
		return err
	}

	// Check if a local branch already exists for the given issue
	branchExists := true
	if isIssueIdProvided {
		formattedIssueID := issue.FormatID()
		currentBranch, branchExists = cpr.Git.BranchExistsContains(fmt.Sprintf("/%s-", formattedIssueID))

		if branchExists && !isInteractive {
			return fmt.Errorf("the branch %s already exists", logging.PaintWarning(currentBranch))
		}

		if branchExists && isInteractive {
			logging.PrintWarn(fmt.Sprintf("there is already a local branch named %s for this issue", logging.PaintInfo(currentBranch)))
		}
	}

	// Confirm usage of the existing branch
	confirmUseExistingBranch := true
	if branchExists && isInteractive {
		confirmUseExistingBranch, err = cpr.UserInteractionProvider.AskUserForConfirmation("Do you want to use this branch to create the pull request", true)
		if err != nil {
			return err
		}

		if !confirmUseExistingBranch {
			if !isIssueIdProvided {
				return nil
			}

		} else {
			if err := cpr.Git.CheckoutBranch(currentBranch); err != nil {
				return fmt.Errorf("could not switch to the branch because %w", err)
			}
		}
	}

	// Create brand new local branch if the branch does not exists or the user wants to create a new branch
	if !branchExists || !confirmUseExistingBranch {
		currentBranch, err = cpr.BranchProvider.GetBranchName(issue, *repo)
		if err != nil {
			return err
		}

		if cpr.Git.RemoteBranchExists(currentBranch) {
			return ErrRemoteBranchAlreadyExists(currentBranch)
		}

		fmt.Printf("\nA new pull request is going to be created from %s to %s branch\n", logging.PaintInfo(currentBranch), logging.PaintInfo(baseBranch))
		if isInteractive {
			confirmed, err := cpr.UserInteractionProvider.AskUserForConfirmation("Do you want to continue?", true)
			if err != nil {
				return err
			}
			if !confirmed {
				return nil
			}
		}

		if err := cpr.createNewLocalBranch(currentBranch, baseBranch); err != nil {
			return err
		}
	}

	// <-------------------------------------------------------------------->

	if cpr.Git.RemoteBranchExists(currentBranch) {
		return ErrRemoteBranchAlreadyExists(currentBranch)
	}

	// 11. CHECK IF BRANCH HAS PENDING COMMITS
	hasPendingCommits, err := cpr.hasPendingCommits(currentBranch)
	if err != nil {
		return err
	}

	if hasPendingCommits {
		logging.PrintWarn("the branch contains commits that have not been pushed yet")

		confirmed, err := cpr.UserInteractionProvider.AskUserForConfirmation("Do you want to continue pushing all pending commits in this branch and create the pull request", true)
		if err != nil {
			return err
		}

		if !confirmed {
			return nil
		}
	} else {
		if err = cpr.createEmptyCommit(); err != nil {
			return fmt.Errorf("could not do the empty commit because %s", err)
		}
	}

	if err := cpr.pushChanges(currentBranch); err != nil {
		return err
	}

	// <-------------------------------------------------------------------->

	// 13. CHECK IF PR DOES ALREADY EXISTS
	pr, err := cpr.PullRequestProvider.GetPullRequestForBranch(currentBranch)
	if err != nil {
		return fmt.Errorf("error while getting pull request for branch: %w", err)
	}

	if pr != nil && !pr.Closed {
		// 14. EXIT
		return fmt.Errorf("a pull request %s for this branch already exists", pr.Url)
	}

	title, body, err := cpr.getPullRequestTitleAndBody(issue)
	if err != nil {
		return err
	}

	labels := []string{}
	typeLabel := issue.TypeLabel()
	if typeLabel != "" {
		labels = append(labels, typeLabel)
	}

	// 16. CREATE PULL REQUEST
	prURL, err := cpr.PullRequestProvider.CreatePullRequest(title, body, baseBranch, currentBranch, cpr.Cfg.DraftPR, labels)
	if err != nil {
		return fmt.Errorf("could not create the pull request because %s", err)
	}

	fmt.Printf("\nThe pull request %s have been created!\nYou are now working on the branch %s\n", logging.PaintInfo(prURL), logging.PaintInfo(currentBranch))
	// 17. EXIT
	return nil
}

func (cpr *CreatePullRequest) extractIssue(currentBranch string) (domain.Issue, error) {
	issueID := cpr.Cfg.IssueID

	if issueID == "" {
		branchNameInfo := branches.ParseBranchName(currentBranch)
		if branchNameInfo == nil || branchNameInfo.IssueId == "" {
			return nil, fmt.Errorf("could not find an issue identifier in the current branch named %s", logging.PaintWarning(currentBranch))
		}

		logging.PrintInfo(fmt.Sprintf("The current branch named %s is available to create a pull request", logging.PaintWarning(currentBranch)))

		issueID = cpr.IssueTrackerProvider.ParseIssueId(branchNameInfo.IssueId)
	}

	return cpr.IssueTrackerProvider.GetIssue(issueID)
}

func (cpr *CreatePullRequest) createEmptyCommit() error {
	return cpr.Git.CommitEmpty("chore: initial commit")
}

func (cpr *CreatePullRequest) createEmptyCommitAndPush(branchName string) (err error) {
	// 18. CREATE EMPTY COMMIT
	if err = cpr.createEmptyCommit(); err != nil {
		return fmt.Errorf("could not do the empty commit because %s", err)
	}

	return cpr.pushChanges(branchName)
}

func (cpr *CreatePullRequest) pushChanges(branchName string) (err error) {
	// 19. PUSH CHANGES
	err = cpr.Git.PushBranch(branchName)
	if err != nil {
		return fmt.Errorf("could not create the remote branch because %s", err)
	}

	return
}

func (cpr *CreatePullRequest) getPullRequestTitleAndBody(issue domain.Issue) (title string, body string, err error) {
	switch issue.TrackerType() {
	case domain.IssueTrackerTypeGithub:
		title = issue.Title()

		keyword := "Closes"
		if !cpr.Cfg.CloseIssue {
			keyword = "Related to"
		}
		body = fmt.Sprintf("%s #%s", keyword, issue.ID())

	case domain.IssueTrackerTypeJira:
		title = fmt.Sprintf("[%s] %s", issue.ID(), issue.Title())

		body = fmt.Sprintf("Relates to [%s](%s)", issue.ID(), issue.URL())
	default:
		err = fmt.Errorf("issue tracker %s is not supported", issue.TrackerType())
	}

	return title, body, err
}

func (cpr *CreatePullRequest) hasPendingCommits(currentBranch string) (bool, error) {
	commitsToPush, err := cpr.Git.GetCommitsToPush(currentBranch)
	if err != nil {
		return false, err
	}

	return len(commitsToPush) > 0, nil
}

func (cpr *CreatePullRequest) pendingCommits(currentBranch string) (canceled bool, err error) {
	// 11. CHECK IF BRANCH HAS PENDING COMMITS
	hasPendingCommits, err := cpr.hasPendingCommits(currentBranch)
	if err != nil {
		return false, err
	}

	if hasPendingCommits {
		logging.PrintWarn("the branch contains commits that have not been pushed yet")

		// 21. ASK USER CONFIRMATION TO PUSH THE COMMITS
		confirmed, err := cpr.UserInteractionProvider.AskUserForConfirmation("Do you want to continue pushing all pending commits in this branch and create the pull request", true)
		if err != nil {
			return false, err
		}

		if !confirmed {
			// 22. EXIT
			return true, nil
		}

		// 19. PUSH CHANGES
		if err := cpr.pushChanges(currentBranch); err != nil {
			return false, err
		}

	} else {
		// 12. DOES THE REMOTE BRANCH EXISTS
		if !cpr.Git.RemoteBranchExists(currentBranch) {
			// 18. & //19.
			if err := cpr.createEmptyCommitAndPush(currentBranch); err != nil {
				return false, err
			}
		}
	}

	return false, nil
}

func (cpr *CreatePullRequest) createNewLocalBranch(currentBranch string, baseBranch string) error {
	// Check if the base branch will be fetched before the new branch is created
	if cpr.Cfg.FetchFromOrigin {
		if err := cpr.Git.FetchBranchFromOrigin(baseBranch); err != nil {
			return fmt.Errorf("could not fetch the changes from base branch because %s", err)
		}
	}

	// Create the new branch from the base branch
	if err := cpr.Git.CheckoutNewBranchFromOrigin(currentBranch, baseBranch); err != nil {
		return fmt.Errorf("could not create the local branch because %s", err)
	}

	return nil
}

func (cpr *CreatePullRequest) createNewUserBranchAndPush(baseBranch string, issue domain.Issue, repo domain.Repository) (branchName string, canceled bool, err error) {
	branchName, err = cpr.BranchProvider.GetBranchName(issue, repo)
	if err != nil {
		return "", false, err
	}

	if cpr.Git.RemoteBranchExists(branchName) {
		return "", true, ErrRemoteBranchAlreadyExists(branchName)
	}

	fmt.Printf("\nA new pull request is going to be created from %s to %s branch\n", logging.PaintInfo(branchName), logging.PaintInfo(baseBranch))
	if cpr.Cfg.IsInteractive {
		confirmed, err := cpr.UserInteractionProvider.AskUserForConfirmation("Do you want to continue?", true)
		if err != nil {
			return "", false, err
		}
		if !confirmed {
			return "", true, nil
		}
	}

	if err = cpr.createNewLocalBranch(branchName, baseBranch); err != nil {
		return
	}

	// 18. && 19.
	if err = cpr.createEmptyCommitAndPush(branchName); err != nil {
		return
	}

	return branchName, false, nil
}
