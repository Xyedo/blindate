// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/xyedo/blindate/pkg/domain/authentication (interfaces: Repository)

// Package mockrepo is a generated GoMock package.
package mockrepo

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuth is a mock of Repository interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// AddRefreshToken mocks base method.
func (m *MockAuth) AddRefreshToken(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRefreshToken", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRefreshToken indicates an expected call of AddRefreshToken.
func (mr *MockAuthMockRecorder) AddRefreshToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRefreshToken", reflect.TypeOf((*MockAuth)(nil).AddRefreshToken), arg0)
}

// DeleteRefreshToken mocks base method.
func (m *MockAuth) DeleteRefreshToken(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRefreshToken", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRefreshToken indicates an expected call of DeleteRefreshToken.
func (mr *MockAuthMockRecorder) DeleteRefreshToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRefreshToken", reflect.TypeOf((*MockAuth)(nil).DeleteRefreshToken), arg0)
}

// VerifyRefreshToken mocks base method.
func (m *MockAuth) VerifyRefreshToken(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyRefreshToken", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyRefreshToken indicates an expected call of VerifyRefreshToken.
func (mr *MockAuthMockRecorder) VerifyRefreshToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyRefreshToken", reflect.TypeOf((*MockAuth)(nil).VerifyRefreshToken), arg0)
}
