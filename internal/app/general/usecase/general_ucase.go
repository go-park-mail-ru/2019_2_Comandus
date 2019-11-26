package general_ucase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
)

type GeneralUsecase struct {
	ManagerClient *clients.ClientManager
}

func NewGeneralUsecase() general.Usecase {
	return &GeneralUsecase{
		ManagerClient: new(clients.ClientManager),
	}
}

func (u *GeneralUsecase) SignUp(data *model.User) (int64, error) {
	user, err := clients.CreateUserOnServer(data)
	if err != nil {
		return 0, errors.Wrap(err, "clients.CreateUserOnServer")
	}

	company, err := clients.CreateCompanyOnServer(user.ID)
	if err != nil {
		return 0, errors.Wrap(err, "clients.CreateCompanyOnServer()")
	}

	_, err = clients.CreateFreelancerOnServer(user.ID)
	if err != nil {
		return 0, errors.Wrap(err, "clients.CreateFreelancerOnServer")
	}

	_, err = u.ManagerClient.CreateManagerOnServer(user.ID, company.ID)
	if err != nil {
		return 0, errors.Wrap(err, "clients.CreateManagerOnServer()")
	}

	return user.ID, nil
}

func (u *GeneralUsecase) VerifyUser(user *model.User) (int64, error) {
	return clients.VerifyUserOnServer(user)
}
