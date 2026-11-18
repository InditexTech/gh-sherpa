package branches

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/logging"
)

// warn prints warnings about the generated branch name. It writes to stderr so it never mixes with
// machine-readable output such as `--output json`. Tests replace it to capture the warnings.
var warn = logging.PrintWarnStderr

// croppedEndingPattern matches the endings that a cropped branch name must not keep: trailing dashes, and
// "/", "." or ".lock", which git does not allow at the end of a branch name.
var croppedEndingPattern = regexp.MustCompile(`(?:-|/|\.|\.lock)+$`)

var patternBranchName = regexp.MustCompile(`^(?:(?P<branch_type>\w*)/)?(?P<issue_id>(?:(?P<issue_key>\w*)-)?(?P<issue_number>\d+))(?:-?(?P<issue_context>[\w\-]*))$`)

type BranchProvider struct {
	cfg             Configuration
	UserInteraction domain.UserInteractionProvider
}

type Configuration struct {
	config.Branches
	IsInteractive     bool
	PreferHotfix      bool
	ForcedBranchType  string
	ForcedDescription string
}

func NewFromConfiguration(globalConfig config.Configuration, userInteractionProvider domain.UserInteractionProvider, isInteractive bool) (*BranchProvider, error) {
	return New(Configuration{
		Branches:      globalConfig.Branches,
		IsInteractive: isInteractive,
	}, userInteractionProvider)
}

func New(cfg Configuration, userInteractionProvider domain.UserInteractionProvider) (*BranchProvider, error) {

	return &BranchProvider{
		cfg:             cfg,
		UserInteraction: userInteractionProvider,
	}, nil
}

type BranchNameInfo struct {
	BranchType   string
	IssueId      string
	IssueContext string
}

func ParseBranchName(branchName string) *BranchNameInfo {
	match := patternBranchName.FindStringSubmatch(branchName)

	if len(match) > 0 {
		branchNameInfo := &BranchNameInfo{}

		branchNameInfo.BranchType = match[1]
		branchNameInfo.IssueId = match[2]
		branchNameInfo.IssueContext = match[5]

		return branchNameInfo
	}

	return nil
}

type branchNameRule struct {
	pattern          regexp.Regexp
	replace          string
	repeatWhileMatch bool
}

var branchNameRules = []branchNameRule{
	{pattern: *regexp.MustCompile(`^/`), replace: ""},                                      // Conventional Git branch naming. See https://git-scm.com/docs/git-check-ref-format
	{pattern: *regexp.MustCompile(`~|\^|:|\?|\*|\[|@\{|\\\\`), replace: ""},                // Conventional Git branch naming.
	{pattern: *regexp.MustCompile(`\/\/| |[\/\.]\.|[[:cntrl:]]`), replace: "-"},            // Conventional Git branch naming.
	{pattern: *regexp.MustCompile(`\.lock$|[\/\.]$`), replace: "", repeatWhileMatch: true}, // Conventional Git branch naming.
	{pattern: *regexp.MustCompile(`[àáâãäåÀÁÂÃÄÅ]`), replace: "a"},                         // Replace characters for Kubernetes compatibility
	{pattern: *regexp.MustCompile(`[èéêëÈÉÊË]`), replace: "e"},                             // Replace characters for Kubernetes compatibility
	{pattern: *regexp.MustCompile(`[ìíîïÌÍÎÏ]`), replace: "i"},                             // Replace characters for Kubernetes compatibility
	{pattern: *regexp.MustCompile(`[òóôõöÒÓÔÕÖ]`), replace: "o"},                           // Replace characters for Kubernetes compatibility
	{pattern: *regexp.MustCompile(`[ùúûüÙÚÛÜ]`), replace: "u"},                             // Replace characters for Kubernetes compatibility
	{pattern: *regexp.MustCompile(`[ñÑ]`), replace: "n"},                                   // Replace characters for Kubernetes compatibility
	{pattern: *regexp.MustCompile(`[çÇ]`), replace: "z"},                                   // Replace characters for Kubernetes compatibility
	{pattern: *regexp.MustCompile(`[^\w\-]`), replace: ""},                                 // Remove any other character for Kubernetes compatibility
}

func normalizeBranch(branchSlug string) string {
	branchSlug = strings.TrimSpace(branchSlug)

	for _, r := range branchNameRules {
		branchSlug = r.pattern.ReplaceAllString(branchSlug, r.replace)

		for r.repeatWhileMatch && r.pattern.MatchString(branchSlug) {
			branchSlug = r.pattern.ReplaceAllString(branchSlug, r.replace)
		}
	}

	branchSlug = strings.ToLower(branchSlug)

	return branchSlug
}

// formatBranchName formats a branch name based on the branch type and the issue identifier.
// It overrides the branch prefix if the issue type is present in the branchPrefixOverride map.
// If the prefix is empty, it uses the branch type as the prefix.
func (b BranchProvider) formatBranchName(repoNameWithOwner string, branchType string, issueId string, issueContext string) (branchName string) {
	branchPrefix := branchType

	for issueType, prefix := range b.cfg.Prefixes {
		if prefix != "" && issueType.String() == branchType {
			branchPrefix = prefix
			break
		}
	}

	branchName = fmt.Sprintf("%s/%s", branchPrefix, issueId)
	prefixAndIssueIdLength := len([]rune(branchName))

	if issueContext != "" {
		branchName = fmt.Sprintf("%s-%s", branchName, issueContext)
	}

	maxBranchNameLength := b.cfg.MaxLength - len([]rune(repoNameWithOwner))
	if maxBranchNameLength > 0 {
		if maxBranchNameLength < prefixAndIssueIdLength {
			if b.cfg.KeepIssueId {
				// Exceed the max length just enough to keep "<prefix>/<issue id>" whole.
				maxBranchNameLength = prefixAndIssueIdLength
			} else {
				warn(fmt.Sprintf("max_length (%d) leaves %d characters for the branch name in %s, so the issue id %s is cropped. Increase max_length or set keep_issue_id to true to keep it whole.",
					b.cfg.MaxLength, maxBranchNameLength, repoNameWithOwner, issueId))
			}
		}

		// Crop by characters, not bytes, so a non-ASCII prefix is never cut in the middle of a character.
		branchNameRunes := []rune(branchName)
		branchName = string(branchNameRunes[:min(maxBranchNameLength, len(branchNameRunes))])
		branchName = croppedEndingPattern.ReplaceAllString(branchName, "")
	}

	// Remove all trailing dashes to ensure branches never end with a dash
	trailingDashPattern := regexp.MustCompile(`-+$`)
	return trailingDashPattern.ReplaceAllString(branchName, "")
}
