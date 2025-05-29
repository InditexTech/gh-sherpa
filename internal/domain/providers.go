package domain

type RepositoryProvider interface {
	GetRepository() (repo *Repository, err error)
}

type PullRequestProvider interface {
	GetPullRequestForBranch(branch string) (pullRequest *PullRequest, err error)
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
	FindBranch(substring string) (branch string, exists bool)
	CheckoutBranch(branch string) (err error)
	GetCommitsToPush(branch string) (commits []string, err error)
	RemoteBranchExists(branch string) (exists bool)
	CommitEmpty(message string) (err error)
	PushBranch(branch string) (err error)
	GetRepositoryRoot() (rootPath string, err error)
}

type BranchProvider interface {
	GetBranchName(issue Issue, repo Repository) (branchName string, err error)
}
