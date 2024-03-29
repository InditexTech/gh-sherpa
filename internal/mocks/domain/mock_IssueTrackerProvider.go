// Code generated by mockery v2.32.4. DO NOT EDIT.

package domain

import (
	domain "github.com/InditexTech/gh-sherpa/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockIssueTrackerProvider is an autogenerated mock type for the IssueTrackerProvider type
type MockIssueTrackerProvider struct {
	mock.Mock
}

type MockIssueTrackerProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIssueTrackerProvider) EXPECT() *MockIssueTrackerProvider_Expecter {
	return &MockIssueTrackerProvider_Expecter{mock: &_m.Mock}
}

// GetIssueTracker provides a mock function with given fields: identifier
func (_m *MockIssueTrackerProvider) GetIssueTracker(identifier string) (domain.IssueTracker, error) {
	ret := _m.Called(identifier)

	var r0 domain.IssueTracker
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (domain.IssueTracker, error)); ok {
		return rf(identifier)
	}
	if rf, ok := ret.Get(0).(func(string) domain.IssueTracker); ok {
		r0 = rf(identifier)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.IssueTracker)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(identifier)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIssueTrackerProvider_GetIssueTracker_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIssueTracker'
type MockIssueTrackerProvider_GetIssueTracker_Call struct {
	*mock.Call
}

// GetIssueTracker is a helper method to define mock.On call
//   - identifier string
func (_e *MockIssueTrackerProvider_Expecter) GetIssueTracker(identifier interface{}) *MockIssueTrackerProvider_GetIssueTracker_Call {
	return &MockIssueTrackerProvider_GetIssueTracker_Call{Call: _e.mock.On("GetIssueTracker", identifier)}
}

func (_c *MockIssueTrackerProvider_GetIssueTracker_Call) Run(run func(identifier string)) *MockIssueTrackerProvider_GetIssueTracker_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockIssueTrackerProvider_GetIssueTracker_Call) Return(issueTracker domain.IssueTracker, err error) *MockIssueTrackerProvider_GetIssueTracker_Call {
	_c.Call.Return(issueTracker, err)
	return _c
}

func (_c *MockIssueTrackerProvider_GetIssueTracker_Call) RunAndReturn(run func(string) (domain.IssueTracker, error)) *MockIssueTrackerProvider_GetIssueTracker_Call {
	_c.Call.Return(run)
	return _c
}

// ParseIssueId provides a mock function with given fields: identifier
func (_m *MockIssueTrackerProvider) ParseIssueId(identifier string) string {
	ret := _m.Called(identifier)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(identifier)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockIssueTrackerProvider_ParseIssueId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ParseIssueId'
type MockIssueTrackerProvider_ParseIssueId_Call struct {
	*mock.Call
}

// ParseIssueId is a helper method to define mock.On call
//   - identifier string
func (_e *MockIssueTrackerProvider_Expecter) ParseIssueId(identifier interface{}) *MockIssueTrackerProvider_ParseIssueId_Call {
	return &MockIssueTrackerProvider_ParseIssueId_Call{Call: _e.mock.On("ParseIssueId", identifier)}
}

func (_c *MockIssueTrackerProvider_ParseIssueId_Call) Run(run func(identifier string)) *MockIssueTrackerProvider_ParseIssueId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockIssueTrackerProvider_ParseIssueId_Call) Return(issueId string) *MockIssueTrackerProvider_ParseIssueId_Call {
	_c.Call.Return(issueId)
	return _c
}

func (_c *MockIssueTrackerProvider_ParseIssueId_Call) RunAndReturn(run func(string) string) *MockIssueTrackerProvider_ParseIssueId_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIssueTrackerProvider creates a new instance of MockIssueTrackerProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIssueTrackerProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIssueTrackerProvider {
	mock := &MockIssueTrackerProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
