package fork

import (
	"fmt"
	"strings"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/logging"
	"github.com/InditexTech/gh-sherpa/internal/utils"
)

type Manager struct {
	cfg                     Configuration
	repositoryProvider      domain.RepositoryProvider
	gitProvider             domain.GitProvider
	userInteractionProvider domain.UserInteractionProvider
	ghCli                   ForkProvider
}

type ForkProvider interface {
	IsRepositoryFork() (bool, error)
	CreateFork(forkName string) error
	ForkExists(forkName string) (bool, error)
	SetDefaultRepository(repo string) error
	GetRemoteConfiguration() (map[string]string, error)
	ConfigureRemotesForExistingFork(forkName string) error
}

func NewManager(
	cfg Configuration,
	repositoryProvider domain.RepositoryProvider,
	gitProvider domain.GitProvider,
	userInteractionProvider domain.UserInteractionProvider,
	ghCli ForkProvider,
) *Manager {
	return &Manager{
		cfg:                     cfg,
		repositoryProvider:      repositoryProvider,
		gitProvider:             gitProvider,
		userInteractionProvider: userInteractionProvider,
		ghCli:                   ghCli,
	}
}

func (m *Manager) DetectForkStatus() (*ForkStatus, error) {
	logging.Debug("Detecting repository fork status...")

	repo, err := m.repositoryProvider.GetRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository information: %w", err)
	}

	remotes, err := m.ghCli.GetRemoteConfiguration()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote configuration: %w", err)
	}

	logging.Debugf("Current repository: %s", repo.NameWithOwner)
	logging.Debugf("Remote configuration: %+v", remotes)

	status := &ForkStatus{}

	// Check if we have both origin and upstream remotes configured
	origin, hasOrigin := remotes["origin"]
	upstream, hasUpstream := remotes["upstream"]

	if hasOrigin && hasUpstream {
		originRepo := utils.ExtractRepoFromURL(origin)
		upstreamRepo := utils.ExtractRepoFromURL(upstream)

		logging.Debugf("Origin repo: %s, Upstream repo: %s", originRepo, upstreamRepo)

		// We're in a properly configured fork setup if:
		// 1. Origin points to a fork (different from the current repo context)
		// 2. Upstream points to the original repo (same as current repo context)
		if originRepo != repo.NameWithOwner && upstreamRepo == repo.NameWithOwner {
			status.IsInFork = true
			status.HasCorrectRemotes = true
			status.ForkName = originRepo
			status.UpstreamName = upstreamRepo
			logging.Debug("Detected properly configured fork setup")
			return status, nil
		}
	}

	// Fallback: check if repository itself is a fork via GitHub API
	isInFork, err := m.ghCli.IsRepositoryFork()
	if err != nil {
		logging.Debugf("Warning: failed to check if repository is a fork via API: %v", err)
		// Don't return error here, continue with remote-based detection
	} else if isInFork {
		status.IsInFork = true
		status.ForkName = repo.NameWithOwner
		logging.Debug("Detected fork via GitHub API, but remotes may not be properly configured")

		// Even if it's a fork via API, check if remotes are configured
		if hasOrigin && hasUpstream {
			originRepo := utils.ExtractRepoFromURL(origin)
			upstreamRepo := utils.ExtractRepoFromURL(upstream)
			if originRepo != repo.NameWithOwner && upstreamRepo == repo.NameWithOwner {
				status.HasCorrectRemotes = true
				status.UpstreamName = upstreamRepo
			}
		}
	}

	return status, nil
}

