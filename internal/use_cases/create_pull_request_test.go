package use_cases_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/config"
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

type CreatePullRequestExecutionTestSuite struct {
	suite.Suite
	defaultBranchName       string
	uc                      use_cases.CreatePullRequest
	gitProvider             *domainMocks.MockGitProvider
	issueTrackerProvider    *domainMocks.MockIssueTrackerProvider
	userInteractionProvider *domainMocks.MockUserInteractionProvider
	pullRequestProvider     *domainMocks.MockPullRequestProvider
	issueTracker            *domainMocks.MockIssueTracker
	branchProvider          *domainMocks.MockBranchProvider
	repositoryProvider      *domainMocks.MockRepositoryProvider
}

type CreateGithubPullRequestExecutionTestSuite struct {
	CreatePullRequestExecutionTestSuite
	gitProvider         *domainFakes.FakeGitProvider
	repositoryProvider  *domainFakes.FakeRepositoryProvider
	pullRequestProvider *domainFakes.FakePullRequestProvider
}

func TestCreatePullRequestExecutionTestSuite(t *testing.T) {
	suite.Run(t, new(CreateGithubPullRequestExecutionTestSuite))
}

func (s *CreateGithubPullRequestExecutionTestSuite) SetupSuite() {
	s.defaultBranchName = "feature/GH-1-sample-issue"
}

