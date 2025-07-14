package fork

import (
	"errors"
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/domain"
)

// Mock implementations for testing
type mockRepositoryProvider struct {
	repo *domain.Repository
	err  error
}

func (m *mockRepositoryProvider) GetRepository() (*domain.Repository, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.repo, nil
}

type mockGitProvider struct {
	fetchError error
}

func (m *mockGitProvider) BranchExists(branch string) bool                         { return false }
func (m *mockGitProvider) FetchBranchFromOrigin(branch string) error              { return m.fetchError }
func (m *mockGitProvider) CheckoutNewBranchFromOrigin(branch, base string) error  { return nil }
func (m *mockGitProvider) GetCurrentBranch() (string, error)                      { return "main", nil }
func (m *mockGitProvider) FindBranch(substring string) (string, bool)             { return "", false }
func (m *mockGitProvider) CheckoutBranch(branch string) error                     { return nil }
func (m *mockGitProvider) GetCommitsToPush(branch string) ([]string, error)       { return nil, nil }
func (m *mockGitProvider) RemoteBranchExists(branch string) bool                  { return false }
func (m *mockGitProvider) CommitEmpty(message string) error                       { return nil }
func (m *mockGitProvider) PushBranch(branch string) error                         { return nil }
func (m *mockGitProvider) GetRepositoryRoot() (string, error)                     { return "/tmp", nil }

type mockUserInteractionProvider struct {
	confirmationResult bool
	confirmationError  error
}

func (m *mockUserInteractionProvider) AskUserForConfirmation(msg string, defaultAnswer bool) (bool, error) {
	if m.confirmationError != nil {
		return false, m.confirmationError
	}
	return m.confirmationResult, nil
}

func (m *mockUserInteractionProvider) SelectOrInputPrompt(message string, validValues []string, variable *string, required bool) error {
	return nil
}

func (m *mockUserInteractionProvider) SelectOrInput(name string, validValues []string, variable *string, required bool) error {
	return nil
}

type mockForkProvider struct {
	isRepositoryFork       bool
	isRepositoryForkError  error
	createForkError        error
	setDefaultRepoError    error
	remoteConfiguration    map[string]string
	remoteConfigError      error
}

func (m *mockForkProvider) IsRepositoryFork() (bool, error) {
	return m.isRepositoryFork, m.isRepositoryForkError
}

func (m *mockForkProvider) CreateFork(forkName string) error {
	return m.createForkError
}

func (m *mockForkProvider) SetDefaultRepository(repo string) error {
	return m.setDefaultRepoError
}

func (m *mockForkProvider) GetRemoteConfiguration() (map[string]string, error) {
	if m.remoteConfigError != nil {
		return nil, m.remoteConfigError
	}
	return m.remoteConfiguration, nil
}

