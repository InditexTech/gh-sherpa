package use_cases_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/config"
	domainFakes "github.com/InditexTech/gh-sherpa/internal/fakes/domain"
	"github.com/InditexTech/gh-sherpa/internal/mocks"
	domainMocks "github.com/InditexTech/gh-sherpa/internal/mocks/domain"
	"github.com/InditexTech/gh-sherpa/internal/use_cases"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CreatePullRequestExecutionTestSuite struct {
	suite.Suite
	defaultBranchName       string
	uc                      use_cases.CreatePullRequest
	gitProvider             *domainFakes.FakeGitProvider
	issueTrackerProvider    *domainFakes.FakeIssueTrackerProvider
	userInteractionProvider *domainMocks.MockUserInteractionProvider
	pullRequestProvider     *domainFakes.FakePullRequestProvider
	issueTracker            *domainFakes.FakeIssueTracker
	branchProvider          *domainMocks.MockBranchProvider
	repositoryProvider      *domainFakes.FakeRepositoryProvider
}

type CreateGithubPullRequestExecutionTestSuite struct {
	CreatePullRequestExecutionTestSuite
}

func (s *CreatePullRequestExecutionTestSuite) setGetBranchName(branchName string) {
	s.branchProvider.EXPECT().GetBranchName(mock.Anything, mock.Anything, mock.Anything).Return(branchName, nil).Once()
}

func TestCreateGitHubPullRequestExecutionTestSuite(t *testing.T) {
	suite.Run(t, new(CreateGithubPullRequestExecutionTestSuite))
}

func (s *CreateGithubPullRequestExecutionTestSuite) SetupSuite() {
	s.defaultBranchName = "feature/GH-1-sample-issue"
}

func (s *CreateGithubPullRequestExecutionTestSuite) SetupSubTest() {
	s.gitProvider = domainFakes.NewFakeGitProvider()
	s.userInteractionProvider = s.initializeUserInteractionProvider()
	s.pullRequestProvider = domainFakes.NewFakePullRequestProvider()
	s.issueTracker = domainFakes.NewFakeGitHubIssueTracker()
	s.issueTrackerProvider = domainFakes.NewFakeIssueTrackerProvider(s.issueTracker)
	s.branchProvider = s.initializeBranchProvider()
	s.repositoryProvider = domainFakes.NewFakeRepositoryProvider()

	defaultConfig := use_cases.CreatePullRequestConfiguration{
		IsInteractive:   true,
		CloseIssue:      true,
		FetchFromOrigin: true,
		DraftPR:         true,
	}
	s.uc = use_cases.CreatePullRequest{
		Cfg:                     defaultConfig,
		Git:                     s.gitProvider,
		IssueTrackerProvider:    s.issueTrackerProvider,
		UserInteractionProvider: s.userInteractionProvider,
		PullRequestProvider:     s.pullRequestProvider,
		BranchProvider:          s.branchProvider,
		RepositoryProvider:      s.repositoryProvider,
	}
}