func (m *Manager) SetupFork(customForkName string) (*ForkSetupResult, error) {
	logging.PrintInfo("Detecting repository setup...")

	status, err := m.DetectForkStatus()
	if err != nil {
		return nil, err
	}

	result := &ForkSetupResult{}

	// If fork is already properly configured, check if it matches the requested fork
	if status.IsInFork && status.HasCorrectRemotes {
		// If a custom fork name is specified, validate it matches the current fork
		if customForkName != "" && customForkName != status.ForkName {
			return nil, fmt.Errorf("fork mismatch: repository is already configured with fork '%s', but you requested '%s'. Please use the existing fork or reconfigure the repository", status.ForkName, customForkName)
		}

		fmt.Printf("Fork already configured, creating branch...\n")
		result.WasAlreadyConfigured = true
		result.ForkName = status.ForkName
		result.UpstreamName = status.UpstreamName

		repo, err := m.repositoryProvider.GetRepository()
		if err != nil {
			logging.PrintWarn("Warning: Could not get repository info to set default repository")
		} else {
			if err := m.ghCli.SetDefaultRepository(repo.NameWithOwner); err != nil {
				logging.PrintWarn(fmt.Sprintf("Warning: Could not set default repository to upstream: %v", err))
				logging.PrintWarn("You may need to run 'gh repo set-default " + repo.NameWithOwner + "' manually")
			} else {
				logging.Debug("✓ Default repository confirmed as upstream: " + repo.NameWithOwner)
			}
		}

		return result, nil
	}

	repo, err := m.repositoryProvider.GetRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository information: %w", err)
	}

	forkName := customForkName
	if forkName == "" && m.cfg.DefaultOrganization != "" {
		forkName = fmt.Sprintf("%s/%s", m.cfg.DefaultOrganization, repo.Name)
	}

	if err := m.handleForkCreation(status, forkName, result); err != nil {
		return nil, err
	}

	fmt.Printf("Setting up remotes (origin: fork, upstream: original)...\n")
	fmt.Printf("Fetching branches from fork...\n")
	if err := m.gitProvider.FetchBranchFromOrigin("main"); err != nil {
		if err := m.gitProvider.FetchBranchFromOrigin("master"); err != nil {
			logging.PrintWarn("Could not fetch main/master branch from fork")
		}
	}

	result.UpstreamName = repo.NameWithOwner
	if result.ForkName == "" {
		result.ForkName = repo.NameWithOwner
	}

	if err := m.ghCli.SetDefaultRepository(repo.NameWithOwner); err != nil {
		logging.PrintWarn(fmt.Sprintf("Warning: Could not set default repository to upstream: %v", err))
		logging.PrintWarn("You may need to run 'gh repo set-default " + repo.NameWithOwner + "' manually")
	} else {
		logging.Debug("✓ Default repository set to upstream: " + repo.NameWithOwner)
	}

	return result, nil
}

func (m *Manager) handleForkCreation(status *ForkStatus, forkName string, result *ForkSetupResult) error {
	if !status.IsInFork {
		return m.createNewFork(forkName, result)
	}

	result.ForkName = status.ForkName
	if !status.HasCorrectRemotes {
		fmt.Printf("Fork detected but remotes need configuration...\n")
	}
	return nil
}

func (m *Manager) createNewFork(forkName string, result *ForkSetupResult) error {
	if forkName != "" {
		return m.createNamedFork(forkName, result)
	}
	return m.createDefaultFork(result)
}

func (m *Manager) createNamedFork(forkName string, result *ForkSetupResult) error {
	exists, err := m.ghCli.ForkExists(forkName)
	if err != nil {
		return fmt.Errorf("failed to check if fork exists: %w", err)
	}

	if exists {
		fmt.Printf("Fork %s already exists, configuring for use...\n", forkName)
		result.ForkCreated = false
		result.ForkName = forkName

		// Actually configure the remotes when fork exists
		if err := m.ghCli.ConfigureRemotesForExistingFork(forkName); err != nil {
			return fmt.Errorf("failed to configure remotes: %w", err)
		}

		return nil
	}

	return m.performForkCreation(forkName, result)
}

func (m *Manager) createDefaultFork(result *ForkSetupResult) error {
	fmt.Printf("No fork detected. Creating fork...")

	if err := m.requestUserConfirmation(); err != nil {
		return err
	}

	return m.executeForkCreation("", result)
}

func (m *Manager) performForkCreation(forkName string, result *ForkSetupResult) error {
	fmt.Printf("No fork detected. Creating fork")
	if forkName != "" {
		fmt.Printf(" %s", logging.PaintInfo(forkName))
	}
	fmt.Println("...")

	if err := m.requestUserConfirmation(); err != nil {
		return err
	}

	return m.executeForkCreation(forkName, result)
}

func (m *Manager) requestUserConfirmation() error {
	if m.cfg.IsInteractive {
		confirmed, err := m.userInteractionProvider.AskUserForConfirmation(
			"Do you want to create a fork and configure it for development?", true)
		if err != nil {
			return err
		}
		if !confirmed {
			return fmt.Errorf("fork creation cancelled by user")
		}
	}
	return nil
}

func (m *Manager) executeForkCreation(forkName string, result *ForkSetupResult) error {
	if err := m.ghCli.CreateFork(forkName); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			logging.PrintWarn("Fork already exists, continuing with setup...")
			result.ForkCreated = false
		} else {
			return fmt.Errorf("failed to create fork: %w", err)
		}
	} else {
		fmt.Printf("✓ Fork created successfully\n")
		result.ForkCreated = true
	}

	if forkName != "" {
		result.ForkName = forkName
	} else {
		return m.updateForkNameFromStatus(result)
	}

	return nil
}

func (m *Manager) updateForkNameFromStatus(result *ForkSetupResult) error {
	updatedStatus, err := m.DetectForkStatus()
	if err != nil {
		return fmt.Errorf("failed to detect fork status after creation: %w", err)
	}
	result.ForkName = updatedStatus.ForkName
	return nil
}
