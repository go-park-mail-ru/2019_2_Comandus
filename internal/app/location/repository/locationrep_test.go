package locationRepository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"reflect"
	"testing"
)

func testCountry(t *testing.T) *model.Country {
	return &model.Country{
		ID:   1,
		Name: "Russia",
	}
}

func testCity(t *testing.T) *model.City {
	return &model.City{
		ID:        1,
		CountryID: 1,
		Name:      "Moscow",
	}
}

func TestLocationRepository_CountryList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "name"})

	expect := []*model.Country{
		testCountry(t),
		testCountry(t),
		testCountry(t),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Name)
	}

	mock.
		ExpectQuery("SELECT id, name FROM country").
		WillReturnRows(rows)

	repo := NewLocationRepository(db)

	countries, err := repo.CountryList()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	for i := 0; i < 3; i++ {
		if countries[i].Name != expect[i].Name || countries[i].ID != expect[i].ID {
			t.Errorf("results not match, want %v, have %v", expect[i], countries[i])
			return
		}
	}
}

func TestLocationRepository_CityListByCountry(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	var countryId int64 = 1
	// good query
	rows := sqlmock.
		NewRows([]string{"id", "country_id", "name"})

	expect := []*model.City{
		testCity(t),
		testCity(t),
		testCity(t),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.CountryID, item.Name)
	}

	mock.
		ExpectQuery("SELECT city.id, region.country_id, city.name" +
		"FROM city " +
		"JOIN region " +
		"ON").
		WithArgs(countryId).
		WillReturnRows(rows)

	repo := NewLocationRepository(db)

	cities, err := repo.CityListByCountry(countryId)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	for i := 0; i < 3; i++ {
		if cities[i].Name != expect[i].Name ||
			cities[i].ID != expect[i].ID ||
			cities[i].CountryID != expect[i].CountryID {
			t.Errorf("results not match, want %v, have %v", expect[i], cities[i])
			return
		}
	}
}

func TestLocationRepository_FindCity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	var cityId int64 = 1

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "country_id", "name"})

	expect := []*model.City{
		testCity(t),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.CountryID, item.Name)
	}

	mock.
		ExpectQuery("SELECT id, country_id, name FROM city WHERE").
		WithArgs(cityId).
		WillReturnRows(rows)

	repo := NewLocationRepository(db)

	item, err := repo.FindCity(cityId)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(item, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], item)
		return
	}

	// query error
	mock.
		ExpectQuery("SELECT id, country_id, name FROM city WHERE").
		WithArgs(cityId).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.FindCity(cityId)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	mock.
		ExpectQuery("SELECT id, country_id, name FROM city WHERE").
		WithArgs(cityId).
		WillReturnRows(rows)

	_, err = repo.FindCity(cityId)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestLocationRepository_FindCountry(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	var countryId int64 = 1

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "name"})

	expect := []*model.Country{
		testCountry(t),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Name)
	}

	mock.
		ExpectQuery("SELECT id, name FROM country WHERE").
		WithArgs(countryId).
		WillReturnRows(rows)

	repo := NewLocationRepository(db)

	item, err := repo.FindCountry(countryId)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(item, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], item)
		return
	}

	// query error
	mock.
		ExpectQuery("SELECT id, name FROM country WHERE").
		WithArgs(countryId).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.FindCountry(countryId)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	mock.
		ExpectQuery("SELECT id, name FROM country WHERE").
		WithArgs(countryId).
		WillReturnRows(rows)

	_, err = repo.FindCountry(countryId)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}