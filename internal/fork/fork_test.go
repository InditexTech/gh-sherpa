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

func (m *mockGitProvider) BranchExists(branch string) bool                       { return false }
func (m *mockGitProvider) FetchBranchFromOrigin(branch string) error             { return m.fetchError }
func (m *mockGitProvider) CheckoutNewBranchFromOrigin(branch, base string) error { return nil }
func (m *mockGitProvider) GetCurrentBranch() (string, error)                     { return "main", nil }
func (m *mockGitProvider) FindBranch(substring string) (string, bool)            { return "", false }
func (m *mockGitProvider) CheckoutBranch(branch string) error                    { return nil }
func (m *mockGitProvider) GetCommitsToPush(branch string) ([]string, error)      { return nil, nil }
func (m *mockGitProvider) RemoteBranchExists(branch string) bool                 { return false }
func (m *mockGitProvider) CommitEmpty(message string) error                      { return nil }
func (m *mockGitProvider) PushBranch(branch string) error                        { return nil }
func (m *mockGitProvider) GetRepositoryRoot() (string, error)                    { return "/tmp", nil }

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
	isRepositoryFork      bool
	isRepositoryForkError error
	createForkError       error
	setDefaultRepoError   error
	remoteConfiguration   map[string]string
	remoteConfigError     error
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
		Owner:            "InditexTech",
		NameWithOwner:    "InditexTech/gh-sherpa",
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
		Owner:            "InditexTech",
		NameWithOwner:    "InditexTech/gh-sherpa",
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

func TestDetectForkStatus_ForkViaAPI_WithCorrectRemotes(t *testing.T) {
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

func TestDetectForkStatus_ForkViaAPI_WithoutCorrectRemotes(t *testing.T) {
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
			"origin": "https://github.com/user/gh-sherpa.git",
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

	if status.HasCorrectRemotes {
		t.Error("Expected HasCorrectRemotes to be false")
	}

	if status.ForkName != "user/gh-sherpa" {
		t.Errorf("Expected ForkName to be 'user/gh-sherpa', got %s", status.ForkName)
	}

	if status.UpstreamName != "" {
		t.Errorf("Expected UpstreamName to be empty, got %s", status.UpstreamName)
	}
}

func TestDetectForkStatus_ForkViaAPI_NoRemotes(t *testing.T) {
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
		isRepositoryFork:    true,
		remoteConfiguration: map[string]string{},
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

	if status.HasCorrectRemotes {
		t.Error("Expected HasCorrectRemotes to be false")
	}

	if status.ForkName != "user/gh-sherpa" {
		t.Errorf("Expected ForkName to be 'user/gh-sherpa', got %s", status.ForkName)
	}

	if status.UpstreamName != "" {
		t.Errorf("Expected UpstreamName to be empty, got %s", status.UpstreamName)
	}
}

func TestDetectForkStatus_APIError_ContinuesWithRemoteDetection(t *testing.T) {
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
		isRepositoryFork:      false,
		isRepositoryForkError: errors.New("API error"),
		remoteConfiguration:   map[string]string{},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	status, err := manager.DetectForkStatus()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if status.IsInFork {
		t.Error("Expected IsInFork to be false when API fails and no remotes configured")
	}

	if status.HasCorrectRemotes {
		t.Error("Expected HasCorrectRemotes to be false")
	}
}

func TestSetupFork_InForkButRemotesNotConfigured(t *testing.T) {
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
			"origin": "https://github.com/user/gh-sherpa.git",
		},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	result, err := manager.SetupFork("")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.WasAlreadyConfigured {
		t.Error("Expected WasAlreadyConfigured to be false when remotes are not configured correctly")
	}

	if result.ForkCreated {
		t.Error("Expected ForkCreated to be false since we're already in a fork")
	}

	if result.ForkName != "user/gh-sherpa" {
		t.Errorf("Expected ForkName to be 'user/gh-sherpa', got %s", result.ForkName)
	}

	if result.UpstreamName != "user/gh-sherpa" {
		t.Errorf("Expected UpstreamName to be set, got %s", result.UpstreamName)
	}
}

