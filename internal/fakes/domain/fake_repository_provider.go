package domain

import (
	"errors"

	"github.com/InditexTech/gh-sherpa/internal/domain"
)

type FakeRepositoryProvider struct {
	Repository *domain.Repository
}

var _ domain.RepositoryProvider = (*FakeRepositoryProvider)(nil)

func NewFakeRepositoryProvider() *FakeRepositoryProvider {
	return &FakeRepositoryProvider{
		Repository: &domain.Repository{
			Name:             "gh-sherpa-test-repo",
			Owner:            "inditextech",
			NameWithOwner:    "inditextech/gh-sherpa-test-repo",
			DefaultBranchRef: "main",
		},
	}
}

var ErrRepositoryNotFound = errors.New("repository not found")

func (f *FakeRepositoryProvider) GetRepository() (repo *domain.Repository, err error) {
	if f.Repository != nil {
		return f.Repository, nil
	}
	return nil, ErrRepositoryNotFound
}
