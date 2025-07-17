package gh

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCli_CreatePullRequest(t *testing.T) {
	type args struct {
		title      string
		body       string
		baseBranch string
		headBranch string
		draft      bool
		labels     []string
	}
	tests := []struct {
		name        string
		args        args
		wantPrURL   string
		wantErr     bool
		executeArgs []string
	}{
		{
			name:        "CreatePR",
			args:        args{title: "title", body: "body", baseBranch: "develop", headBranch: "asd", draft: false, labels: []string{}},
			wantPrURL:   "https://github.com/InditexTech/gh-sherpa/pulls/1",
			wantErr:     false,
			executeArgs: []string{"pr", "create", "-B", "develop", "-H", "asd", "-t", "title", "-b", "body"},
		},
		{
			name:        "CreatePR draft and default",
			args:        args{title: "", body: "", baseBranch: "develop", headBranch: "asd", draft: true, labels: []string{}},
			wantPrURL:   "https://github.com/InditexTech/gh-sherpa/pulls/1",
			wantErr:     false,
			executeArgs: []string{"pr", "create", "-B", "develop", "-H", "asd", "-d", "-f"},
		},
		{
			name:        "CreatePR with labels",
			args:        args{title: "", body: "", baseBranch: "develop", headBranch: "asd", draft: true, labels: []string{"label1", "label2"}},
			wantPrURL:   "https://github.com/InditexTech/gh-sherpa/pulls/1",
			wantErr:     false,
			executeArgs: []string{"pr", "create", "-B", "develop", "-H", "asd", "-d", "-f", "-l", "label1", "-l", "label2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{}

			// Mock the git command execution to avoid fork detection
			originalExecuteGitCommand := executeGitCommand
			defer func() { executeGitCommand = originalExecuteGitCommand }()

			executeGitCommand = func(args ...string) (result string, err error) {
				// Return empty results to simulate no upstream remote (non-fork scenario)
				return "", errors.New("remote not found")
			}

			var executeArgs []string
			ExecuteStringResult = func(args []string) (result string, err error) {
				executeArgs = args
				return "https://github.com/InditexTech/gh-sherpa/pulls/1\n", nil
			}

			gotPrURL, err := c.CreatePullRequest(tt.args.title, tt.args.body, tt.args.baseBranch, tt.args.headBranch, tt.args.draft, tt.args.labels)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.executeArgs, executeArgs)
			assert.Equal(t, tt.wantPrURL, gotPrURL)
		})
	}
}

