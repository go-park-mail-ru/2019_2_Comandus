package location

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Repository interface {
	CountryList() ([]*model.Country, error)
	CityListByCountry(id int64) ([]*model.City, error)
}
