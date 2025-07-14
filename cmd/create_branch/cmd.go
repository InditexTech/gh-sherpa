package create_branch

import (
	"fmt"
	"os"

	"github.com/InditexTech/gh-sherpa/internal/branches"
	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/fork"
	"github.com/InditexTech/gh-sherpa/internal/gh"
	"github.com/InditexTech/gh-sherpa/internal/git"
	"github.com/InditexTech/gh-sherpa/internal/interactive"
	"github.com/InditexTech/gh-sherpa/internal/issue_trackers"
	"github.com/InditexTech/gh-sherpa/internal/logging"
	"github.com/InditexTech/gh-sherpa/internal/use_cases"
	"github.com/spf13/cobra"
)

const cmdName = "create-branch"

var Command = &cobra.Command{
	Use:     cmdName,
	Short:   "Create a local branch from an issue type",
	Long:    "Create a local branch according to GitHub or Jira issue type, checking out and fetching the base branch",
	RunE:    runCommand,
	PreRunE: preRunCommand,
	Example: "`gh sherpa " + cmdName + " --issue 1` for GH or `gh sherpa " + cmdName + " --issue PROJECTKEY-1` for Jira",
	Aliases: []string{"cb"},
}

type createBranchFlags struct {
	IssueValue       string
	BaseValue        string
	NoFetchValue     bool
	UseDefaultValues bool
	ForkValue        bool
	ForkNameValue    string
}

var flags = createBranchFlags{}

func init() {
	Command.PersistentFlags().StringVarP(&flags.IssueValue, "issue", "i", "", "issue identifier")
	err := Command.MarkPersistentFlagRequired("issue")

	if err != nil {
		logging.Errorf("error while setting up the command: %s", err)
		os.Exit(1)
	}

	Command.PersistentFlags().StringVarP(&flags.BaseValue, "base", "b", "", "base branch for checkout. Use the default branch of the repository if it is not set")
	Command.PersistentFlags().BoolVar(&flags.NoFetchValue, "no-fetch", false, "does not fetch the base branch")
	Command.PersistentFlags().BoolVar(&flags.ForkValue, "fork", false, "automatically set up fork for external contributors")
	Command.PersistentFlags().StringVar(&flags.ForkNameValue, "fork-name", "", "specify custom fork organization/user (e.g. MyOrg/gh-sherpa)")
}

func runCommand(cmd *cobra.Command, _ []string) (err error) {
	logging.PrintCommandHeader(cmdName)

	cfg := config.GetConfig()

	issueTrackers, err := issue_trackers.NewFromConfiguration(cfg)
	if err != nil {
		return err
	}

	userInteraction := &interactive.UserInteractionProvider{}

	isInteractive := !flags.UseDefaultValues

	ghCli := &gh.Cli{}

	if flags.ForkValue {
		if err := setupFork(cfg, ghCli, userInteraction, isInteractive); err != nil {
			return err
		}
	}

	branchProviderCfg := branches.Configuration{
		Branches:      cfg.Branches,
		IsInteractive: isInteractive,
	}
	branchProvider, err := branches.New(branchProviderCfg, userInteraction)
	if err != nil {
		return err
	}

	createBranchConfig := use_cases.CreateBranchConfiguration{
		IssueID:         flags.IssueValue,
		BaseBranch:      flags.BaseValue,
		FetchFromOrigin: !flags.NoFetchValue,
		IsInteractive:   isInteractive,
	}
	createBranch := use_cases.CreateBranch{
		Cfg:                     createBranchConfig,
		Git:                     &git.Provider{},
		RepositoryProvider:      ghCli,
		IssueTrackerProvider:    issueTrackers,
		UserInteractionProvider: userInteraction,
		BranchProvider:          branchProvider,
	}

	err = createBranch.Execute()

	if err != nil {
		return err
	}

	return
}

func setupFork(cfg config.Configuration, ghCli *gh.Cli, userInteraction domain.UserInteractionProvider, isInteractive bool) error {
	forkName := flags.ForkNameValue
	if forkName == "" && cfg.Github.ForkOrganization != "" {
		// We'll construct the fork name once we know the repo name
		// This will be handled by the fork manager
	}

	forkCfg := fork.Configuration{
		DefaultOrganization: cfg.Github.ForkOrganization,
		IsInteractive:       isInteractive,
	}

	forkManager := fork.NewManager(
		forkCfg,
		ghCli,
		&git.Provider{},
		userInteraction,
		ghCli,
	)

	result, err := forkManager.SetupFork(forkName)
	if err != nil {
		return err
	}

	if result.WasAlreadyConfigured {
		return nil
	}

	if result.ForkCreated {
		fmt.Printf("✓ Ready to start working on issue #%s!\n", flags.IssueValue)
	}

	return nil
}

func preRunCommand(cmd *cobra.Command, _ []string) error {
	if cmd.Flags().Lookup("no-fetch").Changed {
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
