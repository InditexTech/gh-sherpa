package validator

import (
	"fmt"
	"reflect"

	"github.com/InditexTech/gh-sherpa/internal/domain/issue_types"
	govalidator "github.com/go-playground/validator/v10"
	"golang.org/x/exp/slices"
)

// uniqueMapValues validates that the values of a map are unique across all keys.
// It will only iterate over the second level of the map if the value is a slice.
func uniqueMapValues(fl govalidator.FieldLevel) bool {
	fieldType := fl.Field().Type()
	fieldValue := fl.Field()
	kind := fieldType.Kind()
	if kind != reflect.Map {
		panic(fmt.Sprintf("Invalid type %T. uniqueMapValues only works with map", fieldValue.Interface()))
	}

	mapKeys := fieldValue.MapKeys()
	seen := make(map[any]bool)
	for _, k := range mapKeys {
		v := fieldValue.MapIndex(k)
		if v.Kind() == reflect.Slice {
			for i := 0; i < v.Len(); i++ {
				s := v.Index(i).Interface()
				if _, ok := seen[s]; ok {
					return false
				}
				seen[s] = true
			}
		} else {
			s := v.Interface()
			if _, ok := seen[s]; ok {
				return false
			}
			seen[s] = true
		}
	}

	return true
}

// validIssueTypeKeys validates that the keys of a map are valid issue types.
func validIssueTypeKeys(fl govalidator.FieldLevel) bool {
	fieldType := fl.Field().Type()
	fieldValue := fl.Field()
	kind := fieldType.Kind()
	if kind != reflect.Map {
		panic(fmt.Sprintf("Invalid type %T. uniqueMapValues only works with map", fieldValue.Interface()))
	}

	validIssueTypes := issue_types.GetValidIssueTypes()
	mapKeys := fieldValue.MapKeys()
	for _, k := range mapKeys {
		it, ok := k.Interface().(issue_types.IssueType)
		if !ok {
			return false
		}
		if !slices.Contains(validIssueTypes, it) {
			return false
		}
	}

	return true
}
