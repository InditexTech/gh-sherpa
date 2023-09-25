package branches

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
)

var patternBranchName = regexp.MustCompile(`^(?:(?P<branch_type>\w*)/)?(?P<issue_id>(?:(?P<issue_key>\w*)-)?(?P<issue_number>\d+))(?:-?(?P<issue_context>[\w\-]*))$`)

type BranchProvider struct {
	cfg             Configuration
	UserInteraction domain.UserInteractionProvider
}

type Configuration struct {
	config.BranchPrefixOverrides
}

func (c Configuration) Validate() (err error) {
	//TODO: Validate configuration
	return nil
}

func NewFromConfiguration(globalConfig config.Configuration, userInteractionProvider domain.UserInteractionProvider) (*BranchProvider, error) {
	return New(Configuration{
		BranchPrefixOverrides: globalConfig.BranchPrefixOverrides,
	}, userInteractionProvider)
}

func New(cfg Configuration, userInteractionProvider domain.UserInteractionProvider) (*BranchProvider, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

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

type issueContextRule struct {
	pattern          regexp.Regexp
	replace          string
	repeatWhileMatch bool
}

var issueContextRules = []issueContextRule{
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

func parseIssueContext(issueContext string) string {
	issueContext = strings.TrimSpace(issueContext)

	for _, r := range issueContextRules {
		issueContext = r.pattern.ReplaceAllString(issueContext, r.replace)

		for r.repeatWhileMatch && r.pattern.MatchString(issueContext) {
			issueContext = r.pattern.ReplaceAllString(issueContext, r.replace)
		}
	}

	issueContext = strings.ToLower(issueContext)

	return issueContext
}

// formatBranchName formats a branch name based on the issue type and the issue identifier.
// It overrides the branch prefix if the issue type is present in the branchPrefixOverride map.
// If the prefix is empty, it uses the branch type as the prefix.
func (b BranchProvider) formatBranchName(repoNameWithOwner string, branchType string, issueId string, issueContext string) (branchName string) {
	branchPrefix := branchType

	for issueType, prefix := range b.cfg.BranchPrefixOverrides {
		if prefix != "" && issueType.String() == branchType {
			branchPrefix = prefix
			break
		}
	}

	branchName = fmt.Sprintf("%s/%s", branchPrefix, issueId)

	if issueContext != "" {
		branchName = fmt.Sprintf("%s-%s", branchName, issueContext)
	}

	return strings.TrimSuffix(branchName[:min(63-len([]rune(repoNameWithOwner)), len([]rune(branchName)))], "-")
}

// min returns the smallest of x or y.
func min(x, y int) int {
	if x > y {
		return y
	}

	return x
}
