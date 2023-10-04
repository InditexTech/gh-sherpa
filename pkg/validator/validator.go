package validator

import (
	"errors"

	govalidator "github.com/go-playground/validator/v10"
)

func init() {
	// https://pkg.go.dev/github.com/go-playground/validator/v10#readme-special-notes
	validate = govalidator.New(govalidator.WithRequiredStructEnabled())

	validate.RegisterValidation("uniqueMapValues", uniqueMapValues)
	validate.RegisterValidation("validIssueTypeKeys", validIssueTypeKeys)

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

	// https://github.com/go-playground/validator#error-return-value
	validationErrors := err.(govalidator.ValidationErrors)

	return errors.New(getPrettyErrors(validationErrors))
}
