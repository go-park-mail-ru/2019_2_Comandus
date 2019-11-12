package companyUsecase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
)

type CompanyUsecase struct {
	companyrep company.Repository
	managerep  manager.Repository
}

func NewCompanyUsecase(c company.Repository, m manager.Repository) company.Usecase {
	return &CompanyUsecase{
		companyrep: c,
		managerep:  m,
	}
}

func (u *CompanyUsecase) Find(id int64) (*model.Company, error) {
	c, err := u.companyrep.Find(id)
	if err != nil {
		return nil, errors.Wrapf(err, "HandleEditCompany<-Find: ")
	}
	return c, nil
}

func (u *CompanyUsecase) Edit(user *model.User, company *model.Company) error {
	companyID, err := u.managerep.GetCompanyIDByUserID(user.ID)
	if err != nil {
		return errors.Wrapf(err, "HandleEditCompany<-GetCompanyIDByUserID: ")
	}
	//if companyID != company.ID {
	//	err = errors.New("No access to this company")
	//	return errors.Wrapf(err, "HandleEditCompany<-")
	//}
	company.ID = companyID
	if err := u.companyrep.Edit(company); err != nil {
		return errors.Wrapf(err, "HandleEditCompany<-Edit: ")
	}
	return nil
}
