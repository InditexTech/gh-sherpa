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

	isInFork, err := m.ghCli.IsRepositoryFork()
	if err != nil {
		return nil, fmt.Errorf("failed to check if repository is a fork: %w", err)
	}

	remotes, err := m.ghCli.GetRemoteConfiguration()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote configuration: %w", err)
	}

	status := &ForkStatus{
		IsInFork: isInFork,
	}

	if isInFork {
		status.ForkName = repo.NameWithOwner

		origin, hasOrigin := remotes["origin"]
		upstream, hasUpstream := remotes["upstream"]

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

	if status.IsInFork && status.HasCorrectRemotes {
		logging.PrintInfo("Fork already configured, proceeding...")
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
			return nil, fmt.Errorf("failed to create fork: %w", err)
		}

		fmt.Printf("✓ Fork created successfully\n")
		result.ForkCreated = true

		if forkName == "" {
			updatedStatus, err := m.DetectForkStatus()
			if err != nil {
				return nil, fmt.Errorf("failed to detect fork status after creation: %w", err)
			}
			result.ForkName = updatedStatus.ForkName
		} else {
			result.ForkName = forkName
		}
	}

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

	result.UpstreamName = repo.NameWithOwner
	if result.ForkName == "" {
		result.ForkName = repo.NameWithOwner // This should be updated after fork creation
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
