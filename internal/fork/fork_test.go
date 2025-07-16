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

func TestDetectForkStatus_ForkViaAPI_WithCorrectRemotes(t *testing.T) {
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
			// No origin/upstream in first check, so it falls back to API
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
			// Only origin, no upstream - should trigger fallback to API
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
		remoteConfiguration: map[string]string{
			// No remotes at all - should trigger fallback to API
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
		remoteConfiguration:   map[string]string{
			// No remotes configured
		},
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
			// Only origin, no upstream - so remotes are not correctly configured
			"origin": "https://github.com/user/gh-sherpa.git",
		},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	result, err := manager.SetupFork("")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should not be marked as already configured since remotes are incorrect
	if result.WasAlreadyConfigured {
		t.Error("Expected WasAlreadyConfigured to be false when remotes are not configured correctly")
	}

	// Should not have created a fork since we're already in one
	if result.ForkCreated {
		t.Error("Expected ForkCreated to be false since we're already in a fork")
	}

	// Should have the correct fork name
	if result.ForkName != "user/gh-sherpa" {
		t.Errorf("Expected ForkName to be 'user/gh-sherpa', got %s", result.ForkName)
	}

	// Should have set the upstream name from the original repo
	if result.UpstreamName != "user/gh-sherpa" {
		t.Errorf("Expected UpstreamName to be set, got %s", result.UpstreamName)
	}
}

func TestDetectForkStatus_ForkViaAPI_RemotesConfiguredInSecondCheck(t *testing.T) {
	// Este test cubre específicamente la rama dentro de "else if isInFork"
	// donde se verifica si los remotes están configurados correctamente
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
		isRepositoryFork: true, // API dice que SÍ es un fork
		remoteConfiguration: map[string]string{
			// La configuración es correcta: origin apunta al fork, upstream al original
			// pero para que llegue al bloque else if, necesitamos que la primera verificación falle
			// esto puede pasar si el extractRepoFromURL no funciona correctamente en la primera pasada
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

	// Debe detectar que es un fork
	if !status.IsInFork {
		t.Error("Expected IsInFork to be true")
	}

	// Debe detectar que los remotes están correctamente configurados
	if !status.HasCorrectRemotes {
		t.Error("Expected HasCorrectRemotes to be true")
	}

	// Debe tener el nombre del fork correcto
	if status.ForkName != "user/gh-sherpa" {
		t.Errorf("Expected ForkName to be 'user/gh-sherpa', got %s", status.ForkName)
	}

	// Debe tener el nombre upstream correcto
	if status.UpstreamName != "InditexTech/gh-sherpa" {
		t.Errorf("Expected UpstreamName to be 'InditexTech/gh-sherpa', got %s", status.UpstreamName)
	}
}

func TestDetectForkStatus_ForkViaAPI_CorrectRemotesButWrongConfiguration(t *testing.T) {
	// Este test verifica el caso donde tenemos origin y upstream, pero la configuración
	// no es la correcta para un fork (por example, upstream apunta al mismo repo que origin)
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
		isRepositoryFork: true, // API dice que SÍ es un fork
		remoteConfiguration: map[string]string{
			// Tenemos ambos remotes, pero upstream apunta al mismo repo (configuración incorrecta)
			"origin":   "https://github.com/user/gh-sherpa.git",
			"upstream": "https://github.com/user/gh-sherpa.git", // ¡Misma URL! Configuración incorrecta
		},
	}

	cfg := Configuration{IsInteractive: true}
	manager := NewManager(cfg, repoProvider, gitProvider, userProvider, forkProvider)

	status, err := manager.DetectForkStatus()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Debe detectar que es un fork (via API)
	if !status.IsInFork {
		t.Error("Expected IsInFork to be true")
	}

	// NO debe detectar que los remotes están correctos (ambos apuntan al mismo repo)
	if status.HasCorrectRemotes {
		t.Error("Expected HasCorrectRemotes to be false when both remotes point to same repo")
	}

	// Debe tener el nombre del fork correcto
	if status.ForkName != "user/gh-sherpa" {
		t.Errorf("Expected ForkName to be 'user/gh-sherpa', got %s", status.ForkName)
	}

	// UpstreamName debe estar vacío porque no se detectó configuración correcta
	if status.UpstreamName != "" {
		t.Errorf("Expected UpstreamName to be empty, got %s", status.UpstreamName)
	}
}
