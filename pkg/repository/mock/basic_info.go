// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/xyedo/blindate/pkg/repository (interfaces: BasicInfo)

// Package mockrepo is a generated GoMock package.
package mockrepo

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/xyedo/blindate/pkg/domain/entity"
)

// MockBasicInfo is a mock of BasicInfo interface.
type MockBasicInfo struct {
	ctrl     *gomock.Controller
	recorder *MockBasicInfoMockRecorder
}

// MockBasicInfoMockRecorder is the mock recorder for MockBasicInfo.
type MockBasicInfoMockRecorder struct {
	mock *MockBasicInfo
}

// NewMockBasicInfo creates a new mock instance.
func NewMockBasicInfo(ctrl *gomock.Controller) *MockBasicInfo {
	mock := &MockBasicInfo{ctrl: ctrl}
	mock.recorder = &MockBasicInfoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBasicInfo) EXPECT() *MockBasicInfoMockRecorder {
	return m.recorder
}

// GetBasicInfoByUserId mocks base method.
func (m *MockBasicInfo) GetBasicInfoByUserId(arg0 string) (entity.BasicInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBasicInfoByUserId", arg0)
	ret0, _ := ret[0].(entity.BasicInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBasicInfoByUserId indicates an expected call of GetBasicInfoByUserId.
func (mr *MockBasicInfoMockRecorder) GetBasicInfoByUserId(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBasicInfoByUserId", reflect.TypeOf((*MockBasicInfo)(nil).GetBasicInfoByUserId), arg0)
}

// InsertBasicInfo mocks base method.
func (m *MockBasicInfo) InsertBasicInfo(arg0 entity.BasicInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertBasicInfo", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertBasicInfo indicates an expected call of InsertBasicInfo.
func (mr *MockBasicInfoMockRecorder) InsertBasicInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertBasicInfo", reflect.TypeOf((*MockBasicInfo)(nil).InsertBasicInfo), arg0)
}

// UpdateBasicInfo mocks base method.
func (m *MockBasicInfo) UpdateBasicInfo(arg0 entity.BasicInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBasicInfo", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBasicInfo indicates an expected call of UpdateBasicInfo.
func (mr *MockBasicInfoMockRecorder) UpdateBasicInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBasicInfo", reflect.TypeOf((*MockBasicInfo)(nil).UpdateBasicInfo), arg0)
}
