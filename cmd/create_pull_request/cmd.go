package create_pull_request

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/InditexTech/gh-sherpa/cmd/common"
	"github.com/InditexTech/gh-sherpa/internal/branches"
	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/gh"
	"github.com/InditexTech/gh-sherpa/internal/git"
	"github.com/InditexTech/gh-sherpa/internal/interactive"
	"github.com/InditexTech/gh-sherpa/internal/issue_trackers"
	"github.com/InditexTech/gh-sherpa/internal/logging"
	"github.com/InditexTech/gh-sherpa/internal/use_cases"
	"github.com/spf13/cobra"
)

const cmdName = "create-pr"

var Command = &cobra.Command{
	Use:     cmdName,
	Short:   "Create a pull request from the current local branch or issue type", // by creating and checkout a branch and pushing an empty commit",
	Long:    "Create a pull request in draft mode from the current local branch or create one based on the type of GitHub or Jira issue, pushing all pending local commits or creating an empty one.",
	RunE:    runCommand,
	Example: "`gh sherpa " + cmdName + " --issue 1` for GH or `gh sherpa " + cmdName + " --issue PROJECTKEY-1` for Jira",
	Aliases: []string{"cpr"},
	PreRunE: preRunCommand,
}

type createPullRequestFlags struct {
	IssueID             string
	BaseBranch          string
	NoFetch             bool
	NoDraft             bool
	NoCloseIssue        bool
	UseDefaultValues    bool
	TemplatePath        string
	ForkValue           bool
	ForkNameValue       string
	PreferHotfix        bool
	BranchType          string
	BranchDescription   string
	BranchName          string
	DryRun              bool
	OutputFormat        string
	PRTitle             string
	PRBody              string
	PRBodyFile          string
	NoUseExistingBranch bool
	ExtraLabels         []string
	Reviewers           []string
	Assignees           []string
}

var flags createPullRequestFlags

func init() {
	Command.PersistentFlags().StringVarP(&flags.IssueID, "issue", "i", "", "issue identifier")
	Command.PersistentFlags().StringVarP(&flags.BaseBranch, "base", "b", "", "base branch for checkout. Use the default branch of the repository if it is not set")
	Command.PersistentFlags().BoolVar(&flags.NoFetch, "no-fetch", false, "does not fetch the base branch")
	Command.PersistentFlags().BoolVar(&flags.NoDraft, "no-draft", false, "create the pull request in ready for review mode")
	Command.PersistentFlags().BoolVarP(&flags.NoCloseIssue, "no-close-issue", "n", false, "do not close the GitHub issue after merging the pull request")
	Command.PersistentFlags().StringVar(&flags.TemplatePath, "template", "", "path to a pull request template file")
	Command.PersistentFlags().BoolVar(&flags.ForkValue, "fork", false, "automatically set up fork for external contributors")
	Command.PersistentFlags().StringVar(&flags.ForkNameValue, "fork-name", "", "specify custom fork organization/user (e.g. MyOrg/gh-sherpa)")
	Command.PersistentFlags().BoolVar(&flags.PreferHotfix, "prefer-hotfix", false, "prefer hotfix branch prefix for bug issues when using non-interactive mode")
	Command.PersistentFlags().StringVar(&flags.BranchType, "branch-type", "", "force a specific branch type prefix (e.g. feature, bugfix, hotfix)")
	Command.PersistentFlags().StringVar(&flags.BranchDescription, "branch-description", "", "force a specific branch description slug instead of deriving it from the issue title")
	Command.PersistentFlags().StringVar(&flags.BranchName, "branch-name", "", "use exactly this branch name instead of auto-generating one")
	Command.PersistentFlags().BoolVar(&flags.DryRun, "dry-run", false, "print what would happen without actually creating the PR")
	Command.PersistentFlags().StringVar(&flags.OutputFormat, "output", "", "output format: '' (default human-readable) or 'json'")
	Command.PersistentFlags().StringVar(&flags.PRTitle, "pr-title", "", "override the auto-generated PR title")
	Command.PersistentFlags().StringVar(&flags.PRBody, "pr-body", "", "override the auto-generated PR body")
	Command.PersistentFlags().StringVar(&flags.PRBodyFile, "pr-body-file", "", "read PR body from file (overrides --pr-body and template)")
	Command.PersistentFlags().BoolVar(&flags.NoUseExistingBranch, "no-use-existing-branch", false, "fail if a branch for this issue already exists (default is to reuse it)")
	Command.PersistentFlags().StringArrayVar(&flags.ExtraLabels, "label", []string{}, "additional label to apply to the PR (can be repeated)")
	Command.PersistentFlags().StringArrayVar(&flags.Reviewers, "reviewer", []string{}, "request a review from this user or team (can be repeated)")
	Command.PersistentFlags().StringArrayVar(&flags.Assignees, "assignee", []string{}, "assign this user to the PR (can be repeated)")
}

