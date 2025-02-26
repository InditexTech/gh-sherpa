package gh

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCli_GetRepository(t *testing.T) {
	// Save the original function to restore it later
	originalExecute := Execute

	// Restore the original function after all tests
	defer func() {
		Execute = originalExecute
	}()

	tests := []struct {
		name         string
		mockResponse string
		mockError    error
		expectedRepo *domain.Repository
		expectError  bool
	}{
		{
			name: "Success",
			mockResponse: `{
				"name": "gh-sherpa",
				"owner": {
					"login": "InditexTech"
				},
				"nameWithOwner": "InditexTech/gh-sherpa",
				"defaultBranchRef": {
					"name": "main"
				}
			}`,
			mockError: nil,
			expectedRepo: &domain.Repository{
				Name:             "gh-sherpa",
				Owner:            "InditexTech",
				NameWithOwner:    "InditexTech/gh-sherpa",
				DefaultBranchRef: "main",
			},
			expectError: false,
		},
		{
			name:         "Error from gh CLI",
			mockResponse: "",
			mockError:    fmt.Errorf("gh CLI error"),
			expectedRepo: nil,
			expectError:  true,
		},
		{
			name:         "Error in stderr",
			mockResponse: "",
			mockError:    nil,
			expectedRepo: nil,
			expectError:  true,
		},
		{
			name:         "Invalid JSON response",
			mockResponse: "invalid json",
			mockError:    nil,
			expectedRepo: nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{}

			// Mock the Execute function
			Execute = func(args ...string) (stdout, stderr bytes.Buffer, err error) {
				if tt.mockError != nil {
					return bytes.Buffer{}, bytes.Buffer{}, tt.mockError
				}

				if tt.name == "Error in stderr" {
					stderr.WriteString("error message")
					return stdout, stderr, nil
				}

				stdout.WriteString(tt.mockResponse)
				return stdout, stderr, nil
			}

			repo, err := c.GetRepository()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, repo)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRepo, repo)
			}
		})
	}
}

