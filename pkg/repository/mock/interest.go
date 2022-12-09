// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/xyedo/blindate/pkg/repository (interfaces: Interest)

// Package mockrepo is a generated GoMock package.
package mockrepo

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/xyedo/blindate/pkg/domain"
)

// MockInterest is a mock of Interest interface.
type MockInterest struct {
	ctrl     *gomock.Controller
	recorder *MockInterestMockRecorder
}

// MockInterestMockRecorder is the mock recorder for MockInterest.
type MockInterestMockRecorder struct {
	mock *MockInterest
}

// NewMockInterest creates a new mock instance.
func NewMockInterest(ctrl *gomock.Controller) *MockInterest {
	mock := &MockInterest{ctrl: ctrl}
	mock.recorder = &MockInterestMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterest) EXPECT() *MockInterestMockRecorder {
	return m.recorder
}

// DeleteInterestHobbies mocks base method.
func (m *MockInterest) DeleteInterestHobbies(arg0 string, arg1 []string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteInterestHobbies", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteInterestHobbies indicates an expected call of DeleteInterestHobbies.
func (mr *MockInterestMockRecorder) DeleteInterestHobbies(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteInterestHobbies", reflect.TypeOf((*MockInterest)(nil).DeleteInterestHobbies), arg0, arg1)
}

// DeleteInterestMovieSeries mocks base method.
func (m *MockInterest) DeleteInterestMovieSeries(arg0 string, arg1 []string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteInterestMovieSeries", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteInterestMovieSeries indicates an expected call of DeleteInterestMovieSeries.
func (mr *MockInterestMockRecorder) DeleteInterestMovieSeries(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteInterestMovieSeries", reflect.TypeOf((*MockInterest)(nil).DeleteInterestMovieSeries), arg0, arg1)
}

// DeleteInterestSports mocks base method.
func (m *MockInterest) DeleteInterestSports(arg0 string, arg1 []string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteInterestSports", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteInterestSports indicates an expected call of DeleteInterestSports.
func (mr *MockInterestMockRecorder) DeleteInterestSports(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteInterestSports", reflect.TypeOf((*MockInterest)(nil).DeleteInterestSports), arg0, arg1)
}

// DeleteInterestTraveling mocks base method.
func (m *MockInterest) DeleteInterestTraveling(arg0 string, arg1 []string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteInterestTraveling", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteInterestTraveling indicates an expected call of DeleteInterestTraveling.
func (mr *MockInterestMockRecorder) DeleteInterestTraveling(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteInterestTraveling", reflect.TypeOf((*MockInterest)(nil).DeleteInterestTraveling), arg0, arg1)
}

// GetInterest mocks base method.
func (m *MockInterest) GetInterest(arg0 string) (domain.Interest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInterest", arg0)
	ret0, _ := ret[0].(domain.Interest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInterest indicates an expected call of GetInterest.
func (mr *MockInterestMockRecorder) GetInterest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInterest", reflect.TypeOf((*MockInterest)(nil).GetInterest), arg0)
}

// InsertInterestBio mocks base method.
func (m *MockInterest) InsertInterestBio(arg0 *domain.Bio) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertInterestBio", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertInterestBio indicates an expected call of InsertInterestBio.
func (mr *MockInterestMockRecorder) InsertInterestBio(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertInterestBio", reflect.TypeOf((*MockInterest)(nil).InsertInterestBio), arg0)
}

// InsertInterestHobbies mocks base method.
func (m *MockInterest) InsertInterestHobbies(arg0 string, arg1 []domain.Hobbie) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertInterestHobbies", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertInterestHobbies indicates an expected call of InsertInterestHobbies.
func (mr *MockInterestMockRecorder) InsertInterestHobbies(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertInterestHobbies", reflect.TypeOf((*MockInterest)(nil).InsertInterestHobbies), arg0, arg1)
}

// InsertInterestMovieSeries mocks base method.
func (m *MockInterest) InsertInterestMovieSeries(arg0 string, arg1 []domain.MovieSerie) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertInterestMovieSeries", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertInterestMovieSeries indicates an expected call of InsertInterestMovieSeries.
func (mr *MockInterestMockRecorder) InsertInterestMovieSeries(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertInterestMovieSeries", reflect.TypeOf((*MockInterest)(nil).InsertInterestMovieSeries), arg0, arg1)
}

// InsertInterestSports mocks base method.
func (m *MockInterest) InsertInterestSports(arg0 string, arg1 []domain.Sport) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertInterestSports", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertInterestSports indicates an expected call of InsertInterestSports.
func (mr *MockInterestMockRecorder) InsertInterestSports(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertInterestSports", reflect.TypeOf((*MockInterest)(nil).InsertInterestSports), arg0, arg1)
}

// InsertInterestTraveling mocks base method.
func (m *MockInterest) InsertInterestTraveling(arg0 string, arg1 []domain.Travel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertInterestTraveling", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertInterestTraveling indicates an expected call of InsertInterestTraveling.
func (mr *MockInterestMockRecorder) InsertInterestTraveling(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertInterestTraveling", reflect.TypeOf((*MockInterest)(nil).InsertInterestTraveling), arg0, arg1)
}

// InsertNewStats mocks base method.
func (m *MockInterest) InsertNewStats(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertNewStats", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertNewStats indicates an expected call of InsertNewStats.
func (mr *MockInterestMockRecorder) InsertNewStats(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertNewStats", reflect.TypeOf((*MockInterest)(nil).InsertNewStats), arg0)
}

// SelectInterestBio mocks base method.
func (m *MockInterest) SelectInterestBio(arg0 string) (domain.Bio, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectInterestBio", arg0)
	ret0, _ := ret[0].(domain.Bio)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectInterestBio indicates an expected call of SelectInterestBio.
func (mr *MockInterestMockRecorder) SelectInterestBio(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectInterestBio", reflect.TypeOf((*MockInterest)(nil).SelectInterestBio), arg0)
}

// UpdateInterestBio mocks base method.
func (m *MockInterest) UpdateInterestBio(arg0 domain.Bio) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInterestBio", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInterestBio indicates an expected call of UpdateInterestBio.
func (mr *MockInterestMockRecorder) UpdateInterestBio(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInterestBio", reflect.TypeOf((*MockInterest)(nil).UpdateInterestBio), arg0)
}

// UpdateInterestHobbies mocks base method.
func (m *MockInterest) UpdateInterestHobbies(arg0 string, arg1 []domain.Hobbie) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInterestHobbies", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInterestHobbies indicates an expected call of UpdateInterestHobbies.
func (mr *MockInterestMockRecorder) UpdateInterestHobbies(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInterestHobbies", reflect.TypeOf((*MockInterest)(nil).UpdateInterestHobbies), arg0, arg1)
}

// UpdateInterestMovieSeries mocks base method.
func (m *MockInterest) UpdateInterestMovieSeries(arg0 string, arg1 []domain.MovieSerie) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInterestMovieSeries", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInterestMovieSeries indicates an expected call of UpdateInterestMovieSeries.
func (mr *MockInterestMockRecorder) UpdateInterestMovieSeries(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInterestMovieSeries", reflect.TypeOf((*MockInterest)(nil).UpdateInterestMovieSeries), arg0, arg1)
}

// UpdateInterestSport mocks base method.
func (m *MockInterest) UpdateInterestSport(arg0 string, arg1 []domain.Sport) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInterestSport", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInterestSport indicates an expected call of UpdateInterestSport.
func (mr *MockInterestMockRecorder) UpdateInterestSport(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInterestSport", reflect.TypeOf((*MockInterest)(nil).UpdateInterestSport), arg0, arg1)
}

// UpdateInterestTraveling mocks base method.
func (m *MockInterest) UpdateInterestTraveling(arg0 string, arg1 []domain.Travel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInterestTraveling", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInterestTraveling indicates an expected call of UpdateInterestTraveling.
func (mr *MockInterestMockRecorder) UpdateInterestTraveling(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInterestTraveling", reflect.TypeOf((*MockInterest)(nil).UpdateInterestTraveling), arg0, arg1)
}
