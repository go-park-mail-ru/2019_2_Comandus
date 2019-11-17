package companyUsecase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
)

type CompanyUsecase struct {
	companyRep company.Repository
	manageRep  manager.Repository
}

func NewCompanyUsecase(c company.Repository, m manager.Repository) company.Usecase {
	return &CompanyUsecase{
		companyRep: c,
		manageRep:  m,
	}
}

func (u *CompanyUsecase) Create(c *model.Company) error {
	if err := u.companyRep.Create(c); err != nil {
		return errors.Wrap(err, "Create")
	}
	return nil
}

func (u *CompanyUsecase) Find(id int64) (*model.Company, error) {
	c, err := u.companyRep.Find(id)
	if err != nil {
		return nil, errors.Wrapf(err, "Find")
	}
	return c, nil
}

func (u *CompanyUsecase) Edit(user *model.User, company *model.Company) error {
	if !user.IsManager() {
		return errors.New("only manager can edit company")
	}

	m, err := u.manageRep.FindByUser(user.ID)
	if err != nil {
		return errors.Wrapf(err, "GetCompanyIDByUserID")
	}
	//if companyID != company.ID {
	//	err = errors.New("No access to this company")
	//	return errors.Wrapf(err, "HandleEditCompany<-")
	//}
	company.ID = m.CompanyID
	if err := u.companyRep.Edit(company); err != nil {
		return errors.Wrapf(err, "Edit")
	}
	return nil
}
