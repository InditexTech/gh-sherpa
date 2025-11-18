package branches

import (
	"errors"
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	"github.com/InditexTech/gh-sherpa/internal/interactive"
	"github.com/InditexTech/gh-sherpa/internal/logging"
)

// ErrUndeterminedIssueType is returned when the issue type can't be determined
var ErrUndeterminedIssueType = errors.New("undetermined issue type")

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

		truncatePrompt := ""
		maxContextLen := b.calcIssueContextMaxLen(repo.NameWithOwner, branchType, formattedID)
		if maxContextLen > 0 {
			truncatePrompt = fmt.Sprintf(" Truncate to %d chars", maxContextLen)
		}

		promptIssueContext := "additional description (optional)." + truncatePrompt
		err = b.UserInteraction.SelectOrInput(promptIssueContext, []string{}, &issueSlug, false)
		if err != nil {
			return "", err
		}

		issueSlug = normalizeBranch(issueSlug)

	} else {
		// Check for kind/bug label to determine if prefer-hotfix should apply
		// This works regardless of the issue's determined type or label position
		hasBugLabel := issue.HasLabel("kind/bug")
		
		if hasBugLabel || issueType == issue_types.Bug || issueType == issue_types.Bugfix {
			if b.cfg.PreferHotfix && hasBugLabel {
				// If prefer-hotfix is enabled and kind/bug label is present, use hotfix
				branchType = b.getHotfixBranchType()
				issueType = issue_types.Hotfix
			} else {
				// Otherwise use bugfix
				branchType = b.getBugFixBranchType()
				issueType = issue_types.Bugfix
			}
		}

		if !issueType.Valid() || issueType == issue_types.Other || issueType == issue_types.Unknown {
			return "", ErrUndeterminedIssueType
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

func (b BranchProvider) getHotfixBranchType() (branchType string) {
	if b.cfg.Prefixes[issue_types.Hotfix] != "" {
		branchType = b.cfg.Prefixes[issue_types.Hotfix]
	} else {
		branchType = issue_types.Hotfix.String()
	}

	return branchType
}

func (b BranchProvider) calcIssueContextMaxLen(repository string, branchType string, issueId string) (lenIssueContext int) {
	preBranchName := fmt.Sprintf("%s/%s-", branchType, issueId)

	if lenIssueContext = b.cfg.MaxLength - (len([]rune(repository)) + len([]rune(preBranchName))); lenIssueContext < 0 {
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

	if branchType == issue_types.Other.String() || branchType == issue_types.Unknown.String() {
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
