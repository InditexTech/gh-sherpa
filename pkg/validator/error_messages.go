package validator

import (
	"bytes"
	"fmt"

	govalidator "github.com/go-playground/validator/v10"
)

const fallbackErrMessage = "Field validation for '%s' failed on the '%s' tag"

var validationErrorMessages = map[string]string{
	"required":           "Required field",
	"url":                "Must be a valid URL",
	"validIssueTypeKeys": "Keys must be a valid issue type. Check the documentation for the list of valid issue types",
	"uniqueMapValues":    "Values must be unique across all keys. Check the default values for possible collisions",
}
var validationErrorMessagesWithParam = map[string]string{
	"gte": "Must be greater than or equal to %s",
}

func getPrettyErrors(validationErrors govalidator.ValidationErrors) string {
	var buffer bytes.Buffer

	for _, fieldErr := range validationErrors {
		errKey := fieldErr.Namespace()
		errMsg, ok := validationErrorMessages[fieldErr.Tag()]
		if !ok {
			errMsg, ok = validationErrorMessagesWithParam[fieldErr.Tag()]
			if !ok {
				errMsg = fmt.Sprintf(fallbackErrMessage, fieldErr.Field(), fieldErr.Tag())
			} else {
				errMsg = fmt.Sprintf(errMsg, fieldErr.Param())
			}
		}

		buffer.WriteString(fmt.Sprintf("- %s: %s\n", errKey, errMsg))
	}

	return buffer.String()
}
