// SPDX-FileCopyrightText: 2023 INDITEX S.A
//
// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/stretchr/testify/mock"
)

// UnsetExpectedCall unsets the expected call for the given method.
// `method` should be passed as the method value from the mock.
//
// Example:
//
//	mocks.UnsetExpectedCall(&mockedProvider.Mock, mockedProvider.MockedMethod)
func UnsetExpectedCall(m *mock.Mock, method any) {
	methodName := runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name()
	for _, c := range m.ExpectedCalls {
		// We use this to avoid matching wrong methods.
		// More info -> https://stackoverflow.com/a/33325345
		if strings.Contains(methodName, fmt.Sprintf(".%s-fm", c.Method)) {
			c.Unset()
		}
	}
}
