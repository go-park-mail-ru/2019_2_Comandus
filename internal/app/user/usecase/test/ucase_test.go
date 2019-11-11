package test

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func initUserRep(t * testing.T) *MockRepository{
	t.Helper()
	ctrl := gomock.NewController(t)
	return NewMockRepository(ctrl)
}

func TestUcase_CreateUser(t *testing.T) {

}
