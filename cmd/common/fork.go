package common

import (
	"fmt"

	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/fork"
	"github.com/InditexTech/gh-sherpa/internal/gh"
	"github.com/InditexTech/gh-sherpa/internal/git"
)

// SetupForkForCommand handles the common fork setup logic for both create-branch and create-pr commands
func SetupForkForCommand(
	cfg config.Configuration,
	forkNameValue string,
	issueValue string,
	ghCli *gh.Cli,
	userInteraction domain.UserInteractionProvider,
	isInteractive bool,
	messageType string, // "issue" or "pull request"
) error {
	forkName := forkNameValue

	forkCfg := domain.ForkConfiguration{
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
		if messageType == "issue" && issueValue != "" {
			fmt.Printf("✓ Ready to start working on issue #%s!\n", issueValue)
		} else {
			fmt.Printf("✓ Ready to start working on the %s!\n", messageType)
		}
	}

	return nil
}
