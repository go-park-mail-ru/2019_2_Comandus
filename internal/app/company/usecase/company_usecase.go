package companyUsecase

import (
	server_clients "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients/server-clients"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
)

type CompanyUsecase struct {
	companyRep 		company.Repository
	grpcClients		*server_clients.ServerClients
}

func NewCompanyUsecase(c company.Repository, clients *server_clients.ServerClients) company.Usecase {
		return &CompanyUsecase{
			companyRep: c,
			grpcClients: clients,
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
	m, err := u.grpcClients.ManagerClient.GetManagerByUserFromServer(userId)
	if err != nil {
		return errors.Wrapf(err, "client.FindByUser()")
	}

	company.ID = m.CompanyId
	if err := u.companyRep.Edit(company); err != nil {
		return errors.Wrapf(err, "companyRep.Edit()")
	}
	return nil
}
