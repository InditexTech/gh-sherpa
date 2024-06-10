package interactive

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/InditexTech/gh-sherpa/internal/domain"
)

var ErrOpCanceled error = fmt.Errorf("operation canceled by the user")

type UserInteractionProvider struct{}

func (u UserInteractionProvider) AskUserForConfirmation(message string, defaultValue bool) (bool, error) {
	return AskUserForConfirmation(message, defaultValue)
}

// SelectOrInput Prompt a select if valid values are provided, prompt a simple input otherwise.
// Checks that the user has selected/written a value if required.
func (u UserInteractionProvider) SelectOrInput(name string, validValues []string, variable *string, required bool) error {
	var opts []survey.AskOpt
	if required {
		opts = append(opts, survey.WithValidator(survey.Required))
	}

	var prompt survey.Prompt
	if len(validValues) != 0 {
		prompt = &survey.Select{
			Message:  fmt.Sprintf("Please, select one %v:", name),
			Options:  validValues,
			Default:  *variable,
			PageSize: 10,
		}
	} else {
		prompt = &survey.Input{
			Message: fmt.Sprintf("Please, write one %v:", name),
			Default: *variable,
		}
	}

	err := survey.AskOne(prompt, variable, opts...)
	return handleSurveyError(err)
}

// SelectOrInputPrompt Prompt a select if valid values are provided, prompt a simple input otherwise.
// Checks that the user has selected/written a value if required.
func (u UserInteractionProvider) SelectOrInputPrompt(message string, validValues []string, variable *string, required bool) error {
	var opts []survey.AskOpt
	if required {
		opts = append(opts, survey.WithValidator(survey.Required))
	}

	var prompt survey.Prompt
	if len(validValues) != 0 {
		prompt = &survey.Select{
			Message:  message,
			Options:  validValues,
			Default:  *variable,
			PageSize: 10,
		}
	} else {
		prompt = &survey.Input{
			Message: message,
			Default: *variable,
		}
	}

	err := survey.AskOne(prompt, variable, opts...)
	return handleSurveyError(err)
}

// TODO: Do not use "kind/*" here, use the actual config to retrieve the label
func GetPromptMessageBranchType(branchType string, issueTrackerType domain.IssueTrackerType) string {
	if issueTrackerType == domain.IssueTrackerTypeJira {
		return fmt.Sprintf("Issue type '%s' found. What type of branch name do you want to create?", branchType)
	} else {
		return fmt.Sprintf("Label 'kind/%s' found. What type of branch name do you want to create?", branchType)
	}
}

func SelectPrompt(message string, options []string, defaultOption string, selected *string, required bool) (err error) {
	var opts []survey.AskOpt
	if required {
		opts = append(opts, survey.WithValidator(survey.Required))
	}

	var prompt = &survey.Select{
		Message:  message,
		Options:  options,
		Default:  defaultOption,
		PageSize: 10,
	}

	err = survey.AskOne(prompt, selected, opts...)
	return handleSurveyError(err)
}

func InputPrompt(message string, defaultOption string, selected *string, password bool, required bool) (err error) {
	var opts []survey.AskOpt
	if required {
		opts = append(opts, survey.WithValidator(survey.Required))
	}

	var prompt survey.Prompt = &survey.Input{
		Message: message,
		Default: defaultOption,
	}

	if password {
		prompt = &survey.Password{
			Message: message,
		}
	}

	err = survey.AskOne(prompt, selected, opts...)
	return handleSurveyError(err)
}

func AskUserForConfirmation(promptMessage string, defaultOption bool) (yes bool, err error) {
	yes = false
	prompt := &survey.Confirm{
		Message: promptMessage,
		Default: defaultOption,
	}
	err = survey.AskOne(prompt, &yes)
	err = handleSurveyError(err)

	return
}

func handleSurveyError(err error) error {
	if err == terminal.InterruptErr {
		return ErrOpCanceled
	}

	return err
}

func AskUserForJiraInputs(defaultHost string) (host, pat, username, password, name string, err error) {

	host = defaultHost
	err = InputPrompt("Enter Jira Host", host, &host, false, true)

	if err != nil {
		err = handleSurveyError(err)
		return
	}

	hasAToken, err := AskUserForConfirmation("Do you have a valid PAT to use?", false)

	if err != nil {
		err = handleSurveyError(err)
		return
	}

	if hasAToken {
		if err = InputPrompt("Enter Jira PAT", "", &pat, true, true); err != nil {
			err = handleSurveyError(err)
			return
		}

	} else {
		if err = InputPrompt("Enter Jira Username", "", &username, false, true); err != nil {
			err = handleSurveyError(err)
			return
		}

		if err = InputPrompt("Enter Jira Password", "", &password, true, true); err != nil {
			err = handleSurveyError(err)
			return
		}

		if err = InputPrompt("Enter Jira PAT name", "", &name, false, true); err != nil {
			err = handleSurveyError(err)
			return
		}
	}

	return host, pat, username, password, name, nil
}
