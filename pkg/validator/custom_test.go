// SPDX-FileCopyrightText: 2023 INDITEX S.A
//
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"testing"

	govalidator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestUniqueMapValues(t *testing.T) {

	v := govalidator.New()
	v.RegisterValidation("uniqueMapValues", uniqueMapValues)

	t.Run("should return true if all values are unique with slice", func(t *testing.T) {
		tc := struct {
			M map[string][]string `validate:"uniqueMapValues"`
		}{
			M: map[string][]string{
				"foo": {"val1", "val2"},
				"bar": {"val3", "val4"},
			},
		}

		err := v.Struct(tc)
		assert.NoError(t, err)
	})

	t.Run("should return true if all values are unique with non slice", func(t *testing.T) {
		tc := struct {
			M map[string]string `validate:"uniqueMapValues"`
		}{
			M: map[string]string{
				"foo": "foo",
				"bar": "bar",
			},
		}

		err := v.Struct(tc)
		assert.NoError(t, err)
	})

	t.Run("Should return error if values are not unique with slice", func(t *testing.T) {
		tc := struct {
			M map[string][]string `validate:"uniqueMapValues"`
		}{
			M: map[string][]string{
				"foo": {"val1", "val2"},
				"bar": {"val1", "val4"},
			},
		}

		err := v.Struct(tc)
		assert.Error(t, err)
	})

	t.Run("Should return error if values are not unique with non slice", func(t *testing.T) {
		tc := struct {
			M map[string]string `validate:"uniqueMapValues"`
		}{
			M: map[string]string{
				"foo": "baz",
				"bar": "baz",
			},
		}

		err := v.Struct(tc)
		assert.Error(t, err)
	})

	t.Run("Should panic if not a map", func(t *testing.T) {
		tc := struct {
			M string `validate:"uniqueMapValues"`
		}{
			M: "not a map",
		}

		assert.Panics(t, func() {
			v.Struct(tc)
		})
	})

}
