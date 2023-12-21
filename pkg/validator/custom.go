// SPDX-FileCopyrightText: 2023 INDITEX S.A
//
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	govalidator "github.com/go-playground/validator/v10"
)

// uniqueMapValues validates that the values of a map are unique across all keys.
// It will only iterate over the second level of the map if the value is a slice.
func uniqueMapValues(fl govalidator.FieldLevel) bool {
	field := fl.Field()
	if field.Type().Kind() != reflect.Map {
		panic(fmt.Sprintf("Invalid type %T. uniqueMapValues only works with map", field.Interface()))
	}

	mapKeys := field.MapKeys()
	seen := make(map[any]bool)
	for _, key := range mapKeys {
		value := field.MapIndex(key)
		switch value.Kind() {
		case reflect.Slice:
			for i := 0; i < value.Len(); i++ {
				element := value.Index(i).Interface()
				if _, ok := seen[element]; ok {
					return false
				}
				seen[element] = true
			}
		default:
			element := value.Interface()
			if _, ok := seen[element]; ok {
				return false
			}
			seen[element] = true
		}
	}

	return true
}

// validIssueTypeKeys validates that the keys of a map are valid issue types.
func validIssueTypeKeys(fl govalidator.FieldLevel) bool {
	field := fl.Field()
	if field.Type().Kind() != reflect.Map {
		panic(fmt.Sprintf("Invalid type %T. uniqueMapValues only works with map", field.Interface()))
	}

	validIssueTypes := issue_types.GetValidIssueTypes()
	mapKeys := field.MapKeys()
	for _, key := range mapKeys {
		issueType, ok := key.Interface().(issue_types.IssueType)
		if !ok {
			return false
		}
		if !slices.Contains(validIssueTypes, issueType) {
			return false
		}
	}

	return true
}
