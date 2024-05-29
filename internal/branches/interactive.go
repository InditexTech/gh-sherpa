package branches

import (
	"errors"
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	"github.com/InditexTech/gh-sherpa/internal/interactive"
	"github.com/InditexTech/gh-sherpa/internal/logging"
)

// GetBranchName asks the user for a branch name in an interactive way
func (b BranchProvider) GetBranchName(issue domain.Issue, repo domain.Repository) (branchName string, err error) {
	issueType := issue.Type()
	branchType := issueType.String()

	formattedID := issue.FormatID()

	issueSlug := normalizeBranch(issue.Title())

	issueTrackerType := issue.TrackerType()

	if b.cfg.IsInteractive {
		branchType, err = b.getBranchType(issueType, issueTrackerType)
		if err != nil {
			return "", err
		}

		maxContextLen := calcIssueContextMaxLen(repo.NameWithOwner, branchType, formattedID)
		promptIssueContext := fmt.Sprintf("additional description (optional). Truncate to %d chars", maxContextLen)
		err = b.UserInteraction.SelectOrInput(promptIssueContext, []string{}, &issueSlug, false)
		if err != nil {
			return "", err
		}

		issueSlug = normalizeBranch(issueSlug)

	} else {
		if issueType == issue_types.Other || issueType == issue_types.Unknown {
			return "", errors.New("undetermined issue type")
		}

		// remap bug to bugfix
		if issueType == issue_types.Bug {
			branchType = b.getBugFixBranchType()
		}
	}

	branchName = b.formatBranchName(repo.NameWithOwner, branchType, formattedID, issueSlug)

	return branchName, nil
}

func (b BranchProvider) getBugFixBranchType() (branchType string) {
	if b.cfg.Prefixes[issue_types.Bugfix] != "" {
		branchType = b.cfg.Prefixes[issue_types.Bugfix]
	} else {
		branchType = issue_types.Bugfix.String()
	}

	return branchType
}

func calcIssueContextMaxLen(repository string, branchType string, issueId string) (lenIssueContext int) {
	preBranchName := fmt.Sprintf("%s/%s-", branchType, issueId)

	if lenIssueContext = 63 - (len([]rune(repository)) + len([]rune(preBranchName))); lenIssueContext < 0 {
		lenIssueContext = 0
	}

	return
}

func (b BranchProvider) getBranchType(issueType issue_types.IssueType, issueTrackerType domain.IssueTrackerType) (branchType string, err error) {
	branchType = issueType.String()

	if issueType == issue_types.Bug || issueType == issue_types.Bugfix {
		err = askBranchTypeBug(&branchType, issueTrackerType, b.UserInteraction)
	} else if issueType != issue_types.Other && issueType != issue_types.Unknown {
		err = askBranchType(&branchType, issueTrackerType, b.UserInteraction)
	} else {
		logging.PrintWarn("undetermined issue type")
	}

	if err != nil {
		return
	}

	if issueType == issue_types.Other || issueType == issue_types.Unknown {
		err = askBranchTypeOther(&branchType, b.UserInteraction)
		if err != nil {
			return
		}
	}

	return
}

func askBranchTypeBug(branchType *string, issueTrackerType domain.IssueTrackerType, interactionProvider domain.UserInteractionProvider) error {
	bugValues := issue_types.GetBugValues()
	bugValuesStr := make([]string, len(bugValues))
	for i, branchType := range bugValues {
		bugValuesStr[i] = branchType.String()
	}
	*branchType = bugValuesStr[0]

	promptMessage := interactive.GetPromptMessageBranchType(*branchType, issueTrackerType)
	if err := interactionProvider.SelectOrInputPrompt(promptMessage, bugValuesStr, branchType, true); err != nil {
		return err
	}

	return nil
}

func askBranchType(branchType *string, issueTrackerType domain.IssueTrackerType, interactionProvider domain.UserInteractionProvider) (err error) {
	branchPrefixes := []string{*branchType, issue_types.Other.String()}

	promptMessage := interactive.GetPromptMessageBranchType(*branchType, issueTrackerType)
	if err := interactionProvider.SelectOrInputPrompt(promptMessage, branchPrefixes, branchType, true); err != nil {
		return err
	}

	return nil
}

func askBranchTypeOther(branchType *string, interactionProvider domain.UserInteractionProvider) error {
	validIssueTypes := issue_types.GetValidIssueTypes()
	branchTypes := make([]string, len(validIssueTypes))
	for i, branchType := range validIssueTypes {
		branchTypes[i] = branchType.String()
	}
	*branchType = branchTypes[0]

	if err := interactionProvider.SelectOrInput("branch type", branchTypes, branchType, true); err != nil {
		return err
	}

	return nil
}
