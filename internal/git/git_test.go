// SPDX-FileCopyrightText: 2023 INDITEX S.A
//
// SPDX-License-Identifier: Apache-2.0

package git

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitCheckoutNewBranchFromOrigin(t *testing.T) {
	provider := Provider{}

	t.Run("GitCheckoutNewBranchFromOrigin should checkout a new branch from origin", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			argsSent = args
			return
		}

		err := provider.CheckoutNewBranchFromOrigin("my-branch", "main")

		assert.NoError(t, err)
		assert.Equal(t, []string{"checkout", "--no-track", "-b", "my-branch", "origin/main"}, argsSent)
	})

	t.Run("GitCheckoutNewBranchFromOrigin should return an error if the branch is not found", func(t *testing.T) {
		runGitCommand = func(args ...string) (out string, err error) {
			err = fmt.Errorf("Failed to run Git command (%w)\n\nDetails:\n%s", err, "foo")
			return
		}
		err := provider.CheckoutNewBranchFromOrigin("my-branch", "main")

		assert.Error(t, err)
	})
}

func TestGitFetchBranchFromOrigin(t *testing.T) {
	provider := Provider{}
	t.Run("GitFetchBranchFromOrigin should fetch a branch from origin", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			argsSent = args
			return
		}

		err := provider.FetchBranchFromOrigin("my-branch")

		assert.NoError(t, err)
		assert.Equal(t, []string{"fetch", "origin", "my-branch"}, argsSent)
	})

	t.Run("GitFetchBranchFromOrigin should return an error if the branch is not found", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			argsSent = args
			err = fmt.Errorf("Failed to run Git command (%w)\n\nDetails:\n%s", err, "foo")
			return
		}

		err := provider.FetchBranchFromOrigin("my-branch")

		assert.Error(t, err)
		assert.Equal(t, []string{"fetch", "origin", "my-branch"}, argsSent)
	})
}

func TestGitCheckoutBranch(t *testing.T) {
	provider := Provider{}
	t.Run("GitCheckoutBranch should checkout a branch", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			argsSent = args
			return
		}

		err := provider.CheckoutBranch("my-branch")

		assert.NoError(t, err)
		assert.Equal(t, []string{"checkout", "my-branch"}, argsSent)
	})

	t.Run("GitCheckoutBranch should return an error if the branch is not found", func(t *testing.T) {
		runGitCommand = func(args ...string) (out string, err error) {
			err = fmt.Errorf("Failed to run Git command (%w)\n\nDetails:\n%s", err, "foo")
			return
		}
		err := provider.CheckoutBranch("foo")

		assert.Error(t, err)
	})
}

func TestGitBranchExists(t *testing.T) {
	provider := Provider{}
	t.Run("GitBranchExists should return true if the branch exists", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			argsSent = args
			return
		}

		exists := provider.BranchExists("my-branch")

		assert.True(t, exists)
		assert.Equal(t, []string{"show-ref", "--verify", "refs/heads/my-branch"}, argsSent)
	})

	t.Run("GitBranchExists should return false if the branch does not exist", func(t *testing.T) {
		runGitCommand = func(args ...string) (out string, err error) {
			err = fmt.Errorf("Failed to run Git command (%w)\n\nDetails:\n%s", err, "foo")
			return
		}

		exists := provider.BranchExists("foo")

		assert.False(t, exists)
	})
}

func TestGitPushBranch(t *testing.T) {
	provider := Provider{}
	t.Run("GitPush should push the branch to origin", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			argsSent = args
			return
		}

		err := provider.PushBranch("my-branch")

		assert.NoError(t, err)
		assert.Equal(t, []string{"push", "-u", "origin", "my-branch"}, argsSent)
	})

	t.Run("GitPush should return an error if the branch is not found", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			argsSent = args
			err = fmt.Errorf("Failed to run Git command (%w)\n\nDetails:\n%s", err, "foo")
			return
		}

		err := provider.PushBranch("my-branch")

		assert.Error(t, err)
		assert.Equal(t, []string{"push", "-u", "origin", "my-branch"}, argsSent)
	})
}
