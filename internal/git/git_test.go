package git

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitCheckoutNewBranchFromOrigin(t *testing.T) {
	provider := Provider{}

	t.Run("GitCheckoutNewBranchFromOrigin should checkout a new branch from origin when no upstream", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			// Mock git remote get-url upstream to return error (no upstream)
			if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "upstream" {
				return "", fmt.Errorf("no such remote")
			}
			// Capture the checkout command
			argsSent = args
			return "", nil
		}

		err := provider.CheckoutNewBranchFromOrigin("my-branch", "main")

		assert.NoError(t, err)
		assert.Equal(t, []string{"checkout", "--no-track", "-b", "my-branch", "origin/main"}, argsSent)
	})

	t.Run("GitCheckoutNewBranchFromOrigin should checkout a new branch from upstream when upstream exists", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			// Mock git remote get-url upstream to return success (upstream exists)
			if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "upstream" {
				return "https://github.com/upstream/repo.git", nil
			}
			// Capture the checkout command
			argsSent = args
			return "", nil
		}

		err := provider.CheckoutNewBranchFromOrigin("my-branch", "main")

		assert.NoError(t, err)
		assert.Equal(t, []string{"checkout", "--no-track", "-b", "my-branch", "upstream/main"}, argsSent)
	})

	t.Run("GitCheckoutNewBranchFromOrigin should return an error if the branch is not found", func(t *testing.T) {
		runGitCommand = func(args ...string) (out string, err error) {
			// Mock git remote get-url upstream to return error (no upstream)
			if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "upstream" {
				return "", fmt.Errorf("no such remote")
			}
			err = fmt.Errorf("Failed to run Git command (%w)\n\nDetails:\n%s", err, "foo")
			return
		}
		err := provider.CheckoutNewBranchFromOrigin("my-branch", "main")

		assert.Error(t, err)
	})
}

func TestGitFetchBranchFromOrigin(t *testing.T) {
	provider := Provider{}
	t.Run("GitFetchBranchFromOrigin should fetch a branch from origin when no upstream", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			// Mock git remote get-url upstream to return error (no upstream)
			if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "upstream" {
				return "", fmt.Errorf("no such remote")
			}
			// Capture the fetch command
			argsSent = args
			return "", nil
		}

		err := provider.FetchBranchFromOrigin("my-branch")

		assert.NoError(t, err)
		assert.Equal(t, []string{"fetch", "origin", "my-branch"}, argsSent)
	})

	t.Run("GitFetchBranchFromOrigin should fetch a branch from upstream when upstream exists", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			// Mock git remote get-url upstream to return success (upstream exists)
			if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "upstream" {
				return "https://github.com/upstream/repo.git", nil
			}
			// Capture the fetch command
			argsSent = args
			return "", nil
		}

		err := provider.FetchBranchFromOrigin("my-branch")

		assert.NoError(t, err)
		assert.Equal(t, []string{"fetch", "upstream", "my-branch"}, argsSent)
	})

	t.Run("GitFetchBranchFromOrigin should return an error if the branch is not found", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			// Mock git remote get-url upstream to return error (no upstream)
			if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "upstream" {
				return "", fmt.Errorf("no such remote")
			}
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

func TestGitCreateWorktree(t *testing.T) {
	provider := Provider{}

	t.Run("should create worktree from origin when no upstream", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			// Mock git remote get-url upstream to return error (no upstream)
			if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "upstream" {
				return "", fmt.Errorf("no such remote")
			}
			// Capture the worktree add command
			argsSent = args
			return "", nil
		}

		err := provider.CreateWorktree("../my-worktree", "my-branch", "main")

		assert.NoError(t, err)
		assert.Equal(t, []string{"worktree", "add", "../my-worktree", "-b", "my-branch", "origin/main"}, argsSent)
	})

	t.Run("should create worktree from upstream when upstream exists", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			// Mock git remote get-url upstream to return success (upstream exists)
			if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "upstream" {
				return "https://github.com/upstream/repo.git", nil
			}
			// Capture the worktree add command
			argsSent = args
			return "", nil
		}

		err := provider.CreateWorktree("../my-worktree", "my-branch", "main")

		assert.NoError(t, err)
		assert.Equal(t, []string{"worktree", "add", "../my-worktree", "-b", "my-branch", "upstream/main"}, argsSent)
	})

	t.Run("should return error if worktree creation fails", func(t *testing.T) {
		runGitCommand = func(args ...string) (out string, err error) {
			// Mock git remote get-url upstream to return error (no upstream)
			if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "upstream" {
				return "", fmt.Errorf("no such remote")
			}
			return "", fmt.Errorf("path already exists")
		}

		err := provider.CreateWorktree("../my-worktree", "my-branch", "main")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create worktree")
	})
}

func TestGitListWorktrees(t *testing.T) {
	provider := Provider{}

	t.Run("should parse worktree list output", func(t *testing.T) {
		runGitCommand = func(args ...string) (out string, err error) {
			if args[0] == "worktree" && args[1] == "list" {
				out = `worktree /path/to/main
HEAD abc123def456
branch refs/heads/main

worktree /path/to/feature
HEAD def789abc123
branch refs/heads/feature/test

worktree /path/to/detached
HEAD 111222333444
detached
`
				return out, nil
			}
			return "", nil
		}

		worktrees, err := provider.ListWorktrees()

		assert.NoError(t, err)
		assert.Len(t, worktrees, 3)
		assert.Equal(t, "/path/to/main", worktrees[0].Path)
		assert.Equal(t, "refs/heads/main", worktrees[0].Branch)
		assert.Equal(t, "abc123def456", worktrees[0].Commit)
		assert.Equal(t, "/path/to/feature", worktrees[1].Path)
		assert.Equal(t, "refs/heads/feature/test", worktrees[1].Branch)
	})

	t.Run("should handle empty worktree list", func(t *testing.T) {
		runGitCommand = func(args ...string) (out string, err error) {
			if args[0] == "worktree" && args[1] == "list" {
				return "", nil
			}
			return "", nil
		}

		worktrees, err := provider.ListWorktrees()

		assert.NoError(t, err)
		assert.Len(t, worktrees, 0)
	})

	t.Run("should return error if command fails", func(t *testing.T) {
		runGitCommand = func(args ...string) (out string, err error) {
			return "", fmt.Errorf("git command failed")
		}

		_, err := provider.ListWorktrees()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to list worktrees")
	})
}

func TestGitRemoveWorktree(t *testing.T) {
	provider := Provider{}

	t.Run("should remove worktree successfully", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			argsSent = args
			return "", nil
		}

		err := provider.RemoveWorktree("../my-worktree")

		assert.NoError(t, err)
		assert.Equal(t, []string{"worktree", "remove", "../my-worktree"}, argsSent)
	})

	t.Run("should return error if removal fails", func(t *testing.T) {
		runGitCommand = func(args ...string) (out string, err error) {
			return "", fmt.Errorf("worktree not found")
		}

		err := provider.RemoveWorktree("../nonexistent")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to remove worktree")
	})
}
