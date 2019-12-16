package freelancerUcase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients/interfaces"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
)

type FreelancerUsecase struct {
	freelancerRep  freelancer.Repository
	locationClient clients.LocationClient
}

func NewFreelancerUsecase(f freelancer.Repository, c clients.LocationClient) freelancer.Usecase {
	return &FreelancerUsecase{
		freelancerRep:  f,
		locationClient: c,
	}
}

func (u *FreelancerUsecase) Create(userId int64) (*model.Freelancer, error) {
	f := &model.Freelancer{
		AccountId: userId,
		Country:0,
		City:1,
	}

	if err := u.freelancerRep.Create(f); err != nil {
		return nil, errors.Wrap(err, "freelancerRep.Create()")
	}

	return f, nil
}

func (u *FreelancerUsecase) InsertLocation(freelancer *model.Freelancer) (*model.FreelancerOutput, error) {
	country, err := u.locationClient.GetCountry(freelancer.Country)
	if err != nil {
		return nil, errors.Wrap(err, "clients.GetCountry()")
	}

	city, err := u.locationClient.GetCity(freelancer.City)
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
		Avatar:            freelancer.Avatar,
	}
	return res, nil
}

func (u *FreelancerUsecase) FindByUser(userId int64) (*model.FreelancerOutput, error) {
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

func (u *FreelancerUsecase) Find(id int64) (*model.ExtendedOutputFreelancer, error) {
	f, err := u.freelancerRep.Find(id)
	if err != nil {
		return nil, errors.Wrapf(err, "freelancerRep.Find()")
	}

	ouF, err := u.InsertLocation(f.F)
	exOuFreelancer := &model.ExtendedOutputFreelancer{
		OuFreel: ouF,
		FirstName:  f.FirstName,
		SecondName: f.SecondName,
	}
	if err != nil {
		return nil, errors.Wrap(err, "InsertLocation()")
	}

	return exOuFreelancer, nil
}

func (u *FreelancerUsecase) Edit(userID int64, new *model.Freelancer) error {
	freelancer, err := u.FindByUser(userID)
	if err != nil {
		return err
	}

	new.ID = freelancer.ID
	new.AccountId = freelancer.ID

	if err := u.freelancerRep.Edit(new); err != nil {
		return errors.Wrapf(err, "freelancerRep.Edit()")
	}
	return nil
}

func (u *FreelancerUsecase) PatternSearch(pattern string, params model.SearchParams) ([]model.ExtendFreelancer, error) {
	if (params.MaxGrade < params.MinGrade) ||
		(params.Limit < 0) {
		return nil, errors.New("wrong input parameters")
	}

	exFreelancers, err := u.freelancerRep.ListOnPattern(pattern, params)
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

func (u *FreelancerUsecase) FindNoLocation(id int64) (*model.ExtendFreelancer, error) {
	f, err := u.freelancerRep.Find(id)
	if err != nil {
		return nil, errors.Wrapf(err, "freelancerRep.Find()")
	}
	return f, nil
}

func (u *FreelancerUsecase) FindNoLocationByUser(id int64) (*model.Freelancer, error) {
	f, err := u.freelancerRep.FindByUser(id)
	if err != nil {
		return nil, errors.Wrapf(err, "freelancerRep.Find()")
	}
	return f, nil
}

func (u *FreelancerUsecase) GetRating(id int64) (int64, error) {
	return 0, nil
}
