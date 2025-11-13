package domain

type PullRequest struct {
	Title       string
	Number      int64
	State       string
	Closed      bool
	Url         string
	HeadRefName string
	BaseRefName string
	Labels      []Label
	Body        string
}
