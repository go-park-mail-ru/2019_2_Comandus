package companyUsecase

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/delivery/grpc/manager_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)

type CompanyUsecase struct {
	companyRep company.Repository
}

func NewCompanyUsecase(c company.Repository) company.Usecase {
	return &CompanyUsecase{
		companyRep: c,
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

func (u *CompanyUsecase) Edit(user *model.User, company *model.Company) error {
	conn, err := grpc.Dial(":8084", grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "grpc.Dial()")
	}

	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := manager_grpc.NewManagerHandlerClient(conn)

	userIdMes := &manager_grpc.UserID{
		ID:		user.ID,
	}
	m, err := client.FindByUser(context.Background(), userIdMes)
	if err != nil {
		return errors.Wrapf(err, "client.FindByUser()")
	}

	company.ID = m.CompanyId
	if err := u.companyRep.Edit(company); err != nil {
		return errors.Wrapf(err, "HandleEditCompany<-Edit: ")
	}
	return nil
}
