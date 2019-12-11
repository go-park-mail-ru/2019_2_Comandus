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
	locationClient	clients.LocationClient
}

func NewCompanyUsecase(c company.Repository, mClient clients.ManagerClient, cClient clients.LocationClient) company.Usecase {
	return &CompanyUsecase{
		companyRep: c,
		managerClient: mClient,
		locationClient: cClient,
	}
}

func (u * CompanyUsecase) Create() (*model.Company, error) {
	c := &model.Company{}

	if err := u.companyRep.Create(c); err != nil {
		return nil, errors.Wrap(err, "companyRep.Create()")
	}

	return c, nil
}

func (u *CompanyUsecase) InsertLocation(company *model.Company) (*model.CompanyOutput, error) {
	country, err := u.locationClient.GetCountry(company.Country)
	if err != nil {
		return nil, errors.Wrap(err, "clients.GetCountry()")
	}

	city, err := u.locationClient.GetCity(company.City)
	if err != nil {
		return nil, errors.Wrap(err, "clients.GetCity()")
	}

	res := &model.CompanyOutput{
		ID:          company.ID,
		CompanyName: company.CompanyName,
		Site:        company.Site,
		TagLine:     company.TagLine,
		Description: company.Description,
		Country:     country.Name,
		City:        city.Name,
		Address:     company.Address,
		Phone:       company.Phone,
	}
	return res, nil
}

func (u *CompanyUsecase) Find(id int64) (*model.CompanyOutput, error) {
	c, err := u.companyRep.Find(id)
	if err != nil {
		return nil, errors.Wrapf(err, "companyRep.Find()")
	}

	res, err := u.InsertLocation(c)
	if err != nil {
		return nil, errors.Wrap(err, "InsertLocation()")
	}

	return res, nil
}

func (u *CompanyUsecase) Edit(userID int64, company *model.Company) (*model.CompanyOutput, error) {
	m, err := u.managerClient.GetManagerByUserFromServer(userID)
	if err != nil {
		return nil, errors.Wrapf(err, "client.FindByUser()")
	}

	company.ID = m.CompanyId
	if err := u.companyRep.Edit(company); err != nil {
		return nil, errors.Wrapf(err, "companyRep.Edit()")
	}

	res, err := u.InsertLocation(company)
	if err != nil {
		return nil, errors.Wrap(err, "InsertLocation()")
	}
	return res, nil
}
