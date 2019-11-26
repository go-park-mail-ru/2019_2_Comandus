// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/manager/repository.go

// Package repository_mocks is a generated GoMock package.
package repository_mocks

import (
	model "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockManagerRepository is a mock of Repository interface
type MockManagerRepository struct {
	ctrl     *gomock.Controller
	recorder *MockManagerRepositoryMockRecorder
}

// MockManagerRepositoryMockRecorder is the mock recorder for MockManagerRepository
type MockManagerRepositoryMockRecorder struct {
	mock *MockManagerRepository
}

// NewMockManagerRepository creates a new mock instance
func NewMockManagerRepository(ctrl *gomock.Controller) *MockManagerRepository {
	mock := &MockManagerRepository{ctrl: ctrl}
	mock.recorder = &MockManagerRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockManagerRepository) EXPECT() *MockManagerRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockManagerRepository) Create(manager *model.HireManager) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", manager)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockManagerRepositoryMockRecorder) Create(manager interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockManagerRepository)(nil).Create), manager)
}

// Find mocks base method
func (m *MockManagerRepository) Find(arg0 int64) (*model.HireManager, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0)
	ret0, _ := ret[0].(*model.HireManager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockManagerRepositoryMockRecorder) Find(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockManagerRepository)(nil).Find), arg0)
}

// FindByUser mocks base method
func (m *MockManagerRepository) FindByUser(arg0 int64) (*model.HireManager, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUser", arg0)
	ret0, _ := ret[0].(*model.HireManager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUser indicates an expected call of FindByUser
func (mr *MockManagerRepositoryMockRecorder) FindByUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUser", reflect.TypeOf((*MockManagerRepository)(nil).FindByUser), arg0)
}

// Edit mocks base method
func (m *MockManagerRepository) Edit(manager *model.HireManager) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Edit", manager)
	ret0, _ := ret[0].(error)
	return ret0
}

// Edit indicates an expected call of Edit
func (mr *MockManagerRepositoryMockRecorder) Edit(manager interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Edit", reflect.TypeOf((*MockManagerRepository)(nil).Edit), manager)
}