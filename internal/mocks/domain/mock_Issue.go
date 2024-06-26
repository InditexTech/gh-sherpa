// Code generated by mockery v2.32.4. DO NOT EDIT.

package domain

import (
	domain "github.com/InditexTech/gh-sherpa/internal/domain"
	issue_types "github.com/InditexTech/gh-sherpa/internal/domain/issue_types"

	mock "github.com/stretchr/testify/mock"
)

// MockIssue is an autogenerated mock type for the Issue type
type MockIssue struct {
	mock.Mock
}

type MockIssue_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIssue) EXPECT() *MockIssue_Expecter {
	return &MockIssue_Expecter{mock: &_m.Mock}
}

// Body provides a mock function with given fields:
func (_m *MockIssue) Body() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockIssue_Body_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Body'
type MockIssue_Body_Call struct {
	*mock.Call
}

// Body is a helper method to define mock.On call
func (_e *MockIssue_Expecter) Body() *MockIssue_Body_Call {
	return &MockIssue_Body_Call{Call: _e.mock.On("Body")}
}

func (_c *MockIssue_Body_Call) Run(run func()) *MockIssue_Body_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockIssue_Body_Call) Return(_a0 string) *MockIssue_Body_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIssue_Body_Call) RunAndReturn(run func() string) *MockIssue_Body_Call {
	_c.Call.Return(run)
	return _c
}

// FormatID provides a mock function with given fields:
func (_m *MockIssue) FormatID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockIssue_FormatID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FormatID'
type MockIssue_FormatID_Call struct {
	*mock.Call
}

// FormatID is a helper method to define mock.On call
func (_e *MockIssue_Expecter) FormatID() *MockIssue_FormatID_Call {
	return &MockIssue_FormatID_Call{Call: _e.mock.On("FormatID")}
}

func (_c *MockIssue_FormatID_Call) Run(run func()) *MockIssue_FormatID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockIssue_FormatID_Call) Return(_a0 string) *MockIssue_FormatID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIssue_FormatID_Call) RunAndReturn(run func() string) *MockIssue_FormatID_Call {
	_c.Call.Return(run)
	return _c
}

// ID provides a mock function with given fields:
func (_m *MockIssue) ID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockIssue_ID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ID'
type MockIssue_ID_Call struct {
	*mock.Call
}

// ID is a helper method to define mock.On call
func (_e *MockIssue_Expecter) ID() *MockIssue_ID_Call {
	return &MockIssue_ID_Call{Call: _e.mock.On("ID")}
}

func (_c *MockIssue_ID_Call) Run(run func()) *MockIssue_ID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockIssue_ID_Call) Return(_a0 string) *MockIssue_ID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIssue_ID_Call) RunAndReturn(run func() string) *MockIssue_ID_Call {
	_c.Call.Return(run)
	return _c
}

// Title provides a mock function with given fields:
func (_m *MockIssue) Title() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockIssue_Title_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Title'
type MockIssue_Title_Call struct {
	*mock.Call
}

// Title is a helper method to define mock.On call
func (_e *MockIssue_Expecter) Title() *MockIssue_Title_Call {
	return &MockIssue_Title_Call{Call: _e.mock.On("Title")}
}

func (_c *MockIssue_Title_Call) Run(run func()) *MockIssue_Title_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockIssue_Title_Call) Return(_a0 string) *MockIssue_Title_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIssue_Title_Call) RunAndReturn(run func() string) *MockIssue_Title_Call {
	_c.Call.Return(run)
	return _c
}

// TrackerType provides a mock function with given fields:
func (_m *MockIssue) TrackerType() domain.IssueTrackerType {
	ret := _m.Called()

	var r0 domain.IssueTrackerType
	if rf, ok := ret.Get(0).(func() domain.IssueTrackerType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.IssueTrackerType)
	}

	return r0
}

// MockIssue_TrackerType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TrackerType'
type MockIssue_TrackerType_Call struct {
	*mock.Call
}

// TrackerType is a helper method to define mock.On call
func (_e *MockIssue_Expecter) TrackerType() *MockIssue_TrackerType_Call {
	return &MockIssue_TrackerType_Call{Call: _e.mock.On("TrackerType")}
}

func (_c *MockIssue_TrackerType_Call) Run(run func()) *MockIssue_TrackerType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockIssue_TrackerType_Call) Return(_a0 domain.IssueTrackerType) *MockIssue_TrackerType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIssue_TrackerType_Call) RunAndReturn(run func() domain.IssueTrackerType) *MockIssue_TrackerType_Call {
	_c.Call.Return(run)
	return _c
}

// Type provides a mock function with given fields:
func (_m *MockIssue) Type() issue_types.IssueType {
	ret := _m.Called()

	var r0 issue_types.IssueType
	if rf, ok := ret.Get(0).(func() issue_types.IssueType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(issue_types.IssueType)
	}

	return r0
}

// MockIssue_Type_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Type'
type MockIssue_Type_Call struct {
	*mock.Call
}

// Type is a helper method to define mock.On call
func (_e *MockIssue_Expecter) Type() *MockIssue_Type_Call {
	return &MockIssue_Type_Call{Call: _e.mock.On("Type")}
}

func (_c *MockIssue_Type_Call) Run(run func()) *MockIssue_Type_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockIssue_Type_Call) Return(_a0 issue_types.IssueType) *MockIssue_Type_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIssue_Type_Call) RunAndReturn(run func() issue_types.IssueType) *MockIssue_Type_Call {
	_c.Call.Return(run)
	return _c
}

// TypeLabel provides a mock function with given fields:
func (_m *MockIssue) TypeLabel() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockIssue_TypeLabel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TypeLabel'
type MockIssue_TypeLabel_Call struct {
	*mock.Call
}

// TypeLabel is a helper method to define mock.On call
func (_e *MockIssue_Expecter) TypeLabel() *MockIssue_TypeLabel_Call {
	return &MockIssue_TypeLabel_Call{Call: _e.mock.On("TypeLabel")}
}

func (_c *MockIssue_TypeLabel_Call) Run(run func()) *MockIssue_TypeLabel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockIssue_TypeLabel_Call) Return(_a0 string) *MockIssue_TypeLabel_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIssue_TypeLabel_Call) RunAndReturn(run func() string) *MockIssue_TypeLabel_Call {
	_c.Call.Return(run)
	return _c
}

// URL provides a mock function with given fields:
func (_m *MockIssue) URL() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockIssue_URL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'URL'
type MockIssue_URL_Call struct {
	*mock.Call
}

// URL is a helper method to define mock.On call
func (_e *MockIssue_Expecter) URL() *MockIssue_URL_Call {
	return &MockIssue_URL_Call{Call: _e.mock.On("URL")}
}

func (_c *MockIssue_URL_Call) Run(run func()) *MockIssue_URL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockIssue_URL_Call) Return(_a0 string) *MockIssue_URL_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIssue_URL_Call) RunAndReturn(run func() string) *MockIssue_URL_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIssue creates a new instance of MockIssue. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIssue(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIssue {
	mock := &MockIssue{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
