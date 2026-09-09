package use_cases_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	domainFakes "github.com/InditexTech/gh-sherpa/internal/fakes/domain"
	"github.com/InditexTech/gh-sherpa/internal/mocks"
	domainMocks "github.com/InditexTech/gh-sherpa/internal/mocks/domain"
	"github.com/InditexTech/gh-sherpa/internal/use_cases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CreateBranchExecutionTestSuite struct {
	suite.Suite
	defaultBranchName       string
	uc                      use_cases.CreateBranch
	gitProvider             *domainFakes.FakeGitProvider
	issueTrackerProvider    *domainFakes.FakeIssueTrackerProvider
	userInteractionProvider *domainMocks.MockUserInteractionProvider
	branchProvider          *domainFakes.FakeBranchProvider
	repositoryProvider      *domainFakes.FakeRepositoryProvider
}

type CreateGithubBranchExecutionTestSuite struct {
	CreateBranchExecutionTestSuite
}

func TestCreateGithubBranchExecutionTestSuite(t *testing.T) {
	suite.Run(t, new(CreateGithubBranchExecutionTestSuite))
}

func (s *CreateGithubBranchExecutionTestSuite) SetupSuite() {
	s.defaultBranchName = "feature/GH-1-sample-issue"
}

func (s *CreateGithubBranchExecutionTestSuite) SetupSubTest() {
	s.gitProvider = domainFakes.NewFakeGitProvider()

	s.issueTrackerProvider = domainFakes.NewFakeIssueTrackerProvider()
	issue1 := domainFakes.NewFakeIssue("1", issue_types.Feature, domain.IssueTrackerTypeGithub)
	s.issueTrackerProvider.AddIssue(issue1)
	issue3 := domainFakes.NewFakeIssue("3", issue_types.Documentation, domain.IssueTrackerTypeGithub)
	s.issueTrackerProvider.AddIssue(issue3)
	issue6 := domainFakes.NewFakeIssue("6", issue_types.Refactoring, domain.IssueTrackerTypeGithub)
	s.issueTrackerProvider.AddIssue(issue6)

	s.userInteractionProvider = s.initializeUserInteractionProvider()

	s.branchProvider = domainFakes.NewFakeBranchProvider()
	s.branchProvider.SetBranchName(s.defaultBranchName)

	s.repositoryProvider = domainFakes.NewRepositoryProvider()

	defaultConfig := use_cases.CreateBranchConfiguration{
		FetchFromOrigin: true,
		IsInteractive:   true,
	}
	s.uc = use_cases.CreateBranch{
		Cfg:                     defaultConfig,
		Git:                     s.gitProvider,
		IssueTrackerProvider:    s.issueTrackerProvider,
		UserInteractionProvider: s.userInteractionProvider,
		BranchProvider:          s.branchProvider,
		RepositoryProvider:      s.repositoryProvider,
	}
}

