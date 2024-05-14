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
	gitProvider             *domainMocks.MockGitProvider
	issueTrackerProvider    *domainFakes.FakeIssueTrackerProvider
	issueTracker            *domainFakes.FakeIssueTracker
	userInteractionProvider *domainMocks.MockUserInteractionProvider
	branchProvider          *domainMocks.MockBranchProvider
	repositoryProvider      *domainFakes.FakeRepositoryProvider
}

func (*CreateBranchExecutionTestSuite) newFakeIssueTrackerProvider(issueTracker domain.IssueTracker) *domainFakes.FakeIssueTrackerProvider {
	return &domainFakes.FakeIssueTrackerProvider{
		IssueTracker: issueTracker,
	}
}

func (*CreateBranchExecutionTestSuite) newFakeIssueTracker() *domainFakes.FakeIssueTracker {
	return &domainFakes.FakeIssueTracker{
		Configurations: map[string]domainFakes.FakeIssueTrackerConfiguration{
			"1": {
				IssueTrackerIdentifier: domain.IssueTrackerTypeGithub.String(),
				IssueType:              issue_types.Feature,
				IssueTypeLabel:         "kind/feature",
				Issue: domain.Issue{
					ID:           "1",
					IssueTracker: domain.IssueTrackerTypeGithub,
					Title:        "Sample issue",
					Body:         "Sample issue body",
					Labels: []domain.Label{
						{
							Id:   "kind/feature",
							Name: "kind/feature",
						},
					},
					Url: "https://github.com/InditexTech/gh-sherpa-repo-test/issues/1",
				},
			},
			"PROJECTKEY-1": {
				IssueTrackerIdentifier: domain.IssueTrackerTypeJira.String(),
				IssueType:              issue_types.Feature,
				IssueTypeLabel:         "kind/feature",
				Issue: domain.Issue{
					ID:           "PROJECTKEY-1",
					IssueTracker: domain.IssueTrackerTypeJira,
					Title:        "Sample issue",
					Body:         "Sample issue body",
					Labels: []domain.Label{
						{
							Id:   "kind/feature",
							Name: "kind/feature",
						},
					},
					Url: "https://sample.jira.com/PROJECTKEY-1",
				},
			},
			"3": {
				IssueTrackerIdentifier: domain.IssueTrackerTypeGithub.String(),
				IssueType:              issue_types.Feature,
				IssueTypeLabel:         "kind/documentation",
				Issue: domain.Issue{
					ID:           "3",
					IssueTracker: domain.IssueTrackerTypeGithub,
					Title:        "Sample documentation issue",
					Body:         "Sample documentation issue body",
					Labels: []domain.Label{
						{
							Id:   "kind/documentation",
							Name: "kind/documentation",
						},
					},
					Url: "https://github.com/InditexTech/gh-sherpa-repo-test/issues/3",
				},
			},
			"PROJECTKEY-3": {
				IssueTrackerIdentifier: domain.IssueTrackerTypeJira.String(),
				IssueType:              issue_types.Feature,
				IssueTypeLabel:         "kind/documentation",
				Issue: domain.Issue{
					ID:           "PROJECTKEY-3",
					IssueTracker: domain.IssueTrackerTypeJira,
					Title:        "Sample documentation issue",
					Body:         "Sample documentation issue body",
					Labels: []domain.Label{
						{
							Id:   "kind/documentation",
							Name: "kind/documentation",
						},
					},
					Url: "https://sample.jira.com/PROJECTKEY-3",
				},
			},
			"6": {
				IssueTrackerIdentifier: domain.IssueTrackerTypeGithub.String(),
				IssueType:              issue_types.Feature,
				IssueTypeLabel:         "kind/refactoring",
				Issue: domain.Issue{
					ID:           "6",
					IssueTracker: domain.IssueTrackerTypeJira,
					Title:        "Sample refactoring issue",
					Body:         "Sample refactoring issue body",
					Labels: []domain.Label{
						{
							Id:   "kind/refactoring",
							Name: "kind/refactoring",
						},
					},
					Url: "https://github.com/InditexTech/gh-sherpa-repo-test/issues/6",
				},
			},
			"PROJECTKEY-6": {
				IssueTrackerIdentifier: domain.IssueTrackerTypeJira.String(),
				IssueType:              issue_types.Feature,
				IssueTypeLabel:         "kind/refactoring",
				Issue: domain.Issue{
					ID:           "PROJECTKEY-6",
					IssueTracker: domain.IssueTrackerTypeJira,
					Title:        "Sample refactoring issue",
					Body:         "Sample refactoring issue body",
					Labels: []domain.Label{
						{
							Id:   "kind/refactoring",
							Name: "kind/refactoring",
						},
					},
					Url: "https://sample.jira.com/PROJECTKEY-6",
				},
			},
		},
	}
}

