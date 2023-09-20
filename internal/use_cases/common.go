package use_cases

import (
	"errors"
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/branches"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	"github.com/InditexTech/gh-sherpa/internal/interactive"
	"github.com/InditexTech/gh-sherpa/internal/logging"
)

type providers struct {
	UserInteraction domain.UserInteractionProvider
}

func askBranchName(branchName *string, issueTracker domain.IssueTracker, issueIdentifier string, repo domain.Repository, useDefaultValues bool, prov providers) (err error) {
	issue, err := issueTracker.GetIssue(issueIdentifier)
	if err != nil {
		return err
	}

	issueType, err := issueTracker.GetIssueType(issue)
	if err != nil {
		return err
	}
	branchPrefix := issueType.String()

	issueID := issueTracker.FormatIssueId(issue.ID)

	issueSlug := branches.ParseIssueContext(issue.Title)

	issueTrackerType := issueTracker.GetIssueTrackerType()

	if !useDefaultValues {
		branchPrefix, err = getBranchPrefix(issueType, issueTrackerType, prov)
		if err != nil {
			return err
		}

		promptIssueContext := fmt.Sprintf("additional description (optional). Truncate to %d chars", calcIssueContextMaxLen(repo.NameWithOwner, issueType.String(), issueID))
		err = prov.UserInteraction.SelectOrInput(promptIssueContext, []string{}, &issueSlug, false)
		if err != nil {
			return err
		}

		issueSlug = branches.ParseIssueContext(issueSlug)

	} else {
		if issueType == issue_types.Other || issueType == issue_types.Unknown {
			return errors.New("undetermined issue type")
		}

		if issueType == issue_types.Bug {
			branchPrefix = issue_types.Bugfix.String()
		}
	}

	*branchName = branches.FormatBranchName(repo.NameWithOwner, branchPrefix, issueID, issueSlug)

	return nil
}

func askBranchPrefixBug(branchPrefix *string, issueTrackerType domain.IssueTrackerType, interactionProvider domain.UserInteractionProvider) error {
	bugValues := issue_types.GetBugValues()
	bugValuesStr := make([]string, len(bugValues))
	for i, branchType := range bugValues {
		bugValuesStr[i] = branchType.String()
	}
	*branchPrefix = bugValuesStr[0]

	promptMessage := interactive.GetPromptMessageBranchType(*branchPrefix, issueTrackerType)
	if err := interactionProvider.SelectOrInputPrompt(promptMessage, bugValuesStr, branchPrefix, true); err != nil {
		return err
	}

	return nil
}

func askBranchPrefix(branchPrefix *string, issueTrackerType domain.IssueTrackerType, interactionProvider domain.UserInteractionProvider) (err error) {
	branchPrefixes := []string{*branchPrefix, issue_types.Other.String()}

	promptMessage := interactive.GetPromptMessageBranchType(*branchPrefix, issueTrackerType)
	if err := interactionProvider.SelectOrInputPrompt(promptMessage, branchPrefixes, branchPrefix, true); err != nil {
		return err
	}

	return nil
}

func askBranchPrefixOther(branchPrefix *string, interactionProvider domain.UserInteractionProvider) error {
	validIssueTypes := issue_types.GetValidIssueTypes()
	branchTypes := make([]string, len(validIssueTypes))
	for i, branchType := range validIssueTypes {
		branchTypes[i] = branchType.String()
	}
	*branchPrefix = branchTypes[0]

	if err := interactionProvider.SelectOrInput("branch type", branchTypes, branchPrefix, true); err != nil {
		return err
	}

	return nil
}

func calcIssueContextMaxLen(repository string, branchType string, issueId string) (lenIssueContext int) {
	preBranchName := fmt.Sprintf("%s/%s-", branchType, issueId)

	if lenIssueContext = 63 - (len([]rune(repository)) + len([]rune(preBranchName))); lenIssueContext < 0 {
		lenIssueContext = 0
	}

	return
}

func getBranchPrefix(issueType issue_types.IssueType, issueTrackerType domain.IssueTrackerType, prov providers) (branchPrefix string, err error) {
	branchPrefix = issueType.String()

	if issueType == issue_types.Bug || issueType == issue_types.Bugfix {
		err = askBranchPrefixBug(&branchPrefix, issueTrackerType, prov.UserInteraction)
	} else if issueType != issue_types.Other && issueType != issue_types.Unknown {
		err = askBranchPrefix(&branchPrefix, issueTrackerType, prov.UserInteraction)
	} else {
		logging.PrintWarn("undetermined issue type")
	}

	if err != nil {
		return
	}

	if issueType == issue_types.Other || issueType == issue_types.Unknown {
		err = askBranchPrefixOther(&branchPrefix, prov.UserInteraction)
		if err != nil {
			return
		}
	}

	return
}
