package domain

import "github.com/InditexTech/gh-sherpa/internal/domain"

type FakeGhCli struct{}

var _ domain.GhCli = (*FakeGhCli)(nil)

func NewFakeGhCli() FakeGhCli {
	return FakeGhCli{}
}

func (f FakeGhCli) Execute(result any, args []string) (err error) {
	return nil
}
