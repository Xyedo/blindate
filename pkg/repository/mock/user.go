// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/xyedo/blindate/pkg/repository (interfaces: User)

// Package mockrepo is a generated GoMock package.
package mockrepo

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/xyedo/blindate/pkg/domain"
	entity "github.com/xyedo/blindate/pkg/entity"
)

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// CreateProfilePicture mocks base method.
func (m *MockUser) CreateProfilePicture(arg0, arg1 string, arg2 bool) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProfilePicture", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProfilePicture indicates an expected call of CreateProfilePicture.
func (mr *MockUserMockRecorder) CreateProfilePicture(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProfilePicture", reflect.TypeOf((*MockUser)(nil).CreateProfilePicture), arg0, arg1, arg2)
}

// GetUserByEmail mocks base method.
func (m *MockUser) GetUserByEmail(arg0 string) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", arg0)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserMockRecorder) GetUserByEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUser)(nil).GetUserByEmail), arg0)
}

// GetUserById mocks base method.
func (m *MockUser) GetUserById(arg0 string) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", arg0)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockUserMockRecorder) GetUserById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockUser)(nil).GetUserById), arg0)
}

// InsertUser mocks base method.
func (m *MockUser) InsertUser(arg0 *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertUser indicates an expected call of InsertUser.
func (mr *MockUserMockRecorder) InsertUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertUser", reflect.TypeOf((*MockUser)(nil).InsertUser), arg0)
}

// ProfilePicSelectedToFalse mocks base method.
func (m *MockUser) ProfilePicSelectedToFalse(arg0 string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProfilePicSelectedToFalse", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProfilePicSelectedToFalse indicates an expected call of ProfilePicSelectedToFalse.
func (mr *MockUserMockRecorder) ProfilePicSelectedToFalse(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProfilePicSelectedToFalse", reflect.TypeOf((*MockUser)(nil).ProfilePicSelectedToFalse), arg0)
}

// SelectProfilePicture mocks base method.
func (m *MockUser) SelectProfilePicture(arg0 string, arg1 *entity.ProfilePicQuery) ([]domain.ProfilePicture, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectProfilePicture", arg0, arg1)
	ret0, _ := ret[0].([]domain.ProfilePicture)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectProfilePicture indicates an expected call of SelectProfilePicture.
func (mr *MockUserMockRecorder) SelectProfilePicture(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectProfilePicture", reflect.TypeOf((*MockUser)(nil).SelectProfilePicture), arg0, arg1)
}

// UpdateUser mocks base method.
func (m *MockUser) UpdateUser(arg0 *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserMockRecorder) UpdateUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUser)(nil).UpdateUser), arg0)
}
