package clients

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/delivery/grpc/auth_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
)

type AuthClient interface {
	VerifyUserOnServer(user *model.User) (int64, error)
	CreateUserOnServer(data *model.User) (*auth_grpc.User, error)
}
