package locationUcase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/location"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
)

type LocationUsecase struct {
	locationRep location.Repository
}

func NewLocationUsecase(r location.Repository) location.Usecase {
	return &LocationUsecase{
		locationRep: r,
	}
}

func (u *LocationUsecase) CountryList() ([]*model.Country, error) {
	list, err := u.locationRep.CountryList()
	if err != nil {
		return nil, errors.Wrap(err, "locationRep.CountryList()")
	}
	return list, nil
}

func (u *LocationUsecase) CityListByCountry(id int64) ([]*model.City, error) {
	list, err := u.locationRep.CityListByCountry(id)
	if err != nil {
		return nil, errors.Wrap(err, "locationRep.CityListByCountry()")
	}
	return list, nil
}

func (u *LocationUsecase) GetCountry(id int64) (*model.Country, error) {
	country, err := u.locationRep.FindCountry(id)
	if err != nil {
		return nil, errors.Wrap(err, "localRep.FindCountry()")
	}
	return country, nil
}

func (u *LocationUsecase) GetCity(id int64) (*model.City, error) {
	city, err := u.locationRep.FindCity(id)
	if err != nil {
		return nil, errors.Wrap(err, "localRep.FindCity()")
	}
	return city, nil
}