func (*CreateBranchExecutionTestSuite) newFakeGitProvider() *domainFakes.FakeGitProvider {
	return &domainFakes.FakeGitProvider{
		CurrentBranch: "main",
		RemoteBranches: []string{
			"main",
			"develop",
		},
		LocalBranches: []string{
			"main",
			"develop",
			"feature/GH-3-local-branch",
			"feature/PROJECTKEY-3-local-branch",
		},
		CommitsToPush: map[string][]string{},
		BranchWithCommitError: map[string]error{
			"feature/GH-4-with-commit-error":         domainFakes.ErrGetCommitsToPush,
			"feature/PROJECTKEY-4-with-commit-error": domainFakes.ErrGetCommitsToPush,
		},
	}
}

func (*CreateBranchExecutionTestSuite) newFakeRepositoryProvider() *domainFakes.FakeRepositoryProvider {
	return &domainFakes.FakeRepositoryProvider{
		Repository: &domain.Repository{
			Name:             "gh-sherpa-test-repo",
			Owner:            "inditextech",
			NameWithOwner:    "inditextech/gh-sherpa-test-repo",
			DefaultBranchRef: "main",
		},
	}
}

func (s *CreateBranchExecutionTestSuite) setGetBranchName(branch string) {
	s.branchProvider.EXPECT().GetBranchName(mock.Anything, mock.Anything, mock.Anything).Return(branch, nil).Once()
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
	gitProvider *domainFakes.FakeGitProvider
}

func (s *CreateGithubBranchExecutionTestSuite) newFakeGitHubIssueTracker() *domainFakes.FakeIssueTracker {
	issueTracker := s.newFakeIssueTracker()
	issueTracker.IssueTrackerType = domain.IssueTrackerTypeGithub
	return issueTracker
}

func TestCreateGithubBranchExecutionTestSuite(t *testing.T) {
	suite.Run(t, new(CreateGithubBranchExecutionTestSuite))
}

func (s *CreateGithubBranchExecutionTestSuite) SetupSuite() {
	s.defaultBranchName = "feature/GH-1-sample-issue"
}