func (s *CreateGithubPullRequestExecutionTestSuite) SetupSubTest() {
	s.gitProvider = domainFakes.NewFakeGitProvider()
	s.issueTrackerProvider = s.initializeIssueTrackerProvider()
	s.userInteractionProvider = s.initializeUserInteractionProvider()
	s.pullRequestProvider = domainFakes.NewFakePullRequestProvider()
	s.issueTracker = s.initializeIssueTracker()
	s.branchProvider = s.initializeBranchProvider()
	s.repositoryProvider = domainFakes.NewFakeRepositoryProvider()

	mocks.UnsetExpectedCall(&s.issueTrackerProvider.Mock, s.issueTrackerProvider.GetIssueTracker)
	s.issueTrackerProvider.EXPECT().GetIssueTracker(mock.Anything).Return(s.issueTracker, nil).Maybe()

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
		s.gitProvider.CurrentBranch = "feature/GH-3-local-branch"

		s.repositoryProvider.Repository = nil

		err := s.uc.Execute()

		s.Error(err)
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should error if could not get current branch", func() {
		s.gitProvider.CurrentBranch = ""

		err := s.uc.Execute()

		s.ErrorContains(err, "could not get the current branch name")
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should error if no issue could be identified", func() {
		branchName := "branch-with-no-issue-name"

		s.gitProvider.CurrentBranch = branchName

		err := s.uc.Execute()

		expectedError := fmt.Sprintf("could not find an issue identifier in the current branch named %s", branchName)
		s.EqualError(err, expectedError)
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should exit if user does not confirm current branch", func() {
		s.gitProvider.CurrentBranch = "feature/GH-3-local-branch"

		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation(mock.Anything, mock.Anything).Return(false, nil).Once()

		err := s.uc.Execute()

		s.NoError(err)
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should return error if pr already exists", func() {
		s.gitProvider.CurrentBranch = "feature/GH-3-pull-request-sample"
		s.gitProvider.LocalBranches = append(s.gitProvider.LocalBranches, "feature/GH-3-pull-request-sample")
		s.gitProvider.CommitsToPush["feature/GH-3-pull-request-sample"] = []string{}

		err := s.uc.Execute()

		s.ErrorContains(err, "pull request")
		s.ErrorContains(err, "already exists")
	})

	s.Run("should not ask the user for branch confirmation if default flag is used", func() {
		s.gitProvider.CurrentBranch = "feature/GH-3-local-branch"
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertNotCalled(s.T(), "AskUserForConfirmation")
	})

	s.Run("should create and push empty commit if remote branch nor pull request exists", func() {
		s.gitProvider.RemoteBranches = []string{"main", "develop"}
		s.gitProvider.LocalBranches = append(s.gitProvider.LocalBranches, "feature/GH-99-sample-issue")
		s.gitProvider.CurrentBranch = "feature/GH-99-sample-issue"
		s.gitProvider.CommitsToPush["feature/GH-99-sample-issue"] = []string{}
		s.pullRequestProvider.PullRequests["feature/GH-99-sample-issue"] = nil

		err := s.uc.Execute()

		s.NoError(err)
	})

	s.Run("should return error if remote branch already exists", func() {
		s.gitProvider.CurrentBranch = "feature/GH-1-sample-issue"

		err := s.uc.Execute()

		s.ErrorContains(err, use_cases.ErrRemoteBranchAlreadyExists(s.defaultBranchName).Error())
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should exit if user does not confirm the commit push when default flag is not used", func() {
		s.gitProvider.CurrentBranch = "feature/GH-3-local-branch"
		s.gitProvider.CommitsToPush["feature/GH-3-local-branch"] = []string{"commit 1", "commit 2"}

		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue pushing all pending commits in this branch and create the pull request", true).Return(false, nil).Once()

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should error if could not create empty commit", func() {
		s.gitProvider.RemoteBranches = []string{"main", "develop"}
		s.gitProvider.LocalBranches = append(s.gitProvider.LocalBranches, "feature/GH-99-sample-issue")
		s.gitProvider.CurrentBranch = "feature/GH-99-sample-issue"

		err := s.uc.Execute()

		s.ErrorIs(err, domainFakes.ErrGetCommitsToPush)
	})

	s.Run("should error if could not push branch", func() {
		s.gitProvider.RemoteBranches = []string{"main", "develop"}
		s.gitProvider.LocalBranches = []string{"main", "develop"}
		s.gitProvider.CurrentBranch = "feature/GH-99-sample-issue"
		s.gitProvider.CommitsToPush["feature/GH-99-sample-issue"] = []string{}

		err := s.uc.Execute()

		s.ErrorContains(err, "could not create the remote branch because")
	})

	s.Run("should error if could not create pull request", func() {
		s.gitProvider.RemoteBranches = []string{"main", "develop"}
		s.gitProvider.LocalBranches = append(s.gitProvider.LocalBranches, "feature/GH-99-sample-issue")
		s.gitProvider.CurrentBranch = "feature/GH-99-sample-issue"
		s.gitProvider.CommitsToPush["feature/GH-99-sample-issue"] = []string{}

		err := s.uc.Execute()

		s.ErrorContains(err, "could not create the pull request because")
	})

	s.Run("should checkout local branch if branch exists and user confirms branch usage without default flag and issue flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Once()

		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
	})

	s.Run("should error if branch already exists when using default and issue flags", func() {
		s.uc.Cfg.IsInteractive = false
		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.ErrorContains(err, "the branch feature/GH-1-sample-issue already exists")
	})

	s.Run("should abort execution if remote branch already exists when using issue flags", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(false, nil).Once()

		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.ErrorContains(err, use_cases.ErrRemoteBranchAlreadyExists("feature/GH-1-sample-issue").Error())
		s.userInteractionProvider.AssertExpectations(s.T())
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should create new branch name if user doesn't confirm default branch name when using issue flags", func() {
		s.gitProvider.LocalBranches = []string{"main", "develop"}
		s.gitProvider.RemoteBranches = []string{"main", "develop"}

		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Once()

		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
	})

	s.Run("should abort execution if user doesn't confirm branch name when using issue flags", func() {
		s.gitProvider.LocalBranches = []string{"main", "develop"}
		s.gitProvider.RemoteBranches = []string{"main", "develop"}

		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(false, nil).Once()

		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
		s.assertCreatePullRequestNotCalled()
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
		s.gitProvider.CurrentBranch = "feature/GH-1-sample-issue"

		err := s.uc.Execute()

		s.NoError(err)
		s.assertCreatePullRequestCalled()
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
		s.assertCreatePullRequestCalled()
	})

	s.Run("should create pull request with no close issue flag", func() {
		s.gitProvider.RemoteBranches = []string{"main", "develop"}
		s.gitProvider.LocalBranches = append(s.gitProvider.LocalBranches, "feature/GH-99-sample-issue")
		s.gitProvider.CurrentBranch = "feature/GH-99-sample-issue"
		s.gitProvider.CommitsToPush["feature/GH-99-sample-issue"] = []string{}
		s.pullRequestProvider.PullRequests["feature/GH-99-sample-issue"] = nil

		s.expectNoPrFound()

		s.uc.Cfg.CloseIssue = false

		err := s.uc.Execute()

		s.NoError(err)
	})

	s.Run("should error if could not get issue", func() {
		s.gitProvider.RemoteBranches = []string{"main", "develop"}
		s.gitProvider.LocalBranches = append(s.gitProvider.LocalBranches, "feature/GH-99-sample-issue")
		s.gitProvider.CurrentBranch = "feature/GH-99-sample-issue"
		s.gitProvider.CommitsToPush["feature/GH-99-sample-issue"] = []string{}

		mocks.UnsetExpectedCall(&s.issueTracker.Mock, s.issueTracker.GetIssue)
		s.issueTracker.EXPECT().GetIssue(mock.Anything).Return(domain.Issue{}, assert.AnError).Once()

		err := s.uc.Execute()

		s.Error(err)
		s.issueTracker.AssertExpectations(s.T())
		s.assertCreatePullRequestNotCalled()
	})
}

