// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/clients/interfaces/responseClient.go

// Package client_mocks is a generated GoMock package.
package client_mocks

import (
	response_grpc "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/delivery/grpc/response_grpc"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockClientResponse is a mock of ClientResponse interface
type MockClientResponse struct {
	ctrl     *gomock.Controller
	recorder *MockClientResponseMockRecorder
}

// MockClientResponseMockRecorder is the mock recorder for MockClientResponse
type MockClientResponseMockRecorder struct {
	mock *MockClientResponse
}

// NewMockClientResponse creates a new mock instance
func NewMockClientResponse(ctrl *gomock.Controller) *MockClientResponse {
	mock := &MockClientResponse{ctrl: ctrl}
	mock.recorder = &MockClientResponseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClientResponse) EXPECT() *MockClientResponseMockRecorder {
	return m.recorder
}

// GetResponseFromServer mocks base method
func (m *MockClientResponse) GetResponseFromServer(id int64) (*response_grpc.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResponseFromServer", id)
	ret0, _ := ret[0].(*response_grpc.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResponseFromServer indicates an expected call of GetResponseFromServer
func (mr *MockClientResponseMockRecorder) GetResponseFromServer(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResponseFromServer", reflect.TypeOf((*MockClientResponse)(nil).GetResponseFromServer), id)
}
