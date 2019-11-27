// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/company/usecase.go

// Package ucase_mocks is a generated GoMock package.
package ucase_mocks

import (
	model "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCompanyUsecase is a mock of Usecase interface
type MockCompanyUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockCompanyUsecaseMockRecorder
}

// MockCompanyUsecaseMockRecorder is the mock recorder for MockCompanyUsecase
type MockCompanyUsecaseMockRecorder struct {
	mock *MockCompanyUsecase
}

// NewMockCompanyUsecase creates a new mock instance
func NewMockCompanyUsecase(ctrl *gomock.Controller) *MockCompanyUsecase {
	mock := &MockCompanyUsecase{ctrl: ctrl}
	mock.recorder = &MockCompanyUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCompanyUsecase) EXPECT() *MockCompanyUsecaseMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockCompanyUsecase) Create() (*model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create")
	ret0, _ := ret[0].(*model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockCompanyUsecaseMockRecorder) Create() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCompanyUsecase)(nil).Create))
}

// Find mocks base method
func (m *MockCompanyUsecase) Find(id int64) (*model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", id)
	ret0, _ := ret[0].(*model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockCompanyUsecaseMockRecorder) Find(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockCompanyUsecase)(nil).Find), id)
}

// Edit mocks base method
func (m *MockCompanyUsecase) Edit(userId int64, company *model.Company) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Edit", userId, company)
	ret0, _ := ret[0].(error)
	return ret0
}

// Edit indicates an expected call of Edit
func (mr *MockCompanyUsecaseMockRecorder) Edit(userId, company interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Edit", reflect.TypeOf((*MockCompanyUsecase)(nil).Edit), userId, company)
}
