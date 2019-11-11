package companyRepository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"reflect"
	"testing"
)

func testCompany(t *testing.T) *model.Company {
	t.Helper()
	return &model.Company{
		ID:          1,
		CompanyName: "mocks company",
		Site:        "www.testcompany.com",
		TagLine:     "no tags",
		Description: "no description",
		Country:     "russia",
		City:        "moscow",
		Address:     "baumanskaya ul",
		Phone:       "88888888888",
	}
}

func TestCompanyRep_Create(t *testing.T) {
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

	repo := NewCompanyRepository(db)

	rows := sqlmock.
		NewRows([]string{"id"})

	var elemID int64 = 1
	expect := []*model.Company{
		{ID: elemID},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	c := testCompany(t)

	// TODO: uncomment when validation will be implemented
	/*if err := c.Validate(); err != nil {
		t.Fatal()
	}*/

	//ok query
	mock.
		ExpectQuery(`INSERT INTO companies`).
		WithArgs(c.CompanyName, c.Site, c.TagLine, c.Description, c.Country, c.City, c.Address, c.Phone).
		WillReturnRows(rows)

	err = repo.Create(c)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if c.ID != 1 {
		t.Errorf("bad id: want %v, have %v", c.ID, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectQuery(`INSERT INTO companies`).
		WithArgs(c.CompanyName, c.Site, c.TagLine, c.Description, c.Country, c.City, c.Address, c.Phone).
		WillReturnError(fmt.Errorf("bad query"))

	err = repo.Create(c)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCompanyRep_Find(t *testing.T) {
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
	repo := NewCompanyRepository(db)

	var elemID int64 = 1

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "companyName", "site", "tagLine", "description", "country", "city", "address",
			"phone" })

	expect := []*model.Company{
		testCompany(t),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.CompanyName, item.Site, item.TagLine, item.Description, item.Country,
			item.City, item.Address, item.Phone)
	}

	mock.
		ExpectQuery("SELECT id, companyName, site, tagLine, description, country, city, address, " +
			"phone FROM companies WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	item, err := repo.Find(elemID)
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
		ExpectQuery("SELECT id, companyName, site, tagLine, description, country, city, address, " +
			"phone FROM companies WHERE").
		WithArgs(elemID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.Find(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	expect = []*model.Company{
		testCompany(t),
	}

	mock.
		ExpectQuery("SELECT id, companyName, site, tagLine, description, country, city, address, " +
			"phone FROM companies WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	_, err = repo.Find(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestCompanyRep_Edit(t *testing.T) {
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

	repo := NewCompanyRepository(db)

	rows := sqlmock.
		NewRows([]string{"id"})

	var elemID int64 = 1
	expect := []*model.Company{
		{ ID: elemID },
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	c := testCompany(t)

	// TODO: uncomment when validation will be implemented
	/*if err := f.Validate(); err != nil {
		t.Fatal()
	}*/

	//ok query
	c.Country = "England"
	c.City = "London"
	mock.
		ExpectQuery(`UPDATE companies SET`).
		WithArgs(c.CompanyName, c.Site, c.TagLine, c.Description, c.Country, c.City, c.Address, c.Phone, c.ID).
		WillReturnRows(rows)

	err = repo.Edit(c)
	if err != nil {
		t.Fatal(err)
	}

	if c.ID != 1 {
		t.Errorf("bad id: want %v, have %v", c.ID, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}