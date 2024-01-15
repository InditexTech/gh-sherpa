// Code generated by mockery v2.32.4. DO NOT EDIT.

package domain

import (
	domain "github.com/InditexTech/gh-sherpa/internal/domain"
	issue_types "github.com/InditexTech/gh-sherpa/internal/domain/issue_types"

	mock "github.com/stretchr/testify/mock"
)

// MockIssueTracker is an autogenerated mock type for the IssueTracker type
type MockIssueTracker struct {
	mock.Mock
}

type MockIssueTracker_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIssueTracker) EXPECT() *MockIssueTracker_Expecter {
	return &MockIssueTracker_Expecter{mock: &_m.Mock}
}

// FormatIssueId provides a mock function with given fields: issueId
func (_m *MockIssueTracker) FormatIssueId(issueId string) string {
	ret := _m.Called(issueId)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(issueId)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockIssueTracker_FormatIssueId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FormatIssueId'
type MockIssueTracker_FormatIssueId_Call struct {
	*mock.Call
}

// FormatIssueId is a helper method to define mock.On call
//   - issueId string
func (_e *MockIssueTracker_Expecter) FormatIssueId(issueId interface{}) *MockIssueTracker_FormatIssueId_Call {
	return &MockIssueTracker_FormatIssueId_Call{Call: _e.mock.On("FormatIssueId", issueId)}
}

func (_c *MockIssueTracker_FormatIssueId_Call) Run(run func(issueId string)) *MockIssueTracker_FormatIssueId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockIssueTracker_FormatIssueId_Call) Return(formattedIssueId string) *MockIssueTracker_FormatIssueId_Call {
	_c.Call.Return(formattedIssueId)
	return _c
}

func (_c *MockIssueTracker_FormatIssueId_Call) RunAndReturn(run func(string) string) *MockIssueTracker_FormatIssueId_Call {
	_c.Call.Return(run)
	return _c
}

// GetIssue provides a mock function with given fields: identifier
func (_m *MockIssueTracker) GetIssue(identifier string) (domain.Issue, error) {
	ret := _m.Called(identifier)

	var r0 domain.Issue
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (domain.Issue, error)); ok {
		return rf(identifier)
	}
	if rf, ok := ret.Get(0).(func(string) domain.Issue); ok {
		r0 = rf(identifier)
	} else {
		r0 = ret.Get(0).(domain.Issue)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(identifier)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIssueTracker_GetIssue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIssue'
type MockIssueTracker_GetIssue_Call struct {
	*mock.Call
}

// GetIssue is a helper method to define mock.On call
//   - identifier string
func (_e *MockIssueTracker_Expecter) GetIssue(identifier interface{}) *MockIssueTracker_GetIssue_Call {
	return &MockIssueTracker_GetIssue_Call{Call: _e.mock.On("GetIssue", identifier)}
}

func (_c *MockIssueTracker_GetIssue_Call) Run(run func(identifier string)) *MockIssueTracker_GetIssue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockIssueTracker_GetIssue_Call) Return(issue domain.Issue, err error) *MockIssueTracker_GetIssue_Call {
	_c.Call.Return(issue, err)
	return _c
}

func (_c *MockIssueTracker_GetIssue_Call) RunAndReturn(run func(string) (domain.Issue, error)) *MockIssueTracker_GetIssue_Call {
	_c.Call.Return(run)
	return _c
}

// GetIssueTrackerType provides a mock function with given fields:
func (_m *MockIssueTracker) GetIssueTrackerType() domain.IssueTrackerType {
	ret := _m.Called()

	var r0 domain.IssueTrackerType
	if rf, ok := ret.Get(0).(func() domain.IssueTrackerType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.IssueTrackerType)
	}

	return r0
}

// MockIssueTracker_GetIssueTrackerType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIssueTrackerType'
type MockIssueTracker_GetIssueTrackerType_Call struct {
	*mock.Call
}

// GetIssueTrackerType is a helper method to define mock.On call
func (_e *MockIssueTracker_Expecter) GetIssueTrackerType() *MockIssueTracker_GetIssueTrackerType_Call {
	return &MockIssueTracker_GetIssueTrackerType_Call{Call: _e.mock.On("GetIssueTrackerType")}
}

func (_c *MockIssueTracker_GetIssueTrackerType_Call) Run(run func()) *MockIssueTracker_GetIssueTrackerType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockIssueTracker_GetIssueTrackerType_Call) Return(_a0 domain.IssueTrackerType) *MockIssueTracker_GetIssueTrackerType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIssueTracker_GetIssueTrackerType_Call) RunAndReturn(run func() domain.IssueTrackerType) *MockIssueTracker_GetIssueTrackerType_Call {
	_c.Call.Return(run)
	return _c
}

