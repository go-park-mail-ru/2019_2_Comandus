package generalUsecase

import (
	clients "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients/interfaces"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/delivery/grpc/auth_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"github.com/suggest-go/suggest/pkg/suggest"
)

type GeneralUsecase struct {
	authClient       clients.AuthClient
	freelancerClient clients.ClientFreelancer
	managerClient    clients.ManagerClient
	companyClient    clients.CompanyClient
	suggestService   *suggest.Service
}

func (u *GeneralUsecase) VerifyUser(user *model.User) (int64, error) {
	id, err := u.authClient.VerifyUserOnServer(user)
	if err != nil {
		return 0, errors.Wrap(err, "AuthClient.VerifyUserOnServer()")
	}
	return id, nil
}

func NewGeneralUsecase(a clients.AuthClient, f clients.ClientFreelancer, m clients.ManagerClient,
	c clients.CompanyClient) general.Usecase {
	return &GeneralUsecase{
		authClient:       a,
		freelancerClient: f,
		managerClient:    m,
		companyClient:    c,
	}
}

func (u *GeneralUsecase) CreateUser(newUser *model.User) (*auth_grpc.User, error) {
	user, err := u.authClient.CreateUserOnServer(newUser)
	if err != nil {
		return nil, errors.Wrap(err, "AuthClient.CreateUserOnServer()")
	}

	company, err := u.companyClient.CreateCompanyOnServer(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "clients.CreateCompanyOnServer()")
	}

	freelancer, err := u.freelancerClient.CreateFreelancerOnServer(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "clients.CreateFreelancerOnServer")
	}

	manager, err := u.managerClient.CreateManagerOnServer(user.ID, company.ID)
	if err != nil {
		return nil, errors.Wrap(err, "clients.CreateManagerOnServer()")
	}

	user.CompanyId = company.ID
	user.FreelancerId = freelancer.ID
	user.HireManagerId = manager.ID
	return user, nil
}

func (u *GeneralUsecase) GetSuggest(query string, update bool, dict string) ([]string, error) {
	if update {
		var err error
		u.suggestService, err = NewSuggestService()
		if err != nil {
			return nil, errors.Wrap(err, "NewSuggestService()")
		}
	}

	res, err := GetSuggest(u.suggestService, query, dict)
	if err != nil {
		return nil, errors.Wrap(err, "GetSuggest()")
	}

	return res, nil
}
