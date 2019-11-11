// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/user-job/repository.go

// Package mock_user_job is a generated GoMock package.
package mock_user_job

import (
	model "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m_2 *MockRepository) Create(j *model.Job, m *model.HireManager) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Create", j, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockRepositoryMockRecorder) Create(j, m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), j, m)
}

// Find mocks base method
func (m *MockRepository) Find(arg0 int64) (*model.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0)
	ret0, _ := ret[0].(*model.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockRepositoryMockRecorder) Find(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockRepository)(nil).Find), arg0)
}

// Edit mocks base method
func (m *MockRepository) Edit(job *model.Job) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Edit", job)
	ret0, _ := ret[0].(error)
	return ret0
}

// Edit indicates an expected call of Edit
func (mr *MockRepositoryMockRecorder) Edit(job interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Edit", reflect.TypeOf((*MockRepository)(nil).Edit), job)
}

// List mocks base method
func (m *MockRepository) List() ([]model.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]model.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockRepositoryMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRepository)(nil).List))
}
