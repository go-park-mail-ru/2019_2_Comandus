package clients

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/delivery/grpc/user_grpc"
)

type ClientUser interface {
	GetUserFromServer(userID *user_grpc.UserID) (*user_grpc.User, error)
}
