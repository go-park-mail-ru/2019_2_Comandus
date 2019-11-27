package clients

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/delivery/grpc/user_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
)

type ClientUser interface {
	CreateUserOnServer(data *model.User) (*user_grpc.User, error)
	VerifyUserOnServer(user *model.User) (int64, error)
	GetUserFromServer(userID *user_grpc.UserID) (*user_grpc.User, error)
}