func runCommand(cmd *cobra.Command, _ []string) error {
	isIssueIDFlagUsed := cmd.Flags().Lookup("issue").Changed

	if isIssueIDFlagUsed && flags.IssueID == "" {
		return fmt.Errorf("sherpa needs an valid issue identifier")
	}

	if flags.OutputFormat != "json" {
		logging.PrintCommandHeader(cmdName)
	}

	cfg := config.GetConfig()

	issueTrackers, err := issue_trackers.NewFromConfiguration(cfg)
	if err != nil {
		return err
	}

	userInteraction := &interactive.UserInteractionProvider{}

	isInteractive := !flags.UseDefaultValues
	// --output json implies non-interactive to prevent stdin prompts from corrupting JSON output.
	if flags.OutputFormat == "json" {
		isInteractive = false
	}

	branchProviderCfg := branches.Configuration{
		Branches:          cfg.Branches,
		IsInteractive:     isInteractive,
		PreferHotfix:      flags.PreferHotfix,
		ForcedBranchType:  flags.BranchType,
		ForcedDescription: flags.BranchDescription,
	}
	branchProvider, err := branches.New(branchProviderCfg, userInteraction)
	if err != nil {
		return err
	}

	ghCliProvider := &gh.Cli{}

	if flags.ForkValue {
		if err := common.SetupForkForCommand(cfg, flags.ForkNameValue, flags.IssueID, ghCliProvider, userInteraction, isInteractive, "pull request"); err != nil {
			return err
		}
	}

	createPullRequestConfig := use_cases.CreatePullRequestConfiguration{
		IssueID:             flags.IssueID,
		BaseBranch:          flags.BaseBranch,
		FetchFromOrigin:     !flags.NoFetch,
		IsInteractive:       isInteractive,
		DraftPR:             !flags.NoDraft,
		CloseIssue:          !flags.NoCloseIssue,
		TemplatePath:        flags.TemplatePath,
		BranchName:          flags.BranchName,
		DryRun:              flags.DryRun,
		OutputFormat:        flags.OutputFormat,
		PRTitle:             flags.PRTitle,
		PRBody:              flags.PRBody,
		PRBodyFile:          flags.PRBodyFile,
		NoUseExistingBranch: flags.NoUseExistingBranch,
		ExtraLabels:         flags.ExtraLabels,
		Reviewers:           flags.Reviewers,
		Assignees:           flags.Assignees,
	}
	createPullRequestUseCase := use_cases.CreatePullRequest{
		Cfg:                     createPullRequestConfig,
		Git:                     &git.Provider{},
		RepositoryProvider:      ghCliProvider,
		IssueTrackerProvider:    issueTrackers,
		UserInteractionProvider: userInteraction,
		PullRequestProvider:     ghCliProvider,
		BranchProvider:          branchProvider,
	}

	_, err = createPullRequestUseCase.Execute()
	if err != nil && flags.OutputFormat == "json" {
		errJSON, _ := json.Marshal(map[string]string{"error": err.Error()})
		fmt.Fprintln(os.Stderr, string(errJSON))
		os.Exit(1)
	}
	return err
}

func preRunCommand(cmd *cobra.Command, _ []string) error {
	if cmd.Flags().Lookup("no-fetch").Changed {
		logging.Debug("Flag no-fetch used found, marking issue flag as required...")
		if err := cmd.MarkFlagRequired("issue"); err != nil {
			return err
		}
	}

	yesFlag := cmd.Flags().Lookup("yes")
	if yesFlag != nil {
		flags.UseDefaultValues = yesFlag.Changed
	}

	return nil
}
