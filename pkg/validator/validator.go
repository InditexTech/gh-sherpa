package validator

import (
	govalidator "github.com/go-playground/validator/v10"
)

func init() {
	// https://pkg.go.dev/github.com/go-playground/validator/v10#readme-special-notes
	validate = govalidator.New(govalidator.WithRequiredStructEnabled())

	validate.RegisterValidation("uniqueMapValues", uniqueMapValues)
	validate.RegisterValidation("validIssueTypeKeys", validIssueTypeKeys)
}

var validate *govalidator.Validate

// Struct validates a struct
func Struct(s any) error {
	return validate.Struct(s)
}
