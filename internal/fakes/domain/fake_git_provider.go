package domain

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/InditexTech/gh-sherpa/internal/domain"
)

type FakeGitProvider struct {
	RemoteBranches []string
	LocalBranches  []string
	CurrentBranch  string
	CommitsToPush  map[string][]string
}

var _ domain.GitProvider = (*FakeGitProvider)(nil)

func NewFakeGitProvider() *FakeGitProvider {
	return &FakeGitProvider{
		RemoteBranches: []string{
			"main",
			"develop",
			"feature/GH-1-sample-issue",
			"feature/GH-2-remote-branch",
			"feature/PROJECTKEY-1-sample-issue",
			"feature/PROJECTKEY-2-remote-branch",
		},
		LocalBranches: []string{
			"main",
			"develop",
			"feature/GH-1-sample-issue",
			"feature/PROJECTKEY-1-sample-issue",
			"feature/GH-3-local-branch",
			"feature/PROJECTKEY-3-local-branch",
		},
		CommitsToPush: map[string][]string{
			"main":                              {},
			"develop":                           {},
			"feature/GH-1-sample-issue":         {},
			"feature/PROJECTKEY-1-sample-issue": {},
			"feature/GH-3-local-branch":         {},
			"feature/PROJECTKEY-3-local-branch": {},
		},
		CurrentBranch: "main",
	}
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
	commits, ok := f.CommitsToPush[branch]
	if !ok {
		return []string{}, ErrGetCommitsToPush
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

func (f *FakeGitProvider) PushBranch(branch string) (err error) {
	if !slices.Contains(f.LocalBranches, branch) {
		return fmt.Errorf("local branch %s not found", branch)
	}
	if slices.Contains(f.RemoteBranches, branch) {
		return fmt.Errorf("remote branch %s already exists", branch)
	}
	f.RemoteBranches = append(f.RemoteBranches, branch)
	return nil
}