func TestDetectForkStatus_ForkViaAPI_RemotesConfiguredInSecondCheck(t *testing.T) {
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

func TestDetectForkStatus_ForkViaAPI_CorrectRemotesButWrongConfiguration(t *testing.T) {
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
			"upstream": "https://github.com/user/gh-sherpa.git",
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

	if status.HasCorrectRemotes {
		t.Error("Expected HasCorrectRemotes to be false when both remotes point to same repo")
	}

	if status.ForkName != "user/gh-sherpa" {
		t.Errorf("Expected ForkName to be 'user/gh-sherpa', got %s", status.ForkName)
	}

	if status.UpstreamName != "" {
		t.Errorf("Expected UpstreamName to be empty, got %s", status.UpstreamName)
	}
}

func TestSetupFork_CreateForkAlreadyExists(t *testing.T) {
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
		createForkError:     errors.New("fork already exists for user"),
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	result, err := manager.SetupFork("")

	if err != nil {
		t.Errorf("Expected no error when fork already exists, got %v", err)
	}

	if result.ForkCreated {
		t.Error("Expected ForkCreated to be false when fork already exists")
	}

	if result.WasAlreadyConfigured {
		t.Error("Expected WasAlreadyConfigured to be false")
	}

	if result.UpstreamName != "InditexTech/gh-sherpa" {
		t.Errorf("Expected UpstreamName to be 'InditexTech/gh-sherpa', got %s", result.UpstreamName)
	}
}

func TestSetupFork_CreateForkAlreadyExistsNonInteractive(t *testing.T) {
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
		createForkError:     errors.New("repository already exists in the destination"),
	}

	cfg := Configuration{IsInteractive: false}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	result, err := manager.SetupFork("user/gh-sherpa")

	if err != nil {
		t.Errorf("Expected no error when fork already exists, got %v", err)
	}

	if result.ForkCreated {
		t.Error("Expected ForkCreated to be false when fork already exists")
	}

	if result.ForkName != "user/gh-sherpa" {
		t.Errorf("Expected ForkName to be 'user/gh-sherpa', got %s", result.ForkName)
	}

	if result.UpstreamName != "InditexTech/gh-sherpa" {
		t.Errorf("Expected UpstreamName to be 'InditexTech/gh-sherpa', got %s", result.UpstreamName)
	}
}

