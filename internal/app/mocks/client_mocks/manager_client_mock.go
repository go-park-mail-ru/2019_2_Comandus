// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/clients/interfaces/managerClient.go

// Package client_mocks is a generated GoMock package.
package client_mocks

import (
	manager_grpc "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/delivery/grpc/manager_grpc"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockManagerClient is a mock of ManagerClient interface
type MockManagerClient struct {
	ctrl     *gomock.Controller
	recorder *MockManagerClientMockRecorder
}

// MockManagerClientMockRecorder is the mock recorder for MockManagerClient
type MockManagerClientMockRecorder struct {
	mock *MockManagerClient
}

// NewMockManagerClient creates a new mock instance
func NewMockManagerClient(ctrl *gomock.Controller) *MockManagerClient {
	mock := &MockManagerClient{ctrl: ctrl}
	mock.recorder = &MockManagerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockManagerClient) EXPECT() *MockManagerClientMockRecorder {
	return m.recorder
}

// CreateManagerOnServer mocks base method
func (m *MockManagerClient) CreateManagerOnServer(userId, companyId int64) (*manager_grpc.Manager, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateManagerOnServer", userId, companyId)
	ret0, _ := ret[0].(*manager_grpc.Manager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateManagerOnServer indicates an expected call of CreateManagerOnServer
func (mr *MockManagerClientMockRecorder) CreateManagerOnServer(userId, companyId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateManagerOnServer", reflect.TypeOf((*MockManagerClient)(nil).CreateManagerOnServer), userId, companyId)
}

// GetManagerByUserFromServer mocks base method
func (m *MockManagerClient) GetManagerByUserFromServer(id int64) (*manager_grpc.Manager, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManagerByUserFromServer", id)
	ret0, _ := ret[0].(*manager_grpc.Manager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetManagerByUserFromServer indicates an expected call of GetManagerByUserFromServer
func (mr *MockManagerClientMockRecorder) GetManagerByUserFromServer(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManagerByUserFromServer", reflect.TypeOf((*MockManagerClient)(nil).GetManagerByUserFromServer), id)
}

// GetManagerFromServer mocks base method
func (m *MockManagerClient) GetManagerFromServer(id int64) (*manager_grpc.Manager, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManagerFromServer", id)
	ret0, _ := ret[0].(*manager_grpc.Manager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetManagerFromServer indicates an expected call of GetManagerFromServer
func (mr *MockManagerClientMockRecorder) GetManagerFromServer(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManagerFromServer", reflect.TypeOf((*MockManagerClient)(nil).GetManagerFromServer), id)
}
