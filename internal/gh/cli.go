package gh

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/cli/go-gh/v2"
)

type Cli struct{}

func (c *Cli) GetRepository() (repo *domain.Repository, err error) {
	baseCommand := []string{"repo", "view", "--json", "name,owner,nameWithOwner,defaultBranchRef"}

	stdout, stderr, err := gh.Exec(baseCommand...)

	if stderr.String() != "" {
		return nil, errors.New(stderr.String())
	}

	if err != nil {
		return
	}

	apiResponse := struct {
		Name             string
		Owner            struct{ Login string }
		NameWithOwner    string
		DefaultBranchRef struct{ Name string }
	}{}

	err = json.Unmarshal(stdout.Bytes(), &apiResponse)

	if err != nil {
		return
	}

	repo = &domain.Repository{
		Name:             apiResponse.Name,
		Owner:            apiResponse.Owner.Login,
		NameWithOwner:    apiResponse.NameWithOwner,
		DefaultBranchRef: apiResponse.DefaultBranchRef.Name,
	}

	return
}

func (c *Cli) Execute(result any, args []string) (err error) {
	stdout, stderr, err := Execute(args...)

	if stderr.String() != "" {
		return errors.New(stderr.String())
	}

	if err != nil {
		return
	}

	err = json.Unmarshal(stdout.Bytes(), &result)

	if err != nil {
		return
	}

	return
}

var Execute = func(args ...string) (stdout, stderr bytes.Buffer, err error) {
	return gh.Exec(args...)
}

var ExecuteStringResult = func(args []string) (result string, err error) {
	stdout, stderr, err := gh.Exec(args...)
	if err != nil {

		err = fmt.Errorf("failed to run GitHub CLI command (%w)\n\nDetails:\n%s", err, stderr.String())
	}

	result = stdout.String()
	return
}

func (c *Cli) CreatePullRequest(title string, body string, baseBranch string, headBranch string, draft bool, labels []string) (prURL string, err error) {
	args := []string{"pr", "create"}

	if baseBranch != "" {
		args = append(args, "-B", baseBranch)
	}

	if headBranch != "" {
		args = append(args, "-H", headBranch)
	}

	if draft {
		args = append(args, "-d")
	}

	if title == "" && body == "" {
		args = append(args, "-f")
	} else {
		args = append(args, "-t", title)
		args = append(args, "-b", body)
	}

	for _, label := range labels {
		args = append(args, "-l", label)
	}

	result, err := ExecuteStringResult(args)

	if err != nil {
		return
	}

	prURL = strings.Split(result, "\n")[0]
	prURL = strings.TrimSpace(prURL)

	return
}

func (c *Cli) GetPullRequestForBranch(branchName string) (*domain.PullRequest, error) {
	command := []string{"pr", "view", branchName, "--json", "closed,number,state,title,url"}

	stdout, stderr, err := gh.Exec(command...)
	if strings.Contains(stderr.String(), "no pull requests found") {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if stderr.String() != "" {
		return nil, fmt.Errorf("error while executing the command: %s", stderr.String())
	}

	var pr domain.PullRequest
	if err := json.Unmarshal(stdout.Bytes(), &pr); err != nil {
		return nil, err
	}

	return &pr, nil
}

func (c *Cli) GetPullRequestTemplate() (template string, err error) {
	templatePaths := []string{
		".github/pull_request_template.md",
		".github/PULL_REQUEST_TEMPLATE.md",
		"docs/pull_request_template.md",
		"docs/PULL_REQUEST_TEMPLATE.md",
		"pull_request_template.md",
		"PULL_REQUEST_TEMPLATE.md",
	}

	args := []string{"api", "/repos/{owner}/{repo}/contents/.github/PULL_REQUEST_TEMPLATE"}
	stdout, stderr, err := gh.Exec(args...)
	if err == nil && stderr.String() == "" {
		var contents []struct {
			Name     string `json:"name"`
			Path     string `json:"path"`
			Content  string `json:"content"`
			Encoding string `json:"encoding"`
		}
		if err := json.Unmarshal(stdout.Bytes(), &contents); err == nil {
			for _, content := range contents {
				if strings.HasSuffix(content.Name, ".md") {
					templatePaths = append(templatePaths, content.Path)
				}
			}
		}
	}

	for _, path := range templatePaths {
		args := []string{"api", fmt.Sprintf("/repos/{owner}/{repo}/contents/%s", path)}
		stdout, stderr, err := gh.Exec(args...)
		if err == nil && stderr.String() == "" {
			var response struct {
				Content  string `json:"content"`
				Encoding string `json:"encoding"`
			}

			if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
				continue
			}

			if response.Encoding == "base64" {
				decoded, err := base64.StdEncoding.DecodeString(response.Content)
				if err != nil {
					continue
				}
				return string(decoded), nil
			}
			return response.Content, nil
		}
	}

	return "", nil
}