func TestDetectForkStatus_NotAFork(t *testing.T) {
	repo := &domain.Repository{
		Name:             "gh-sherpa",
		Owner:            "InditexTech",
		NameWithOwner:    "InditexTech/gh-sherpa",
		DefaultBranchRef: "main",
	}

	repoProvider := &mockRepositoryProvider{repo: repo}
	gitProvider := &mockGitProvider{}
	userProvider := &mockUserInteractionProvider{}
	forkProvider := &mockForkProvider{
		isRepositoryFork:    false,
		remoteConfiguration: map[string]string{},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	status, err := manager.DetectForkStatus()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if status.IsInFork {
		t.Error("Expected IsInFork to be false")
	}

	if status.HasCorrectRemotes {
		t.Error("Expected HasCorrectRemotes to be false")
	}
}

func TestDetectForkStatus_IsFork(t *testing.T) {
	repo := &domain.Repository{
		Name:             "gh-sherpa",
		Owner:            "user",
		NameWithOwner:    "user/gh-sherpa",
		DefaultBranchRef: "main",
	}

	repoProvider := &mockRepositoryProvider{repo: repo}
	gitProvider := &mockGitProvider{}
	userProvider := &mockUserInteractionProvider{}
	forkProvider := &mockForkProvider{
		isRepositoryFork: true,
		remoteConfiguration: map[string]string{
			"origin":   "https://github.com/user/gh-sherpa.git",
			"upstream": "https://github.com/InditexTech/gh-sherpa.git",
		},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	status, err := manager.DetectForkStatus()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !status.IsInFork {
		t.Error("Expected IsInFork to be true")
	}

	if !status.HasCorrectRemotes {
		t.Error("Expected HasCorrectRemotes to be true")
	}

	if status.ForkName != "user/gh-sherpa" {
		t.Errorf("Expected ForkName to be 'user/gh-sherpa', got %s", status.ForkName)
	}

	if status.UpstreamName != "InditexTech/gh-sherpa" {
		t.Errorf("Expected UpstreamName to be 'InditexTech/gh-sherpa', got %s", status.UpstreamName)
	}
}

func TestSetupFork_AlreadyConfigured(t *testing.T) {
	repo := &domain.Repository{
		Name:             "gh-sherpa",
		Owner:            "user",
		NameWithOwner:    "user/gh-sherpa",
		DefaultBranchRef: "main",
	}

	repoProvider := &mockRepositoryProvider{repo: repo}
	gitProvider := &mockGitProvider{}
	userProvider := &mockUserInteractionProvider{}
	forkProvider := &mockForkProvider{
		isRepositoryFork: true,
		remoteConfiguration: map[string]string{
			"origin":   "https://github.com/user/gh-sherpa.git",
			"upstream": "https://github.com/InditexTech/gh-sherpa.git",
		},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	result, err := manager.SetupFork("")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result.WasAlreadyConfigured {
		t.Error("Expected WasAlreadyConfigured to be true")
	}

	if result.ForkCreated {
		t.Error("Expected ForkCreated to be false")
	}
}

func TestSetupFork_CreateNewFork(t *testing.T) {
	repo := &domain.Repository{
		Name:             "gh-sherpa",
		Owner:            "InditexTech",
		NameWithOwner:    "InditexTech/gh-sherpa",
		DefaultBranchRef: "main",
	}

	repoProvider := &mockRepositoryProvider{repo: repo}
	gitProvider := &mockGitProvider{}
	userProvider := &mockUserInteractionProvider{confirmationResult: true}
	forkProvider := &mockForkProvider{
		isRepositoryFork:    false,
		remoteConfiguration: map[string]string{},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	result, err := manager.SetupFork("")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.WasAlreadyConfigured {
		t.Error("Expected WasAlreadyConfigured to be false")
	}

	if !result.ForkCreated {
		t.Error("Expected ForkCreated to be true")
	}
}

func TestSetupFork_UserDeclinesToCreateFork(t *testing.T) {
	repo := &domain.Repository{
		Name:             "gh-sherpa",
		Owner:            "InditexTech",
		NameWithOwner:    "InditexTech/gh-sherpa",
		DefaultBranchRef: "main",
	}

	repoProvider := &mockRepositoryProvider{repo: repo}
	gitProvider := &mockGitProvider{}
	userProvider := &mockUserInteractionProvider{confirmationResult: false}
	forkProvider := &mockForkProvider{
		isRepositoryFork:    false,
		remoteConfiguration: map[string]string{},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	_, err := manager.SetupFork("")

	if err == nil {
		t.Error("Expected error when user declines to create fork")
	}

	if err.Error() != "fork creation cancelled by user" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestSetupFork_CreateForkError(t *testing.T) {
	repo := &domain.Repository{
		Name:             "gh-sherpa",
		Owner:            "InditexTech",
		NameWithOwner:    "InditexTech/gh-sherpa",
		DefaultBranchRef: "main",
	}

	repoProvider := &mockRepositoryProvider{repo: repo}
	gitProvider := &mockGitProvider{}
	userProvider := &mockUserInteractionProvider{confirmationResult: true}
	forkProvider := &mockForkProvider{
		isRepositoryFork:    false,
		remoteConfiguration: map[string]string{},
		createForkError:     errors.New("fork creation failed"),
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	_, err := manager.SetupFork("")

	if err == nil {
		t.Error("Expected error when fork creation fails")
	}

	if err.Error() != "failed to create fork: fork creation failed" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}
