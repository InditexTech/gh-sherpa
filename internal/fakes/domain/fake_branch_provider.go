package domain

import (
	"errors"

	"github.com/InditexTech/gh-sherpa/internal/domain"
)

type FakeBranchProvider struct {
	BranchName string
}

var _ domain.BranchProvider = (*FakeBranchProvider)(nil)

func NewFakeBranchProvider() *FakeBranchProvider {
	return &FakeBranchProvider{}
}

func (f *FakeBranchProvider) SetBranchName(branchName string) {
	f.BranchName = branchName
}

var ErrGetBranchName = errors.New("get branch name error")

func (f *FakeBranchProvider) GetBranchName(_ domain.Issue, _ domain.Repository) (branchName string, err error) {

	if f.BranchName == "" {
		return "", ErrGetBranchName
	}

	return f.BranchName, nil
}
