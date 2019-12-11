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

func (u *FreelancerUsecase) Create(userId int64) (*model.Freelancer, error) {
	f := &model.Freelancer{
		AccountId: userId,
	}

	if err := u.freelancerRep.Create(f); err != nil {
		return nil, errors.Wrap(err, "freelancerRep.Create()")
	}

	return f, nil
}

func (u *FreelancerUsecase) FindByUser(userId int64) (*model.Freelancer, error) {
	f, err := u.freelancerRep.FindByUser(userId)
	if err != nil {
		return nil, errors.Wrapf(err, "freelancerRep.FindByUser()")
	}
	return f, nil
}

func (u *FreelancerUsecase) Find(id int64) (*model.Freelancer, error) {
	f, err := u.freelancerRep.Find(id)
	if err != nil {
		return nil, errors.Wrapf(err, "freelancerRep.Find()")
	}
	return f, nil
}

func (u *FreelancerUsecase) Edit(new *model.Freelancer, old *model.Freelancer) error {
	if new.ID != old.ID {
		return errors.New("can't change ID")
	}

	if new.AccountId != old.AccountId {
		return errors.New("can't change user associated with")
	}

	if err := u.freelancerRep.Edit(new); err != nil {
		return errors.Wrapf(err, "freelancerRep.Edit()")
	}
	return nil
}

func (u *FreelancerUsecase) PatternSearch(pattern string) ([]model.ExtendFreelancer, error) {
	exFreelancers, err := u.freelancerRep.ListOnPattern(pattern)
	if err != nil {
		return nil, errors.Wrap(err, "freelancerRep.PatternSearch()")
	}
	return exFreelancers, nil
}

func (u *FreelancerUsecase) FindPart(offset int, limit int) ([]model.ExtendFreelancer, error) {
	exFreelancers, err := u.freelancerRep.FindPartByTime(offset, limit)
	if err != nil {
		return nil, errors.Wrap(err, "freelancerRep.FindAll()")
	}
	return exFreelancers, nil
}