package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/InditexTech/gh-sherpa/internal/domain"
	"github.com/InditexTech/gh-sherpa/internal/logging"
)

const gitBin = "git"

const DRY_RUN_ENV = "SHERPA_DRY_RUN"

type Provider struct{}

var _ domain.GitProvider = (*Provider)(nil)

var runGitCommand = func(args ...string) (out string, err error) {
	logging.Debugf("Running git command: %s %v", gitBin, strings.Join(args, " "))

	if os.Getenv(DRY_RUN_ENV) != "" {
		fmt.Printf("DRY RUN: \"%s %v\"\n", gitBin, strings.Join(args, " "))
		return
	}

	cmd := exec.Command(gitBin, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf(stderr.String())
		return
	}

	out = stdout.String()
	return
}

func (p *Provider) BranchExists(branch string) bool {
	args := []string{"show-ref", "--verify", "refs/heads/" + branch}

	_, err := runGitCommand(args...)

	if os.Getenv(DRY_RUN_ENV) != "" {
		return false
	}

	return err == nil
}

func (p *Provider) FetchBranchFromOrigin(branch string) (err error) {
	args := []string{"fetch", "origin", branch}

	_, err = runGitCommand(args...)

	return
}

func (p *Provider) CheckoutNewBranchFromOrigin(branch string, base string) (err error) {
	args := []string{"checkout", "--no-track", "-b", branch, "origin/" + base}

	_, err = runGitCommand(args...)

	if err != nil {
		err = fmt.Errorf("failed to checkout the new branch.\n\nDetails:\n%s", err)

		return
	}

	return
}

func (p *Provider) GetCurrentBranch() (branchName string, err error) {
	args := []string{"rev-parse", "--abbrev-ref", "HEAD"}

	out, err := runGitCommand(args...)

	if err != nil {
		err = fmt.Errorf("failed to get the current branch.\n\nDetails:\n%s", err)

		return
	}

	if out != "" {
		branchName = strings.Split(out, "\n")[0]
		branchName = strings.TrimSpace(branchName)
	}

	return
}

func (p *Provider) FindBranch(substring string) (branch string, exists bool) {
	args := []string{"rev-parse", "--abbrev-ref", "--branches=*" + substring + "*"}

	out, _ := runGitCommand(args...)

	branch = strings.Split(out, "\n")[0]
	branch = strings.TrimSpace(branch)

	if os.Getenv(DRY_RUN_ENV) != "" {
		return "", false
	}

	return branch, branch != ""
}

func (p *Provider) CheckoutBranch(branch string) (err error) {
	args := []string{"checkout", branch}

	_, err = runGitCommand(args...)

	if err != nil {
		err = fmt.Errorf("failed to checkout the branch.\n\nDetails:\n%s", err)

		return
	}

	return
}

func (p *Provider) GetCommitsToPush(branch string) ([]string, error) {
	commits := []string{}

	args := []string{"log", "--pretty=format:'%h %s'", branch, "--not", "--remotes=origin"}

	out, err := runGitCommand(args...)
	if err != nil {
		return commits, err
	}

	if out != "" {
		commits = strings.Split(out, "\n")
	}

	return commits, nil
}

func (p *Provider) RemoteBranchExists(branch string) (exists bool) {
	args := []string{"show-ref", "--verify", "refs/remotes/origin/" + branch}

	_, err := runGitCommand(args...)

	return err == nil
}

func (p *Provider) CommitEmpty(message string) (err error) {
	signing := CommitSigningEnabled()

	args := []string{"commit", "--allow-empty", "-m", message}

	if signing {
		args = append(args, "-S")
	}

	_, err = runGitCommand(args...)

	if err != nil {
		err = fmt.Errorf("failed to commit.\n\nDetails:\n%s", err)

		return
	}

	return
}

func (p *Provider) PushBranch(branch string) (err error) {
	args := []string{"push", "-u", "origin", branch}

	_, err = runGitCommand(args...)

	if err != nil {
		err = fmt.Errorf("failed to push the branch.\n\nDetails:\n%s", err)

		return
	}

	return
}

func CommitSigningEnabled() bool {
	args := []string{"config", "--get", "commit.gpgsign"}

	stdout, err := runGitCommand(args...)

	return err == nil && strings.Contains(stdout, "true")
}

func (p *Provider) GetRepositoryRoot() (rootPath string, err error) {
	args := []string{"rev-parse", "--show-toplevel"}

	out, err := runGitCommand(args...)
	if err != nil {
		return "", fmt.Errorf("failed to get repository root: %w", err)
	}

	return strings.TrimSpace(out), nil
}
