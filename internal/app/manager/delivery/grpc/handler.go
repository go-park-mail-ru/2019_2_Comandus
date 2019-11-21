package mgrpc

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/delivery/grpc/manager_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"time"
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

	date := &timestamp.Timestamp{
		Seconds:              manager.RegistrationDate.Unix(),
		Nanos:                int32(manager.RegistrationDate.UnixNano()),
	}

	res := &manager_grpc.Manager{
		ID:                   manager.ID,
		AccountId:            manager.AccountID,
		RegistrationDate:     date,
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
		RegistrationDate: time.Time{},
		Location:         manager.Location,
		CompanyID:        manager.CompanyId,
	}
	return res
}

func (s *ManagerServer) CreateManager(context context.Context, info *manager_grpc.Info) (*manager_grpc.Manager, error) {
	log.Println("Manager Server", info.UserID, info.UserID)
	newManager, err := s.Ucase.Create(info.UserID, info.CompanyID)
	if err != nil {
		return nil, errors.Wrap(err, "UserUcase.Create")
	}

	res := s.TransformManagerRPC(newManager)
	return res, nil
}