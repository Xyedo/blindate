// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/xyedo/blindate/pkg/repository (interfaces: Location)

// Package mockrepo is a generated GoMock package.
package mockrepo

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/xyedo/blindate/pkg/domain"
	entity "github.com/xyedo/blindate/pkg/entity"
)

// MockLocation is a mock of Location interface.
type MockLocation struct {
	ctrl     *gomock.Controller
	recorder *MockLocationMockRecorder
}

// MockLocationMockRecorder is the mock recorder for MockLocation.
type MockLocationMockRecorder struct {
	mock *MockLocation
}

// NewMockLocation creates a new mock instance.
func NewMockLocation(ctrl *gomock.Controller) *MockLocation {
	mock := &MockLocation{ctrl: ctrl}
	mock.recorder = &MockLocationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLocation) EXPECT() *MockLocationMockRecorder {
	return m.recorder
}

// GetClosestUser mocks base method.
func (m *MockLocation) GetClosestUser(arg0 string, arg1 int) ([]domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClosestUser", arg0, arg1)
	ret0, _ := ret[0].([]domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClosestUser indicates an expected call of GetClosestUser.
func (mr *MockLocationMockRecorder) GetClosestUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClosestUser", reflect.TypeOf((*MockLocation)(nil).GetClosestUser), arg0, arg1)
}

// GetLocationByUserId mocks base method.
func (m *MockLocation) GetLocationByUserId(arg0 string) (*entity.Location, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocationByUserId", arg0)
	ret0, _ := ret[0].(*entity.Location)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLocationByUserId indicates an expected call of GetLocationByUserId.
func (mr *MockLocationMockRecorder) GetLocationByUserId(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocationByUserId", reflect.TypeOf((*MockLocation)(nil).GetLocationByUserId), arg0)
}

// InsertNewLocation mocks base method.
func (m *MockLocation) InsertNewLocation(arg0 *entity.Location) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertNewLocation", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertNewLocation indicates an expected call of InsertNewLocation.
func (mr *MockLocationMockRecorder) InsertNewLocation(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertNewLocation", reflect.TypeOf((*MockLocation)(nil).InsertNewLocation), arg0)
}

// UpdateLocation mocks base method.
func (m *MockLocation) UpdateLocation(arg0 *entity.Location) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLocation", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateLocation indicates an expected call of UpdateLocation.
func (mr *MockLocationMockRecorder) UpdateLocation(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLocation", reflect.TypeOf((*MockLocation)(nil).UpdateLocation), arg0)
}