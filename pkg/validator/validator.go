package validator

import (
	"bytes"
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	govalidator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func init() {

	en := en.New()
	uni = ut.New(en, en)

	trans, _ = uni.GetTranslator("en")

	// https://pkg.go.dev/github.com/go-playground/validator/v10#readme-special-notes
	validate = govalidator.New(govalidator.WithRequiredStructEnabled())

	en_translations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterValidation("uniqueMapValues", uniqueMapValues)
	validate.RegisterTranslation("uniqueMapValues", trans, func(ut ut.Translator) error {
		return ut.Add("uniqueMapValues", "{0} must have unique values across all keys. Check the default values for collisions.", true)
	}, func(ut ut.Translator, fe govalidator.FieldError) string {
		t, _ := ut.T("uniqueMapValues", fe.Field())
		return t
	})

	validate.RegisterValidation("validIssueTypeKeys", validIssueTypeKeys)
	validate.RegisterTranslation("validIssueTypeKeys", trans, func(ut ut.Translator) error {
		return ut.Add("validIssueTypeKeys", "{0} must have valid GH Sherpa issue types as keys. Check the documentation for the issue types.", true)
	}, func(ut ut.Translator, fe govalidator.FieldError) string {
		t, _ := ut.T("validIssueTypeKeys", fe.Field())
		return t
	})

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
	uni      *ut.UniversalTranslator
	trans    ut.Translator
)

// Struct validates a struct
func Struct(s any) error {
	return handleValidationError(validate.Struct(s))
}

func TranslateError(err ValidationErrors) ValidationErrorsTranslations {
	return ValidationErrorsTranslations{err.Translate(trans)}
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
