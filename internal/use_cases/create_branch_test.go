package use_cases_test

import (
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
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
	gitProvider             *domainMocks.MockGitProvider
	issueTrackerProvider    *domainMocks.MockIssueTrackerProvider
	issueTracker            *domainMocks.MockIssueTracker
	userInteractionProvider *domainMocks.MockUserInteractionProvider
	branchProvider          *domainMocks.MockBranchProvider
	repositoryProvider      *domainMocks.MockRepositoryProvider
}

func (s *CreateBranchExecutionTestSuite) expectCreateBranchNotCalled() {
	mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.CheckoutNewBranchFromOrigin)
	// s.gitProvider.EXPECT().CheckoutNewBranchFromOrigin(mock.Anything, mock.Anything).Times(0)
}

func (s *CreateBranchExecutionTestSuite) assertCreateBranchNotCalled() {
	s.gitProvider.AssertNotCalled(s.T(), "CheckoutNewBranchFromOrigin")
}

func (s *CreateBranchExecutionTestSuite) assertCreateBranchCalled(issueTrackerType domain.IssueTrackerType) {
	var branchName string
	switch issueTrackerType {
	case domain.IssueTrackerTypeJira:
		branchName = "feature/PROJECTKEY-1-sample-issue"
	case domain.IssueTrackerTypeGithub:
		fallthrough
	default:
		branchName = "feature/GH-1-sample-issue"
	}
	s.gitProvider.AssertCalled(s.T(), "CheckoutNewBranchFromOrigin", branchName, "main")
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
	s.gitProvider = s.initializeGitProvider()
	s.issueTrackerProvider = s.initializeIssueTrackerProvider()
	s.issueTracker = s.initializeIssueTracker()
	s.userInteractionProvider = s.initializeUserInteractionProvider()
	s.branchProvider = s.initializeBranchProvider()
	s.repositoryProvider = s.initializeRepositoryProvider()

	mocks.UnsetExpectedCall(&s.issueTrackerProvider.Mock, s.issueTrackerProvider.GetIssueTracker)
	s.issueTrackerProvider.EXPECT().GetIssueTracker(mock.Anything).Return(s.issueTracker, nil).Maybe()

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
		mocks.UnsetExpectedCall(&s.repositoryProvider.Mock, s.repositoryProvider.GetRepository)
		s.repositoryProvider.EXPECT().GetRepository().Return(nil, assert.AnError).Once()

		s.expectCreateBranchNotCalled()

		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.Error(err)
		s.gitProvider.AssertExpectations(s.T())
		s.assertCreateBranchNotCalled()
	})

	s.Run("should error if no issue flag is provided", func() {
		s.expectCreateBranchNotCalled()

		err := s.uc.Execute()

		s.ErrorContains(err, "sherpa needs an valid issue identifier")
		s.gitProvider.AssertExpectations(s.T())
		s.assertCreateBranchNotCalled()
	})

	s.Run("should error if branch already exists with default flag", func() {
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.BranchExists)
		s.gitProvider.EXPECT().BranchExists("feature/GH-1-sample-issue").Return(true).Maybe()

		s.expectCreateBranchNotCalled()

		s.uc.Cfg.IssueID = "1"
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.ErrorContains(err, "a local branch with the name feature/GH-1-sample-issue already exists")
		s.assertCreateBranchNotCalled()
	})

	s.Run("should create branch if branch doesn't exists with default flag", func() {

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.BranchExists)
		s.gitProvider.EXPECT().BranchExists("feature/GH-1-sample-issue").Return(false).Maybe()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.CheckoutNewBranchFromOrigin)
		s.gitProvider.EXPECT().CheckoutNewBranchFromOrigin("feature/GH-1-sample-issue", "main").Return(nil).Maybe()

		s.uc.Cfg.IssueID = "1"
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.NoError(err)
		s.gitProvider.AssertExpectations(s.T())
		s.assertCreateBranchCalled(domain.IssueTrackerTypeGithub)
	})

	s.Run("should create branch if not exists without default flag", func() {

		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Maybe()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.BranchExists)
		s.gitProvider.EXPECT().BranchExists("feature/GH-1-sample-issue").Return(false).Maybe()

		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.NoError(err)
		s.assertCreateBranchCalled(domain.IssueTrackerTypeGithub)
	})

	s.Run("should error if branch already exists without default flag", func() {

		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Maybe()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.BranchExists)
		s.gitProvider.EXPECT().BranchExists("feature/GH-1-sample-issue").Return(true).Maybe()

		s.expectCreateBranchNotCalled()

		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.ErrorContains(err, "a local branch with the name feature/GH-1-sample-issue already exists")
		s.assertCreateBranchNotCalled()
	})

}

