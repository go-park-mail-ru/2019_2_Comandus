package freelancerUcase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
)

type FreelancerUsecase struct {
	freelancerRep freelancer.Repository
}

func NewFreelancerUsecase(f freelancer.Repository) freelancer.Usecase {
	return &FreelancerUsecase{
		freelancerRep: f,
	}
}

func (u *FreelancerUsecase) FindByUser(user *model.User) (*model.Freelancer, error) {
	f, err := u.freelancerRep.FindByUser(user.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "HandleEditFreelancer<-FindByUser: ")
	}
	return f, nil
}

func (u *FreelancerUsecase) Find(id int64) (*model.Freelancer, error) {
	f, err := u.freelancerRep.Find(id)
	if err != nil {
		return nil, errors.Wrapf(err, "HandleEditFreelancer<-FindByUser: ")
	}
	return f, nil
}

func (u *FreelancerUsecase) Edit(user *model.User, freelancer *model.Freelancer) error {
	f, err := u.freelancerRep.FindByUser(user.ID)
	if err != nil {
		return errors.Wrapf(err, "HandleEditFreelancer<-FindByUser: ")
	}
	freelancer.ID = f.ID
	freelancer.AccountId = f.AccountId
	freelancer.RegistrationDate = f.RegistrationDate
	if err := u.freelancerRep.Edit(freelancer); err != nil {
		return errors.Wrapf(err, "HandleEditFreelancer<-Edit: ")
	}
	return nil
}
