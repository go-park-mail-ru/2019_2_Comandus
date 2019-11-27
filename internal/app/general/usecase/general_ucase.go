package general_ucase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients/interfaces"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
)

type GeneralUsecase struct {
	UserClient clients.ClientUser
	ManagerClient clients.ManagerClient
	CompanyClient clients.CompanyClient
	FreelancerClient clients.ClientFreelancer
}

func NewGeneralUsecase(UClient clients.ClientUser, MClient clients.ManagerClient,
	CClient clients.CompanyClient, FClient clients.ClientFreelancer) general.Usecase {
	return &GeneralUsecase{
		UserClient: UClient,
		ManagerClient: MClient,
		CompanyClient: CClient,
		FreelancerClient: FClient,
	}
}

func (u *GeneralUsecase) SignUp(data *model.User) (int64, error) {
	user, err := u.UserClient.CreateUserOnServer(data)
	if err != nil {
		return 0, errors.Wrap(err, "clients.CreateUserOnServer")
	}

	company, err := u.CompanyClient.CreateCompanyOnServer(user.ID)
	if err != nil {
		return 0, errors.Wrap(err, "clients.CreateCompanyOnServer()")
	}

	_, err = u.FreelancerClient.CreateFreelancerOnServer(user.ID)
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
	return u.UserClient.VerifyUserOnServer(user)
}