func (s *CreateGithubBranchExecutionTestSuite) initializeGitProvider() *domainMocks.MockGitProvider {
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

func (s *CreateGithubBranchExecutionTestSuite) initializeUserInteractionProvider() *domainMocks.MockUserInteractionProvider {
	userInteractionProvider := &domainMocks.MockUserInteractionProvider{}

	userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInputPrompt("Label 'kind/feature' found. What type of branch name do you want to create?", []string{"feature", "other"}, mock.Anything, true).Return(nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInput(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	return userInteractionProvider
}

func (s *CreateGithubBranchExecutionTestSuite) initializeIssueTrackerProvider() *domainMocks.MockIssueTrackerProvider {
	issueTrackerProvider := &domainMocks.MockIssueTrackerProvider{}

	// issueTrackerProvider.EXPECT().GetIssueTracker(mock.Anything).Return(GetDefaultIssueTracker(), nil).Maybe()
	issueTrackerProvider.EXPECT().ParseIssueId(mock.Anything).Return("1").Maybe()

	return issueTrackerProvider
}

func (s *CreateGithubBranchExecutionTestSuite) initializeBranchProvider() *domainMocks.MockBranchProvider {
	branchProvider := &domainMocks.MockBranchProvider{}

	branchProvider.EXPECT().GetBranchName(mock.Anything, mock.Anything, mock.Anything).Return("feature/GH-1-sample-issue", nil).Maybe()

	return branchProvider
}

func (s *CreateGithubBranchExecutionTestSuite) initializeRepositoryProvider() *domainMocks.MockRepositoryProvider {
	repositoryProvider := &domainMocks.MockRepositoryProvider{}

	repositoryProvider.EXPECT().GetRepository().Return(&domain.Repository{
		Owner:            "inditex",
		Name:             "gh-sherpa",
		NameWithOwner:    "InditexTech/gh-sherpa",
		DefaultBranchRef: "main",
	}, nil).Maybe()

	return repositoryProvider
}

func (s *CreateGithubBranchExecutionTestSuite) initializeIssueTracker() *domainMocks.MockIssueTracker {
	issueTracker := &domainMocks.MockIssueTracker{}

	issueTracker.EXPECT().FormatIssueId(mock.Anything).Return("GH-1").Maybe()
	issueTracker.EXPECT().GetIssue(mock.Anything).Return(domain.Issue{
		ID:           "1",
		Title:        "Sample issue",
		Body:         "Sample issue body",
		Labels:       []domain.Label{},
		IssueTracker: domain.IssueTrackerTypeGithub,
	}, nil).Maybe()
	issueTracker.EXPECT().GetIssueType(mock.Anything).Return(issue_types.Feature).Maybe()
	issueTracker.EXPECT().GetIssueTrackerType().Return(domain.IssueTrackerTypeGithub).Maybe()

	return issueTracker
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
	s.gitProvider = s.initializeGitProvider()
	s.issueTrackerProvider = s.initializeIssueTrackerProvider()
	s.issueTracker = s.initializeIssueTracker()
	s.userInteractionProvider = s.initializeUserInteractionProvider()
	s.branchProvider = s.initializeBranchProvider()
	s.repositoryProvider = s.initializeRepositoryProvider()

	mocks.UnsetExpectedCall(&s.issueTrackerProvider.Mock, s.issueTrackerProvider.GetIssueTracker)
	s.issueTrackerProvider.EXPECT().GetIssueTracker(mock.Anything).Return(s.issueTracker, nil).Maybe()

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
		mocks.UnsetExpectedCall(&s.repositoryProvider.Mock, s.repositoryProvider.GetRepository)
		s.repositoryProvider.EXPECT().GetRepository().Return(nil, assert.AnError).Once()

		s.expectCreateBranchNotCalled()

		s.uc.Cfg.IssueID = issueID

		err := s.uc.Execute()

		s.Error(err)
		s.gitProvider.AssertExpectations(s.T())
		s.assertCreateBranchNotCalled()
	})

	s.Run("should error if no issue flag is provided", func() {
		s.expectCreateBranchNotCalled()

		err := s.uc.Execute()

		s.ErrorContains(err, "sherpa needs an valid issue identifier")
		s.gitProvider.AssertExpectations(s.T())
		s.assertCreateBranchNotCalled()
	})

	s.Run("should error if branch already exists with default flag", func() {
		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.BranchExists)
		s.gitProvider.EXPECT().BranchExists("feature/PROJECTKEY-1-sample-issue").Return(true).Maybe()

		s.expectCreateBranchNotCalled()

		s.uc.Cfg.IssueID = issueID
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.ErrorContains(err, "a local branch with the name feature/PROJECTKEY-1-sample-issue already exists")
		s.assertCreateBranchNotCalled()
	})

	s.Run("should create branch if branch doesn't exists with default flag", func() {

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.BranchExists)
		s.gitProvider.EXPECT().BranchExists("feature/PROJECTKEY-1-sample-issue").Return(false).Maybe()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.CheckoutNewBranchFromOrigin)
		s.gitProvider.EXPECT().CheckoutNewBranchFromOrigin("feature/PROJECTKEY-1-sample-issue", "main").Return(nil).Maybe()

		s.uc.Cfg.IssueID = issueID
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.NoError(err)
		s.gitProvider.AssertExpectations(s.T())
		s.assertCreateBranchCalled(domain.IssueTrackerTypeJira)
	})

	s.Run("should create branch if not exists without default flag", func() {

		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Maybe()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.BranchExists)
		s.gitProvider.EXPECT().BranchExists("feature/PROJECTKEY-1-sample-issue").Return(false).Maybe()

		s.uc.Cfg.IssueID = issueID

		err := s.uc.Execute()

		s.NoError(err)
		s.assertCreateBranchCalled(domain.IssueTrackerTypeJira)
	})

	s.Run("should error if branch already exists without default flag", func() {

		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Maybe()

		mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.BranchExists)
		s.gitProvider.EXPECT().BranchExists("feature/PROJECTKEY-1-sample-issue").Return(true).Maybe()

		s.expectCreateBranchNotCalled()

		s.uc.Cfg.IssueID = issueID

		err := s.uc.Execute()

		s.ErrorContains(err, "a local branch with the name feature/PROJECTKEY-1-sample-issue already exists")
		s.assertCreateBranchNotCalled()
	})

}

