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
	issueTracker            *domainFakes.FakeIssueTracker
	userInteractionProvider *domainMocks.MockUserInteractionProvider
	branchProvider          *domainMocks.MockBranchProvider
	repositoryProvider      *domainFakes.FakeRepositoryProvider
}

func (s *CreateBranchExecutionTestSuite) setGetBranchName(branch string) {
	s.branchProvider.EXPECT().GetBranchName(mock.Anything, mock.Anything, mock.Anything).Return(branch, nil).Once()
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

	s.issueTracker = domainFakes.NewFakeIssueTracker()
	s.issueTracker.IssueTrackerType = domain.IssueTrackerTypeGithub
	s.issueTracker.AddIssue("1", issue_types.Feature)
	s.issueTracker.AddIssue("3", issue_types.Documentation)
	s.issueTracker.AddIssue("6", issue_types.Refactoring)

	s.issueTrackerProvider = domainFakes.NewFakeIssueTrackerProvider(s.issueTracker)

	s.userInteractionProvider = s.initializeUserInteractionProvider()
	s.branchProvider = s.initializeBranchProvider()
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
		s.setGetBranchName(s.defaultBranchName)
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
		s.setGetBranchName(branchName)

		s.uc.Cfg.IssueID = "3"
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.ErrorContains(err, fmt.Sprintf("a local branch with the name %s already exists", branchName))
	})

	s.Run("should create branch if branch doesn't exists with default flag", func() {
		s.setGetBranchName(s.defaultBranchName)

		s.uc.Cfg.IssueID = "1"
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.NoError(err)
		s.True(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should create branch if not exists without default flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Maybe()

		s.setGetBranchName(s.defaultBranchName)

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
		s.setGetBranchName(branchName)

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

func (s *CreateGithubBranchExecutionTestSuite) initializeBranchProvider() *domainMocks.MockBranchProvider {
	branchProvider := &domainMocks.MockBranchProvider{}

	return branchProvider
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

	s.issueTracker = domainFakes.NewFakeIssueTracker()
	s.issueTracker.IssueTrackerType = domain.IssueTrackerTypeJira
	s.issueTracker.AddIssue("PROJECTKEY-1", issue_types.Feature)
	s.issueTracker.AddIssue("PROJECTKEY-3", issue_types.Documentation)
	s.issueTracker.AddIssue("PROJECTKEY-6", issue_types.Refactoring)

	s.issueTrackerProvider = domainFakes.NewFakeIssueTrackerProvider(s.issueTracker)

	s.userInteractionProvider = s.initializeUserInteractionProvider()
	s.branchProvider = s.initializeBranchProvider()
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
		s.setGetBranchName(s.defaultBranchName)
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
		s.setGetBranchName(branchName)

		s.uc.Cfg.IssueID = "PROJECTKEY-3"
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.ErrorContains(err, fmt.Sprintf("a local branch with the name %s already exists", branchName))
	})

	s.Run("should create branch if branch doesn't exists with default flag", func() {
		s.setGetBranchName(s.defaultBranchName)

		s.uc.Cfg.IssueID = issueID
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.NoError(err)
		s.True(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should create branch if not exists without default flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Maybe()

		s.setGetBranchName(s.defaultBranchName)

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
		s.setGetBranchName(branchName)

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

func (s *CreateJiraBranchExecutionTestSuite) initializeBranchProvider() *domainMocks.MockBranchProvider {
	branchProvider := &domainMocks.MockBranchProvider{}

	return branchProvider
}
