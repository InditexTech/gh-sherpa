package branches

import (
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	domainFakes "github.com/InditexTech/gh-sherpa/internal/fakes/domain"
	domainMocks "github.com/InditexTech/gh-sherpa/internal/mocks/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BranchTestSuite struct {
	suite.Suite
	b                       *BranchProvider
	userInteractionProvider *domainMocks.MockUserInteractionProvider
	fakeIssue               *domainFakes.FakeIssue
	defaultRepository       *domain.Repository
}

func TestBranchTestSuite(t *testing.T) {
	suite.Run(t, new(BranchTestSuite))
}

func (s *BranchTestSuite) SetupSubTest() {
	s.userInteractionProvider = &domainMocks.MockUserInteractionProvider{}

	s.b = &BranchProvider{
		cfg: Configuration{
			Branches: config.Branches{
				Prefixes: map[issue_types.IssueType]string{
					issue_types.Bug: "bugfix",
				},
				MaxLength: 0,
			},
		},
		UserInteraction: s.userInteractionProvider,
	}

	s.defaultRepository = &domain.Repository{
		Name:             "test-name",
		Owner:            "test-owner",
		NameWithOwner:    "test-owner/test-name",
		DefaultBranchRef: "main",
	}

	s.fakeIssue = domainFakes.NewFakeIssue("1", issue_types.Bug, domain.IssueTrackerTypeGithub)
}

func (s *BranchTestSuite) TestGetBranchName() {

	s.Run("should return expected branch name", func() {
		expectedBrachName := "bugfix/GH-1-fake-title"

		branchName, err := s.b.GetBranchName(s.fakeIssue, *s.defaultRepository)

		s.NoError(err)
		s.Equal(expectedBrachName, branchName)
	})

	s.Run("should return cropped branch", func() {
		expectedBrachName := "bugfix/GH-1-my-title-is-too-long-and-it-sho"

		s.fakeIssue.SetTitle("my title is too long and it should not matter")

		s.b.cfg.Branches.MaxLength = 63
		branchName, err := s.b.GetBranchName(s.fakeIssue, *s.defaultRepository)

		s.NoError(err)
		s.Equal(expectedBrachName, branchName)
	})

	s.Run("should return expected branch name when interactive", func() {
		expectedBrachName := "bugfix/GH-1-fake-title-from-interactive"

		s.userInteractionProvider.EXPECT().SelectOrInputPrompt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		s.userInteractionProvider.EXPECT().SelectOrInput(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(name string, validValues []string, variable *string, required bool) {
			*variable = "fake title from interactive"
		}).Return(nil).Once()

		s.b.cfg.IsInteractive = true

		branchName, err := s.b.GetBranchName(s.fakeIssue, *s.defaultRepository)

		s.NoError(err)
		s.Equal(expectedBrachName, branchName)
	})

	s.Run("should return cropped branch name when interactive", func() {
		expectedBrachName := "bugfix/GH-1-this-is-a-very-long-fake-title"

		s.userInteractionProvider.EXPECT().SelectOrInputPrompt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		s.userInteractionProvider.EXPECT().SelectOrInput(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(name string, validValues []string, variable *string, required bool) {
			*variable = "this is a very long fake title from interactive"
		}).Return(nil).Once()

		s.b.cfg.IsInteractive = true
		s.b.cfg.Branches.MaxLength = 63

		branchName, err := s.b.GetBranchName(s.fakeIssue, *s.defaultRepository)

		s.NoError(err)
		s.Equal(expectedBrachName, branchName)
	})
}

func TestParseIssueContext(t *testing.T) {
	tests := []struct {
		name  string
		given string
		want  string
	}{
		{
			name:  "issue context with special chars",
			given: "  /hello world//test/.test..test@{test\\\\test.lock  ",
			want:  "hello-world-test-test-testtesttest",
		},
		{
			name:  "issue context with special chars",
			given: "Begoña Caçadora_renombrádo' Él !cÓncepto $de \"bloqueo\" en cc%, úÍ",
			want:  "begona-cazadora_renombrado-el-concepto-de-bloqueo-en-cc-ui",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := normalizeBranch(tt.given)

			assert.Equal(t, tt.want, context)
		})
	}
}

