package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/InditexTech/gh-sherpa/pkg/validator"
)

const fallbackErrMessage = "Field validation for '%s' failed on the '%s' tag"

var validationErrorMessages = map[string]string{
	"required":           "Required field",
	"url":                "Must be a valid URL",
	"validIssueTypeKeys": "Key must be a valid issue type. Check the documentation for the list of valid issue types",
	"uniqueMapValues":    "Values must be unique across all keys. Check the default values for possible collisions",
}

func getPrettyErrors(validationErrors validator.ValidationErrors) string {
	var buffer bytes.Buffer

	for _, fieldErr := range validationErrors.ValidationErrors {
		errKey := formatErrorKey(fieldErr.Namespace())
		errMsg, ok := validationErrorMessages[fieldErr.Tag()]
		if !ok {
			errMsg = fmt.Sprintf(fallbackErrMessage, fieldErr.Field(), fieldErr.Tag())
		}
		buffer.WriteString(fmt.Sprintf("- %s: %s\n", errKey, errMsg))
	}

	return buffer.String()
}

// formatErrorKey formats the error key to remove the "configuration." prefix and lowercase it
func formatErrorKey(key string) string {
	return strings.Join(strings.Split(strings.ToLower(key), "configuration.")[1:], "")
}
