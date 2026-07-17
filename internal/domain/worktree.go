package domain

type Worktree struct {
	Path     string
	Branch   string
	Commit   string
	Prunable bool
}
