package domain

import (
	"errors"

	"github.com/InditexTech/gh-sherpa/internal/domain"
)

type FakeUserInteractionProvider struct {
	// TODO: Map this :)
}

var _ domain.UserInteractionProvider = (*FakeUserInteractionProvider)(nil)

func NewFakeUserInteractionProvider() FakeUserInteractionProvider {
	return FakeUserInteractionProvider{}
}

func (f *FakeUserInteractionProvider) AskUserForConfirmation(msg string, defaultAnswer bool) (answer bool, err error) {
	return false, errors.New("not implemented")
}

func (f *FakeUserInteractionProvider) SelectOrInputPrompt(message string, validValues []string, variable *string, required bool) error {
	return errors.New("not implemented")
}

func (f *FakeUserInteractionProvider) SelectOrInput(name string, validValues []string, variable *string, required bool) error {
	return errors.New("not implemented")
}
