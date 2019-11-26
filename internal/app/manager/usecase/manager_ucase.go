package managerUcase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
)

type ManagerUsecase struct {
	managerRep manager.Repository
}

func NewManagerUsecase(rep manager.Repository) manager.Usecase {
	return &ManagerUsecase{
		managerRep: rep,
	}
}

func (u *ManagerUsecase) Create(userId int64, compId int64) (*model.HireManager, error) {
	m := &model.HireManager{
		AccountID:			userId,
		CompanyID:			compId,
	}

	if err := u.managerRep.Create(m); err != nil {
		return nil, errors.Wrap(err, "CreateUser<-managerRep.Create()")
	}

	return m, nil
}

func (u *ManagerUsecase) FindByUser(id int64) (*model.HireManager, error) {
	m, err := u.managerRep.FindByUser(id)
	if err != nil {
		return nil, errors.Wrap(err, "managerRep.FindByUser()")
	}
	return m, nil
}

func (u *ManagerUsecase) Find(id int64) (*model.HireManager, error) {
	m, err := u.managerRep.Find(id)
	if err != nil {
		return nil, errors.Wrap(err, "managerRep.Find()")
	}
	return m, nil
}