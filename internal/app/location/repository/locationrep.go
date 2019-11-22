package locationRepository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/location"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
)

type LocationRepository struct {
	db	*sql.DB
}

func NewLocationRepository(db *sql.DB) location.Repository {
	return &LocationRepository{db}
}

func (r *LocationRepository) CountryList() ([]*model.Country, error) {
	var countries []*model.Country
	rows, err := r.db.Query("SELECT id, name " +
		"FROM country")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := &model.Country{}
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		countries = append(countries, c)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	return countries, nil
}

func (r *LocationRepository) CityListByCountry(id int64) ([]*model.City, error) {
	var cities []*model.City
	rows, err := r.db.Query("SELECT city.id, region.country_id, city.name" +
		"FROM city " +
		"JOIN region " +
		"ON region.country_id = $1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := &model.City{}
		err := rows.Scan(&c.ID, &c.CountryID, &c.Name)
		if err != nil {
			return nil, err
		}
		cities = append(cities, c)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return cities, nil
}

func (r *LocationRepository) FindCountry(id int64) (*model.Country, error) {
	country := &model.Country{}
	if err := r.db.QueryRow(
		"SELECT id, name FROM country WHERE id = $1",
		id,
	).Scan(
		&country.ID,
		&country.Name,
	); err != nil {
		return nil, err
	}
	return country, nil
}

func (r *LocationRepository) FindCity(id int64) (*model.City, error) {
	city := &model.City{}
	if err := r.db.QueryRow(
		"SELECT id, country_id, name FROM city WHERE id = $1",
		id,
	).Scan(
		&city.ID,
		&city.CountryID,
		&city.Name,
	); err != nil {
		return nil, err
	}
	return city, nil
}