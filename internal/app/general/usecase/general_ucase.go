package generalUsecase

import (
	server_clients "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients/server-clients"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/delivery/grpc/auth_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
)

type GeneralUsecase struct {
	grpcClients		*server_clients.ServerClients
}

func (u *GeneralUsecase) VerifyUser(user *model.User) (int64, error) {
	id, err := u.grpcClients.AuthClient.VerifyUserOnServer(user)
	if err != nil {
		return 0, errors.Wrap(err, "AuthClient.VerifyUserOnServer()")
	}
	return id, nil
}

func NewGeneralUsecase(clients *server_clients.ServerClients) general.Usecase {
	return &GeneralUsecase{
		grpcClients: clients,
	}
}

func (u *GeneralUsecase) CreateUser(newUser *model.User) (*auth_grpc.User, error) {
	user, err := u.grpcClients.AuthClient.CreateUserOnServer(newUser)
	if err != nil {
		return nil, errors.Wrap(err, "AuthClient.CreateUserOnServer()")
	}
	return user, nil
}
