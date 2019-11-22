package ugrpc

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/delivery/grpc/user_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type UserServer struct {
	UserUcase user.Usecase
}

func NewUserServerGrpc(gserver *grpc.Server, userUcase user.Usecase) {
	userServer := &UserServer{
		UserUcase: userUcase,
	}
	user_grpc.RegisterUserHandlerServer(gserver, userServer)
	reflection.Register(gserver)
}

func (s *UserServer) TransformUserRPC(user *model.User) *user_grpc.User {
	if user == nil {
		return nil
	}

	res := &user_grpc.User{
		ID:              	user.ID,
		FirstName:       	user.FirstName,
		SecondName:      	user.SecondName,
		UserName:        	user.UserName,
		Email:           	user.Email,
		Password:        	user.Password,
		EncryptPassword: 	user.EncryptPassword,
		UserType:        	user.UserType,
		FreelancerId:    	user.FreelancerId,
		HireManagerId:   	user.HireManagerId,
		CompanyId:       	user.CompanyId,
		Avatar:				user.Avatar,
	}
	return res
}


func (s *UserServer) TransformUserData(user *user_grpc.User) *model.User {
	res := &model.User{
		ID:              	user.ID,
		FirstName:       	user.FirstName,
		SecondName:      	user.SecondName,
		UserName:        	user.UserName,
		Email:           	user.Email,
		Password:        	user.Password,
		EncryptPassword: 	user.EncryptPassword,
		UserType:        	user.UserType,
		FreelancerId:    	user.FreelancerId,
		HireManagerId:   	user.HireManagerId,
		CompanyId:       	user.CompanyId,
		Avatar:				user.Avatar,
	}
	return res
}

func (s *UserServer) CreateUser(context context.Context,userReq *user_grpc.User) (*user_grpc.User, error) {
	newUser := &model.User{
		Email:		userReq.Email,
		Password:	userReq.Password,
		FirstName: 	userReq.FirstName,
		SecondName:	userReq.SecondName,
		UserType:	userReq.UserType,
	}

	if err := s.UserUcase.CreateUser(newUser); err != nil {
		return nil, errors.Wrap(err, "UserUcase.CreateUser")
	}

	res := s.TransformUserRPC(newUser)
	return res, nil
}

func (s *UserServer) VerifyUser(context context.Context,userReq *user_grpc.UserRequest) (*user_grpc.UserID, error) {
	newUser := &model.User{
		Email:           userReq.Email,
		Password:        userReq.Password,
	}

	id, err := s.UserUcase.VerifyUser(newUser)
	if err != nil {
		return nil, err
	}

	res := &user_grpc.UserID{
		ID:		id,
	}
	return res, nil
}


func (s *UserServer) Find(context context.Context,userId *user_grpc.UserID) (*user_grpc.User, error) {
	currUser, err := s.UserUcase.Find(userId.ID)
	if err != nil {
		return nil, errors.Wrap(err, "UserUcase.Find()")
	}
	res := s.TransformUserRPC(currUser)
	return res, nil
}