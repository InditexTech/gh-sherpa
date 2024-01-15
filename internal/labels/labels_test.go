package labels

import (
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	domainMocks "github.com/InditexTech/gh-sherpa/internal/mocks/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newGHIssue() domain.Issue {
	return domain.Issue{
		ID: "GH-1",
		Labels: []domain.Label{
			{
				Id:          "1",
				Name:        "kind/feat",
				Description: "Feature",
				Color:       "ffffff",
			},
		},
	}
}

func TestGetIssueTypeLabel(t *testing.T) {

	t.Run("Returns the label for the github issue type", func(t *testing.T) {

		labelName := "kind/feat"

		cfg := Configuration{
			IssueLabels: IssueLabelsMap{
				issue_types.Feature: {labelName},
			},
		}
		issue := newGHIssue()

		issueTracker := &domainMocks.MockIssueTracker{}
		issueTracker.EXPECT().GetIssueType(issue).Return(issue_types.Feature, nil)

		issueTrackerProvider := &domainMocks.MockIssueTrackerProvider{}
		issueTrackerProvider.EXPECT().GetIssueTracker("GH-1").Return(issueTracker, nil)

		provider, err := New(cfg, issueTrackerProvider)
		require.NoError(t, err)

		label, err := provider.GetIssueTypeLabel(issue)
		require.NoError(t, err)

		assert.Equal(t, labelName, label)
	})

	t.Run("Returns an empty label if the github issue label type could not be determined", func(t *testing.T) {

		cfg := Configuration{}
		issue := newGHIssue()

		issueTracker := &domainMocks.MockIssueTracker{}
		issueTracker.EXPECT().GetIssueType(issue).Return(issue_types.Unknown, nil)

		issueTrackerProvider := &domainMocks.MockIssueTrackerProvider{}
		issueTrackerProvider.EXPECT().GetIssueTracker("GH-1").Return(issueTracker, nil)

		provider, err := New(cfg, issueTrackerProvider)
		require.NoError(t, err)

		label, err := provider.GetIssueTypeLabel(issue)

		require.NoError(t, err)
		assert.Empty(t, label)
	})
}

func TestGetLabelFromBranchType(t *testing.T) {
	t.Run("Returns the label for the branch type", func(t *testing.T) {

		labelName := "kind/feat"

		cfg := Configuration{
			IssueLabels: IssueLabelsMap{
				issue_types.Feature: {labelName},
			},
		}

		issueTrackerProvider := &domainMocks.MockIssueTrackerProvider{}

		provider, err := New(cfg, issueTrackerProvider)
		require.NoError(t, err)

		label, err := provider.GetLabelFromBranchType("feature")
		require.NoError(t, err)

		assert.Equal(t, labelName, label)
	})

	t.Run("Returns an empty label if the branch type could not be determined", func(t *testing.T) {

		cfg := Configuration{}

		issueTrackerProvider := &domainMocks.MockIssueTrackerProvider{}

		provider, err := New(cfg, issueTrackerProvider)
		require.NoError(t, err)

		label, err := provider.GetLabelFromBranchType("unknown")

		require.NoError(t, err)
		assert.Empty(t, label)
	})
}

func TestNew(t *testing.T) {
	t.Run("Creates a labels provider from configuration", func(t *testing.T) {
		cfg := Configuration{}
		issueTrackerProvider := &domainMocks.MockIssueTrackerProvider{}
		provider, err := New(cfg, issueTrackerProvider)
		require.NoError(t, err)

		assert.NotNil(t, provider)
	})
}
