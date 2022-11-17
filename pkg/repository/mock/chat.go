// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/xyedo/blindate/pkg/repository (interfaces: ChatRepo)

// Package mockrepo is a generated GoMock package.
package mockrepo

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/xyedo/blindate/pkg/entity"
)

// MockChatRepo is a mock of ChatRepo interface.
type MockChatRepo struct {
	ctrl     *gomock.Controller
	recorder *MockChatRepoMockRecorder
}

// MockChatRepoMockRecorder is the mock recorder for MockChatRepo.
type MockChatRepoMockRecorder struct {
	mock *MockChatRepo
}

// NewMockChatRepo creates a new mock instance.
func NewMockChatRepo(ctrl *gomock.Controller) *MockChatRepo {
	mock := &MockChatRepo{ctrl: ctrl}
	mock.recorder = &MockChatRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChatRepo) EXPECT() *MockChatRepoMockRecorder {
	return m.recorder
}

// DeleteChatById mocks base method.
func (m *MockChatRepo) DeleteChatById(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteChatById", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteChatById indicates an expected call of DeleteChatById.
func (mr *MockChatRepoMockRecorder) DeleteChatById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteChatById", reflect.TypeOf((*MockChatRepo)(nil).DeleteChatById), arg0)
}

// InsertNewChat mocks base method.
func (m *MockChatRepo) InsertNewChat(arg0 *entity.Chat) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertNewChat", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertNewChat indicates an expected call of InsertNewChat.
func (mr *MockChatRepoMockRecorder) InsertNewChat(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertNewChat", reflect.TypeOf((*MockChatRepo)(nil).InsertNewChat), arg0)
}

// SelectChat mocks base method.
func (m *MockChatRepo) SelectChat(arg0 string, arg1 entity.ChatFilter) ([]entity.Chat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectChat", arg0, arg1)
	ret0, _ := ret[0].([]entity.Chat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectChat indicates an expected call of SelectChat.
func (mr *MockChatRepoMockRecorder) SelectChat(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectChat", reflect.TypeOf((*MockChatRepo)(nil).SelectChat), arg0, arg1)
}

// UpdateSeenChatById mocks base method.
func (m *MockChatRepo) UpdateSeenChatById(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSeenChatById", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSeenChatById indicates an expected call of UpdateSeenChatById.
func (mr *MockChatRepoMockRecorder) UpdateSeenChatById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSeenChatById", reflect.TypeOf((*MockChatRepo)(nil).UpdateSeenChatById), arg0)
}
