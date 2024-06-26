// Code generated by mockery v2.32.4. DO NOT EDIT.

package domain

import mock "github.com/stretchr/testify/mock"

// MockGitProvider is an autogenerated mock type for the GitProvider type
type MockGitProvider struct {
	mock.Mock
}

type MockGitProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGitProvider) EXPECT() *MockGitProvider_Expecter {
	return &MockGitProvider_Expecter{mock: &_m.Mock}
}

// BranchExists provides a mock function with given fields: branch
func (_m *MockGitProvider) BranchExists(branch string) bool {
	ret := _m.Called(branch)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(branch)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockGitProvider_BranchExists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BranchExists'
type MockGitProvider_BranchExists_Call struct {
	*mock.Call
}

// BranchExists is a helper method to define mock.On call
//   - branch string
func (_e *MockGitProvider_Expecter) BranchExists(branch interface{}) *MockGitProvider_BranchExists_Call {
	return &MockGitProvider_BranchExists_Call{Call: _e.mock.On("BranchExists", branch)}
}

func (_c *MockGitProvider_BranchExists_Call) Run(run func(branch string)) *MockGitProvider_BranchExists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGitProvider_BranchExists_Call) Return(_a0 bool) *MockGitProvider_BranchExists_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGitProvider_BranchExists_Call) RunAndReturn(run func(string) bool) *MockGitProvider_BranchExists_Call {
	_c.Call.Return(run)
	return _c
}

// CheckoutBranch provides a mock function with given fields: branch
func (_m *MockGitProvider) CheckoutBranch(branch string) error {
	ret := _m.Called(branch)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(branch)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockGitProvider_CheckoutBranch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckoutBranch'
type MockGitProvider_CheckoutBranch_Call struct {
	*mock.Call
}

// CheckoutBranch is a helper method to define mock.On call
//   - branch string
func (_e *MockGitProvider_Expecter) CheckoutBranch(branch interface{}) *MockGitProvider_CheckoutBranch_Call {
	return &MockGitProvider_CheckoutBranch_Call{Call: _e.mock.On("CheckoutBranch", branch)}
}

func (_c *MockGitProvider_CheckoutBranch_Call) Run(run func(branch string)) *MockGitProvider_CheckoutBranch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGitProvider_CheckoutBranch_Call) Return(err error) *MockGitProvider_CheckoutBranch_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockGitProvider_CheckoutBranch_Call) RunAndReturn(run func(string) error) *MockGitProvider_CheckoutBranch_Call {
	_c.Call.Return(run)
	return _c
}

// CheckoutNewBranchFromOrigin provides a mock function with given fields: branch, base
func (_m *MockGitProvider) CheckoutNewBranchFromOrigin(branch string, base string) error {
	ret := _m.Called(branch, base)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(branch, base)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockGitProvider_CheckoutNewBranchFromOrigin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckoutNewBranchFromOrigin'
type MockGitProvider_CheckoutNewBranchFromOrigin_Call struct {
	*mock.Call
}

// CheckoutNewBranchFromOrigin is a helper method to define mock.On call
//   - branch string
//   - base string
func (_e *MockGitProvider_Expecter) CheckoutNewBranchFromOrigin(branch interface{}, base interface{}) *MockGitProvider_CheckoutNewBranchFromOrigin_Call {
	return &MockGitProvider_CheckoutNewBranchFromOrigin_Call{Call: _e.mock.On("CheckoutNewBranchFromOrigin", branch, base)}
}

func (_c *MockGitProvider_CheckoutNewBranchFromOrigin_Call) Run(run func(branch string, base string)) *MockGitProvider_CheckoutNewBranchFromOrigin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockGitProvider_CheckoutNewBranchFromOrigin_Call) Return(err error) *MockGitProvider_CheckoutNewBranchFromOrigin_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockGitProvider_CheckoutNewBranchFromOrigin_Call) RunAndReturn(run func(string, string) error) *MockGitProvider_CheckoutNewBranchFromOrigin_Call {
	_c.Call.Return(run)
	return _c
}

