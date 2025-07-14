package fork

import (
	"fmt"
	"strings"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/logging"
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
	SetDefaultRepository(repo string) error
	GetRemoteConfiguration() (map[string]string, error)
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
		originRepo := extractRepoFromURL(origin)
		upstreamRepo := extractRepoFromURL(upstream)

		logging.Debugf("Origin repo: %s, Upstream repo: %s", originRepo, upstreamRepo)

		// We're in a properly configured fork setup if:
		// 1. Origin points to current repo (the fork)
		// 2. Upstream points to a different repo (the original)
		if originRepo == repo.NameWithOwner && upstreamRepo != repo.NameWithOwner {
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
			originRepo := extractRepoFromURL(origin)
			upstreamRepo := extractRepoFromURL(upstream)
			if originRepo == repo.NameWithOwner && upstreamRepo != repo.NameWithOwner {
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

	// If fork is already properly configured, return early
	if status.IsInFork && status.HasCorrectRemotes {
		fmt.Printf("Fork already configured, creating branch...\n")
		result.WasAlreadyConfigured = true
		result.ForkName = status.ForkName
		result.UpstreamName = status.UpstreamName
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

	// Only try to create fork if we're not in a fork yet
	if !status.IsInFork {
		fmt.Printf("No fork detected. Creating fork")
		if forkName != "" {
			fmt.Printf(" %s", logging.PaintInfo(forkName))
		}
		fmt.Println("...")

		if m.cfg.IsInteractive {
			confirmed, err := m.userInteractionProvider.AskUserForConfirmation(
				"Do you want to create a fork and configure it for development?", true)
			if err != nil {
				return nil, err
			}
			if !confirmed {
				return nil, fmt.Errorf("fork creation cancelled by user")
			}
		}

		if err := m.ghCli.CreateFork(forkName); err != nil {
			// Check if error is about fork already existing
			if strings.Contains(err.Error(), "already exists") {
				logging.PrintWarn("Fork already exists, continuing with setup...")
				result.ForkCreated = false
			} else {
				return nil, fmt.Errorf("failed to create fork: %w", err)
			}
		} else {
			fmt.Printf("✓ Fork created successfully\n")
			result.ForkCreated = true
		}

		// Determine the fork name for result
		if forkName == "" {
			updatedStatus, err := m.DetectForkStatus()
			if err != nil {
				return nil, fmt.Errorf("failed to detect fork status after creation: %w", err)
			}
			result.ForkName = updatedStatus.ForkName
		} else {
			result.ForkName = forkName
		}
	} else {
		// We're in a fork but remotes might not be configured correctly
		result.ForkName = status.ForkName
		if !status.HasCorrectRemotes {
			fmt.Printf("Fork detected but remotes need configuration...\n")
		}
	}

	// Only do remote setup if remotes are not correctly configured
	if !status.HasCorrectRemotes {
		fmt.Printf("Setting up remotes (origin: fork, upstream: original)...\n")

		fmt.Printf("Setting default repository to upstream...\n")
		if err := m.ghCli.SetDefaultRepository(repo.NameWithOwner); err != nil {
			return nil, fmt.Errorf("failed to set default repository: %w", err)
		}

		fmt.Printf("Fetching branches from fork...\n")
		if err := m.gitProvider.FetchBranchFromOrigin("main"); err != nil {
			if err := m.gitProvider.FetchBranchFromOrigin("master"); err != nil {
				logging.PrintWarn("Could not fetch main/master branch from fork")
			}
		}
	}

	result.UpstreamName = repo.NameWithOwner
	if result.ForkName == "" {
		result.ForkName = repo.NameWithOwner
	}

	return result, nil
}

func extractRepoFromURL(url string) string {

	if strings.Contains(url, "github.com") {
		parts := strings.Split(url, "/")
		if len(parts) >= 2 {
			owner := parts[len(parts)-2]
			repo := parts[len(parts)-1]

			repo = strings.TrimSuffix(repo, ".git")

			if strings.Contains(owner, ":") {
				owner = strings.Split(owner, ":")[1]
			}

			return fmt.Sprintf("%s/%s", owner, repo)
		}
	}

	return url
}
