package gen_grpc

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/delivery/grpc/general_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GeneralServer struct {
	UserUcase user.Usecase
}

func (s *GeneralServer) GetToken(context.Context, *general_grpc.UserID) (*general_grpc.Token, error) {
	panic("implement me")
}

func NewGeneralServerGrpc(gserver *grpc.Server, userUcase user.Usecase) {
	generalServer := &GeneralServer{
		UserUcase: userUcase,
	}
	general_grpc.RegisterGeneralHandlerServer(gserver, generalServer)
	reflection.Register(gserver)
}
