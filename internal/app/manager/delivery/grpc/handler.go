package mgrpc

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/delivery/grpc/manager_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ManagerServer struct {
	Ucase manager.Usecase
}

func NewManagerServerGrpc(gserver *grpc.Server, managerUcase manager.Usecase) {
	server := &ManagerServer{
		Ucase: managerUcase,
	}
	manager_grpc.RegisterManagerHandlerServer(gserver, server)
	reflection.Register(gserver)
}

func (s *ManagerServer) TransformManagerRPC(manager *model.HireManager) *manager_grpc.Manager {
	if manager == nil {
		return nil
	}

	res := &manager_grpc.Manager{
		ID:                   manager.ID,
		AccountId:            manager.AccountID,
		Location:             manager.Location,
		CompanyId:            manager.CompanyID,
	}
	return res
}


func (s *ManagerServer) TransformManagerData(manager *manager_grpc.Manager) *model.HireManager {
	// TODO: fix date
	res := &model.HireManager{
		ID:               manager.ID,
		AccountID:        manager.AccountId,
		Location:         manager.Location,
		CompanyID:        manager.CompanyId,
	}
	return res
}

func (s *ManagerServer) CreateManager(context context.Context, info *manager_grpc.Info) (*manager_grpc.Manager, error) {
	newManager, err := s.Ucase.Create(info.UserID, info.CompanyID)
	if err != nil {
		return nil, errors.Wrap(err, "UserUcase.Create")
	}

	res := s.TransformManagerRPC(newManager)
	return res, nil
}

func (s *ManagerServer) FindByUser(context context.Context, user *manager_grpc.UserID) (*manager_grpc.Manager, error) {
	newManager, err := s.Ucase.FindByUser(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "ManagerUcase.FindByUser")
	}
	res := s.TransformManagerRPC(newManager)
	return res, nil
}

func (s *ManagerServer) Find(context context.Context, manager *manager_grpc.ManagerID) (*manager_grpc.Manager, error) {
	newManager, err := s.Ucase.Find(manager.ID)
	if err != nil {
		return nil, errors.Wrap(err, "ManagerUcase.FindByUser")
	}
	res := s.TransformManagerRPC(newManager)
	return res, nil
}