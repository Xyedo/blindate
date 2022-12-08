// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/xyedo/blindate/pkg/repository (interfaces: Online)

// Package mockrepo is a generated GoMock package.
package mockrepo

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/xyedo/blindate/pkg/domain"
)

// MockOnline is a mock of Online interface.
type MockOnline struct {
	ctrl     *gomock.Controller
	recorder *MockOnlineMockRecorder
}

// MockOnlineMockRecorder is the mock recorder for MockOnline.
type MockOnlineMockRecorder struct {
	mock *MockOnline
}

// NewMockOnline creates a new mock instance.
func NewMockOnline(ctrl *gomock.Controller) *MockOnline {
	mock := &MockOnline{ctrl: ctrl}
	mock.recorder = &MockOnlineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOnline) EXPECT() *MockOnlineMockRecorder {
	return m.recorder
}

// InsertNewOnline mocks base method.
func (m *MockOnline) InsertNewOnline(arg0 domain.Online) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertNewOnline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertNewOnline indicates an expected call of InsertNewOnline.
func (mr *MockOnlineMockRecorder) InsertNewOnline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertNewOnline", reflect.TypeOf((*MockOnline)(nil).InsertNewOnline), arg0)
}

// SelectOnline mocks base method.
func (m *MockOnline) SelectOnline(arg0 string) (domain.Online, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectOnline", arg0)
	ret0, _ := ret[0].(domain.Online)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectOnline indicates an expected call of SelectOnline.
func (mr *MockOnlineMockRecorder) SelectOnline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectOnline", reflect.TypeOf((*MockOnline)(nil).SelectOnline), arg0)
}

// UpdateOnline mocks base method.
func (m *MockOnline) UpdateOnline(arg0 string, arg1 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOnline", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOnline indicates an expected call of UpdateOnline.
func (mr *MockOnlineMockRecorder) UpdateOnline(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOnline", reflect.TypeOf((*MockOnline)(nil).UpdateOnline), arg0, arg1)
}