// GetIssueType provides a mock function with given fields: issue
func (_m *MockIssueTracker) GetIssueType(issue domain.Issue) (issue_types.IssueType, error) {
	ret := _m.Called(issue)

	var r0 issue_types.IssueType
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.Issue) (issue_types.IssueType, error)); ok {
		return rf(issue)
	}
	if rf, ok := ret.Get(0).(func(domain.Issue) issue_types.IssueType); ok {
		r0 = rf(issue)
	} else {
		r0 = ret.Get(0).(issue_types.IssueType)
	}

	if rf, ok := ret.Get(1).(func(domain.Issue) error); ok {
		r1 = rf(issue)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIssueTracker_GetIssueType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIssueType'
type MockIssueTracker_GetIssueType_Call struct {
	*mock.Call
}

// GetIssueType is a helper method to define mock.On call
//   - issue domain.Issue
func (_e *MockIssueTracker_Expecter) GetIssueType(issue interface{}) *MockIssueTracker_GetIssueType_Call {
	return &MockIssueTracker_GetIssueType_Call{Call: _e.mock.On("GetIssueType", issue)}
}

func (_c *MockIssueTracker_GetIssueType_Call) Run(run func(issue domain.Issue)) *MockIssueTracker_GetIssueType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(domain.Issue))
	})
	return _c
}

func (_c *MockIssueTracker_GetIssueType_Call) Return(issueType issue_types.IssueType, err error) *MockIssueTracker_GetIssueType_Call {
	_c.Call.Return(issueType, err)
	return _c
}

func (_c *MockIssueTracker_GetIssueType_Call) RunAndReturn(run func(domain.Issue) (issue_types.IssueType, error)) *MockIssueTracker_GetIssueType_Call {
	_c.Call.Return(run)
	return _c
}

// GetIssueTypeLabel provides a mock function with given fields: issue
func (_m *MockIssueTracker) GetIssueTypeLabel(issue domain.Issue) (string, error) {
	ret := _m.Called(issue)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.Issue) (string, error)); ok {
		return rf(issue)
	}
	if rf, ok := ret.Get(0).(func(domain.Issue) string); ok {
		r0 = rf(issue)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(domain.Issue) error); ok {
		r1 = rf(issue)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIssueTracker_GetIssueTypeLabel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIssueTypeLabel'
type MockIssueTracker_GetIssueTypeLabel_Call struct {
	*mock.Call
}

// GetIssueTypeLabel is a helper method to define mock.On call
//   - issue domain.Issue
func (_e *MockIssueTracker_Expecter) GetIssueTypeLabel(issue interface{}) *MockIssueTracker_GetIssueTypeLabel_Call {
	return &MockIssueTracker_GetIssueTypeLabel_Call{Call: _e.mock.On("GetIssueTypeLabel", issue)}
}

func (_c *MockIssueTracker_GetIssueTypeLabel_Call) Run(run func(issue domain.Issue)) *MockIssueTracker_GetIssueTypeLabel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(domain.Issue))
	})
	return _c
}

func (_c *MockIssueTracker_GetIssueTypeLabel_Call) Return(_a0 string, _a1 error) *MockIssueTracker_GetIssueTypeLabel_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIssueTracker_GetIssueTypeLabel_Call) RunAndReturn(run func(domain.Issue) (string, error)) *MockIssueTracker_GetIssueTypeLabel_Call {
	_c.Call.Return(run)
	return _c
}

// IdentifyIssue provides a mock function with given fields: identifier
func (_m *MockIssueTracker) IdentifyIssue(identifier string) bool {
	ret := _m.Called(identifier)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(identifier)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockIssueTracker_IdentifyIssue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IdentifyIssue'
type MockIssueTracker_IdentifyIssue_Call struct {
	*mock.Call
}

// IdentifyIssue is a helper method to define mock.On call
//   - identifier string
func (_e *MockIssueTracker_Expecter) IdentifyIssue(identifier interface{}) *MockIssueTracker_IdentifyIssue_Call {
	return &MockIssueTracker_IdentifyIssue_Call{Call: _e.mock.On("IdentifyIssue", identifier)}
}

func (_c *MockIssueTracker_IdentifyIssue_Call) Run(run func(identifier string)) *MockIssueTracker_IdentifyIssue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockIssueTracker_IdentifyIssue_Call) Return(_a0 bool) *MockIssueTracker_IdentifyIssue_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIssueTracker_IdentifyIssue_Call) RunAndReturn(run func(string) bool) *MockIssueTracker_IdentifyIssue_Call {
	_c.Call.Return(run)
	return _c
}

// ParseRawIssueId provides a mock function with given fields: identifier
func (_m *MockIssueTracker) ParseRawIssueId(identifier string) string {
	ret := _m.Called(identifier)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(identifier)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockIssueTracker_ParseRawIssueId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ParseRawIssueId'
type MockIssueTracker_ParseRawIssueId_Call struct {
	*mock.Call
}

// ParseRawIssueId is a helper method to define mock.On call
//   - identifier string
func (_e *MockIssueTracker_Expecter) ParseRawIssueId(identifier interface{}) *MockIssueTracker_ParseRawIssueId_Call {
	return &MockIssueTracker_ParseRawIssueId_Call{Call: _e.mock.On("ParseRawIssueId", identifier)}
}

func (_c *MockIssueTracker_ParseRawIssueId_Call) Run(run func(identifier string)) *MockIssueTracker_ParseRawIssueId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockIssueTracker_ParseRawIssueId_Call) Return(issueId string) *MockIssueTracker_ParseRawIssueId_Call {
	_c.Call.Return(issueId)
	return _c
}

func (_c *MockIssueTracker_ParseRawIssueId_Call) RunAndReturn(run func(string) string) *MockIssueTracker_ParseRawIssueId_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIssueTracker creates a new instance of MockIssueTracker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIssueTracker(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIssueTracker {
	mock := &MockIssueTracker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