func (s *CreateGithubPullRequestExecutionTestSuite) assertCreatePullRequestNotCalled() {
	branch := s.gitProvider.CurrentBranch
	if pr := s.pullRequestProvider.PullRequests[branch]; pr != nil {
		s.Failf("pull request exists for branch %s", branch)
	}
}

func (s *CreateGithubPullRequestExecutionTestSuite) assertCreatePullRequestCalled() {
	branch := s.gitProvider.CurrentBranch
	if _, ok := s.pullRequestProvider.PullRequests[branch]; !ok {
		s.Failf("pull request not found for branch %s", branch)
	}
}

func (s *CreateGithubPullRequestExecutionTestSuite) expectNoPrFound() {
	branch := s.gitProvider.CurrentBranch
	if pr := s.pullRequestProvider.PullRequests[branch]; pr != nil {
		s.Failf("pull request exists for branch %s", branch)
	}
}

func (s *CreateGithubPullRequestExecutionTestSuite) initializeGitProvider() *domainMocks.MockGitProvider {
	gitProvider := &domainMocks.MockGitProvider{}

	gitProvider.EXPECT().GetCurrentBranch().Return(s.defaultBranchName, nil).Maybe()
	gitProvider.EXPECT().GetCommitsToPush(s.defaultBranchName).Return([]string{}, nil).Maybe()
	gitProvider.EXPECT().RemoteBranchExists(s.defaultBranchName).Return(true).Maybe()
	gitProvider.EXPECT().BranchExistsContains("/GH-1-").Return("feature/GH-1-sample-issue", true).Maybe()
	gitProvider.EXPECT().BranchExists("/GH-1-").Return(true).Maybe()

	gitProvider.EXPECT().CommitEmpty(mock.Anything).Return(nil).Maybe()
	gitProvider.EXPECT().PushBranch(mock.Anything).Return(nil).Maybe()
	gitProvider.EXPECT().FetchBranchFromOrigin("main").Return(nil).Maybe()
	gitProvider.EXPECT().CheckoutNewBranchFromOrigin("feature/GH-1-sample-issue", "main").Return(nil).Maybe()

	return gitProvider
}

