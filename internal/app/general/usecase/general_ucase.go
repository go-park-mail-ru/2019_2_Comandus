package general_ucase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
)

type GeneralUsecase struct {
}

func NewGeneralUsecase() general.Usecase {
	return &GeneralUsecase{
	}
}

func (u *GeneralUsecase) SignUp(data *model.User) error {
	user, err := clients.CreateUserOnServer(data)
	if err != nil {
		return errors.Wrap(err, "clients.CreateUserOnServer")
	}

	company, err := clients.CreateCompanyOnServer(user.ID)
	if err != nil {
		return errors.Wrap(err, "clients.CreateCompanyOnServer()")
	}

	_, err = clients.CreateFreelancerOnServer(user.ID)
	if err != nil {
		return errors.Wrap(err, "clients.CreateFreelancerOnServer")
	}

	_, err = clients.CreateManagerOnServer(user.ID, company.ID)
	if err != nil {
		return errors.Wrap(err, "clients.CreateManagerOnServer()")
	}

	return nil
}

func (u *GeneralUsecase) VerifyUser(user *model.User) (int64, error) {
	return clients.VerifyUserOnServer(user)
}