func (s *CreateGithubPullRequestExecutionTestSuite) TestCreatePullRequestExecution() {
	s.Run("should error if could not get git repository", func() {
		branchName := "feature/GH-3-local-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		s.repositoryProvider.Repository = nil

		err := s.uc.Execute()

		s.Error(err)
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(branchName))
	})

	s.Run("should error if could not get current branch", func() {
		s.gitProvider.CurrentBranch = ""

		err := s.uc.Execute()

		s.ErrorContains(err, "could not get the current branch name")
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(s.gitProvider.CurrentBranch))
	})

	s.Run("should error if no issue could be identified", func() {
		branchName := "branch-with-no-issue-name"

		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		err := s.uc.Execute()

		expectedError := fmt.Sprintf("could not find an issue identifier in the current branch named %s", branchName)
		s.EqualError(err, expectedError)
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(branchName))
	})

	s.Run("should exit if user does not confirm current branch", func() {
		branchName := "feature/GH-3-local-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation(mock.Anything, mock.Anything).Return(false, nil).Once()

		err := s.uc.Execute()

		s.NoError(err)
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(branchName))
	})

	s.Run("should return error if pr already exists", func() {
		branchName := "feature/GH-3-pull-request-sample"
		s.gitProvider.CurrentBranch = branchName
		s.gitProvider.LocalBranches = append(s.gitProvider.LocalBranches, branchName)
		s.gitProvider.CommitsToPush[branchName] = []string{}
		s.setGetBranchName(branchName)

		err := s.uc.Execute()

		s.ErrorContains(err, "pull request")
		s.ErrorContains(err, "already exists")
	})

	s.Run("should not ask the user for branch confirmation if default flag is used", func() {
		branchName := "feature/GH-3-local-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertNotCalled(s.T(), "AskUserForConfirmation")
	})

	s.Run("should create and push empty commit if remote branch nor pull request exists", func() {
		branchName := "feature/GH-3-local-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		err := s.uc.Execute()

		s.NoError(err)
	})

	s.Run("should return error if remote branch already exists", func() {
		branchName := "feature/GH-1-sample-issue"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		err := s.uc.Execute()

		s.ErrorContains(err, use_cases.ErrRemoteBranchAlreadyExists(s.defaultBranchName).Error())
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(branchName))
	})

	s.Run("should exit if user does not confirm the commit push when default flag is not used", func() {
		branchName := "feature/GH-3-local-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)
		s.gitProvider.CommitsToPush[branchName] = []string{"commit 1", "commit 2"}

		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue pushing all pending commits in this branch and create the pull request", true).Return(false, nil).Once()

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(branchName))
	})

	s.Run("should error if could not create empty commit", func() {
		branchName := "feature/GH-4-commit-error"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		err := s.uc.Execute()

		s.ErrorIs(err, domainFakes.ErrGetCommitsToPush)
	})

	s.Run("should error if could not push branch", func() {
		branchName := "feature/GH-5-with-no-local-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		err := s.uc.Execute()

		s.ErrorContains(err, "could not create the remote branch because")
	})

	s.Run("should error if could not create pull request", func() {
		branchName := "feature/GH-6-with-no-remote-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		err := s.uc.Execute()

		s.ErrorContains(err, "could not create the pull request because")
	})

	s.Run("should checkout local branch if branch exists and user confirms branch usage without default flag and issue flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Once()

		s.setGetBranchName(s.defaultBranchName)

		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
	})

	s.Run("should error if branch already exists when using default and issue flags", func() {
		s.setGetBranchName(s.defaultBranchName)

		s.uc.Cfg.IsInteractive = false
		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.ErrorContains(err, "the branch feature/GH-1-sample-issue already exists")
	})

	s.Run("should return error if remote branch already exists when using issue flags", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(false, nil).Once()

		s.setGetBranchName(s.defaultBranchName)

		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.ErrorContains(err, use_cases.ErrRemoteBranchAlreadyExists("feature/GH-1-sample-issue").Error())
		s.userInteractionProvider.AssertExpectations(s.T())
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(s.gitProvider.CurrentBranch))
	})

	s.Run("should create new branch name if user doesn't confirm default branch name when using issue flags", func() {
		s.gitProvider.LocalBranches = []string{"main", "develop"}
		s.gitProvider.RemoteBranches = []string{"main", "develop"}

		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Once()

		s.setGetBranchName(s.defaultBranchName)

		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
	})

	s.Run("should return error if user doesn't confirm branch name when using issue flags", func() {
		s.gitProvider.LocalBranches = []string{"main", "develop"}
		s.gitProvider.RemoteBranches = []string{"main", "develop"}

		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(false, nil).Once()

		s.setGetBranchName(s.defaultBranchName)

		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(s.gitProvider.CurrentBranch))
	})

	s.Run("should checkout branch if user confirms branch usage with issue flag and no default flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Once()

		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
	})

	s.Run("should not error if pull request is created correctly", func() {
		s.gitProvider.RemoteBranches = []string{"main", "develop"}
		branchName := "feature/GH-1-sample-issue"
		s.gitProvider.CurrentBranch = branchName

		s.setGetBranchName(branchName)

		err := s.uc.Execute()

		s.NoError(err)
		s.Assert().True(s.pullRequestProvider.HasPullRequestForBranch(branchName))
	})

	s.Run("should create branch and pull request if local branch doesn't exists with issue flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.SelectOrInputPrompt)
		s.userInteractionProvider.EXPECT().SelectOrInputPrompt("Label 'kind/feature' found. What type of branch name do you want to create?", []string{"feature", "other"}, mock.Anything, true).Return(nil).Once()
		s.userInteractionProvider.EXPECT().SelectOrInput("additional description (optional). Truncate to 29 chars", []string{}, mock.Anything, false).Return(nil).Once()
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Once()
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Once()

		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.NoError(err)
		s.Assert().True(s.pullRequestProvider.HasPullRequestForBranch(s.gitProvider.CurrentBranch))
	})

	s.Run("should create pull request with no close issue flag", func() {
		branchName := "feature/GH-3-local-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		s.uc.Cfg.CloseIssue = false

		err := s.uc.Execute()

		s.NoError(err)
		s.Assert().True(s.pullRequestProvider.HasPullRequestForBranch(branchName))
	})

	s.Run("should error if could not get issue", func() {
		branchName := "feature/GH-6-with-no-remote-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		err := s.uc.Execute()

		s.Error(err)
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(branchName))
	})
}