func (s *CreateGithubPullRequestExecutionTestSuite) initializeUserInteractionProvider() *domainMocks.MockUserInteractionProvider {
	userInteractionProvider := &domainMocks.MockUserInteractionProvider{}

	userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInputPrompt("Label 'kind/feature' found. What type of branch name do you want to create?", []string{"feature", "other"}, mock.Anything, true).Return(nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInput(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	return userInteractionProvider
}

func (s *CreateGithubPullRequestExecutionTestSuite) initializePullRequestProvider() *domainMocks.MockPullRequestProvider {
	pullRequestProvider := &domainMocks.MockPullRequestProvider{}

	pullRequestProvider.EXPECT().GetPullRequestForBranch(s.defaultBranchName).Return(&domain.PullRequest{
		Title:  "Sample issue",
		Number: 1,
		State:  "OPEN",
		Closed: false,
		Url:    "https://example.com",
		Labels: []domain.Label{},
	}, nil).Maybe()

	pullRequestProvider.EXPECT().CreatePullRequest(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("https://example.com", nil).Maybe()

	return pullRequestProvider
}

func (s *CreateGithubPullRequestExecutionTestSuite) initializeIssueTrackerProvider() *domainMocks.MockIssueTrackerProvider {
	issueTrackerProvider := &domainMocks.MockIssueTrackerProvider{}

	// issueTrackerProvider.EXPECT().GetIssueTracker(mock.Anything).Return(GetDefaultIssueTracker(), nil).Maybe()
	issueTrackerProvider.EXPECT().ParseIssueId(mock.Anything).Return("1").Maybe()

	return issueTrackerProvider
}

func (s *CreateGithubPullRequestExecutionTestSuite) initializeRepositoryProvider() *domainMocks.MockRepositoryProvider {
	repositoryProvider := &domainMocks.MockRepositoryProvider{}

	repositoryProvider.EXPECT().GetRepository().Return(&domain.Repository{
		Owner:            "inditex",
		Name:             "gh-sherpa",
		NameWithOwner:    "InditexTech/gh-sherpa",
		DefaultBranchRef: "main",
	}, nil).Maybe()

	return repositoryProvider
}

func (s *CreateGithubPullRequestExecutionTestSuite) initializeIssueTracker() *domainMocks.MockIssueTracker {
	issueTracker := &domainMocks.MockIssueTracker{}

	issueTracker.EXPECT().FormatIssueId(mock.Anything).Return("GH-1").Maybe()
	issueTracker.EXPECT().GetIssue(mock.Anything).Return(domain.Issue{
		ID:           "1",
		Title:        "Sample issue",
		Body:         "Sample issue body",
		Labels:       []domain.Label{},
		IssueTracker: domain.IssueTrackerTypeGithub,
		Url:          "https://github.com/InditexTech/gh-sherpa/issues/1",
	}, nil).Maybe()
	issueTracker.EXPECT().GetIssueType(mock.Anything).Return(issue_types.Feature).Maybe()
	issueTracker.EXPECT().GetIssueTrackerType().Return(domain.IssueTrackerTypeGithub).Maybe()
	issueTracker.EXPECT().GetIssueTypeLabel(mock.Anything).Return("kind/feature").Maybe()

	return issueTracker
}

func (s *CreateGithubPullRequestExecutionTestSuite) initializeBranchProvider() *domainMocks.MockBranchProvider {
	branchProvider := &domainMocks.MockBranchProvider{}

	branchProvider.EXPECT().GetBranchName(mock.Anything, mock.Anything, mock.Anything).Return("feature/GH-1-sample-issue", nil).Maybe()

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
	s.gitProvider = s.initializeGitProvider()
	s.issueTrackerProvider = s.initializeIssueTrackerProvider()
	s.userInteractionProvider = s.initializeUserInteractionProvider()
	s.pullRequestProvider = s.initializePullRequestProvider()
	s.issueTracker = s.initializeIssueTracker()
	s.branchProvider = s.initializeBranchProvider()
	s.repositoryProvider = s.initializeRepositoryProvider()

	mocks.UnsetExpectedCall(&s.issueTrackerProvider.Mock, s.issueTrackerProvider.GetIssueTracker)
	s.issueTrackerProvider.EXPECT().GetIssueTracker(mock.Anything).Return(s.issueTracker, nil).Maybe()

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
		mocks.UnsetExpectedCall(&s.repositoryProvider.Mock, s.repositoryProvider.GetRepository)
		s.repositoryProvider.EXPECT().GetRepository().Return(nil, assert.AnError).Once()

		s.expectCreatePullRequestNotCalled()

		err := s.uc.Execute()

		s.Error(err)
		s.gitProvider.AssertExpectations(s.T())
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should error if could not get current branch", func() {
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.GetCurrentBranch)
		s.gitProvider.EXPECT().GetCurrentBranch().Return("", assert.AnError).Once()

		s.expectCreatePullRequestNotCalled()

		err := s.uc.Execute()

		s.ErrorContains(err, "could not get the current branch name")
		s.gitProvider.AssertExpectations(s.T())
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should error if no issue could be identified", func() {
		branchName := "branch-with-no-issue-name"

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.GetCurrentBranch)
		s.gitProvider.EXPECT().GetCurrentBranch().Return(branchName, nil).Once()

		s.expectCreatePullRequestNotCalled()

		err := s.uc.Execute()

		expectedError := fmt.Sprintf("could not find an issue identifier in the current branch named %s", branchName)
		s.EqualError(err, expectedError)
		s.gitProvider.AssertExpectations(s.T())
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should exit if user does not confirm current branch", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation(mock.Anything, mock.Anything).Return(false, nil).Once()

		s.expectCreatePullRequestNotCalled()

		err := s.uc.Execute()

		s.NoError(err)
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should return error if pr already exists", func() {
		s.expectCreatePullRequestNotCalled()

		err := s.uc.Execute()

		s.ErrorContains(err, use_cases.ErrRemoteBranchAlreadyExists(s.defaultBranchName).Error())
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should not ask the user for branch confirmation if default flag is used", func() {
		s.uc.Cfg.IsInteractive = false

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.RemoteBranchExists)
		s.gitProvider.EXPECT().RemoteBranchExists(s.defaultBranchName).Return(false).Times(2)

		mocks.UnsetExpectedCall(&s.pullRequestProvider.Mock, s.pullRequestProvider.GetPullRequestForBranch)
		s.pullRequestProvider.EXPECT().GetPullRequestForBranch(mock.Anything).Return(nil, nil).Once()

		err := s.uc.Execute()

		s.NoError(err)
		s.gitProvider.AssertExpectations(s.T())
		s.userInteractionProvider.AssertNotCalled(s.T(), "AskUserForConfirmation")
	})

	s.Run("should return error if remote branch already exists", func() {
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.RemoteBranchExists)
		s.gitProvider.EXPECT().RemoteBranchExists(s.defaultBranchName).Return(true).Once()

		s.expectCreatePullRequestNotCalled()

		err := s.uc.Execute()

		s.ErrorContains(err, use_cases.ErrRemoteBranchAlreadyExists(s.defaultBranchName).Error())
		s.gitProvider.AssertExpectations(s.T())
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should create and push empty commit if remote branch does not exists", func() {
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.RemoteBranchExists)
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.CommitEmpty)
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.PushBranch)
		s.gitProvider.EXPECT().RemoteBranchExists(s.defaultBranchName).Return(false).Times(2)
		s.gitProvider.EXPECT().CommitEmpty(mock.Anything).Return(nil).Once()
		s.gitProvider.EXPECT().PushBranch(s.defaultBranchName).Return(nil).Once()

		mocks.UnsetExpectedCall(&s.pullRequestProvider.Mock, s.pullRequestProvider.GetPullRequestForBranch)
		s.pullRequestProvider.EXPECT().GetPullRequestForBranch(mock.Anything).Return(nil, nil).Once()

		err := s.uc.Execute()

		s.NoError(err)
		s.gitProvider.AssertExpectations(s.T())
	})

	s.Run("should exit if user does not confirm the commit push when default flag is not used", func() {
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.GetCommitsToPush)
		s.gitProvider.EXPECT().GetCommitsToPush(s.defaultBranchName).Return([]string{"chore: update docs", "refactor: method"}, nil).Once()

		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue pushing all pending commits in this branch and create the pull request", true).Return(false, nil).Once()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.PushBranch)
		s.gitProvider.EXPECT().PushBranch(s.defaultBranchName).Once()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.RemoteBranchExists)
		s.gitProvider.EXPECT().RemoteBranchExists(s.defaultBranchName).Return(false).Once()

		s.expectCreatePullRequestNotCalled()

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should error if could not create empty commit", func() {
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.RemoteBranchExists)
		s.gitProvider.EXPECT().RemoteBranchExists(s.defaultBranchName).Return(false).Times(2)
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.CommitEmpty)
		s.gitProvider.EXPECT().CommitEmpty(mock.Anything).Return(assert.AnError).Once()
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.PushBranch)

		err := s.uc.Execute()

		s.ErrorContains(err, "could not do the empty commit because")
		s.gitProvider.AssertNotCalled(s.T(), "PushBranch")
		s.gitProvider.AssertExpectations(s.T())
	})

	s.Run("should error if could not push branch", func() {
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.RemoteBranchExists)
		s.gitProvider.EXPECT().RemoteBranchExists(s.defaultBranchName).Return(false).Times(2)
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.PushBranch)
		s.gitProvider.EXPECT().PushBranch(s.defaultBranchName).Return(assert.AnError).Once()

		err := s.uc.Execute()

		s.ErrorContains(err, "could not create the remote branch because")
		s.gitProvider.AssertExpectations(s.T())
	})

	s.Run("should error if could not create pull request", func() {
		mocks.UnsetExpectedCall(&s.pullRequestProvider.Mock, s.pullRequestProvider.GetPullRequestForBranch)
		s.pullRequestProvider.EXPECT().GetPullRequestForBranch(mock.Anything).Return(nil, nil).Once()

		mocks.UnsetExpectedCall(&s.pullRequestProvider.Mock, s.pullRequestProvider.CreatePullRequest)
		s.pullRequestProvider.EXPECT().CreatePullRequest(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", assert.AnError).Once()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.RemoteBranchExists)
		s.gitProvider.EXPECT().RemoteBranchExists(mock.Anything).Return(false).Times(2)

		err := s.uc.Execute()

		s.ErrorContains(err, "could not create the pull request because")
		s.pullRequestProvider.AssertExpectations(s.T())
	})

	s.Run("should checkout local branch if branch exists and user confirms branch usage without default flag and issue flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Once()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.CheckoutBranch)
		s.gitProvider.EXPECT().CheckoutBranch("feature/PROJECTKEY-1-sample-issue").Return(nil).Once()

		mocks.UnsetExpectedCall(&s.pullRequestProvider.Mock, s.pullRequestProvider.GetPullRequestForBranch)
		s.pullRequestProvider.EXPECT().GetPullRequestForBranch(mock.Anything).Return(nil, nil).Once()

		s.uc.Cfg.IssueID = "PROJECTKEY-1"

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
	})

	s.Run("should error if branch already exists when using default and issue flags", func() {
		s.uc.Cfg.IssueID = "1"
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.ErrorContains(err, "the branch feature/PROJECTKEY-1-sample-issue already exists")
	})

	s.Run("should return error if remote branch already exists when using issue flags", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(false, nil).Once()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.RemoteBranchExists)
		s.gitProvider.EXPECT().RemoteBranchExists("feature/PROJECTKEY-1-sample-issue").Return(true).Once()

		s.uc.Cfg.IssueID = "PROJECTKEY-1"

		err := s.uc.Execute()

		s.ErrorContains(err, use_cases.ErrRemoteBranchAlreadyExists("feature/PROJECTKEY-1-sample-issue").Error())
		s.userInteractionProvider.AssertExpectations(s.T())
		s.gitProvider.AssertExpectations(s.T())
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should create new branch name if user doesn't confirm default branch name when using issue flags", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(false, nil).Once()
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Once()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.FetchBranchFromOrigin)
		s.gitProvider.EXPECT().FetchBranchFromOrigin("main").Return(nil).Once()
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.CheckoutNewBranchFromOrigin)
		s.gitProvider.EXPECT().CheckoutNewBranchFromOrigin("feature/PROJECTKEY-1-sample-issue", "main").Return(nil).Once()
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.RemoteBranchExists)
		s.gitProvider.EXPECT().RemoteBranchExists("feature/PROJECTKEY-1-sample-issue").Return(false).Once()

		mocks.UnsetExpectedCall(&s.pullRequestProvider.Mock, s.pullRequestProvider.GetPullRequestForBranch)
		s.pullRequestProvider.EXPECT().GetPullRequestForBranch(mock.Anything).Return(nil, nil).Once()

		s.uc.Cfg.IssueID = "PROJECTKEY-1"

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
		s.gitProvider.AssertExpectations(s.T())
	})

	s.Run("should abort execution if user doesn't confirm branch name when using issue flags", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(false, nil).Once()
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(false, nil).Once()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.RemoteBranchExists)
		s.gitProvider.EXPECT().RemoteBranchExists("feature/PROJECTKEY-1-sample-issue").Return(false).Once()

		s.expectCreatePullRequestNotCalled()

		s.uc.Cfg.IssueID = "PROJECTKEY-1"

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
		s.gitProvider.AssertExpectations(s.T())
		s.assertCreatePullRequestNotCalled()
	})

	s.Run("should checkout branch if user confirms branch usage with issue flag and no default flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Once()

		s.gitProvider.EXPECT().CheckoutBranch("feature/PROJECTKEY-1-sample-issue").Return(nil).Once()

		mocks.UnsetExpectedCall(&s.pullRequestProvider.Mock, s.pullRequestProvider.GetPullRequestForBranch)
		s.pullRequestProvider.EXPECT().GetPullRequestForBranch(mock.Anything).Return(nil, nil).Once()

		s.uc.Cfg.IssueID = "PROJECTKEY-1"

		err := s.uc.Execute()

		s.NoError(err)
		s.userInteractionProvider.AssertExpectations(s.T())
	})

	s.Run("should not error if pull request is created correctly", func() {
		mocks.UnsetExpectedCall(&s.pullRequestProvider.Mock, s.pullRequestProvider.GetPullRequestForBranch)
		s.pullRequestProvider.EXPECT().GetPullRequestForBranch(mock.Anything).Return(nil, nil).Once()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.RemoteBranchExists)
		s.gitProvider.EXPECT().RemoteBranchExists(mock.Anything).Return(false).Times(2)

		err := s.uc.Execute()

		s.NoError(err)
		s.pullRequestProvider.AssertExpectations(s.T())
		s.assertCreatePullRequestCalled()
	})

	s.Run("should create branch and pull request if local branch doesn't exists with issue flag", func() {
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.BranchExistsContains)
		s.gitProvider.EXPECT().BranchExistsContains("/PROJECTKEY-1-").Return("", false)
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.RemoteBranchExists)
		s.gitProvider.EXPECT().RemoteBranchExists("feature/PROJECTKEY-1-sample-issue").Return(false).Once()

		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.SelectOrInputPrompt)
		s.userInteractionProvider.EXPECT().SelectOrInputPrompt("Issue type 'feature' found. What type of branch name do you want to create?", []string{"feature", "other"}, mock.Anything, true).Return(nil).Once()
		s.userInteractionProvider.EXPECT().SelectOrInput("additional description (optional). Truncate to 21 chars", []string{}, mock.Anything, false).Return(nil).Once()
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Once()
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Once()

		mocks.UnsetExpectedCall(&s.pullRequestProvider.Mock, s.pullRequestProvider.GetPullRequestForBranch)
		s.pullRequestProvider.EXPECT().GetPullRequestForBranch(mock.Anything).Return(nil, nil).Once()

		s.uc.Cfg.IssueID = "PROJECTKEY-1"

		err := s.uc.Execute()

		s.NoError(err)
		s.pullRequestProvider.AssertExpectations(s.T())
		s.assertCreatePullRequestCalled()
	})

	s.Run("should create pull request with no close issue flag", func() {
		mocks.UnsetExpectedCall(&s.pullRequestProvider.Mock, s.pullRequestProvider.CreatePullRequest)
		s.pullRequestProvider.EXPECT().CreatePullRequest("[PROJECTKEY-1] Sample issue", "Relates to [PROJECTKEY-1](https://jira.example.com/browse/PROJECTKEY-1)", "main", "feature/PROJECTKEY-1-sample-issue", true, []string{"kind/feature"}).Return("https://example.com", nil).Once()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.RemoteBranchExists)
		s.gitProvider.EXPECT().RemoteBranchExists(s.defaultBranchName).Return(false).Times(2)

		s.expectNoPrFound()

		s.uc.Cfg.CloseIssue = false

		err := s.uc.Execute()

		s.NoError(err)
		s.pullRequestProvider.AssertExpectations(s.T())
	})

	s.Run("should error if could not get issue", func() {
		mocks.UnsetExpectedCall(&s.pullRequestProvider.Mock, s.pullRequestProvider.GetPullRequestForBranch)
		s.pullRequestProvider.EXPECT().GetPullRequestForBranch(mock.Anything).Return(nil, nil).Once()

		mocks.UnsetExpectedCall(&s.issueTracker.Mock, s.issueTracker.GetIssue)
		s.issueTracker.EXPECT().GetIssue(mock.Anything).Return(domain.Issue{}, assert.AnError).Once()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.RemoteBranchExists)
		s.gitProvider.EXPECT().RemoteBranchExists(s.defaultBranchName).Return(false).Times(2)

		s.expectCreatePullRequestNotCalled()

		err := s.uc.Execute()

		s.Error(err)
		s.issueTracker.AssertExpectations(s.T())
		s.assertCreatePullRequestNotCalled()
	})
}

