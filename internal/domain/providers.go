// SPDX-FileCopyrightText: 2023 INDITEX S.A
//
// SPDX-License-Identifier: Apache-2.0

package domain

type GhCli interface {
	Execute(result any, args []string) (err error)
}

type RepositoryProvider interface {
	GetRepository() (repo *Repository, err error)
}

type PullRequestProvider interface {
	GetPullRequestForBranch(string) (pullRequest *PullRequest, err error)
	CreatePullRequest(title string, body string, baseBranch string, headBranch string, draft bool, labels []string) (prUrl string, err error)
}

type UserInteractionProvider interface {
	AskUserForConfirmation(msg string, defaultAnswer bool) (answer bool, err error)
	SelectOrInputPrompt(message string, validValues []string, variable *string, required bool) error
	SelectOrInput(name string, validValues []string, variable *string, required bool) error
}

type GitProvider interface {
	BranchExists(branch string) bool
	FetchBranchFromOrigin(branch string) (err error)
	CheckoutNewBranchFromOrigin(branch string, base string) (err error)
	GetCurrentBranch() (branchName string, err error)
	BranchExistsContains(branch string) (name string, exists bool)
	CheckoutBranch(branch string) (err error)
	GetCommitsToPush(branch string) (commits []string, err error)
	RemoteBranchExists(branch string) (exists bool)
	CommitEmpty(message string) (err error)
	PushBranch(branch string) (err error)
}

type BranchProvider interface {
	GetBranchName(issueTracker IssueTracker, issueIdentifier string, repo Repository) (branchName string, err error)
}