func (s *CreateGithubPullRequestExecutionTestSuite) initializeUserInteractionProvider() *domainMocks.MockUserInteractionProvider {
	userInteractionProvider := &domainMocks.MockUserInteractionProvider{}

	userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInputPrompt("Label 'kind/feature' found. What type of branch name do you want to create?", []string{"feature", "other"}, mock.Anything, true).Return(nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInput(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	return userInteractionProvider
}

func (s *CreateGithubPullRequestExecutionTestSuite) initializeBranchProvider() *domainMocks.MockBranchProvider {
	branchProvider := &domainMocks.MockBranchProvider{}

	return branchProvider
}

type CreateJiraPullRequestExecutionTestSuite struct {
	CreatePullRequestExecutionTestSuite
	getConfigFile func() (config.ConfigFile, error)
}

func TestCreateJiraPullRequestExecutionTestSuite(t *testing.T) {
	suite.Run(t, new(CreateJiraPullRequestExecutionTestSuite))
}

func (s *CreateJiraPullRequestExecutionTestSuite) SetupSuite() {
	s.defaultBranchName = "feature/PROJECTKEY-1-sample-issue"
}

func (s *CreateJiraPullRequestExecutionTestSuite) SetupTest() {
	s.getConfigFile = config.GetConfigFile
	dir, _ := os.Getwd()
	path := filepath.Join(dir, "testdata")
	config.GetConfigFile = func() (config.ConfigFile, error) {
		return config.ConfigFile{
			Path: path,
			Name: "config",
			Type: "yml",
		}, nil
	}
}

func (s *CreateJiraPullRequestExecutionTestSuite) TeardownTest() {
	config.GetConfigFile = s.getConfigFile
}

func (s *CreateJiraPullRequestExecutionTestSuite) SetupSubTest() {
	s.gitProvider = domainFakes.NewFakeGitProvider()
	s.userInteractionProvider = s.initializeUserInteractionProvider()
	s.pullRequestProvider = domainFakes.NewFakePullRequestProvider()
	s.issueTracker = domainFakes.NewFakeJiraIssueTracker()
	s.issueTrackerProvider = domainFakes.NewFakeIssueTrackerProvider(s.issueTracker)
	s.branchProvider = s.initializeBranchProvider()
	s.repositoryProvider = domainFakes.NewFakeRepositoryProvider()

	defaultConfig := use_cases.CreatePullRequestConfiguration{
		IsInteractive:   true,
		CloseIssue:      true,
		FetchFromOrigin: true,
		DraftPR:         true,
	}
	s.uc = use_cases.CreatePullRequest{
		Cfg:                     defaultConfig,
		Git:                     s.gitProvider,
		IssueTrackerProvider:    s.issueTrackerProvider,
		UserInteractionProvider: s.userInteractionProvider,
		PullRequestProvider:     s.pullRequestProvider,
		BranchProvider:          s.branchProvider,
		RepositoryProvider:      s.repositoryProvider,
	}
}

func (s *CreateJiraPullRequestExecutionTestSuite) TestCreatePullRequestExecution() {
	s.Run("should error if could not get git repository", func() {
		s.gitProvider.CurrentBranch = "feature/PROJECTKEY-3-local-branch"
		s.repositoryProvider.Repository = nil

		err := s.uc.Execute()

		s.Error(err)
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(s.gitProvider.CurrentBranch))
	})

	s.Run("should error if could not get current branch", func() {
		s.gitProvider.CurrentBranch = ""

		err := s.uc.Execute()

		s.ErrorContains(err, "could not get the current branch name")
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(s.gitProvider.CurrentBranch))
	})

	s.Run("should error if no issue could be identified", func() {
		branchName := "branch-with-no-issue-name"

		s.gitProvider.CurrentBranch = branchName

		err := s.uc.Execute()

		expectedError := fmt.Sprintf("could not find an issue identifier in the current branch named %s", branchName)
		s.EqualError(err, expectedError)
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(branchName))
	})

	s.Run("should exit if user does not confirm current branch", func() {
		branchName := "feature/PROJECTKEY-3-local-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation(mock.Anything, mock.Anything).Return(false, nil).Once()

		err := s.uc.Execute()

		s.NoError(err)
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(branchName))
	})

	s.Run("should return error if pr already exists", func() {
		branchName := "feature/PROJECTKEY-3-pull-request-sample"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)
		s.gitProvider.LocalBranches = append(s.gitProvider.LocalBranches, branchName)
		s.gitProvider.CommitsToPush[branchName] = []string{}

		err := s.uc.Execute()

		s.ErrorContains(err, "pull request")
		s.ErrorContains(err, "already exists")
	})

	s.Run("should not ask the user for branch confirmation if default flag is used", func() {
		branchName := "feature/PROJECTKEY-3-local-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertNotCalled(s.T(), "AskUserForConfirmation")
	})

	s.Run("should create and push empty commit if remote branch does not exists", func() {
		branchName := "feature/PROJECTKEY-3-local-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		err := s.uc.Execute()

		s.NoError(err)
	})

	s.Run("should return error if remote branch already exists", func() {
		branchName := "feature/PROJECTKEY-1-sample-issue"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		err := s.uc.Execute()

		s.ErrorContains(err, use_cases.ErrRemoteBranchAlreadyExists(s.defaultBranchName).Error())
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(branchName))
	})

	s.Run("should exit if user does not confirm the commit push when default flag is not used", func() {
		branchName := "feature/PROJECTKEY-3-local-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)
		s.gitProvider.CommitsToPush[branchName] = []string{"commit 1", "commit 2"}

		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue pushing all pending commits in this branch and create the pull request", true).Return(false, nil).Once()

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(branchName))
	})

	s.Run("should error if could not create empty commit", func() {
		branchName := "feature/PROJECTKEY-4-with-commit-error"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		err := s.uc.Execute()

		s.ErrorIs(err, domainFakes.ErrGetCommitsToPush)
	})

	s.Run("should error if could not push branch", func() {
		branchName := "feature/PROJECTKEY-5-with-no-local-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		err := s.uc.Execute()

		s.ErrorContains(err, "could not create the remote branch because")
	})

	s.Run("should error if could not create pull request", func() {
		branchName := "feature/PROJECTKEY-6-with-no-remote-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		err := s.uc.Execute()

		s.ErrorContains(err, "could not create the pull request because")
	})

	s.Run("should checkout local branch if branch exists and user confirms branch usage without default flag and issue flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Once()

		s.setGetBranchName(s.defaultBranchName)

		s.uc.Cfg.IssueID = "PROJECTKEY-1"

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
	})

	s.Run("should error if branch already exists when using default and issue flags", func() {
		s.uc.Cfg.IssueID = "PROJECTKEY-1"
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.ErrorContains(err, "the branch feature/PROJECTKEY-1-sample-issue already exists")
	})

	s.Run("should return error if remote branch already exists when using issue flags", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(false, nil).Once()

		s.setGetBranchName(s.defaultBranchName)

		s.uc.Cfg.IssueID = "PROJECTKEY-1"

		err := s.uc.Execute()

		s.ErrorContains(err, use_cases.ErrRemoteBranchAlreadyExists("feature/PROJECTKEY-1-sample-issue").Error())
		s.userInteractionProvider.AssertExpectations(s.T())
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(s.gitProvider.CurrentBranch))
	})

	s.Run("should create new branch name if user doesn't confirm default branch name when using issue flags", func() {
		s.gitProvider.LocalBranches = []string{"main", "develop"}
		s.gitProvider.RemoteBranches = []string{"main", "develop"}

		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Once()

		s.setGetBranchName(s.defaultBranchName)

		s.uc.Cfg.IssueID = "PROJECTKEY-1"

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
	})

	s.Run("should return error if user doesn't confirm branch name when using issue flags", func() {
		s.gitProvider.LocalBranches = []string{"main", "develop"}
		s.gitProvider.RemoteBranches = []string{"main", "develop"}

		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(false, nil).Once()

		s.setGetBranchName(s.defaultBranchName)

		s.uc.Cfg.IssueID = "PROJECTKEY-1"

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(s.gitProvider.CurrentBranch))
	})

	s.Run("should checkout branch if user confirms branch usage with issue flag and no default flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Once()

		s.uc.Cfg.IssueID = "PROJECTKEY-1"

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
	})

	s.Run("should not error if pull request is created correctly", func() {
		s.gitProvider.RemoteBranches = []string{"main", "develop"}
		s.gitProvider.CurrentBranch = "feature/PROJECTKEY-1-sample-issue"

		err := s.uc.Execute()

		s.NoError(err)
		s.Assert().True(s.pullRequestProvider.HasPullRequestForBranch(s.gitProvider.CurrentBranch))
	})

	s.Run("should create branch and pull request if local branch doesn't exists with issue flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.SelectOrInputPrompt)
		s.userInteractionProvider.EXPECT().SelectOrInputPrompt("Issue type 'feature' found. What type of branch name do you want to create?", []string{"feature", "other"}, mock.Anything, true).Return(nil).Once()
		s.userInteractionProvider.EXPECT().SelectOrInput("additional description (optional). Truncate to 21 chars", []string{}, mock.Anything, false).Return(nil).Once()
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Once()
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Once()

		s.uc.Cfg.IssueID = "PROJECTKEY-1"

		err := s.uc.Execute()

		s.NoError(err)
		s.Assert().True(s.pullRequestProvider.HasPullRequestForBranch(s.gitProvider.CurrentBranch))
	})

	s.Run("should create pull request with no close issue flag", func() {
		branchName := "feature/PROJECTKEY-3-local-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		s.expectNoPrFound()

		s.uc.Cfg.CloseIssue = false

		err := s.uc.Execute()

		s.NoError(err)
	})

	s.Run("should error if could not get issue", func() {
		branchName := "feature/PROJECTKEY-6-with-no-remote-branch"
		s.gitProvider.CurrentBranch = branchName
		s.setGetBranchName(branchName)

		err := s.uc.Execute()

		s.Error(err)
		s.Assert().False(s.pullRequestProvider.HasPullRequestForBranch(branchName))
	})
}

func (s *CreateJiraPullRequestExecutionTestSuite) expectNoPrFound() {
	branch := s.gitProvider.CurrentBranch
	if pr := s.pullRequestProvider.PullRequests[branch]; pr != nil {
		s.Failf("pull request exists for branch %s", branch)
	}
}

func (s *CreateJiraPullRequestExecutionTestSuite) initializeUserInteractionProvider() *domainMocks.MockUserInteractionProvider {
	userInteractionProvider := &domainMocks.MockUserInteractionProvider{}

	userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInputPrompt("Issue type 'feature' found. What type of branch name do you want to create?", []string{"feature", "other"}, mock.Anything, true).Return(nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInput(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	return userInteractionProvider
}

func (s *CreateJiraPullRequestExecutionTestSuite) initializeBranchProvider() *domainMocks.MockBranchProvider {
	branchProvider := &domainMocks.MockBranchProvider{}

	return branchProvider
}