func (s *CreateJiraPullRequestExecutionTestSuite) expectCreatePullRequestNotCalled() {
	mocks.UnsetExpectedCall(&s.pullRequestProvider.Mock, s.pullRequestProvider.CreatePullRequest)
	s.pullRequestProvider.EXPECT().CreatePullRequest(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(0)
}

func (s *CreateJiraPullRequestExecutionTestSuite) assertCreatePullRequestNotCalled() {
	s.pullRequestProvider.AssertNotCalled(s.T(), "CreatePullRequest")
}

func (s *CreateJiraPullRequestExecutionTestSuite) assertCreatePullRequestCalled() {
	s.pullRequestProvider.AssertCalled(s.T(), "CreatePullRequest", "[PROJECTKEY-1] Sample issue", "Relates to [PROJECTKEY-1](https://jira.example.com/browse/PROJECTKEY-1)", "main", "feature/PROJECTKEY-1-sample-issue", true, []string{"kind/feature"})
}

func (s *CreateJiraPullRequestExecutionTestSuite) expectNoPrFound() {
	mocks.UnsetExpectedCall(&s.pullRequestProvider.Mock, s.pullRequestProvider.GetPullRequestForBranch)
	s.pullRequestProvider.EXPECT().GetPullRequestForBranch(mock.Anything).Return(nil, nil).Once()
}

func (s *CreateJiraPullRequestExecutionTestSuite) initializeGitProvider() *domainMocks.MockGitProvider {
	gitProvider := &domainMocks.MockGitProvider{}

	gitProvider.EXPECT().GetCurrentBranch().Return(s.defaultBranchName, nil).Maybe()
	gitProvider.EXPECT().GetCommitsToPush(s.defaultBranchName).Return([]string{}, nil).Maybe()
	gitProvider.EXPECT().RemoteBranchExists(s.defaultBranchName).Return(true).Maybe()
	gitProvider.EXPECT().BranchExistsContains("/PROJECTKEY-1-").Return("feature/PROJECTKEY-1-sample-issue", true).Maybe()
	gitProvider.EXPECT().BranchExists("/PROJECTKEY-1-").Return(true).Maybe()

	gitProvider.EXPECT().CommitEmpty(mock.Anything).Return(nil).Maybe()
	gitProvider.EXPECT().PushBranch(mock.Anything).Return(nil).Maybe()
	gitProvider.EXPECT().FetchBranchFromOrigin("main").Return(nil).Maybe()
	gitProvider.EXPECT().CheckoutNewBranchFromOrigin("feature/PROJECTKEY-1-sample-issue", "main").Return(nil).Maybe()

	return gitProvider
}

func (s *CreateJiraPullRequestExecutionTestSuite) initializeUserInteractionProvider() *domainMocks.MockUserInteractionProvider {
	userInteractionProvider := &domainMocks.MockUserInteractionProvider{}

	userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInputPrompt("Issue type 'feature' found. What type of branch name do you want to create?", []string{"feature", "other"}, mock.Anything, true).Return(nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInput(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	return userInteractionProvider
}

func (s *CreateJiraPullRequestExecutionTestSuite) initializePullRequestProvider() *domainMocks.MockPullRequestProvider {
	pullRequestProvider := &domainMocks.MockPullRequestProvider{}

	pullRequestProvider.EXPECT().GetPullRequestForBranch(s.defaultBranchName).Return(&domain.PullRequest{
		Title:  "Sample issue",
		Number: 1,
		State:  "OPEN",
		Closed: false,
		Url:    "https://example.com",
		Labels: []domain.Label{},
	}, nil).Maybe()

	pullRequestProvider.EXPECT().CreatePullRequest(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("https://example.com", nil).Maybe()

	return pullRequestProvider
}

func (s *CreateJiraPullRequestExecutionTestSuite) initializeIssueTrackerProvider() *domainMocks.MockIssueTrackerProvider {
	issueTrackerProvider := &domainMocks.MockIssueTrackerProvider{}

	// issueTrackerProvider.EXPECT().GetIssueTracker(mock.Anything).Return(GetDefaultIssueTracker(), nil).Maybe()
	issueTrackerProvider.EXPECT().ParseIssueId(mock.Anything).Return("1").Maybe()

	return issueTrackerProvider
}

func (s *CreateJiraPullRequestExecutionTestSuite) initializeRepositoryProvider() *domainMocks.MockRepositoryProvider {
	repositoryProvider := &domainMocks.MockRepositoryProvider{}

	repositoryProvider.EXPECT().GetRepository().Return(&domain.Repository{
		Owner:            "inditex",
		Name:             "gh-sherpa",
		NameWithOwner:    "InditexTech/gh-sherpa",
		DefaultBranchRef: "main",
	}, nil).Maybe()

	return repositoryProvider
}

func (s *CreateJiraPullRequestExecutionTestSuite) initializeIssueTracker() *domainMocks.MockIssueTracker {
	issueTracker := &domainMocks.MockIssueTracker{}

	issueTracker.EXPECT().FormatIssueId(mock.Anything).Return("PROJECTKEY-1").Maybe()
	issueTracker.EXPECT().GetIssue(mock.Anything).Return(domain.Issue{
		ID:           "PROJECTKEY-1",
		Title:        "Sample issue",
		Body:         "Sample issue body",
		Labels:       []domain.Label{},
		IssueTracker: domain.IssueTrackerTypeJira,
		Type: domain.IssueType{
			Id:          "3",
			Name:        "feature",
			Description: "A new feature of the product, which has to be developed and tested.",
		},
		Url: "https://jira.example.com/browse/PROJECTKEY-1",
	}, nil).Maybe()
	issueTracker.EXPECT().GetIssueType(mock.Anything).Return(issue_types.Feature).Maybe()
	issueTracker.EXPECT().GetIssueTrackerType().Return(domain.IssueTrackerTypeJira).Maybe()
	issueTracker.EXPECT().GetIssueTypeLabel(mock.Anything).Return("kind/feature").Maybe()

	return issueTracker
}

func (s *CreateJiraPullRequestExecutionTestSuite) initializeBranchProvider() *domainMocks.MockBranchProvider {
	branchProvider := &domainMocks.MockBranchProvider{}

	branchProvider.EXPECT().GetBranchName(mock.Anything, mock.Anything, mock.Anything).Return("feature/PROJECTKEY-1-sample-issue", nil).Maybe()

	return branchProvider
}