// CommitEmpty provides a mock function with given fields: message
func (_m *MockGitProvider) CommitEmpty(message string) error {
	ret := _m.Called(message)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockGitProvider_CommitEmpty_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CommitEmpty'
type MockGitProvider_CommitEmpty_Call struct {
	*mock.Call
}

// CommitEmpty is a helper method to define mock.On call
//   - message string
func (_e *MockGitProvider_Expecter) CommitEmpty(message interface{}) *MockGitProvider_CommitEmpty_Call {
	return &MockGitProvider_CommitEmpty_Call{Call: _e.mock.On("CommitEmpty", message)}
}

func (_c *MockGitProvider_CommitEmpty_Call) Run(run func(message string)) *MockGitProvider_CommitEmpty_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGitProvider_CommitEmpty_Call) Return(err error) *MockGitProvider_CommitEmpty_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockGitProvider_CommitEmpty_Call) RunAndReturn(run func(string) error) *MockGitProvider_CommitEmpty_Call {
	_c.Call.Return(run)
	return _c
}

// FetchBranchFromOrigin provides a mock function with given fields: branch
func (_m *MockGitProvider) FetchBranchFromOrigin(branch string) error {
	ret := _m.Called(branch)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(branch)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockGitProvider_FetchBranchFromOrigin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchBranchFromOrigin'
type MockGitProvider_FetchBranchFromOrigin_Call struct {
	*mock.Call
}

// FetchBranchFromOrigin is a helper method to define mock.On call
//   - branch string
func (_e *MockGitProvider_Expecter) FetchBranchFromOrigin(branch interface{}) *MockGitProvider_FetchBranchFromOrigin_Call {
	return &MockGitProvider_FetchBranchFromOrigin_Call{Call: _e.mock.On("FetchBranchFromOrigin", branch)}
}

func (_c *MockGitProvider_FetchBranchFromOrigin_Call) Run(run func(branch string)) *MockGitProvider_FetchBranchFromOrigin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGitProvider_FetchBranchFromOrigin_Call) Return(err error) *MockGitProvider_FetchBranchFromOrigin_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockGitProvider_FetchBranchFromOrigin_Call) RunAndReturn(run func(string) error) *MockGitProvider_FetchBranchFromOrigin_Call {
	_c.Call.Return(run)
	return _c
}

// FindBranch provides a mock function with given fields: substring
func (_m *MockGitProvider) FindBranch(substring string) (string, bool) {
	ret := _m.Called(substring)

	var r0 string
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) (string, bool)); ok {
		return rf(substring)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(substring)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(substring)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// MockGitProvider_FindBranch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindBranch'
type MockGitProvider_FindBranch_Call struct {
	*mock.Call
}

// FindBranch is a helper method to define mock.On call
//   - substring string
func (_e *MockGitProvider_Expecter) FindBranch(substring interface{}) *MockGitProvider_FindBranch_Call {
	return &MockGitProvider_FindBranch_Call{Call: _e.mock.On("FindBranch", substring)}
}

func (_c *MockGitProvider_FindBranch_Call) Run(run func(substring string)) *MockGitProvider_FindBranch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGitProvider_FindBranch_Call) Return(branch string, exists bool) *MockGitProvider_FindBranch_Call {
	_c.Call.Return(branch, exists)
	return _c
}

func (_c *MockGitProvider_FindBranch_Call) RunAndReturn(run func(string) (string, bool)) *MockGitProvider_FindBranch_Call {
	_c.Call.Return(run)
	return _c
}

// GetCommitsToPush provides a mock function with given fields: branch
func (_m *MockGitProvider) GetCommitsToPush(branch string) ([]string, error) {
	ret := _m.Called(branch)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]string, error)); ok {
		return rf(branch)
	}
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(branch)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(branch)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitProvider_GetCommitsToPush_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCommitsToPush'
type MockGitProvider_GetCommitsToPush_Call struct {
	*mock.Call
}

// GetCommitsToPush is a helper method to define mock.On call
//   - branch string
func (_e *MockGitProvider_Expecter) GetCommitsToPush(branch interface{}) *MockGitProvider_GetCommitsToPush_Call {
	return &MockGitProvider_GetCommitsToPush_Call{Call: _e.mock.On("GetCommitsToPush", branch)}
}

