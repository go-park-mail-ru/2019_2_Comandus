package freelancerUcase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients"
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
		AccountId:         userId,
	}

	if err := u.freelancerRep.Create(f); err != nil {
		return nil, errors.Wrap(err, "Create<-freelancerRep.Create()")
	}

	return f, nil
}

func (u *FreelancerUsecase) InsertLocation(freelancer *model.Freelancer) (*model.FreelancerOutput, error) {
	country, err := clients.GetCountry(freelancer.Country)
	if err != nil {
		return nil, errors.Wrap(err, "clients.GetCountry()")
	}

	city, err := clients.GetCity(freelancer.City)
	if err != nil {
		return nil, errors.Wrap(err, "clients.GetCity()")
	}

	res := &model.FreelancerOutput{
		ID:                freelancer.ID,
		AccountId:         freelancer.AccountId,
		Country:           country.Name,
		City:              city.Name,
		Address:           freelancer.Address,
		Phone:             freelancer.Phone,
		TagLine:           freelancer.TagLine,
		Overview:          freelancer.Overview,
		ExperienceLevelId: freelancer.ExperienceLevelId,
		SpecialityId:      freelancer.SpecialityId,
	}
	return res, nil
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
		return nil, errors.Wrapf(err, "HandleEditFreelancer<-FindByUser: ")
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
		return errors.Wrapf(err, "HandleEditFreelancer<-Edit: ")
	}
	return nil
}

func (u *FreelancerUsecase) FindWithLocation(id int64) (*model.FreelancerOutput, error) {
	f, err := u.freelancerRep.Find(id)
	if err != nil {
		return nil, errors.Wrapf(err, "freelancerRep.Find()")
	}

	res, err := u.InsertLocation(f)
	if err != nil {
		return nil, errors.Wrap(err, "InsertLocation()")
	}
	return res, nil
}

func (u *FreelancerUsecase) FindByUserWithLocation(userId int64) (*model.FreelancerOutput, error) {
	f, err := u.freelancerRep.FindByUser(userId)
	if err != nil {
		return nil, errors.Wrapf(err, "freelancerRep.FindByUser()")
	}

	res, err := u.InsertLocation(f)
	if err != nil {
		return nil, errors.Wrap(err, "InsertLocation()")
	}
	return res, nil
}
