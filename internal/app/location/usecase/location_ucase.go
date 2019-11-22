package locationUcase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/location"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
)

type LocationUsecase struct {
	locationRep   location.Repository
}

func NewLocationUsecase(r location.Repository) location.Usecase {
	return &LocationUsecase{
		locationRep:   r,
	}
}

func (u *LocationUsecase)CountryList() ([]*model.Country, error) {
	list, err := u.locationRep.CountryList()
	if err != nil {
		return nil, errors.Wrap(err, "locationRep.CountryList()")
	}
	return list, nil
}

func (u *LocationUsecase)CityListByCountry(id int64) ([]*model.City, error) {
	list, err := u.locationRep.CityListByCountry(id)
	if err != nil {
		return nil, errors.Wrap(err, "locationRep.CityListByCountry()")
	}
	return list, nil
}