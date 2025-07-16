package gh

import (
	"bytes"
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

func (c *Cli) IsRepositoryFork() (bool, error) {
	command := []string{"repo", "view", "--json", "isFork"}

	jsonResult, err := ExecuteStringResult(command)
	if err != nil {
		return false, err
	}

	var result struct {
		IsFork bool `json:"isFork"`
	}

	if err := json.Unmarshal([]byte(jsonResult), &result); err != nil {
		return false, err
	}

	return result.IsFork, nil
}

func (c *Cli) CreateFork(forkName string) error {
	args := []string{"repo", "fork", "--remote"}

	if forkName != "" {
		parts := strings.Split(forkName, "/")
		if len(parts) == 2 {
			args = append(args, "--org", parts[0])
		}
	}

	_, err := ExecuteStringResult(args)
	if err != nil {
		return fmt.Errorf("error creating fork: %s", err.Error())
	}

	return nil
}

func (c *Cli) SetDefaultRepository(repo string) error {
	args := []string{"repo", "set-default", repo}

	_, err := ExecuteStringResult(args)
	if err != nil {
		return fmt.Errorf("error setting default repository: %s", err.Error())
	}

	return err
}

func (c *Cli) GetRemoteConfiguration() (map[string]string, error) {
	remotes := make(map[string]string)

	originResult, err := ExecuteStringResult([]string{"git", "remote", "get-url", "origin"})
	if err == nil {
		remotes["origin"] = strings.TrimSpace(originResult)
	}

	upstreamResult, err := ExecuteStringResult([]string{"git", "remote", "get-url", "upstream"})
	if err == nil {
		remotes["upstream"] = strings.TrimSpace(upstreamResult)
	}

	return remotes, nil
}
