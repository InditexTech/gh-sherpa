package labels

import (
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newIssue() domain.Issue {
	return domain.Issue{}
}

func TestGetIssueTypeLabel(t *testing.T) {

	t.Run("Returns the label for the issue type", func(t *testing.T) {

		provider, err := New(Configuration{})
		require.NoError(t, err)

		issue := newIssue()
		label, err := provider.GetIssueTypeLabel(issue)
		require.NoError(t, err)

		assert.Equal(t, "kind/feature", label)
	})
}

func TestNew(t *testing.T) {
	t.Run("Creates a labels provider from configuration", func(t *testing.T) {
		cfg := Configuration{}
		provider, err := New(cfg)
		require.NoError(t, err)

		assert.NotNil(t, provider)
	})
}