func (s *CreateJiraBranchExecutionTestSuite) initializeGitProvider() *domainMocks.MockGitProvider {
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

func (s *CreateJiraBranchExecutionTestSuite) initializeUserInteractionProvider() *domainMocks.MockUserInteractionProvider {
	userInteractionProvider := &domainMocks.MockUserInteractionProvider{}

	userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to use this branch to create the pull request", true).Return(true, nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInputPrompt("Issue type 'feature' found. What type of branch name do you want to create?", []string{"feature", "other"}, mock.Anything, true).Return(nil).Maybe()
	userInteractionProvider.EXPECT().SelectOrInput(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	return userInteractionProvider
}

func (s *CreateJiraBranchExecutionTestSuite) initializeIssueTrackerProvider() *domainMocks.MockIssueTrackerProvider {
	issueTrackerProvider := &domainMocks.MockIssueTrackerProvider{}

	// issueTrackerProvider.EXPECT().GetIssueTracker(mock.Anything).Return(GetDefaultIssueTracker(), nil).Maybe()
	issueTrackerProvider.EXPECT().ParseIssueId(mock.Anything).Return("1").Maybe()

	return issueTrackerProvider
}

func (s *CreateJiraBranchExecutionTestSuite) initializeRepositoryProvider() *domainMocks.MockRepositoryProvider {
	repositoryProvider := &domainMocks.MockRepositoryProvider{}

	repositoryProvider.EXPECT().GetRepository().Return(&domain.Repository{
		Owner:            "inditex",
		Name:             "gh-sherpa",
		NameWithOwner:    "InditexTech/gh-sherpa",
		DefaultBranchRef: "main",
	}, nil).Maybe()

	return repositoryProvider
}

func (s *CreateJiraBranchExecutionTestSuite) initializeIssueTracker() *domainMocks.MockIssueTracker {
	issueTracker := &domainMocks.MockIssueTracker{}

	issueTracker.EXPECT().FormatIssueId(mock.Anything).Return("PROJECTKEY-1").Maybe()
	issueTracker.EXPECT().GetIssue(mock.Anything).Return(domain.Issue{
		ID:           "1",
		Title:        "Sample issue",
		Body:         "Sample issue body",
		Labels:       []domain.Label{},
		IssueTracker: domain.IssueTrackerTypeJira,
		Type: domain.IssueType{
			Id:          "3",
			Name:        "feature",
			Description: "A new feature of the product, which has to be developed and tested.",
		},
	}, nil).Maybe()
	issueTracker.EXPECT().GetIssueType(mock.Anything).Return(issue_types.Feature).Maybe()
	issueTracker.EXPECT().GetIssueTrackerType().Return(domain.IssueTrackerTypeJira).Maybe()

	return issueTracker
}

func (s *CreateJiraBranchExecutionTestSuite) initializeBranchProvider() *domainMocks.MockBranchProvider {
	branchProvider := &domainMocks.MockBranchProvider{}

	branchProvider.EXPECT().GetBranchName(mock.Anything, mock.Anything, mock.Anything).Return("feature/PROJECTKEY-1-sample-issue", nil).Maybe()

	return branchProvider
}