func (_c *MockGitProvider_GetCommitsToPush_Call) Run(run func(branch string)) *MockGitProvider_GetCommitsToPush_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGitProvider_GetCommitsToPush_Call) Return(commits []string, err error) *MockGitProvider_GetCommitsToPush_Call {
	_c.Call.Return(commits, err)
	return _c
}

func (_c *MockGitProvider_GetCommitsToPush_Call) RunAndReturn(run func(string) ([]string, error)) *MockGitProvider_GetCommitsToPush_Call {
	_c.Call.Return(run)
	return _c
}

// GetCurrentBranch provides a mock function with given fields:
func (_m *MockGitProvider) GetCurrentBranch() (string, error) {
	ret := _m.Called()

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitProvider_GetCurrentBranch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCurrentBranch'
type MockGitProvider_GetCurrentBranch_Call struct {
	*mock.Call
}

// GetCurrentBranch is a helper method to define mock.On call
func (_e *MockGitProvider_Expecter) GetCurrentBranch() *MockGitProvider_GetCurrentBranch_Call {
	return &MockGitProvider_GetCurrentBranch_Call{Call: _e.mock.On("GetCurrentBranch")}
}

func (_c *MockGitProvider_GetCurrentBranch_Call) Run(run func()) *MockGitProvider_GetCurrentBranch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockGitProvider_GetCurrentBranch_Call) Return(branchName string, err error) *MockGitProvider_GetCurrentBranch_Call {
	_c.Call.Return(branchName, err)
	return _c
}

func (_c *MockGitProvider_GetCurrentBranch_Call) RunAndReturn(run func() (string, error)) *MockGitProvider_GetCurrentBranch_Call {
	_c.Call.Return(run)
	return _c
}

// PushBranch provides a mock function with given fields: branch
func (_m *MockGitProvider) PushBranch(branch string) error {
	ret := _m.Called(branch)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(branch)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockGitProvider_PushBranch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PushBranch'
type MockGitProvider_PushBranch_Call struct {
	*mock.Call
}

// PushBranch is a helper method to define mock.On call
//   - branch string
func (_e *MockGitProvider_Expecter) PushBranch(branch interface{}) *MockGitProvider_PushBranch_Call {
	return &MockGitProvider_PushBranch_Call{Call: _e.mock.On("PushBranch", branch)}
}

func (_c *MockGitProvider_PushBranch_Call) Run(run func(branch string)) *MockGitProvider_PushBranch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGitProvider_PushBranch_Call) Return(err error) *MockGitProvider_PushBranch_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockGitProvider_PushBranch_Call) RunAndReturn(run func(string) error) *MockGitProvider_PushBranch_Call {
	_c.Call.Return(run)
	return _c
}

// RemoteBranchExists provides a mock function with given fields: branch
func (_m *MockGitProvider) RemoteBranchExists(branch string) bool {
	ret := _m.Called(branch)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(branch)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockGitProvider_RemoteBranchExists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoteBranchExists'
type MockGitProvider_RemoteBranchExists_Call struct {
	*mock.Call
}

// RemoteBranchExists is a helper method to define mock.On call
//   - branch string
func (_e *MockGitProvider_Expecter) RemoteBranchExists(branch interface{}) *MockGitProvider_RemoteBranchExists_Call {
	return &MockGitProvider_RemoteBranchExists_Call{Call: _e.mock.On("RemoteBranchExists", branch)}
}

func (_c *MockGitProvider_RemoteBranchExists_Call) Run(run func(branch string)) *MockGitProvider_RemoteBranchExists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGitProvider_RemoteBranchExists_Call) Return(exists bool) *MockGitProvider_RemoteBranchExists_Call {
	_c.Call.Return(exists)
	return _c
}

func (_c *MockGitProvider_RemoteBranchExists_Call) RunAndReturn(run func(string) bool) *MockGitProvider_RemoteBranchExists_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGitProvider creates a new instance of MockGitProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGitProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGitProvider {
	mock := &MockGitProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