func TestCli_GetPullRequestTemplate(t *testing.T) {
	// Save the original function to restore it later
	originalExecute := Execute

	// Restore the original function after all tests
	defer func() {
		Execute = originalExecute
	}()

	tests := []struct {
		name           string
		apiResponses   map[string]string
		apiErrors      map[string]error
		expectedResult string
		expectError    bool
	}{
		{
			name: "Template found in .github/pull_request_template.md",
			apiResponses: map[string]string{
				"/repos/{owner}/{repo}/contents/.github/pull_request_template.md": `{
					"content": "IyMgRGVzY3JpcHRpb24KClBsZWFzZSBkZXNjcmliZSB5b3VyIGNoYW5nZXMuLi4=",
					"encoding": "base64"
				}`,
			},
			expectedResult: "## Description\n\nPlease describe your changes...",
			expectError:    false,
		},
		{
			name: "Template found in .github/PULL_REQUEST_TEMPLATE directory",
			apiResponses: map[string]string{
				"/repos/{owner}/{repo}/contents/.github/PULL_REQUEST_TEMPLATE": `[
					{
						"name": "default.md",
						"path": ".github/PULL_REQUEST_TEMPLATE/default.md"
					}
				]`,
				"/repos/{owner}/{repo}/contents/.github/PULL_REQUEST_TEMPLATE/default.md": `{
					"content": "IyMgRGVzY3JpcHRpb24KClBsZWFzZSBkZXNjcmliZSB5b3VyIGNoYW5nZXMuLi4=",
					"encoding": "base64"
				}`,
			},
			expectedResult: "## Description\n\nPlease describe your changes...",
			expectError:    false,
		},
		{
			name:         "No template found",
			apiResponses: map[string]string{},
			apiErrors: map[string]error{
				"/repos/{owner}/{repo}/contents/.github/pull_request_template.md": fmt.Errorf("not found"),
				"/repos/{owner}/{repo}/contents/.github/PULL_REQUEST_TEMPLATE.md": fmt.Errorf("not found"),
				"/repos/{owner}/{repo}/contents/docs/pull_request_template.md":    fmt.Errorf("not found"),
				"/repos/{owner}/{repo}/contents/docs/PULL_REQUEST_TEMPLATE.md":    fmt.Errorf("not found"),
				"/repos/{owner}/{repo}/contents/pull_request_template.md":         fmt.Errorf("not found"),
				"/repos/{owner}/{repo}/contents/PULL_REQUEST_TEMPLATE.md":         fmt.Errorf("not found"),
				"/repos/{owner}/{repo}/contents/.github/PULL_REQUEST_TEMPLATE":    fmt.Errorf("not found"),
			},
			expectedResult: "",
			expectError:    false,
		},
		{
			name: "Template with non-base64 encoding",
			apiResponses: map[string]string{
				"/repos/{owner}/{repo}/contents/.github/pull_request_template.md": `{
					"content": "## Description\n\nPlease describe your changes...",
					"encoding": "utf-8"
				}`,
			},
			expectedResult: "## Description\n\nPlease describe your changes...",
			expectError:    false,
		},
		{
			name: "Invalid JSON response",
			apiResponses: map[string]string{
				"/repos/{owner}/{repo}/contents/.github/pull_request_template.md": `invalid json`,
			},
			expectedResult: "",
			expectError:    false, // The function handles errors internally and continues searching
		},
		{
			name: "Invalid base64 content",
			apiResponses: map[string]string{
				"/repos/{owner}/{repo}/contents/.github/pull_request_template.md": `{
					"content": "invalid base64!@#$",
					"encoding": "base64"
				}`,
			},
			expectedResult: "",
			expectError:    false, // The function handles errors internally and continues searching
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{}

			// Mock the Execute function
			Execute = func(args ...string) (stdout, stderr bytes.Buffer, err error) {
				if len(args) < 2 {
					return bytes.Buffer{}, bytes.Buffer{}, fmt.Errorf("invalid args")
				}

				path := args[1]
				if response, ok := tt.apiResponses[path]; ok {
					stdout.WriteString(response)
					return stdout, stderr, nil
				}

				if err, ok := tt.apiErrors[path]; ok {
					return stdout, stderr, err
				}

				return stdout, stderr, fmt.Errorf("unexpected path: %s", path)
			}

			result, err := c.GetPullRequestTemplate()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestCli_GetPullRequestForBranch(t *testing.T) {
	// Save the original function to restore it later
	originalExecute := Execute

	// Restore the original function after all tests
	defer func() {
		Execute = originalExecute
	}()

	tests := []struct {
		name         string
		branchName   string
		mockResponse string
		mockStderr   string
		mockError    error
		expectedPR   *domain.PullRequest
		expectError  bool
	}{
		{
			name:       "Success",
			branchName: "feature/GH-123-test",
			mockResponse: `{
				"closed": false,
				"number": 123,
				"state": "open",
				"title": "Test PR",
				"url": "https://github.com/InditexTech/gh-sherpa/pull/123"
			}`,
			mockStderr: "",
			mockError:  nil,
			expectedPR: &domain.PullRequest{
				Closed: false,
				Number: 123,
				State:  "open",
				Title:  "Test PR",
				Url:    "https://github.com/InditexTech/gh-sherpa/pull/123",
			},
			expectError: false,
		},
		{
			name:         "No PR found",
			branchName:   "feature/GH-456-no-pr",
			mockResponse: "",
			mockStderr:   "no pull requests found",
			mockError:    nil,
			expectedPR:   nil,
			expectError:  false,
		},
		{
			name:         "Error from gh CLI",
			branchName:   "feature/GH-789-error",
			mockResponse: "",
			mockStderr:   "",
			mockError:    fmt.Errorf("gh CLI error"),
			expectedPR:   nil,
			expectError:  true,
		},
		{
			name:         "Error in stderr",
			branchName:   "feature/GH-101-stderr-error",
			mockResponse: "",
			mockStderr:   "error message",
			mockError:    nil,
			expectedPR:   nil,
			expectError:  true,
		},
		{
			name:         "Invalid JSON response",
			branchName:   "feature/GH-102-invalid-json",
			mockResponse: "invalid json",
			mockStderr:   "",
			mockError:    nil,
			expectedPR:   nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{}

			// Mock the Execute function
			Execute = func(args ...string) (stdout, stderr bytes.Buffer, err error) {
				if tt.mockError != nil {
					return bytes.Buffer{}, bytes.Buffer{}, tt.mockError
				}

				if tt.mockStderr != "" {
					stderr.WriteString(tt.mockStderr)
				}

				stdout.WriteString(tt.mockResponse)
				return stdout, stderr, nil
			}

			pr, err := c.GetPullRequestForBranch(tt.branchName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, pr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPR, pr)
			}
		})
	}
}

func TestCli_Execute(t *testing.T) {
	// Save the original function to restore it later
	originalExecute := Execute

	// Restore the original function after all tests
	defer func() {
		Execute = originalExecute
	}()

	tests := []struct {
		name         string
		mockResponse string
		mockStderr   string
		mockError    error
		result       interface{}
		expectError  bool
	}{
		{
			name:         "Success",
			mockResponse: `{"key": "value"}`,
			mockStderr:   "",
			mockError:    nil,
			result: &struct {
				Key string `json:"key"`
			}{},
			expectError: false,
		},
		{
			name:         "Error from Execute",
			mockResponse: "",
			mockStderr:   "",
			mockError:    fmt.Errorf("execute error"),
			result:       &struct{}{},
			expectError:  true,
		},
		{
			name:         "Error in stderr",
			mockResponse: "",
			mockStderr:   "error message",
			mockError:    nil,
			result:       &struct{}{},
			expectError:  true,
		},
		{
			name:         "Invalid JSON response",
			mockResponse: "invalid json",
			mockStderr:   "",
			mockError:    nil,
			result:       &struct{}{},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{}

			// Mock the Execute function
			Execute = func(args ...string) (stdout, stderr bytes.Buffer, err error) {
				if tt.mockError != nil {
					return bytes.Buffer{}, bytes.Buffer{}, tt.mockError
				}

				if tt.mockStderr != "" {
					stderr.WriteString(tt.mockStderr)
				}

				stdout.WriteString(tt.mockResponse)
				return stdout, stderr, nil
			}

			err := c.Execute(tt.result, []string{"test"})

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.name == "Success" {
					assert.Equal(t, "value", tt.result.(*struct {
						Key string `json:"key"`
					}).Key)
				}
			}
		})
	}
}

