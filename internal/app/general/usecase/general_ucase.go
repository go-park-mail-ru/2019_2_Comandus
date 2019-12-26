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
	userClient		 clients.ClientUser
	freelancerClient clients.ClientFreelancer
	managerClient    clients.ManagerClient
	companyClient    clients.CompanyClient
	jobClient		 clients.ClientJob
	suggestService   *suggest.Service
}

func (u *GeneralUsecase) VerifyUser(user *model.User) (int64, *model.HttpError) {
	id, err := u.authClient.VerifyUserOnServer(user)
	if err != nil {
		err.LogErr = errors.Wrapf(err.LogErr, "AuthClient.VerifyUserOnServer()")
		return 0, err
	}
	return id, nil
}

func NewGeneralUsecase(a clients.AuthClient, f clients.ClientFreelancer, m clients.ManagerClient,
	c clients.CompanyClient, u clients.ClientUser, j clients.ClientJob) general.Usecase {
	return &GeneralUsecase{
		authClient:       a,
		freelancerClient: f,
		managerClient:    m,
		companyClient:    c,
		userClient:       u,
		jobClient:		  j,
	}
}

func (u *GeneralUsecase) CreateUser(newUser *model.User) (*auth_grpc.User, *model.HttpError) {
	user, err := u.authClient.CreateUserOnServer(newUser)
	if err != nil {
		err.LogErr = errors.Wrap(err.LogErr, "AuthClient.CreateUserOnServer()")
		return nil, err
	}

	company, err1 := u.companyClient.CreateCompanyOnServer(user.ID)
	if err1 != nil {
		return nil, &model.HttpError{ClientErr: errors.Wrap(err1, "clients.CreateCompanyOnServer()")}
	}

	freelancer, err1 := u.freelancerClient.CreateFreelancerOnServer(user.ID)
	if err1 != nil {
		return nil, &model.HttpError{ClientErr: errors.Wrap(err1, "clients.CreateFreelancerOnServer")}
	}

	manager, err1 := u.managerClient.CreateManagerOnServer(user.ID, company.ID)
	if err1 != nil {
		return nil, &model.HttpError{ClientErr: errors.Wrap(err1, "clients.CreateManagerOnServer()")}
	}

	user.CompanyId = company.ID
	user.FreelancerId = freelancer.ID
	user.HireManagerId = manager.ID
	return user, nil
}

func (u *GeneralUsecase) GetSuggest(query string, dict string) ([]string, error) {
	if u.suggestService == nil {
		if err := u.createSuggestService(); err != nil {
			return nil, errors.Wrap(err, "suggest service not implemented")
		}
	}

	res, err := GetSuggest(u.suggestService, query, dict)
	if err != nil {
		return nil, errors.Wrap(err, "GetSuggest()")
	}

	return res, nil
}

func (u *GeneralUsecase) createSuggestService() error {
	suggestService, err := NewSuggestService(u.userClient, u.jobClient)
	if err != nil {
		u.suggestService = nil
		return err
	}
	u.suggestService = suggestService
	return nil
}
