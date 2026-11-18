package git

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestGitCreateWorktree(t *testing.T) {
	provider := Provider{}

	t.Run("creates a worktree from origin and preserves a path with spaces", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "upstream" {
				return "", fmt.Errorf("no such remote")
			}
			argsSent = args
			return "", nil
		}

		err := provider.CreateWorktree("../worktrees/my feature", "feature/my-branch", "main")

		assert.NoError(t, err)
		assert.Equal(t, []string{"worktree", "add", "--no-track", "-b", "feature/my-branch", "../worktrees/my feature", "origin/main"}, argsSent)
	})

	t.Run("creates a worktree from upstream when it exists", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "upstream" {
				return "https://github.com/upstream/repo.git", nil
			}
			argsSent = args
			return "", nil
		}

		err := provider.CreateWorktree("../worktrees/feature", "feature/my-branch", "main")

		assert.NoError(t, err)
		assert.Equal(t, []string{"worktree", "add", "--no-track", "-b", "feature/my-branch", "../worktrees/feature", "upstream/main"}, argsSent)
	})

	t.Run("returns an error when git cannot create the worktree", func(t *testing.T) {
		runGitCommand = func(args ...string) (out string, err error) {
			if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "upstream" {
				return "", fmt.Errorf("no such remote")
			}
			return "", fmt.Errorf("path already exists")
		}

		err := provider.CreateWorktree("../worktrees/feature", "feature/my-branch", "main")

		assert.ErrorContains(t, err, "failed to create the worktree")
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

func TestGitCommitEmpty(t *testing.T) {
	provider := Provider{}

	t.Run("CommitEmpty should include a Signed-off-by trailer when signing is disabled", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			// Mock git config --get commit.gpgsign to return disabled
			if len(args) >= 3 && args[0] == "config" && args[1] == "--get" && args[2] == "commit.gpgsign" {
				return "false\n", nil
			}
			argsSent = args
			return "", nil
		}

		err := provider.CommitEmpty("chore: initial commit")

		assert.NoError(t, err)
		assert.Equal(t, []string{"commit", "--allow-empty", "-s", "-m", "chore: initial commit"}, argsSent)
	})

	t.Run("CommitEmpty should include both Signed-off-by and GPG signing when signing is enabled", func(t *testing.T) {
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			// Mock git config --get commit.gpgsign to return enabled
			if len(args) >= 3 && args[0] == "config" && args[1] == "--get" && args[2] == "commit.gpgsign" {
				return "true\n", nil
			}
			argsSent = args
			return "", nil
		}

		err := provider.CommitEmpty("chore: initial commit")

		assert.NoError(t, err)
		assert.Equal(t, []string{"commit", "--allow-empty", "-s", "-m", "chore: initial commit", "-S"}, argsSent)
	})

	t.Run("CommitEmpty should bypass hooks and preserve GPG signing when skipping Git hooks", func(t *testing.T) {
		provider := Provider{SkipGitHooks: true}
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			if len(args) >= 3 && args[0] == "config" && args[1] == "--get" && args[2] == "commit.gpgsign" {
				return "true\n", nil
			}
			argsSent = args
			return "", nil
		}

		err := provider.CommitEmpty("chore: initial commit")

		assert.NoError(t, err)
		assert.Equal(t, []string{"commit", "--allow-empty", "-s", "-m", "chore: initial commit", "--no-verify", "-S"}, argsSent)
	})

	t.Run("CommitEmpty should return an error if the commit fails", func(t *testing.T) {
		runGitCommand = func(args ...string) (out string, err error) {
			if len(args) >= 3 && args[0] == "config" && args[1] == "--get" && args[2] == "commit.gpgsign" {
				return "false\n", nil
			}
			err = fmt.Errorf("Failed to run Git command (%w)\n\nDetails:\n%s", err, "foo")
			return
		}

		err := provider.CommitEmpty("chore: initial commit")

		assert.Error(t, err)
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

	t.Run("GitPush should bypass the pre-push hook when skipping Git hooks", func(t *testing.T) {
		provider := Provider{SkipGitHooks: true}
		var argsSent []string
		runGitCommand = func(args ...string) (out string, err error) {
			argsSent = args
			return
		}

		err := provider.PushBranch("my-branch")

		assert.NoError(t, err)
		assert.Equal(t, []string{"push", "-u", "origin", "my-branch", "--no-verify"}, argsSent)
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

func TestGitFindBranch(t *testing.T) {
	repoDir := t.TempDir()
	gitCommand := func(args ...string) (string, error) {
		cmd := exec.Command("git", args...)
		cmd.Dir = repoDir
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	for _, args := range [][]string{
		{"-c", "init.templateDir=", "init"},
		{"-c", "user.name=Test", "-c", "user.email=test@example.com", "-c", "commit.gpgsign=false", "commit", "--allow-empty", "-m", "initial"},
		{"branch", "feature/GH-1"},
		{"branch", "feature/GH-10"},
		{"branch", "feature/GH-10-other"},
		{"branch", "feature/GH-2-description"},
		{"branch", "feature/GH-30"},
		{"branch", "feature/GH-30-other"},
	} {
		out, err := gitCommand(args...)
		require.NoError(t, err, "%s", out)
	}

	previousRunGitCommand := runGitCommand
	runGitCommand = gitCommand
	t.Cleanup(func() { runGitCommand = previousRunGitCommand })
	t.Setenv(DRY_RUN_ENV, "")

	for _, tt := range []struct {
		name    string
		issueID string
		want    string
	}{
		{name: "issue-only branch", issueID: "GH-1", want: "feature/GH-1"},
		{name: "branch with description", issueID: "GH-2", want: "feature/GH-2-description"},
		{name: "does not match longer issue numbers", issueID: "GH-3"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			provider := Provider{}
			branch, exists := provider.FindBranch(tt.issueID)
			assert.Equal(t, tt.want, branch)
			assert.Equal(t, tt.want != "", exists)
		})
	}
}