func TestCli_GetRemoteConfiguration(t *testing.T) {
	tests := []struct {
		name             string
		originResponse   string
		originError      error
		upstreamResponse string
		upstreamError    error
		expected         map[string]string
	}{
		{
			name:             "Both remotes exist",
			originResponse:   "https://github.com/user/repo.git\n",
			originError:      nil,
			upstreamResponse: "https://github.com/upstream/repo.git\n",
			upstreamError:    nil,
			expected: map[string]string{
				"origin":   "https://github.com/user/repo.git",
				"upstream": "https://github.com/upstream/repo.git",
			},
		},
		{
			name:             "Only origin exists",
			originResponse:   "https://github.com/user/repo.git\n",
			originError:      nil,
			upstreamResponse: "",
			upstreamError:    errors.New("upstream not found"),
			expected: map[string]string{
				"origin": "https://github.com/user/repo.git",
			},
		},
		{
			name:             "Only upstream exists",
			originResponse:   "",
			originError:      errors.New("origin not found"),
			upstreamResponse: "https://github.com/upstream/repo.git\n",
			upstreamError:    nil,
			expected: map[string]string{
				"upstream": "https://github.com/upstream/repo.git",
			},
		},
		{
			name:             "No remotes exist",
			originResponse:   "",
			originError:      errors.New("origin not found"),
			upstreamResponse: "",
			upstreamError:    errors.New("upstream not found"),
			expected:         map[string]string{},
		},
		{
			name:             "Remotes with whitespace",
			originResponse:   "  https://github.com/user/repo.git  \n",
			originError:      nil,
			upstreamResponse: "\t https://github.com/upstream/repo.git \t\n",
			upstreamError:    nil,
			expected: map[string]string{
				"origin":   "https://github.com/user/repo.git",
				"upstream": "https://github.com/upstream/repo.git",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{}

			originalExecuteGitCommand := executeGitCommand
			defer func() { executeGitCommand = originalExecuteGitCommand }()

			executeGitCommand = func(args ...string) (result string, err error) {
				// Check if this is a git remote get-url origin command
				if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "origin" {
					return tt.originResponse, tt.originError
				}
				// Check if this is a git remote get-url upstream command
				if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "upstream" {
					return tt.upstreamResponse, tt.upstreamError
				}
				return "", errors.New("unexpected command")
			}

			result, err := c.GetRemoteConfiguration()

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCli_CreateFork(t *testing.T) {
	tests := []struct {
		name         string
		forkName     string
		mockResult   string
		mockError    error
		expectedArgs []string
		wantErr      bool
	}{
		{
			name:         "Success - No organization specified",
			forkName:     "",
			mockResult:   "✓ Created fork user/repo\n✓ Added remote \"origin\"",
			mockError:    nil,
			expectedArgs: []string{"repo", "fork", "--remote"},
			wantErr:      false,
		},
		{
			name:         "Success - Organization specified",
			forkName:     "MyOrg/repo-name",
			mockResult:   "✓ Created fork MyOrg/repo-name\n✓ Added remote \"origin\"",
			mockError:    nil,
			expectedArgs: []string{"repo", "fork", "--remote", "--org", "MyOrg"},
			wantErr:      false,
		},
		{
			name:         "Success - Invalid format (single part)",
			forkName:     "just-a-name",
			mockResult:   "✓ Created fork user/repo\n✓ Added remote \"origin\"",
			mockError:    nil,
			expectedArgs: []string{"repo", "fork", "--remote"},
			wantErr:      false,
		},
		{
			name:         "Success - Multiple parts but invalid",
			forkName:     "org/repo/extra",
			mockResult:   "✓ Created fork user/repo\n✓ Added remote \"origin\"",
			mockError:    nil,
			expectedArgs: []string{"repo", "fork", "--remote"},
			wantErr:      false,
		},
		{
			name:         "Error - Repository not found",
			forkName:     "invalid/repo",
			mockResult:   "",
			mockError:    fmt.Errorf("failed to run GitHub CLI command (exit status 1)\n\nDetails:\nrepository not found"),
			expectedArgs: []string{"repo", "fork", "--remote", "--org", "invalid"},
			wantErr:      true,
		},
		{
			name:         "Error - Permission denied",
			forkName:     "private/repo",
			mockResult:   "",
			mockError:    fmt.Errorf("failed to run GitHub CLI command (exit status 1)\n\nDetails:\npermission denied"),
			expectedArgs: []string{"repo", "fork", "--remote", "--org", "private"},
			wantErr:      true,
		},
		{
			name:         "Error - Fork already exists",
			forkName:     "MyOrg/existing-fork",
			mockResult:   "",
			mockError:    fmt.Errorf("failed to run GitHub CLI command (exit status 1)\n\nDetails:\nfork already exists"),
			expectedArgs: []string{"repo", "fork", "--remote", "--org", "MyOrg"},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{}

			var capturedArgs []string
			originalExecuteStringResult := ExecuteStringResult
			defer func() { ExecuteStringResult = originalExecuteStringResult }()

			ExecuteStringResult = func(args []string) (result string, err error) {
				capturedArgs = args
				return tt.mockResult, tt.mockError
			}

			err := c.CreateFork(tt.forkName)

			// Verify the arguments passed to ExecuteStringResult
			assert.Equal(t, tt.expectedArgs, capturedArgs)

			// Verify the error behavior
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCli_SetDefaultRepository(t *testing.T) {
	tests := []struct {
		name         string
		repo         string
		mockResult   string
		mockError    error
		expectedArgs []string
		wantErr      bool
	}{
		{
			name:         "Success - Standard repository",
			repo:         "InditexTech/gh-sherpa",
			mockResult:   "",
			mockError:    nil,
			expectedArgs: []string{"repo", "set-default", "InditexTech/gh-sherpa"},
			wantErr:      false,
		},
		{
			name:         "Success - User repository",
			repo:         "user/another-repo",
			mockResult:   "",
			mockError:    nil,
			expectedArgs: []string{"repo", "set-default", "user/another-repo"},
			wantErr:      false,
		},
		{
			name:         "Error - Repository not found",
			repo:         "invalid/repo",
			mockResult:   "",
			mockError:    fmt.Errorf("failed to run GitHub CLI command (exit status 1)\n\nDetails:\nrepository not found"),
			expectedArgs: []string{"repo", "set-default", "invalid/repo"},
			wantErr:      true,
		},
		{
			name:         "Error - Authentication required",
			repo:         "private/repo",
			mockResult:   "",
			mockError:    fmt.Errorf("failed to run GitHub CLI command (exit status 1)\n\nDetails:\nauthentication required"),
			expectedArgs: []string{"repo", "set-default", "private/repo"},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{}

			var capturedArgs []string
			originalExecuteStringResult := ExecuteStringResult
			defer func() { ExecuteStringResult = originalExecuteStringResult }()

			ExecuteStringResult = func(args []string) (result string, err error) {
				capturedArgs = args
				return tt.mockResult, tt.mockError
			}

			err := c.SetDefaultRepository(tt.repo)

			// Verify the arguments passed to ExecuteStringResult
			assert.Equal(t, tt.expectedArgs, capturedArgs)

			// Verify the error behavior
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCli_IsRepositoryFork(t *testing.T) {
	tests := []struct {
		name         string
		mockResult   string
		mockError    error
		expectedArgs []string
		expectedFork bool
		wantErr      bool
	}{
		{
			name:         "Success - Repository is a fork",
			mockResult:   `{"isFork":true}`,
			mockError:    nil,
			expectedArgs: []string{"repo", "view", "--json", "isFork"},
			expectedFork: true,
			wantErr:      false,
		},
		{
			name:         "Success - Repository is not a fork",
			mockResult:   `{"isFork":false}`,
			mockError:    nil,
			expectedArgs: []string{"repo", "view", "--json", "isFork"},
			expectedFork: false,
			wantErr:      false,
		},
		{
			name:         "Error - Repository not found",
			mockResult:   "",
			mockError:    fmt.Errorf("failed to run GitHub CLI command (exit status 1)\n\nDetails:\nrepository not found"),
			expectedArgs: []string{"repo", "view", "--json", "isFork"},
			expectedFork: false,
			wantErr:      true,
		},
		{
			name:         "Error - Authentication required",
			mockResult:   "",
			mockError:    fmt.Errorf("failed to run GitHub CLI command (exit status 1)\n\nDetails:\nauthentication required"),
			expectedArgs: []string{"repo", "view", "--json", "isFork"},
			expectedFork: false,
			wantErr:      true,
		},
		{
			name:         "Success - Fork status with additional fields",
			mockResult:   `{"isFork":true,"name":"my-repo","owner":{"login":"user"}}`,
			mockError:    nil,
			expectedArgs: []string{"repo", "view", "--json", "isFork"},
			expectedFork: true,
			wantErr:      false,
		},
		{
			name:         "Error - Invalid JSON response",
			mockResult:   `{"isFork":true,invalid}`,
			mockError:    nil,
			expectedArgs: []string{"repo", "view", "--json", "isFork"},
			expectedFork: false,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{}

			var capturedArgs []string
			originalExecuteStringResult := ExecuteStringResult
			defer func() { ExecuteStringResult = originalExecuteStringResult }()

			ExecuteStringResult = func(args []string) (result string, err error) {
				capturedArgs = args
				return tt.mockResult, tt.mockError
			}

			isFork, err := c.IsRepositoryFork()

			// Verify the arguments passed to ExecuteStringResult
			assert.Equal(t, tt.expectedArgs, capturedArgs)

			// Verify the error behavior
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedFork, isFork)
			}
		})
	}
}

func TestGetRemoteConfiguration_Logic(t *testing.T) {
	// This test already works well, so we'll keep a simplified version
	c := &Cli{}

	originalExecuteGitCommand := executeGitCommand
	defer func() { executeGitCommand = originalExecuteGitCommand }()

	executeGitCommand = func(args ...string) (result string, err error) {
		if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "origin" {
			return "https://github.com/user/repo.git\n", nil
		}
		if len(args) >= 3 && args[0] == "remote" && args[1] == "get-url" && args[2] == "upstream" {
			return "https://github.com/upstream/repo.git\n", nil
		}
		return "", nil
	}

	result, err := c.GetRemoteConfiguration()

	assert.NoError(t, err)
	expected := map[string]string{
		"origin":   "https://github.com/user/repo.git",
		"upstream": "https://github.com/upstream/repo.git",
	}
	assert.Equal(t, expected, result)
}
