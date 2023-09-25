package use_cases

import (
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/branches"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	"github.com/InditexTech/gh-sherpa/internal/logging"
)

// CreatePullRequestArgs contains the arguments for the CreatePullRequest use case
type CreatePullRequestArgs struct {
	IssueId          string
	BaseBranch       string
	NoFetch          bool
	NoDraft          bool
	UseDefaultValues bool
	NoCloseIssue     bool
}

type CreatePullRequest struct {
	BranchPrefixOverride    map[issue_types.IssueType]string
	Git                     domain.GitProvider
	GhCli                   domain.GhCli
	IssueTrackerProvider    domain.IssueTrackerProvider
	UserInteractionProvider domain.UserInteractionProvider
	PullRequestProvider     domain.PullRequestProvider
	BranchProvider          domain.BranchProvider
}

// Execute executes the create pull request use case
func (cpr CreatePullRequest) Execute(args CreatePullRequestArgs) error {

	repo, err := cpr.GhCli.GetRepo()
	if err != nil {
		logging.Debugf("error while getting repo: %s", err)
		return err
	}

	baseBranch := args.BaseBranch
	if baseBranch == "" {
		baseBranch = repo.DefaultBranchRef
	}

	currentBranch, err := cpr.Git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("could not get the current branch name because %s", err)
	}

	issueID := args.IssueId
	//1. FLAG ISSUE IS USED
	if issueID != "" {
		issueTracker, err := cpr.IssueTrackerProvider.GetIssueTracker(issueID)
		if err != nil {
			return err
		}

		formattedIssueId := issueTracker.FormatIssueId(issueID)

		//7. CHECK IF A LOCAL BRANCH CONTAINS THIS ISSUE
		name, exists := cpr.Git.BranchExistsContains(fmt.Sprintf("/%s-", formattedIssueId))
		if exists {
			//8. FLAG DEFAULT IS USED
			if args.UseDefaultValues {
				return fmt.Errorf("the branch %s already exists", logging.PaintWarning(name))
			}

			logging.PrintWarn(fmt.Sprintf("there is already a local branch named %s for this issue", logging.PaintInfo(name)))

			//9. CONFIRM USER TO USE THE BRANCH
			confirmed, err := cpr.UserInteractionProvider.AskUserForConfirmation("Do you want to use this branch to create the pull request", true)
			if err != nil {
				return err
			}

			if !confirmed {
				branchName, canceled, err := cpr.createNewUserBranchAndPush(baseBranch, issueTracker, issueID, *repo, args.UseDefaultValues, args.NoFetch)
				if err != nil {
					return err
				}
				currentBranch = branchName
				if canceled {
					return nil
				}

			} else {
				currentBranch = name
				if err := cpr.Git.CheckoutBranch(currentBranch); err != nil {
					return fmt.Errorf("could not switch to the branch because %w", err)
				}
			}

			//11. CHECK IF BRANCH HAS PENDING COMMITS
			canceled, err := cpr.pendingCommits(currentBranch)
			if err != nil {
				return err
			}
			if canceled {
				return nil
			}

		} else {
			branchName, canceled, err := cpr.createNewUserBranchAndPush(baseBranch, issueTracker, issueID, *repo, args.UseDefaultValues, args.NoFetch)
			if err != nil {
				return err
			}
			currentBranch = branchName
			if canceled {
				return nil
			}
		}

	} else {
		//1.NO

		branchNameInfo := branches.ParseBranchName(currentBranch)
		//2. DOES CURRENT BRANCH CONTAINS A ISSUE IN IT
		if branchNameInfo == nil || branchNameInfo.IssueId == "" {
			//3. EXIT
			return fmt.Errorf("could not find an issue identifier in the current branch named %s", logging.PaintWarning(currentBranch))
		}

		logging.PrintInfo(fmt.Sprintf("The current branch named %s is available to create a pull request", logging.PaintWarning(currentBranch)))

		issueID = cpr.IssueTrackerProvider.ParseIssueId(branchNameInfo.IssueId)

		//4. FLAG DEFAULT IS USED
		if !args.UseDefaultValues {
			//5. CONFIG USER TO USE THIS BRANCH
			confirmed, err := cpr.UserInteractionProvider.AskUserForConfirmation("Do you want to use this branch to create the pull request", true)
			if err != nil {
				return err
			}

			if !confirmed {
				return nil
			}

		}

		//11. CHECK IF BRANCH HAS PENDING COMMITS
		canceled, err := cpr.pendingCommits(currentBranch)
		if err != nil {
			return err
		}
		if canceled {
			return nil
		}
	}

	//13. CHECK IF PR DOES ALREADY EXISTS
	pr, err := cpr.PullRequestProvider.GetPullRequestForBranch(currentBranch)
	if err != nil {
		return fmt.Errorf("error while getting pull request for branch: %w", err)
	}

	if pr != nil && !pr.Closed {
		//14. EXIT
		logging.PrintInfo(fmt.Sprintf("\nA pull request %s for this branch already exists", logging.PaintWarning(pr.Url)))
		return nil
	}

	//15. GET ISSUE INFO
	title, body, err := cpr.getPullRequestInfo(issueID, args.NoCloseIssue)
	if err != nil {
		return err
	}

	//16. CREATE PULL REQUEST
	prURL, err := cpr.PullRequestProvider.CreatePullRequest(title, body, baseBranch, currentBranch, !args.NoDraft)
	if err != nil {
		return fmt.Errorf("could not create the pull request because %s", err)
	}

	if !args.UseDefaultValues {
		fmt.Println()
	}

	fmt.Printf("The pull request %s have been created!\nYou are now working on the branch %s", logging.PaintInfo(prURL), logging.PaintInfo(currentBranch))
	fmt.Println()

	// 17. EXIT
	return nil
}

