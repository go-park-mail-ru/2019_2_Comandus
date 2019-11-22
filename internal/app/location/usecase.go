package location

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	CountryList() ([]*model.Country, error)
	CityListByCountry(id int64) ([]*model.City, error)
	GetCountry(id int64) (*model.Country, error)
	GetCity(id int64) (*model.City, error)
}