func (s *CreateGithubBranchExecutionTestSuite) SetupSubTest() {
	s.gitProvider = s.newFakeGitProvider()
	s.issueTracker = s.newFakeGitHubIssueTracker()
	s.issueTrackerProvider = s.newFakeIssueTrackerProvider(s.issueTracker)
	s.userInteractionProvider = s.initializeUserInteractionProvider()
	s.branchProvider = s.initializeBranchProvider()
	s.repositoryProvider = s.newFakeRepositoryProvider()

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
		s.Assert().False(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should error if no issue flag is provided", func() {
		err := s.uc.Execute()

		s.ErrorContains(err, "sherpa needs an valid issue identifier")
		s.Assert().False(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should error if branch already exists with default flag", func() {
		branchName := "feature/GH-3-local-branch"
		s.setGetBranchName(branchName)

		s.uc.Cfg.IssueID = "3"
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.ErrorContains(err, fmt.Sprintf("a local branch with the name %s already exists", branchName))
	})

	s.Run("should create branch if branch doesn't exists with default flag", func() {
		// mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.BranchExists)
		// s.gitProvider.EXPECT().BranchExists("feature/GH-1-sample-issue").Return(false).Maybe()

		// mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.CheckoutNewBranchFromOrigin)
		// s.gitProvider.EXPECT().CheckoutNewBranchFromOrigin("feature/GH-1-sample-issue", "main").Return(nil).Maybe()

		s.setGetBranchName(s.defaultBranchName)

		s.uc.Cfg.IssueID = "1"
		s.uc.Cfg.IsInteractive = false

		err := s.uc.Execute()

		s.NoError(err)
		s.Assert().True(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should create branch if not exists without default flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Maybe()

		// mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.BranchExists)
		// s.gitProvider.EXPECT().BranchExists("feature/GH-1-sample-issue").Return(false).Maybe()

		s.setGetBranchName(s.defaultBranchName)

		s.uc.Cfg.IssueID = "1"

		err := s.uc.Execute()

		s.NoError(err)
		s.Assert().True(s.gitProvider.BranchExists(s.defaultBranchName))
	})

	s.Run("should error if branch already exists without default flag", func() {
		mocks.UnsetExpectedCall(&s.userInteractionProvider.Mock, s.userInteractionProvider.AskUserForConfirmation)
		s.userInteractionProvider.EXPECT().AskUserForConfirmation("Do you want to continue?", true).Return(true, nil).Maybe()

		// mocks.UnsetExpectedCall(&s.gitProvider.Mock, s.gitProvider.BranchExists)
		// s.gitProvider.EXPECT().BranchExists("feature/GH-1-sample-issue").Return(true).Maybe()

		branchName := "feature/GH-3-local-branch"
		s.setGetBranchName(branchName)

		s.uc.Cfg.IssueID = "3"

		err := s.uc.Execute()

		s.ErrorContains(err, fmt.Sprintf("a local branch with the name %s already exists", branchName))
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

func (s *CreateGithubBranchExecutionTestSuite) initializeBranchProvider() *domainMocks.MockBranchProvider {
	branchProvider := &domainMocks.MockBranchProvider{}

	return branchProvider
}

type CreateJiraBranchExecutionTestSuite struct {
	CreateBranchExecutionTestSuite
}

func (s *CreateJiraBranchExecutionTestSuite) newFakeJiraIssueTracker() *domainFakes.FakeIssueTracker {
	issueTracker := s.newFakeIssueTracker()
	issueTracker.IssueTrackerType = domain.IssueTrackerTypeJira
	return issueTracker
}

func TestCreateJiraBranchExecutionTestSuite(t *testing.T) {
	suite.Run(t, new(CreateJiraBranchExecutionTestSuite))
}

func (s *CreateJiraBranchExecutionTestSuite) SetupSuite() {
	s.defaultBranchName = "feature/PROJECTKEY-1-sample-issue"
}

func (s *CreateJiraBranchExecutionTestSuite) SetupSubTest() {
	s.gitProvider = s.initializeGitProvider()
	s.issueTracker = s.newFakeJiraIssueTracker()
	s.issueTrackerProvider = s.newFakeIssueTrackerProvider(s.issueTracker)
	s.userInteractionProvider = s.initializeUserInteractionProvider()
	s.branchProvider = s.initializeBranchProvider()
	s.repositoryProvider = s.newFakeRepositoryProvider()

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
		s.expectCreateBranchNotCalled()

		s.setGetBranchName(s.defaultBranchName)

		s.repositoryProvider.Repository = nil

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

		s.setGetBranchName(s.defaultBranchName)

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

		s.setGetBranchName(s.defaultBranchName)

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

		s.setGetBranchName(s.defaultBranchName)

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

		s.setGetBranchName(s.defaultBranchName)

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

func (s *CreateJiraBranchExecutionTestSuite) initializeBranchProvider() *domainMocks.MockBranchProvider {
	branchProvider := &domainMocks.MockBranchProvider{}

	return branchProvider
}
