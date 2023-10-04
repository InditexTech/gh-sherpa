package validator

import (
	"bytes"
	"fmt"

	govalidator "github.com/go-playground/validator/v10"
)

func init() {
	// https://pkg.go.dev/github.com/go-playground/validator/v10#readme-special-notes
	validate = govalidator.New(govalidator.WithRequiredStructEnabled())

	validate.RegisterValidation("uniqueMapValues", uniqueMapValues)
	validate.RegisterValidation("validIssueTypeKeys", validIssueTypeKeys)

}

type ValidationErrors struct {
	govalidator.ValidationErrors
}

type ValidationErrorsTranslations struct {
	govalidator.ValidationErrorsTranslations
}

func (vet ValidationErrorsTranslations) PrettyPrint() string {
	var buffer bytes.Buffer
	for k, v := range vet.ValidationErrorsTranslations {
		buffer.WriteString(fmt.Sprintf("- %s: %s\n", k, v))
	}
	return buffer.String()
}

var (
	validate *govalidator.Validate
)

// Struct validates a struct
func Struct(s any) error {
	return handleValidationError(validate.Struct(s))
}

func handleValidationError(err error) error {
	if err == nil {
		return nil
	}

	switch err := err.(type) {
	case govalidator.ValidationErrors:
		return ValidationErrors{err}
	default:
		return err
	}
}
