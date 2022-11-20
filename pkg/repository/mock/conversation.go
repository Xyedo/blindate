// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/xyedo/blindate/pkg/repository (interfaces: Conversation)

// Package mockrepo is a generated GoMock package.
package mockrepo

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/xyedo/blindate/pkg/domain"
	entity "github.com/xyedo/blindate/pkg/entity"
)

// MockConversation is a mock of Conversation interface.
type MockConversation struct {
	ctrl     *gomock.Controller
	recorder *MockConversationMockRecorder
}

// MockConversationMockRecorder is the mock recorder for MockConversation.
type MockConversationMockRecorder struct {
	mock *MockConversation
}

// NewMockConversation creates a new mock instance.
func NewMockConversation(ctrl *gomock.Controller) *MockConversation {
	mock := &MockConversation{ctrl: ctrl}
	mock.recorder = &MockConversationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConversation) EXPECT() *MockConversationMockRecorder {
	return m.recorder
}

// DeleteConversationById mocks base method.
func (m *MockConversation) DeleteConversationById(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteConversationById", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteConversationById indicates an expected call of DeleteConversationById.
func (mr *MockConversationMockRecorder) DeleteConversationById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteConversationById", reflect.TypeOf((*MockConversation)(nil).DeleteConversationById), arg0)
}

// InsertConversation mocks base method.
func (m *MockConversation) InsertConversation(arg0, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertConversation", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertConversation indicates an expected call of InsertConversation.
func (mr *MockConversationMockRecorder) InsertConversation(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertConversation", reflect.TypeOf((*MockConversation)(nil).InsertConversation), arg0, arg1)
}

// SelectConversationById mocks base method.
func (m *MockConversation) SelectConversationById(arg0 string) (*domain.Conversation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectConversationById", arg0)
	ret0, _ := ret[0].(*domain.Conversation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectConversationById indicates an expected call of SelectConversationById.
func (mr *MockConversationMockRecorder) SelectConversationById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectConversationById", reflect.TypeOf((*MockConversation)(nil).SelectConversationById), arg0)
}

// SelectConversationByUserId mocks base method.
func (m *MockConversation) SelectConversationByUserId(arg0 string, arg1 *entity.ConvFilter) ([]domain.Conversation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectConversationByUserId", arg0, arg1)
	ret0, _ := ret[0].([]domain.Conversation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectConversationByUserId indicates an expected call of SelectConversationByUserId.
func (mr *MockConversationMockRecorder) SelectConversationByUserId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectConversationByUserId", reflect.TypeOf((*MockConversation)(nil).SelectConversationByUserId), arg0, arg1)
}

// UpdateChatRow mocks base method.
func (m *MockConversation) UpdateChatRow(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateChatRow", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateChatRow indicates an expected call of UpdateChatRow.
func (mr *MockConversationMockRecorder) UpdateChatRow(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateChatRow", reflect.TypeOf((*MockConversation)(nil).UpdateChatRow), arg0)
}

// UpdateDayPass mocks base method.
func (m *MockConversation) UpdateDayPass(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDayPass", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDayPass indicates an expected call of UpdateDayPass.
func (mr *MockConversationMockRecorder) UpdateDayPass(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDayPass", reflect.TypeOf((*MockConversation)(nil).UpdateDayPass), arg0)
}
