package companyUsecase

import (
	clients "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients/interfaces"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
)

type CompanyUsecase struct {
	companyRep 		company.Repository
	managerClient 	clients.ManagerClient
}

func NewCompanyUsecase(c company.Repository, mClient clients.ManagerClient ) company.Usecase {
	return &CompanyUsecase{
		companyRep: c,
		managerClient: mClient,
	}
}

func (u * CompanyUsecase) Create() (*model.Company, error) {
	c := &model.Company{}

	if err := u.companyRep.Create(c); err != nil {
		return nil, errors.Wrap(err, "companyRep.Create()")
	}

	return c, nil
}

func (u *CompanyUsecase) Find(id int64) (*model.Company, error) {
	c, err := u.companyRep.Find(id)
	if err != nil {
		return nil, errors.Wrapf(err, "companyRep.Find()")
	}
	return c, nil
}

func (u *CompanyUsecase) Edit(userId int64, company *model.Company) error {
	m, err := u.managerClient.GetManagerByUserFromServer(userId)
	if err != nil {
		return errors.Wrapf(err, "client.FindByUser()")
	}

	company.ID = m.CompanyId
	if err := u.companyRep.Edit(company); err != nil {
		return errors.Wrapf(err, "companyRep.Edit()")
	}
	return nil
}