func (s *CreateGithubBranchExecutionTestSuite) TestCreateBranchExecution() {
	s.Run("should error if could not get git repository", func() {
		s.repositoryProvider.Repository = nil

		s.uc.Cfg.IssueID = "1"

		_, err := s.uc.Execute()

		s.Error(err)
		s.False(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should error if no issue flag is provided", func() {
		_, err := s.uc.Execute()

		s.ErrorContains(err, "sherpa needs an valid issue identifier")
		s.False(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should error if branch already exists with default flag", func() {
		branchName := "feature/GH-3-local-branch"
		s.gitProvider.AddLocalBranches(branchName)
		s.branchProvider.SetBranchName(branchName)

		s.uc.Cfg.IssueID = "3"
		s.uc.Cfg.IsInteractive = false

		_, err := s.uc.Execute()

		s.ErrorContains(err, fmt.Sprintf("a local branch with the name %s already exists", branchName))
	})

	s.Run("should create branch if branch doesn't exists with default flag", func() {
		s.uc.Cfg.IssueID = "1"
		s.uc.Cfg.IsInteractive = false

		_, err := s.uc.Execute()

		s.NoError(err)
		s.True(s.gitProvider.BranchExists(s.defaultBranchName))
		s.Equal([]string{"main"}, s.gitProvider.FetchedBranches)
	})

	s.Run("should not fetch when fetch from origin is disabled", func() {
		s.uc.Cfg.IssueID = "1"
		s.uc.Cfg.IsInteractive = false
		s.uc.Cfg.FetchFromOrigin = false

		_, err := s.uc.Execute()

		s.NoError(err)
		s.Empty(s.gitProvider.FetchedBranches)
		s.True(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should create branch if not exists without default flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Maybe()
		s.uc.Cfg.IssueID = "1"

		_, err := s.uc.Execute()

		s.NoError(err)
		s.True(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should error if branch already exists without default flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Maybe()

		branchName := "feature/GH-3-local-branch"
		s.gitProvider.AddLocalBranches(branchName)
		s.branchProvider.SetBranchName(branchName)

		s.uc.Cfg.IssueID = "3"

		_, err := s.uc.Execute()

		s.ErrorContains(err, fmt.Sprintf("a local branch with the name %s already exists", branchName))
	})
}

func TestCreateBranchWithWorktree(t *testing.T) {
	newUseCase := func() (use_cases.CreateBranch, *domainFakes.FakeGitProvider) {
		gitProvider := domainFakes.NewFakeGitProvider()
		issueTrackerProvider := domainFakes.NewFakeIssueTrackerProvider()
		issueTrackerProvider.AddIssue(domainFakes.NewFakeIssue("1", issue_types.Feature, domain.IssueTrackerTypeGithub))
		branchProvider := domainFakes.NewFakeBranchProvider()
		branchProvider.SetBranchName("feature/GH-1-sample-issue")

		return use_cases.CreateBranch{
			Cfg: use_cases.CreateBranchConfiguration{
				IssueID:         "1",
				FetchFromOrigin: true,
				UseWorktree:     true,
			},
			Git:                     gitProvider,
			IssueTrackerProvider:    issueTrackerProvider,
			UserInteractionProvider: &domainMocks.MockUserInteractionProvider{},
			BranchProvider:          branchProvider,
			RepositoryProvider:      domainFakes.NewRepositoryProvider(),
		}, gitProvider
	}

	t.Run("creates a worktree at the default path without switching the current branch", func(t *testing.T) {
		uc, gitProvider := newUseCase()

		result, err := uc.Execute()

		assert.NoError(t, err)
		repositoryRoot, err := gitProvider.GetRepositoryRoot()
		assert.NoError(t, err)
		expectedPath, err := filepath.Abs(filepath.Join(repositoryRoot, "..", "worktrees", "feature", "GH-1-sample-issue"))
		assert.NoError(t, err)
		assert.Equal(t, expectedPath, result.WorktreePath)
		assert.Equal(t, "feature/GH-1-sample-issue", gitProvider.Worktrees[expectedPath])
		assert.Equal(t, "main", gitProvider.CurrentBranch)
		assert.Equal(t, []string{"main"}, gitProvider.FetchedBranches)
	})

	t.Run("uses a custom path as the exact worktree destination", func(t *testing.T) {
		uc, gitProvider := newUseCase()
		uc.Cfg.UseWorktree = false
		uc.Cfg.WorktreePath = filepath.Join("testdata", "worktree with spaces")

		result, err := uc.Execute()

		assert.NoError(t, err)
		expectedPath, err := filepath.Abs(uc.Cfg.WorktreePath)
		assert.NoError(t, err)
		assert.Equal(t, expectedPath, result.WorktreePath)
		assert.Equal(t, "feature/GH-1-sample-issue", gitProvider.Worktrees[expectedPath])
	})

	t.Run("does not fetch when fetch from origin is disabled", func(t *testing.T) {
		uc, gitProvider := newUseCase()
		uc.Cfg.FetchFromOrigin = false

		_, err := uc.Execute()

		assert.NoError(t, err)
		assert.Empty(t, gitProvider.FetchedBranches)
		assert.Len(t, gitProvider.Worktrees, 1)
	})

	t.Run("does not create a worktree when fetching the base fails", func(t *testing.T) {
		uc, gitProvider := newUseCase()
		gitProvider.ResetRemoteBranches()
		gitProvider.RemoteBranches = nil

		_, err := uc.Execute()

		assert.ErrorContains(t, err, "error while fetching the branch main")
		assert.Equal(t, []string{"main"}, gitProvider.FetchedBranches)
		assert.Empty(t, gitProvider.Worktrees)
		assert.False(t, gitProvider.BranchExists("feature/GH-1-sample-issue"))
	})

	t.Run("does not create a branch or worktree during dry-run", func(t *testing.T) {
		uc, gitProvider := newUseCase()
		uc.Cfg.DryRun = true
		uc.Cfg.OutputFormat = "json"

		output := captureStdout(t, func() {
			result, err := uc.Execute()
			assert.NoError(t, err)
			assert.NotEmpty(t, result.WorktreePath)
		})

		assert.Contains(t, output, `"branch":"feature/GH-1-sample-issue"`)
		assert.Contains(t, output, `"worktree_path":`)
		assert.Empty(t, gitProvider.FetchedBranches)
		assert.Empty(t, gitProvider.Worktrees)
		assert.False(t, gitProvider.BranchExists("feature/GH-1-sample-issue"))
	})

	t.Run("keeps the existing JSON shape for a regular branch", func(t *testing.T) {
		uc, _ := newUseCase()
		uc.Cfg.UseWorktree = false
		uc.Cfg.DryRun = true
		uc.Cfg.OutputFormat = "json"

		output := captureStdout(t, func() {
			_, err := uc.Execute()
			assert.NoError(t, err)
		})

		assert.Equal(t, "{\"branch\":\"feature/GH-1-sample-issue\"}\n", output)
	})
}

func captureStdout(t *testing.T, run func()) string {
	t.Helper()

	reader, writer, err := os.Pipe()
	assert.NoError(t, err)
	originalStdout := os.Stdout
	os.Stdout = writer
	t.Cleanup(func() { os.Stdout = originalStdout })

	run()
	assert.NoError(t, writer.Close())
	os.Stdout = originalStdout

	output, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.NoError(t, reader.Close())
	return string(output)
}

func (s *CreateGithubBranchExecutionTestSuite) initializeUserInteractionProvider() *domainMocks.MockUserInteractionProvider {
	userInteractionProvider := &domainMocks.MockUserInteractionProvider{}

	userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInputPrompt("Label 'kind/feature' found. What type of branch name do you want to create?", []string{"feature", "other"}, mock.Anything, true).Return(nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInput(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	return userInteractionProvider
}

type CreateJiraBranchExecutionTestSuite struct {
	CreateBranchExecutionTestSuite
}

func TestCreateJiraBranchExecutionTestSuite(t *testing.T) {
	suite.Run(t, new(CreateJiraBranchExecutionTestSuite))
}

func (s *CreateJiraBranchExecutionTestSuite) SetupSuite() {
	s.defaultBranchName = "feature/PROJECTKEY-1-sample-issue"
}

func (s *CreateJiraBranchExecutionTestSuite) SetupSubTest() {
	s.gitProvider = domainFakes.NewFakeGitProvider()

	s.issueTrackerProvider = domainFakes.NewFakeIssueTrackerProvider()
	issue1 := domainFakes.NewFakeIssue("PROJECTKEY-1", issue_types.Feature, domain.IssueTrackerTypeJira)
	s.issueTrackerProvider.AddIssue(issue1)
	issue3 := domainFakes.NewFakeIssue("PROJECTKEY-3", issue_types.Documentation, domain.IssueTrackerTypeJira)
	s.issueTrackerProvider.AddIssue(issue3)
	issue6 := domainFakes.NewFakeIssue("PROJECTKEY-6", issue_types.Refactoring, domain.IssueTrackerTypeJira)
	s.issueTrackerProvider.AddIssue(issue6)

	s.userInteractionProvider = s.initializeUserInteractionProvider()

	s.branchProvider = domainFakes.NewFakeBranchProvider()
	s.branchProvider.SetBranchName(s.defaultBranchName)

	s.repositoryProvider = domainFakes.NewRepositoryProvider()

	s.uc = use_cases.CreateBranch{
		Git:                     s.gitProvider,
		RepositoryProvider:      s.repositoryProvider,
		IssueTrackerProvider:    s.issueTrackerProvider,
		UserInteractionProvider: s.userInteractionProvider,
		BranchProvider:          s.branchProvider,
	}
}

func (s *CreateJiraBranchExecutionTestSuite) TestCreateBranchExecution() {
	issueID := "PROJECTKEY-1"

	s.Run("should error if could not get git repository", func() {
		s.repositoryProvider.Repository = nil

		s.uc.Cfg.IssueID = issueID

		_, err := s.uc.Execute()

		s.Error(err)
		s.False(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should error if no issue flag is provided", func() {
		_, err := s.uc.Execute()

		s.ErrorContains(err, "sherpa needs an valid issue identifier")
		s.False(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should error if branch already exists with default flag", func() {
		branchName := "feature/PROJECTKEY-3-local-branch"
		s.gitProvider.AddLocalBranches(branchName)
		s.branchProvider.SetBranchName(branchName)

		s.uc.Cfg.IssueID = "PROJECTKEY-3"
		s.uc.Cfg.IsInteractive = false

		_, err := s.uc.Execute()

		s.ErrorContains(err, fmt.Sprintf("a local branch with the name %s already exists", branchName))
	})

	s.Run("should create branch if branch doesn't exists with default flag", func() {
		s.uc.Cfg.IssueID = issueID
		s.uc.Cfg.IsInteractive = false

		_, err := s.uc.Execute()

		s.NoError(err)
		s.True(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should create branch if not exists without default flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Maybe()
		s.uc.Cfg.IssueID = issueID

		_, err := s.uc.Execute()

		s.NoError(err)
		s.True(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should error if branch already exists without default flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Maybe()

		branchName := "feature/PROJECTKEY-3-local-branch"
		s.gitProvider.AddLocalBranches(branchName)
		s.branchProvider.SetBranchName(branchName)

		s.uc.Cfg.IssueID = issueID

		_, err := s.uc.Execute()

		s.ErrorContains(err, fmt.Sprintf("a local branch with the name %s already exists", branchName))
	})
}

func (s *CreateJiraBranchExecutionTestSuite) initializeUserInteractionProvider() *domainMocks.MockUserInteractionProvider {
	userInteractionProvider := &domainMocks.MockUserInteractionProvider{}

	userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInputPrompt("Issue type 'feature' found. What type of branch name do you want to create?", []string{"feature", "other"}, mock.Anything, true).Return(nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInput(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	return userInteractionProvider
}