func TestFormatBranchName(t *testing.T) {
	repositoryName := "InditexTech/gh-sherpa"

	type args struct {
		repository           string
		branchType           string
		issueId              string
		issueContext         string
		branchPrefixOverride map[issue_types.IssueType]string
		maxLength            int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Does format branch name",
			args: args{
				repository:   repositoryName,
				branchType:   "feature",
				issueId:      "GH-1",
				issueContext: "my-title",
				maxLength:    63,
			},

			want: "feature/GH-1-my-title",
		},
		{
			name: "Does format branch name with override",
			args: args{
				repository:           repositoryName,
				branchType:           "feature",
				issueId:              "GH-1",
				issueContext:         "my-title",
				branchPrefixOverride: map[issue_types.IssueType]string{issue_types.Feature: "feat"},
				maxLength:            63,
			},
			want: "feat/GH-1-my-title",
		},
		{
			name: "Does format long branch name",
			args: args{
				repository:   repositoryName,
				branchType:   "feature",
				issueId:      "GH-1",
				issueContext: "my-title-is-too-long-and-it-should-be-truncated",
				maxLength:    63,
			},
			want: "feature/GH-1-my-title-is-too-long-and-it-s",
		},
		{
			name: "Does format long branch name with override",
			args: args{
				repository:           repositoryName,
				branchType:           "feature",
				issueId:              "GH-1",
				issueContext:         "my-title-is-too-long-and-it-should-be-truncated",
				branchPrefixOverride: map[issue_types.IssueType]string{issue_types.Feature: "feat"},
				maxLength:            63,
			},
			want: "feat/GH-1-my-title-is-too-long-and-it-shou",
		},
		{
			name: "Does format branch with empty override",
			args: args{
				repository:           repositoryName,
				branchType:           "refactoring",
				issueId:              "GH-1",
				issueContext:         "refactor-issue",
				branchPrefixOverride: map[issue_types.IssueType]string{issue_types.Refactoring: ""},
				maxLength:            63,
			},
			want: "refactoring/GH-1-refactor-issue",
		},
		{
			name: "Does not crop the title if the maxLength is 0",
			args: args{
				repository:   repositoryName,
				branchType:   "feature",
				issueId:      "GH-1",
				issueContext: "this-is-a-very-long-title-that-will-not-be-cropped",
				maxLength:    0,
			},
			want: "feature/GH-1-this-is-a-very-long-title-that-will-not-be-cropped",
		},
		{
			name: "Does not crop the title if the maxLength is negative",
			args: args{
				repository:   repositoryName,
				branchType:   "feature",
				issueId:      "GH-1",
				issueContext: "this-is-a-very-long-title-that-will-not-be-cropped",
				maxLength:    -1,
			},
			want: "feature/GH-1-this-is-a-very-long-title-that-will-not-be-cropped",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			b := BranchProvider{
				cfg: Configuration{
					Branches: config.Branches{
						Prefixes:  tt.args.branchPrefixOverride,
						MaxLength: tt.args.maxLength,
					},
				},
			}
			branchName := b.formatBranchName(tt.args.repository, tt.args.branchType, tt.args.issueId, tt.args.issueContext)

			assert.Equal(t, tt.want, branchName)
		})
	}
}

func TestParseBranchName(t *testing.T) {
	for _, tc := range []struct {
		name       string
		branchName string
		want       *BranchNameInfo
	}{
		{
			branchName: "feature/GH-1-my-title",
			want:       &BranchNameInfo{BranchType: "feature", IssueId: "GH-1", IssueContext: "my-title"},
		},
		{
			branchName: "bugfix/PROJECTKEY-1-my-title",
			want:       &BranchNameInfo{BranchType: "bugfix", IssueId: "PROJECTKEY-1", IssueContext: "my-title"},
		},
		{
			branchName: "feature/GH-1-my-title-is-too-long-and-it-should-not-matter",
			want:       &BranchNameInfo{BranchType: "feature", IssueId: "GH-1", IssueContext: "my-title-is-too-long-and-it-should-not-matter"},
		},
		{
			branchName: "randomprefix/A_PROJECT_KEY-99-issue-tittle-here",
			want:       &BranchNameInfo{BranchType: "randomprefix", IssueId: "A_PROJECT_KEY-99", IssueContext: "issue-tittle-here"},
		},
	} {
		tc := tc
		t.Run(tc.branchName, func(t *testing.T) {
			branchInfo := ParseBranchName(tc.branchName)

			assert.Equal(t, tc.want, branchInfo)
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("Creates a branch provider from configuration", func(t *testing.T) {
		cfg := config.Configuration{}
		userInteraction := new(domainMocks.MockUserInteractionProvider)
		provider, err := NewFromConfiguration(cfg, userInteraction, false)
		require.NoError(t, err)

		assert.NotNil(t, provider)
	})
	//TODO: Add test cases when validation is implemented
}
