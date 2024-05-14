package domain

import (
	"errors"

	"github.com/InditexTech/gh-sherpa/internal/domain"
)

type FakeRepositoryProvider struct {
	Repository *domain.Repository
}

var _ domain.RepositoryProvider = (*FakeRepositoryProvider)(nil)

var ErrRepositoryNotFound = errors.New("repository not found")

func (f *FakeRepositoryProvider) GetRepository() (repo *domain.Repository, err error) {
	if f.Repository != nil {
		return f.Repository, nil
	}
	return nil, ErrRepositoryNotFound
}
