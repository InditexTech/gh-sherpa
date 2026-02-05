package use_cases_test

import (
	"fmt"
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	domainFakes "github.com/InditexTech/gh-sherpa/internal/fakes/domain"
	"github.com/InditexTech/gh-sherpa/internal/mocks"
	domainMocks "github.com/InditexTech/gh-sherpa/internal/mocks/domain"
	"github.com/InditexTech/gh-sherpa/internal/use_cases"
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

		err := s.uc.Execute()

		s.Error(err)
		s.False(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should error if no issue flag is provided", func() {
		err := s.uc.Execute()

		s.ErrorContains(err, "sherpa needs an valid issue identifier")
		s.False(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should error if branch already exists with default flag", func() {
		branchName := "feature/GH-3-local-branch"
		s.gitProvider.AddLocalBranches(branchName)
		s.branchProvider.SetBranchName(branchName)

		s.uc.Cfg.IssueID = "3"
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.ErrorContains(err, fmt.Sprintf("a local branch with the name %s already exists", branchName))
	})

	s.Run("should create branch if branch doesn't exists with default flag", func() {
		s.uc.Cfg.IssueID = "1"
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.NoError(err)
		s.True(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should create branch if not exists without default flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Maybe()
		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

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

		err := s.uc.Execute()

		s.ErrorContains(err, fmt.Sprintf("a local branch with the name %s already exists", branchName))
	})
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

		err := s.uc.Execute()

		s.Error(err)
		s.False(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should error if no issue flag is provided", func() {
		err := s.uc.Execute()

		s.ErrorContains(err, "sherpa needs an valid issue identifier")
		s.False(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should error if branch already exists with default flag", func() {
		branchName := "feature/PROJECTKEY-3-local-branch"
		s.gitProvider.AddLocalBranches(branchName)
		s.branchProvider.SetBranchName(branchName)

		s.uc.Cfg.IssueID = "PROJECTKEY-3"
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.ErrorContains(err, fmt.Sprintf("a local branch with the name %s already exists", branchName))
	})

	s.Run("should create branch if branch doesn't exists with default flag", func() {
		s.uc.Cfg.IssueID = issueID
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.NoError(err)
		s.True(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should create branch if not exists without default flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Maybe()
		s.uc.Cfg.IssueID = issueID

		err := s.uc.Execute()

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

		err := s.uc.Execute()

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

func TestCreateBranchWithWorktree(t *testing.T) {
	suite.Run(t, new(CreateBranchWorktreeTestSuite))
}

type CreateBranchWorktreeTestSuite struct {
	suite.Suite
	gitProvider             *domainFakes.FakeGitProvider
	issueTrackerProvider    *domainFakes.FakeIssueTrackerProvider
	userInteractionProvider *domainMocks.MockUserInteractionProvider
	branchProvider          *domainFakes.FakeBranchProvider
	repositoryProvider      *domainFakes.FakeRepositoryProvider
	uc                      use_cases.CreateBranch
}

func (s *CreateBranchWorktreeTestSuite) SetupTest() {
	s.gitProvider = domainFakes.NewFakeGitProvider()

	s.issueTrackerProvider = domainFakes.NewFakeIssueTrackerProvider()
	issue1 := domainFakes.NewFakeIssue("1", issue_types.Feature, domain.IssueTrackerTypeGithub)
	s.issueTrackerProvider.AddIssue(issue1)

	s.userInteractionProvider = &domainMocks.MockUserInteractionProvider{}
	s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Maybe()

	s.branchProvider = domainFakes.NewFakeBranchProvider()
	s.branchProvider.SetBranchName("feature/GH-1-sample-issue")

	s.repositoryProvider = domainFakes.NewRepositoryProvider()

	s.uc = use_cases.CreateBranch{
		Cfg: use_cases.CreateBranchConfiguration{
			IssueID:         "1",
			FetchFromOrigin: true,
			IsInteractive:   false,
			UseWorktree:     true,
		},
		Git:                     s.gitProvider,
		IssueTrackerProvider:    s.issueTrackerProvider,
		UserInteractionProvider: s.userInteractionProvider,
		BranchProvider:          s.branchProvider,
		RepositoryProvider:      s.repositoryProvider,
	}
}

func (s *CreateBranchWorktreeTestSuite) TestCreateBranchWithWorktree() {
	s.Run("should create worktree with default path", func() {
		err := s.uc.Execute()

		s.NoError(err)
		s.True(s.gitProvider.BranchExists("feature/GH-1-sample-issue"))

		worktrees, _ := s.gitProvider.ListWorktrees()
		s.Len(worktrees, 1)
		s.Equal("feature/GH-1-sample-issue", worktrees[0].Branch)
		s.Contains(worktrees[0].Path, "gh-sherpa-feature/GH-1-sample-issue")
	})

	s.Run("should create worktree with custom path", func() {
		s.uc.Cfg.WorktreePath = "/custom/path/my-worktree"
		s.branchProvider.SetBranchName("feature/GH-1-custom")

		err := s.uc.Execute()

		s.NoError(err)

		worktrees, _ := s.gitProvider.ListWorktrees()
		s.Len(worktrees, 1)
		s.Equal("/custom/path/my-worktree", worktrees[0].Path)
		s.Equal("feature/GH-1-custom", worktrees[0].Branch)
	})

	s.Run("should error if branch already exists", func() {
		s.gitProvider.AddLocalBranches("feature/GH-1-existing")
		s.branchProvider.SetBranchName("feature/GH-1-existing")

		err := s.uc.Execute()

		s.ErrorContains(err, "a local branch with the name feature/GH-1-existing already exists")
	})

	s.Run("should fetch from origin before creating worktree", func() {
		s.gitProvider.ResetRemoteBranches()
		s.gitProvider.AddRemoteBranches("main")
		s.branchProvider.SetBranchName("feature/GH-1-fetch-test")

		err := s.uc.Execute()

		s.NoError(err)
		// Verify branch was created in worktree
		worktrees, _ := s.gitProvider.ListWorktrees()
		s.Len(worktrees, 1)
		s.Equal("feature/GH-1-fetch-test", worktrees[0].Branch)
	})

	s.Run("should not fetch if no-fetch is set", func() {
		s.uc.Cfg.FetchFromOrigin = false
		s.branchProvider.SetBranchName("feature/GH-1-no-fetch")

		err := s.uc.Execute()

		s.NoError(err)
		worktrees, _ := s.gitProvider.ListWorktrees()
		s.Len(worktrees, 1)
		s.Equal("feature/GH-1-no-fetch", worktrees[0].Branch)
	})
}
