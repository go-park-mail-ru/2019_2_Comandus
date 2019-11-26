package fgrpc

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/delivery/grpc/freelancer_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type FreelancerServer struct {
	Ucase freelancer.Usecase
}

func NewFreelancerServerGrpc(gserver *grpc.Server, freelancerUcase freelancer.Usecase) {
	server := &FreelancerServer{
		Ucase: freelancerUcase,
	}
	freelancer_grpc.RegisterFreelancerHandlerServer(gserver, server)
	reflection.Register(gserver)
}

func (s *FreelancerServer) TransformFreelancerRPC(freelancer *model.Freelancer) *freelancer_grpc.Freelancer {
	if freelancer == nil {
		return nil
	}

	res := &freelancer_grpc.Freelancer{
		ID:                   freelancer.ID,
		AccountId:            freelancer.AccountId,
		Country:              freelancer.Country,
		City:                 freelancer.City,
		Address:              freelancer.Address,
		Phone:                freelancer.Phone,
		TagLine:              freelancer.TagLine,
		Overview:             freelancer.Overview,
		ExperienceLevelId:    freelancer.ExperienceLevelId,
		SpecialityId:         freelancer.SpecialityId,
	}
	return res
}


func (s *FreelancerServer) TransformFreelancerData(freelancer *freelancer_grpc.Freelancer) *model.Freelancer {
	res := &model.Freelancer{
		ID:                freelancer.ID,
		AccountId:         freelancer.AccountId,
		Country:           freelancer.Country,
		City:              freelancer.City,
		Address:           freelancer.Address,
		Phone:             freelancer.Phone,
		TagLine:           freelancer.TagLine,
		Overview:          freelancer.Overview,
		ExperienceLevelId: freelancer.ExperienceLevelId,
		SpecialityId:      freelancer.SpecialityId,
	}
	return res
}

func (s *FreelancerServer) CreateFreelancer(context context.Context, userID *freelancer_grpc.UserID) (*freelancer_grpc.Freelancer, error) {
	newFreelancer, err := s.Ucase.Create(userID.ID)
	if err != nil {
		return nil, errors.Wrap(err, "UserUcase.CreateUser")
	}

	res := s.TransformFreelancerRPC(newFreelancer)
	return res, nil
}

func (s *FreelancerServer) FindByUser(context context.Context, userID *freelancer_grpc.UserID) (*freelancer_grpc.Freelancer, error) {
	currFreelancer, err := s.Ucase.FindByUser(userID.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Ucase.FindByUser()")
	}
	res := s.TransformFreelancerRPC(currFreelancer)
	return res, nil
}