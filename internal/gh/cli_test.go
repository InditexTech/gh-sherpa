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

func TestCli_CreatePullRequest_InForkContext(t *testing.T) {
	tests := []struct {
		name             string
		title            string
		body             string
		baseBranch       string
		headBranch       string
		draft            bool
		labels           []string
		originResponse   string
		originError      error
		upstreamResponse string
		upstreamError    error
		expectedArgs     []string
		wantPrURL        string
		wantErr          bool
	}{
		{
			name:             "CreatePR in fork context with upstream repo",
			title:            "Test PR",
			body:             "Test body",
			baseBranch:       "main",
			headBranch:       "feature/test",
			draft:            false,
			labels:           []string{},
			originResponse:   "https://github.com/user/gh-sherpa.git\n",
			originError:      nil,
			upstreamResponse: "https://github.com/InditexTech/gh-sherpa.git\n",
			upstreamError:    nil,
			expectedArgs:     []string{"pr", "create", "-B", "main", "-H", "user:feature/test", "--repo", "InditexTech/gh-sherpa", "-t", "Test PR", "-b", "Test body"},
			wantPrURL:        "https://github.com/InditexTech/gh-sherpa/pulls/1",
			wantErr:          false,
		},
		{
			name:             "CreatePR in fork context with upstream repo error",
			title:            "Test PR",
			body:             "Test body",
			baseBranch:       "main",
			headBranch:       "feature/test",
			draft:            false,
			labels:           []string{},
			originResponse:   "https://github.com/user/gh-sherpa.git\n",
			originError:      nil,
			upstreamResponse: "",
			upstreamError:    errors.New("upstream not found"),
			expectedArgs:     []string{"pr", "create", "-B", "main", "-H", "feature/test", "-t", "Test PR", "-b", "Test body"},
			wantPrURL:        "https://github.com/InditexTech/gh-sherpa/pulls/1",
			wantErr:          false,
		},
		{
			name:             "CreatePR in fork context with empty upstream repo",
			title:            "Test PR",
			body:             "Test body",
			baseBranch:       "main",
			headBranch:       "feature/test",
			draft:            false,
			labels:           []string{},
			originResponse:   "https://github.com/user/gh-sherpa.git\n",
			originError:      nil,
			upstreamResponse: "\n",
			upstreamError:    nil,
			expectedArgs:     []string{"pr", "create", "-B", "main", "-H", "user:feature/test", "-t", "Test PR", "-b", "Test body"},
			wantPrURL:        "https://github.com/InditexTech/gh-sherpa/pulls/1",
			wantErr:          false,
		},
		{
			name:             "CreatePR in fork context with SSH upstream",
			title:            "Test PR",
			body:             "Test body",
			baseBranch:       "main",
			headBranch:       "feature/test",
			draft:            true,
			labels:           []string{"enhancement", "fork"},
			originResponse:   "git@github.com:user/gh-sherpa.git\n",
			originError:      nil,
			upstreamResponse: "git@github.com:InditexTech/gh-sherpa.git\n",
			upstreamError:    nil,
			expectedArgs:     []string{"pr", "create", "-B", "main", "-H", "user:feature/test", "--repo", "InditexTech/gh-sherpa", "-d", "-t", "Test PR", "-b", "Test body", "-l", "enhancement", "-l", "fork"},
			wantPrURL:        "https://github.com/InditexTech/gh-sherpa/pulls/1",
			wantErr:          false,
		},
		{
			name:             "CreatePR non-fork context (no upstream)",
			title:            "Test PR",
			body:             "Test body",
			baseBranch:       "main",
			headBranch:       "feature/test",
			draft:            false,
			labels:           []string{},
			originResponse:   "https://github.com/InditexTech/gh-sherpa.git\n",
			originError:      nil,
			upstreamResponse: "",
			upstreamError:    errors.New("upstream not found"),
			expectedArgs:     []string{"pr", "create", "-B", "main", "-H", "feature/test", "-t", "Test PR", "-b", "Test body"},
			wantPrURL:        "https://github.com/InditexTech/gh-sherpa/pulls/1",
			wantErr:          false,
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

			var executeArgs []string
			originalExecuteStringResult := ExecuteStringResult
			defer func() { ExecuteStringResult = originalExecuteStringResult }()

			ExecuteStringResult = func(args []string) (result string, err error) {
				executeArgs = args
				return "https://github.com/InditexTech/gh-sherpa/pulls/1\n", nil
			}

			gotPrURL, err := c.CreatePullRequest(tt.title, tt.body, tt.baseBranch, tt.headBranch, tt.draft, tt.labels)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedArgs, executeArgs)
			assert.Equal(t, tt.wantPrURL, gotPrURL)
		})
	}
}

