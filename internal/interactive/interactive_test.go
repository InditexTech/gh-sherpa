package interactive_test

import (
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestUserInteractionProvider_AskUserForConfirmation(t *testing.T) {
	program := bubbletea.NewProgram(nil)
	userInteraction := &UserInteractionProvider{program: program}

	t.Run("should return true when user confirms", func(t *testing.T) {
		program.Send(bubbletea.KeyMsg{Type: bubbletea.KeyEnter})
		confirmed, err := userInteraction.AskUserForConfirmation("Do you want to continue?", true)
		assert.NoError(t, err)
		assert.True(t, confirmed)
	})

	t.Run("should return false when user does not confirm", func(t *testing.T) {
		program.Send(bubbletea.KeyMsg{Type: bubbletea.KeyEsc})
		confirmed, err := userInteraction.AskUserForConfirmation("Do you want to continue?", true)
		assert.NoError(t, err)
		assert.False(t, confirmed)
	})
}

func TestUserInteractionProvider_SelectOrInput(t *testing.T) {
	program := bubbletea.NewProgram(nil)
	userInteraction := &UserInteractionProvider{program: program}

	t.Run("should return selected value from options", func(t *testing.T) {
		program.Send(bubbletea.KeyMsg{Type: bubbletea.KeyEnter})
		var selectedValue string
		err := userInteraction.SelectOrInput("Select a value", []string{"option1", "option2"}, &selectedValue, true)
		assert.NoError(t, err)
		assert.Equal(t, "option1", selectedValue)
	})

	t.Run("should return input value when no options are provided", func(t *testing.T) {
		program.Send(bubbletea.KeyMsg{Type: bubbletea.KeyEnter})
		var inputValue string
		err := userInteraction.SelectOrInput("Enter a value", []string{}, &inputValue, true)
		assert.NoError(t, err)
		assert.Equal(t, "input value", inputValue)
	})
}

func TestUserInteractionProvider_SelectOrInputPrompt(t *testing.T) {
	program := bubbletea.NewProgram(nil)
	userInteraction := &UserInteractionProvider{program: program}

	t.Run("should return selected value from options with prompt", func(t *testing.T) {
		program.Send(bubbletea.KeyMsg{Type: bubbletea.KeyEnter})
		var selectedValue string
		err := userInteraction.SelectOrInputPrompt("Select a value", []string{"option1", "option2"}, &selectedValue, true)
		assert.NoError(t, err)
		assert.Equal(t, "option1", selectedValue)
	})

	t.Run("should return input value when no options are provided with prompt", func(t *testing.T) {
		program.Send(bubbletea.KeyMsg{Type: bubbletea.KeyEnter})
		var inputValue string
		err := userInteraction.SelectOrInputPrompt("Enter a value", []string{}, &inputValue, true)
		assert.NoError(t, err)
		assert.Equal(t, "input value", inputValue)
	})
}