func TestExecuteStringResult(t *testing.T) {
	// Save the original function to restore it later
	originalExecute := Execute

	// Restore the original function after all tests
	defer func() {
		Execute = originalExecute
	}()

	tests := []struct {
		name           string
		mockStdout     string
		mockStderr     string
		mockError      error
		expectedResult string
		expectError    bool
	}{
		{
			name:           "Success",
			mockStdout:     "command output",
			mockStderr:     "",
			mockError:      nil,
			expectedResult: "command output",
			expectError:    false,
		},
		{
			name:           "Error",
			mockStdout:     "",
			mockStderr:     "error details",
			mockError:      fmt.Errorf("command error"),
			expectedResult: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the Execute function
			Execute = func(args ...string) (stdout, stderr bytes.Buffer, err error) {
				stdout.WriteString(tt.mockStdout)
				stderr.WriteString(tt.mockStderr)
				return stdout, stderr, tt.mockError
			}

			result, err := ExecuteStringResult([]string{"test"})

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

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
		mockError   error
	}{
		{
			name:        "CreatePR",
			args:        args{title: "title", body: "body", baseBranch: "develop", headBranch: "asd", draft: false, labels: []string{}},
			wantPrURL:   "https://github.com/InditexTech/gh-sherpa/pulls/1",
			wantErr:     false,
			executeArgs: []string{"pr", "create", "-B", "develop", "-H", "asd", "-t", "title", "-b", "body"},
			mockError:   nil,
		},
		{
			name:        "CreatePR draft and default",
			args:        args{title: "", body: "", baseBranch: "develop", headBranch: "asd", draft: true, labels: []string{}},
			wantPrURL:   "https://github.com/InditexTech/gh-sherpa/pulls/1",
			wantErr:     false,
			executeArgs: []string{"pr", "create", "-B", "develop", "-H", "asd", "-d", "-f"},
			mockError:   nil,
		},
		{
			name:        "CreatePR with labels",
			args:        args{title: "", body: "", baseBranch: "develop", headBranch: "asd", draft: true, labels: []string{"label1", "label2"}},
			wantPrURL:   "https://github.com/InditexTech/gh-sherpa/pulls/1",
			wantErr:     false,
			executeArgs: []string{"pr", "create", "-B", "develop", "-H", "asd", "-d", "-f", "-l", "label1", "-l", "label2"},
			mockError:   nil,
		},
		{
			name:        "CreatePR with error",
			args:        args{title: "title", body: "body", baseBranch: "develop", headBranch: "asd", draft: false, labels: []string{}},
			wantPrURL:   "",
			wantErr:     true,
			executeArgs: []string{"pr", "create", "-B", "develop", "-H", "asd", "-t", "title", "-b", "body"},
			mockError:   fmt.Errorf("error creating PR"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{}

			var executeArgs []string
			ExecuteStringResult = func(args []string) (result string, err error) {
				executeArgs = args
				return "https://github.com/InditexTech/gh-sherpa/pulls/1\n", tt.mockError
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