func TestCli_formatHeadBranchForFork(t *testing.T) {
	tests := []struct {
		name             string
		headBranch       string
		originResponse   string
		originError      error
		upstreamResponse string
		upstreamError    error
		expectedResult   string
		wantErr          bool
	}{
		{
			name:             "Success - Fork context with HTTPS origin",
			headBranch:       "feature/new-feature",
			originResponse:   "https://github.com/user/repo.git\n",
			originError:      nil,
			upstreamResponse: "https://github.com/InditexTech/gh-sherpa.git\n",
			upstreamError:    nil,
			expectedResult:   "user:feature/new-feature",
			wantErr:          false,
		},
		{
			name:             "Success - Fork context with SSH origin",
			headBranch:       "bugfix/fix-issue",
			originResponse:   "git@github.com:fork-user/awesome-repo.git\n",
			originError:      nil,
			upstreamResponse: "git@github.com:original/awesome-repo.git\n",
			upstreamError:    nil,
			expectedResult:   "fork-user:bugfix/fix-issue",
			wantErr:          false,
		},
		{
			name:             "Success - Complex organization name in fork",
			headBranch:       "hotfix/urgent-fix",
			originResponse:   "https://github.com/my.fork.org/project.git\n",
			originError:      nil,
			upstreamResponse: "https://github.com/upstream.org/project.git\n",
			upstreamError:    nil,
			expectedResult:   "my.fork.org:hotfix/urgent-fix",
			wantErr:          false,
		},
		{
			name:             "Success - Branch with special characters",
			headBranch:       "feature/user-auth_v2",
			originResponse:   "https://github.com/dev-team/app.git\n",
			originError:      nil,
			upstreamResponse: "https://github.com/main-org/app.git\n",
			upstreamError:    nil,
			expectedResult:   "dev-team:feature/user-auth_v2",
			wantErr:          false,
		},
		{
			name:             "Non-fork context - No upstream remote",
			headBranch:       "feature/new-feature",
			originResponse:   "https://github.com/user/repo.git\n",
			originError:      nil,
			upstreamResponse: "",
			upstreamError:    errors.New("upstream not found"),
			expectedResult:   "feature/new-feature",
			wantErr:          false,
		},
		{
			name:             "Non-fork context - No origin remote",
			headBranch:       "feature/test",
			originResponse:   "",
			originError:      errors.New("origin not found"),
			upstreamResponse: "https://github.com/upstream/repo.git\n",
			upstreamError:    nil,
			expectedResult:   "feature/test",
			wantErr:          false,
		},
		{
			name:             "Fork context - No origin remote",
			headBranch:       "feature/test",
			originResponse:   "",
			originError:      errors.New("origin not found"),
			upstreamResponse: "https://github.com/upstream/repo.git\n",
			upstreamError:    nil,
			expectedResult:   "feature/test",
			wantErr:          false,
		},
		{
			name:             "Success - Origin URL without .git suffix",
			headBranch:       "develop",
			originResponse:   "https://github.com/fork-owner/project\n",
			originError:      nil,
			upstreamResponse: "https://github.com/original-owner/project\n",
			upstreamError:    nil,
			expectedResult:   "fork-owner:develop",
			wantErr:          false,
		},
		{
			name:             "Success - Origin URL with extra whitespace",
			headBranch:       "main",
			originResponse:   "  https://github.com/my-fork/repo.git  \n",
			originError:      nil,
			upstreamResponse: "https://github.com/original/repo.git\n",
			upstreamError:    nil,
			expectedResult:   "my-fork:main",
			wantErr:          false,
		},
		{
			name:             "Edge case - Empty branch name",
			headBranch:       "",
			originResponse:   "https://github.com/user/repo.git\n",
			originError:      nil,
			upstreamResponse: "https://github.com/upstream/repo.git\n",
			upstreamError:    nil,
			expectedResult:   "user:",
			wantErr:          false,
		},
		{
			name:             "Success - Username with hyphens and numbers",
			headBranch:       "feature/api-v2",
			originResponse:   "https://github.com/user-123/my-project.git\n",
			originError:      nil,
			upstreamResponse: "https://github.com/company/my-project.git\n",
			upstreamError:    nil,
			expectedResult:   "user-123:feature/api-v2",
			wantErr:          false,
		},
		{
			name:             "Git command error handling",
			headBranch:       "feature/test",
			originResponse:   "",
			originError:      errors.New("git command failed"),
			upstreamResponse: "",
			upstreamError:    errors.New("git command failed"),
			expectedResult:   "feature/test",
			wantErr:          false,
		},
		{
			name:             "GetRemoteConfiguration error - return original branch",
			headBranch:       "feature/error-handling",
			originResponse:   "",
			originError:      errors.New("fatal: not a git repository"),
			upstreamResponse: "",
			upstreamError:    errors.New("fatal: not a git repository"),
			expectedResult:   "feature/error-handling",
			wantErr:          false,
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

			result, err := c.formatHeadBranchForFork(tt.headBranch)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
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
		name             string
		forkName         string
		forkExistsError  error
		forkExistsResult string
		mockResult       string
		mockError        error
		expectedArgs     []string
		wantErr          bool
	}{
		{
			name:             "Success - No organization specified",
			forkName:         "",
			forkExistsError:  nil,
			forkExistsResult: "",
			mockResult:       "✓ Created fork user/repo\n✓ Added remote \"origin\"",
			mockError:        nil,
			expectedArgs:     []string{"repo", "fork", "--remote"},
			wantErr:          false,
		},
		{
			name:             "Success - Organization specified",
			forkName:         "MyOrg/repo-name",
			forkExistsError:  errors.New("Could not resolve to a Repository"),
			forkExistsResult: "",
			mockResult:       "✓ Created fork MyOrg/repo-name\n✓ Added remote \"origin\"",
			mockError:        nil,
			expectedArgs:     []string{"repo", "fork", "--remote", "--org", "MyOrg"},
			wantErr:          false,
		},
		{
			name:             "Success - Invalid format (single part)",
			forkName:         "just-a-name",
			forkExistsError:  errors.New("Could not resolve to a Repository"),
			forkExistsResult: "",
			mockResult:       "✓ Created fork user/repo\n✓ Added remote \"origin\"",
			mockError:        nil,
			expectedArgs:     []string{"repo", "fork", "--remote"},
			wantErr:          false,
		},
		{
			name:             "Success - Multiple parts but invalid",
			forkName:         "org/repo/extra",
			forkExistsError:  errors.New("Could not resolve to a Repository"),
			forkExistsResult: "",
			mockResult:       "✓ Created fork user/repo\n✓ Added remote \"origin\"",
			mockError:        nil,
			expectedArgs:     []string{"repo", "fork", "--remote"},
			wantErr:          false,
		},
		{
			name:             "Error - Repository not found",
			forkName:         "invalid/repo",
			forkExistsError:  errors.New("Could not resolve to a Repository"),
			forkExistsResult: "",
			mockResult:       "",
			mockError:        fmt.Errorf("failed to run GitHub CLI command (exit status 1)\n\nDetails:\nrepository not found"),
			expectedArgs:     []string{"repo", "fork", "--remote", "--org", "invalid"},
			wantErr:          true,
		},
		{
			name:             "Error - Permission denied",
			forkName:         "private/repo",
			forkExistsError:  fmt.Errorf("failed to run GitHub CLI command (exit status 1)\n\nDetails:\npermission denied"),
			forkExistsResult: "",
			mockResult:       "",
			mockError:        fmt.Errorf("failed to run GitHub CLI command (exit status 1)\n\nDetails:\npermission denied"),
			expectedArgs:     []string{"repo", "fork", "--remote", "--org", "private"},
			wantErr:          true,
		},
		{
			name:             "Fork already exists - Configure remotes",
			forkName:         "MyOrg/existing-fork",
			forkExistsError:  nil,
			forkExistsResult: `{"name":"existing-fork"}`,
			mockResult:       "",
			mockError:        nil,
			expectedArgs:     []string(nil),
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{}

			var capturedArgs []string
			originalExecuteStringResult := ExecuteStringResult
			defer func() { ExecuteStringResult = originalExecuteStringResult }()

			// Mock git command for remote configuration
			originalExecuteGitCommand := executeGitCommand
			defer func() { executeGitCommand = originalExecuteGitCommand }()

			executeGitCommand = func(args ...string) (result string, err error) {
				// Mock git remote commands for ConfigureRemotesForExistingFork
				if len(args) >= 2 && args[0] == "remote" && args[1] == "set-url" {
					return "", nil
				}
				if len(args) >= 2 && args[0] == "remote" && args[1] == "add" {
					return "", nil
				}
				return "", errors.New("unexpected git command")
			}

			callCount := 0
			ExecuteStringResult = func(args []string) (result string, err error) {
				callCount++
				// Handle ForkExists check (repo view) - only when forkName is not empty
				if tt.forkName != "" && len(args) >= 2 && args[0] == "repo" && args[1] == "view" {
					return tt.forkExistsResult, tt.forkExistsError
				}
				// Handle fork creation (repo fork)
				if len(args) >= 2 && args[0] == "repo" && args[1] == "fork" {
					capturedArgs = args
					return tt.mockResult, tt.mockError
				}
				return "", fmt.Errorf("unexpected call: %v", args)
			}

			err := c.CreateFork(tt.forkName)

			// Verify the arguments passed to ExecuteStringResult for fork creation
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

func TestCli_getUpstreamRepository(t *testing.T) {
	tests := []struct {
		name             string
		originResponse   string
		originError      error
		upstreamResponse string
		upstreamError    error
		expectedRepo     string
		wantErr          bool
		expectedErrMsg   string
	}{
		{
			name:             "Success - HTTPS upstream URL",
			originResponse:   "https://github.com/user/repo.git\n",
			originError:      nil,
			upstreamResponse: "https://github.com/InditexTech/gh-sherpa.git\n",
			upstreamError:    nil,
			expectedRepo:     "InditexTech/gh-sherpa",
			wantErr:          false,
		},
		{
			name:             "Success - SSH upstream URL",
			originResponse:   "git@github.com:user/repo.git\n",
			originError:      nil,
			upstreamResponse: "git@github.com:InditexTech/gh-sherpa.git\n",
			upstreamError:    nil,
			expectedRepo:     "InditexTech/gh-sherpa",
			wantErr:          false,
		},
		{
			name:             "Success - upstream URL without .git suffix",
			originResponse:   "https://github.com/user/repo\n",
			originError:      nil,
			upstreamResponse: "https://github.com/InditexTech/gh-sherpa\n",
			upstreamError:    nil,
			expectedRepo:     "InditexTech/gh-sherpa",
			wantErr:          false,
		},
		{
			name:             "Success - upstream URL with extra whitespace",
			originResponse:   "https://github.com/user/repo.git\n",
			originError:      nil,
			upstreamResponse: "  https://github.com/InditexTech/gh-sherpa.git  \n",
			upstreamError:    nil,
			expectedRepo:     "InditexTech/gh-sherpa",
			wantErr:          false,
		},
		{
			name:             "Error - No upstream remote found",
			originResponse:   "https://github.com/user/repo.git\n",
			originError:      nil,
			upstreamResponse: "",
			upstreamError:    errors.New("upstream not found"),
			expectedRepo:     "",
			wantErr:          true,
			expectedErrMsg:   "no upstream remote found",
		},
		{
			name:             "Error - Git command fails for both remotes",
			originResponse:   "",
			originError:      errors.New("origin not found"),
			upstreamResponse: "",
			upstreamError:    errors.New("upstream not found"),
			expectedRepo:     "",
			wantErr:          true,
			expectedErrMsg:   "no upstream remote found",
		},
		{
			name:             "Success - Complex organization names",
			originResponse:   "https://github.com/my-fork-org/repo.git\n",
			originError:      nil,
			upstreamResponse: "https://github.com/my.upstream.org/awesome-repo.git\n",
			upstreamError:    nil,
			expectedRepo:     "my.upstream.org/awesome-repo",
			wantErr:          false,
		},
		{
			name:             "Success - Different user upstream",
			originResponse:   "https://github.com/fork-user/repo.git\n",
			originError:      nil,
			upstreamResponse: "git@github.com:original-author/repo.git\n",
			upstreamError:    nil,
			expectedRepo:     "original-author/repo",
			wantErr:          false,
		},
		{
			name:             "Success - Repository with special characters",
			originResponse:   "https://github.com/user/my-fork.git\n",
			originError:      nil,
			upstreamResponse: "https://github.com/upstream/my_awesome-repo.123.git\n",
			upstreamError:    nil,
			expectedRepo:     "upstream/my_awesome-repo.123",
			wantErr:          false,
		},
		{
			name:             "Error - Only origin exists, no upstream",
			originResponse:   "https://github.com/user/repo.git\n",
			originError:      nil,
			upstreamResponse: "",
			upstreamError:    errors.New("fatal: No such remote 'upstream'"),
			expectedRepo:     "",
			wantErr:          true,
			expectedErrMsg:   "no upstream remote found",
		},
		{
			name:             "Error - GetRemoteConfiguration fails",
			originResponse:   "",
			originError:      errors.New("fatal: not a git repository (or any of the parent directories): .git"),
			upstreamResponse: "",
			upstreamError:    errors.New("fatal: not a git repository (or any of the parent directories): .git"),
			expectedRepo:     "",
			wantErr:          true,
			expectedErrMsg:   "", // No specific error message needed, just that an error occurs
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

			repo, err := c.getUpstreamRepository()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Equal(t, tt.expectedRepo, repo)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRepo, repo)
			}
		})
	}
}

func Test_executeGitCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name:    "Success - git version command",
			args:    []string{"version"},
			wantErr: false,
		},
		{
			name:        "Error - git invalid command",
			args:        []string{"this-is-not-a-valid-git-command"},
			wantErr:     true,
			errContains: "failed to run git command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the actual function (not mocked)
			result, err := executeGitCommand(tt.args...)

			// Verify results
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				// For version command, we expect output. For config, it might be empty if no global config
				if tt.name == "Success - git version command" {
					assert.NotEmpty(t, result)
				}
			}
		})
	}
}
