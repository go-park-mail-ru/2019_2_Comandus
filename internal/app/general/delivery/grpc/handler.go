package authgrpc

import (
	"context"
	auth_grpc2 "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/delivery/grpc/auth_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type AuthServer struct {
	UserUcase user.Usecase
}

func NewAuthServerGrpc(gserver *grpc.Server, userUcase user.Usecase) {
	authServer := &AuthServer{
		UserUcase: userUcase,
	}
	auth_grpc2.RegisterAuthHandlerServer(gserver, authServer)
	reflection.Register(gserver)
}

func (s *AuthServer) TransformUserRPC(user *model.User) *auth_grpc2.User {
	if user == nil {
		return nil
	}

	res := &auth_grpc2.User{
		ID:              user.ID,
		FirstName:       user.FirstName,
		SecondName:      user.SecondName,
		UserName:        user.UserName,
		Email:           user.Email,
		Password:        user.Password,
		EncryptPassword: user.EncryptPassword,
		UserType:        user.UserType,
		FreelancerId:    user.FreelancerId,
		HireManagerId:   user.HireManagerId,
		CompanyId:       user.CompanyId,
		Avatar:          user.Avatar,
	}
	return res
}

func (s *AuthServer) TransformUserData(user *auth_grpc2.User) *model.User {
	res := &model.User{
		ID:              user.ID,
		FirstName:       user.FirstName,
		SecondName:      user.SecondName,
		UserName:        user.UserName,
		Email:           user.Email,
		Password:        user.Password,
		EncryptPassword: user.EncryptPassword,
		UserType:        user.UserType,
		FreelancerId:    user.FreelancerId,
		HireManagerId:   user.HireManagerId,
		CompanyId:       user.CompanyId,
		Avatar:          user.Avatar,
	}
	return res
}

func (s *AuthServer) CreateUser(context context.Context, userReq *auth_grpc2.User) (*auth_grpc2.User, error) {
	newUser := &model.User{
		Email:      userReq.Email,
		Password:   userReq.Password,
		FirstName:  userReq.FirstName,
		SecondName: userReq.SecondName,
		UserType:   userReq.UserType,
	}

	if err := s.UserUcase.CreateUser(newUser); err != nil {
		return nil, errors.Wrap(err, "UserUcase.CreateUser")
	}

	res := s.TransformUserRPC(newUser)
	return res, nil
}

func (s *AuthServer) VerifyUser(context context.Context, userReq *auth_grpc2.UserRequest) (*auth_grpc2.UserID, error) {
	newUser := &model.User{
		Email:    userReq.Email,
		Password: userReq.Password,
	}

	id, err := s.UserUcase.VerifyUser(newUser)
	if err != nil {
		return nil, err
	}

	res := &auth_grpc2.UserID{
		ID: id,
	}
	return res, nil
}