func TestSetupFork_CreateForkOtherError(t *testing.T) {
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
		createForkError:     errors.New("network error"),
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	_, err := manager.SetupFork("")

	if err == nil {
		t.Error("Expected error for network error when creating fork")
	}

	if err.Error() != "failed to create fork: network error" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestDetectForkStatus_RemotesCorrectlyConfigured_ViaRemoteDetection(t *testing.T) {
	// Test the specific condition: originRepo != repo.NameWithOwner && upstreamRepo == repo.NameWithOwner
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
		isRepositoryFork: false, // Testing remote-based detection, not API
		remoteConfiguration: map[string]string{
			"origin":   "https://github.com/user/gh-sherpa.git",        // Different from repo.NameWithOwner
			"upstream": "https://github.com/InditexTech/gh-sherpa.git", // Same as repo.NameWithOwner
		},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	status, err := manager.DetectForkStatus()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !status.IsInFork {
		t.Error("Expected IsInFork to be true when remotes are properly configured")
	}

	if !status.HasCorrectRemotes {
		t.Error("Expected HasCorrectRemotes to be true when origin != repo and upstream == repo")
	}

	if status.ForkName != "user/gh-sherpa" {
		t.Errorf("Expected ForkName to be 'user/gh-sherpa', got %s", status.ForkName)
	}

	if status.UpstreamName != "InditexTech/gh-sherpa" {
		t.Errorf("Expected UpstreamName to be 'InditexTech/gh-sherpa', got %s", status.UpstreamName)
	}
}

func TestDetectForkStatus_RemotesIncorrectlyConfigured_OriginSameAsRepo(t *testing.T) {
	// Test when origin is same as repo.NameWithOwner (should not be detected as properly configured fork)
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
		isRepositoryFork: false,
		remoteConfiguration: map[string]string{
			"origin":   "https://github.com/InditexTech/gh-sherpa.git", // Same as repo.NameWithOwner
			"upstream": "https://github.com/InditexTech/gh-sherpa.git", // Same as repo.NameWithOwner
		},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	status, err := manager.DetectForkStatus()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if status.IsInFork {
		t.Error("Expected IsInFork to be false when origin is same as current repo")
	}

	if status.HasCorrectRemotes {
		t.Error("Expected HasCorrectRemotes to be false when origin is same as current repo")
	}
}

func TestDetectForkStatus_RemotesIncorrectlyConfigured_UpstreamDifferentFromRepo(t *testing.T) {
	// Test when upstream is different from repo.NameWithOwner (should not be detected as properly configured fork)
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
		isRepositoryFork: false,
		remoteConfiguration: map[string]string{
			"origin":   "https://github.com/user/gh-sherpa.git",      // Different from repo.NameWithOwner
			"upstream": "https://github.com/different/gh-sherpa.git", // Different from repo.NameWithOwner
		},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	status, err := manager.DetectForkStatus()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if status.IsInFork {
		t.Error("Expected IsInFork to be false when upstream is different from current repo")
	}

	if status.HasCorrectRemotes {
		t.Error("Expected HasCorrectRemotes to be false when upstream is different from current repo")
	}
}

func TestDetectForkStatus_ForkViaAPI_RemotesCorrectlyConfiguredInSecondaryCheck(t *testing.T) {
	// Test the specific condition in the API fallback section:
	// originRepo != repo.NameWithOwner && upstreamRepo == repo.NameWithOwner
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
		isRepositoryFork: true, // Fork detected via API
		remoteConfiguration: map[string]string{
			"origin":   "https://github.com/user/gh-sherpa.git",        // Different from repo.NameWithOwner
			"upstream": "https://github.com/InditexTech/gh-sherpa.git", // Same as repo.NameWithOwner
		},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	status, err := manager.DetectForkStatus()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !status.IsInFork {
		t.Error("Expected IsInFork to be true when detected via API")
	}

	if !status.HasCorrectRemotes {
		t.Error("Expected HasCorrectRemotes to be true when remotes are correctly configured in secondary check")
	}

	if status.UpstreamName != "InditexTech/gh-sherpa" {
		t.Errorf("Expected UpstreamName to be set to 'InditexTech/gh-sherpa', got %s", status.UpstreamName)
	}
}

func TestSetupFork_ForkMismatch_DifferentCustomForkName(t *testing.T) {
	// Test the specific condition: customForkName != "" && customForkName != status.ForkName
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
		isRepositoryFork: true,
		remoteConfiguration: map[string]string{
			"origin":   "https://github.com/existing-user/gh-sherpa.git", // Current fork
			"upstream": "https://github.com/InditexTech/gh-sherpa.git",
		},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	// Try to setup with a different fork name than the existing one
	_, err := manager.SetupFork("different-user/gh-sherpa")

	if err == nil {
		t.Error("Expected error when custom fork name doesn't match existing fork")
	}

	expectedError := "fork mismatch: repository is already configured with fork 'existing-user/gh-sherpa', but you requested 'different-user/gh-sherpa'. Please use the existing fork or reconfigure the repository"
	if err.Error() != expectedError {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestSetupFork_ForkMatch_SameCustomForkName(t *testing.T) {
	// Test when custom fork name matches existing fork (should succeed)
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
		isRepositoryFork: true,
		remoteConfiguration: map[string]string{
			"origin":   "https://github.com/existing-user/gh-sherpa.git", // Current fork
			"upstream": "https://github.com/InditexTech/gh-sherpa.git",
		},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	// Try to setup with the same fork name as the existing one
	result, err := manager.SetupFork("existing-user/gh-sherpa")

	if err != nil {
		t.Errorf("Expected no error when custom fork name matches existing fork, got %v", err)
	}

	if !result.WasAlreadyConfigured {
		t.Error("Expected WasAlreadyConfigured to be true when fork names match")
	}

	if result.ForkName != "existing-user/gh-sherpa" {
		t.Errorf("Expected ForkName to be 'existing-user/gh-sherpa', got %s", result.ForkName)
	}
}

func TestDetectForkStatus_FallbackBranch_RemotesCorrectlyConfigured(t *testing.T) {
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
		isRepositoryFork: true,
		remoteConfiguration: map[string]string{
			"origin":   "https://github.com/fork-owner/gh-sherpa.git",
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
		t.Error("Expected IsInFork to be true when fork detected via API")
	}

	if !status.HasCorrectRemotes {
		t.Error("Expected HasCorrectRemotes to be true when originRepo != repo.NameWithOwner && upstreamRepo == repo.NameWithOwner in fallback")
	}

	if status.UpstreamName != "InditexTech/gh-sherpa" {
		t.Errorf("Expected UpstreamName to be 'InditexTech/gh-sherpa', got %s", status.UpstreamName)
	}

	if status.ForkName != "fork-owner/gh-sherpa" {
		t.Errorf("Expected ForkName to be 'fork-owner/gh-sherpa' (from remote detection), got %s", status.ForkName)
	}
}

func TestSetupFork_NoCustomForkName_ExistingFork(t *testing.T) {
	// Test when no custom fork name is provided and fork already exists (should succeed)
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
		isRepositoryFork: true,
		remoteConfiguration: map[string]string{
			"origin":   "https://github.com/existing-user/gh-sherpa.git",
			"upstream": "https://github.com/InditexTech/gh-sherpa.git",
		},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	// Setup without custom fork name
	result, err := manager.SetupFork("")

	if err != nil {
		t.Errorf("Expected no error when no custom fork name provided, got %v", err)
	}

	if !result.WasAlreadyConfigured {
		t.Error("Expected WasAlreadyConfigured to be true")
	}

	if result.ForkName != "existing-user/gh-sherpa" {
		t.Errorf("Expected ForkName to be 'existing-user/gh-sherpa', got %s", result.ForkName)
	}
}
