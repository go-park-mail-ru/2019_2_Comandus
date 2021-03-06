// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/company/repository.go

// Package repository_mocks is a generated GoMock package.
package repository_mocks

import (
	model "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCompanyRepository is a mock of Repository interface
type MockCompanyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCompanyRepositoryMockRecorder
}

// MockCompanyRepositoryMockRecorder is the mock recorder for MockCompanyRepository
type MockCompanyRepositoryMockRecorder struct {
	mock *MockCompanyRepository
}

// NewMockCompanyRepository creates a new mock instance
func NewMockCompanyRepository(ctrl *gomock.Controller) *MockCompanyRepository {
	mock := &MockCompanyRepository{ctrl: ctrl}
	mock.recorder = &MockCompanyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCompanyRepository) EXPECT() *MockCompanyRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockCompanyRepository) Create(company *model.Company) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", company)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockCompanyRepositoryMockRecorder) Create(company interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCompanyRepository)(nil).Create), company)
}

// Find mocks base method
func (m *MockCompanyRepository) Find(arg0 int64) (*model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0)
	ret0, _ := ret[0].(*model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockCompanyRepositoryMockRecorder) Find(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockCompanyRepository)(nil).Find), arg0)
}

// Edit mocks base method
func (m *MockCompanyRepository) Edit(company *model.Company) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Edit", company)
	ret0, _ := ret[0].(error)
	return ret0
}

// Edit indicates an expected call of Edit
func (mr *MockCompanyRepositoryMockRecorder) Edit(company interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Edit", reflect.TypeOf((*MockCompanyRepository)(nil).Edit), company)
}
