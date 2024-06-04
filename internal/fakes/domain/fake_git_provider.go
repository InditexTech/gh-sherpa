package domain

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/InditexTech/gh-sherpa/internal/domain"
)

type FakeGitProvider struct {
	RemoteBranches        []string
	LocalBranches         []string
	CurrentBranch         string
	CommitsToPush         map[string][]string
	BranchWithCommitError []string
	BranchWithPushError   []string
}

var _ domain.GitProvider = (*FakeGitProvider)(nil)

func NewFakeGitProvider() *FakeGitProvider {
	return &FakeGitProvider{
		CurrentBranch: "main",
		RemoteBranches: []string{
			"main",
			"develop",
		},
		LocalBranches: []string{
			"main",
			"develop",
		},
		CommitsToPush:         map[string][]string{},
		BranchWithCommitError: []string{},
	}
}

func (f *FakeGitProvider) ResetLocalBranches() {
	f.LocalBranches = []string{"main", "develop"}
}

func (f *FakeGitProvider) ResetRemoteBranches() {
	f.RemoteBranches = []string{"main", "develop"}
}

func (f *FakeGitProvider) AddLocalBranches(branches ...string) {
	f.LocalBranches = append(f.LocalBranches, branches...)
}

func (f *FakeGitProvider) AddRemoteBranches(branches ...string) {
	f.RemoteBranches = append(f.RemoteBranches, branches...)
}

func (f *FakeGitProvider) BranchExists(branch string) bool {
	return slices.Contains(f.LocalBranches, branch)
}

func (f *FakeGitProvider) FetchBranchFromOrigin(branch string) (err error) {
	idx := slices.Index(f.RemoteBranches, branch)
	if idx == -1 {
		return fmt.Errorf("remote branch %s not found", branch)
	}
	return nil
}

func (f *FakeGitProvider) CheckoutNewBranchFromOrigin(branch string, base string) (err error) {
	idx := slices.Index(f.RemoteBranches, base)
	if idx == -1 {
		return fmt.Errorf("remote branch %s not found", base)
	}
	f.LocalBranches = append(f.LocalBranches, branch)
	f.CurrentBranch = branch
	return nil
}

var ErrGetCurrentBranch = errors.New("no current branch")

func (f *FakeGitProvider) GetCurrentBranch() (branch string, err error) {
	if f.CurrentBranch != "" {
		return f.CurrentBranch, nil
	}

	return "", ErrGetCurrentBranch
}

func (f *FakeGitProvider) BranchExistsContains(branch string) (name string, exists bool) {
	for _, b := range f.LocalBranches {
		if strings.Contains(b, branch) {
			return b, true
		}
	}
	for _, b := range f.RemoteBranches {
		if strings.Contains(b, branch) {
			return b, true
		}
	}
	return "", false
}

func (f *FakeGitProvider) CheckoutBranch(branch string) (err error) {
	if !slices.Contains(f.LocalBranches, branch) {
		return fmt.Errorf("local branch %s not found", branch)
	}
	f.CurrentBranch = branch
	return nil
}

var ErrGetCommitsToPush = errors.New("error getting commits to push")

func (f *FakeGitProvider) GetCommitsToPush(branch string) (commits []string, err error) {
	if slices.Contains(f.BranchWithCommitError, branch) {
		return commits, ErrGetCommitsToPush
	}

	commits, ok := f.CommitsToPush[branch]
	if !ok {
		commits = []string{}
	}

	return commits, nil
}

func (f *FakeGitProvider) RemoteBranchExists(branch string) (exists bool) {
	return slices.Contains(f.RemoteBranches, branch)
}

func (f *FakeGitProvider) CommitEmpty(message string) (err error) {
	currentCommits, ok := f.CommitsToPush[f.CurrentBranch]
	if !ok {
		currentCommits = []string{}
	}
	currentCommits = append(currentCommits, message)
	f.CommitsToPush[f.CurrentBranch] = currentCommits
	return nil
}

var ErrPushBranch = errors.New("error pushing branch")

func (f *FakeGitProvider) PushBranch(branch string) (err error) {
	if slices.Contains(f.BranchWithPushError, branch) {
		return ErrPushBranch
	}
	if !slices.Contains(f.LocalBranches, branch) {
		return fmt.Errorf("local branch %s not found", branch)
	}
	if slices.Contains(f.RemoteBranches, branch) {
		return fmt.Errorf("remote branch %s already exists", branch)
	}
	f.RemoteBranches = append(f.RemoteBranches, branch)
	return nil
}

func (f *FakeGitProvider) IsCurrentBranch(branch string) bool {
	return f.CurrentBranch == branch
}