func (cpr *CreatePullRequest) createEmptyCommitAndPush(branchName string) (err error) {
	//18. CREATE EMPTY COMMIT
	err = cpr.Git.CommitEmpty("chore: initial commit")
	if err != nil {
		return fmt.Errorf("could not do the empty commit because %s", err)
	}

	return cpr.pushChanges(branchName)
}

func (cpr *CreatePullRequest) pushChanges(branchName string) (err error) {
	//19. PUSH CHANGES
	err = cpr.Git.PushBranch(branchName)
	if err != nil {
		return fmt.Errorf("could not create the remote branch because %s", err)
	}

	return
}

func (cpr *CreatePullRequest) getPullRequestTitleAndBody(issue domain.Issue, noCloseIssue bool) (title string, body string, err error) {
	switch issue.IssueTracker {
	case domain.IssueTrackerTypeGithub:
		title = issue.Title

		keyword := "Closes"
		if noCloseIssue {
			keyword = "Related to"
		}
		body = fmt.Sprintf("%s #%s", keyword, issue.ID)

	case domain.IssueTrackerTypeJira:
		title = fmt.Sprintf("[%s] %s", issue.ID, issue.Title)

		body = fmt.Sprintf("Relates to [%s](%s)", issue.ID, issue.Url)
	default:
		err = fmt.Errorf("issue tracker %s is not supported", issue.IssueTracker)
	}

	return title, body, err
}

func (cpr *CreatePullRequest) pendingCommits(currentBranch string) (canceled bool, err error) {
	//11. CHECK IF BRANCH HAS PENDING COMMITS
	commitsToPush, err := cpr.Git.GetCommitsToPush(currentBranch)
	if err != nil {
		return false, err
	}

	if len(commitsToPush) > 0 {
		logging.PrintWarn("the branch contains commits that have not been pushed yet")

		//21. ASK USER CONFIRMATION TO PUSH THE COMMITS
		confirmed, err := cpr.UserInteractionProvider.AskUserForConfirmation("Do you want to continue pushing all pending commits in this branch and create the pull request", true)
		if err != nil {
			return false, err
		}

		if !confirmed {
			//22. EXIT
			return true, nil
		}

		//19. PUSH CHANGES
		if err := cpr.pushChanges(currentBranch); err != nil {
			return false, err
		}

	} else {
		//12. DOES THE REMOTE BRANCH EXISTS
		if !cpr.Git.RemoteBranchExists(currentBranch) {
			//18. & //19.
			if err := cpr.createEmptyCommitAndPush(currentBranch); err != nil {
				return false, err
			}
		} else {
			logging.Info("The remote branch already exists, skipping creation...")
		}
	}

	return false, nil
}

func (cpr *CreatePullRequest) getIssueFromIssueID(issueID string) (issue domain.Issue, err error) {
	issueTracker, err := cpr.IssueTrackerProvider.GetIssueTracker(issueID)
	if err != nil {
		return
	}

	return issueTracker.GetIssue(issueID)
}

func (cpr *CreatePullRequest) createNewLocalBranch(currentBranch string, baseBranch string, fetch bool) error {
	// Check if the base branch will be fetched before the new branch is created
	if fetch {
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

func (cpr *CreatePullRequest) getPullRequestInfo(issueID string, noCloseIssue bool) (title string, body string, err error) {
	issue, err := cpr.getIssueFromIssueID(issueID)
	if err != nil {
		return
	}

	title, body, err = cpr.getPullRequestTitleAndBody(issue, noCloseIssue)

	return
}

func (cpr *CreatePullRequest) askUserForNewBranchName(baseBranch string, issueTracker domain.IssueTracker, issueID string, repo domain.Repository, useDefaultValues bool) (newBranchName string, cancelled bool, err error) {

	newBranchName, err = cpr.BranchProvider.AskBranchName(issueTracker, issueID, repo, useDefaultValues)
	if err != nil {
		return "", false, err
	}

	if !useDefaultValues {
		fmt.Println()
		fmt.Printf("A new pull request is going to be created from %s to %s branch", logging.PaintInfo(newBranchName), logging.PaintInfo(baseBranch))
		fmt.Println()

		confirmation, err := cpr.UserInteractionProvider.AskUserForConfirmation("Do you want to continue?", true)
		if err != nil {
			return "", false, err
		}

		if !confirmation {
			return "", true, nil
		}
	}

	return newBranchName, false, nil
}

func (cpr *CreatePullRequest) createNewUserBranchAndPush(baseBranch string, issueTracker domain.IssueTracker, issueID string, repo domain.Repository, useDefaultValues, noFetch bool) (branchName string, canceled bool, err error) {
	//useDefaultValues should always be false in this specific case but we use it anyways
	branchName, canceled, err = cpr.askUserForNewBranchName(baseBranch, issueTracker, issueID, repo, useDefaultValues)
	if err != nil || canceled {
		return
	}

	if err = cpr.createNewLocalBranch(branchName, baseBranch, !noFetch); err != nil {
		return
	}

	//18. && 19.
	if err = cpr.createEmptyCommitAndPush(branchName); err != nil {
		return
	}

	return branchName, false, nil
}
