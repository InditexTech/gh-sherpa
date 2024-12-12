package create_pull_request

import (
	"fmt"

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
	IssueID          string
	BaseBranch       string
	NoFetch          bool
	NoDraft          bool
	NoCloseIssue     bool
	UseDefaultValues bool
	FromTemplate     bool
}

var flags createPullRequestFlags

func init() {
	Command.PersistentFlags().StringVarP(&flags.IssueID, "issue", "i", "", "issue identifier")
	Command.PersistentFlags().StringVarP(&flags.BaseBranch, "base", "b", "", "base branch for checkout. Use the default branch of the repository if it is not set")
	Command.PersistentFlags().BoolVar(&flags.NoFetch, "no-fetch", false, "does not fetch the base branch")
	Command.PersistentFlags().BoolVar(&flags.NoDraft, "no-draft", false, "create the pull request in ready for review mode")
	Command.PersistentFlags().BoolVarP(&flags.NoCloseIssue, "no-close-issue", "n", false, "do not close the GitHub issue after merging the pull request")
	Command.PersistentFlags().BoolVarP(&flags.FromTemplate, "from-template", "t", false, "create the pull request using the repository's pull request template")
}

func runCommand(cmd *cobra.Command, _ []string) error {
	isIssueIDFlagUsed := cmd.Flags().Lookup("issue").Changed

	if isIssueIDFlagUsed && flags.IssueID == "" {
		return fmt.Errorf("sherpa needs an valid issue identifier")
	}

	logging.PrintCommandHeader(cmdName)

	cfg := config.GetConfig()

	issueTrackers, err := issue_trackers.NewFromConfiguration(cfg)
	if err != nil {
		return err
	}

	userInteraction := &interactive.UserInteractionProvider{}

	isInteractive := !flags.UseDefaultValues

	branchProviderCfg := branches.Configuration{
		Branches:      cfg.Branches,
		IsInteractive: isInteractive,
	}
	branchProvider, err := branches.New(branchProviderCfg, userInteraction)
	if err != nil {
		return err
	}

	ghCliProvider := &gh.Cli{}

	createPullRequestConfig := use_cases.CreatePullRequestConfiguration{
		IssueID:         flags.IssueID,
		BaseBranch:      flags.BaseBranch,
		FetchFromOrigin: !flags.NoFetch,
		IsInteractive:   isInteractive,
		DraftPR:         !flags.NoDraft,
		CloseIssue:      !flags.NoCloseIssue,
		FromTemplate:    flags.FromTemplate,
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

	return createPullRequestUseCase.Execute()
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
