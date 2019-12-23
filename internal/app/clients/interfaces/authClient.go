package clients

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/delivery/grpc/auth_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/location/delivery/grpc/location_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
)

type AuthClient interface {
	VerifyUserOnServer(user *model.User) (int64, *model.HttpError)
	CreateUserOnServer(data *model.User) (*auth_grpc.User, *model.HttpError)
}

type LocationClient interface {
	GetCountry(id int64) (*location_grpc.Country, error)
	GetCity(id int64) (*location_grpc.City, error)
}
