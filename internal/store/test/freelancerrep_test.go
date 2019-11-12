package test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
	"reflect"
	"testing"
)

func TestFreelancerRep_Create(t *testing.T) {
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

	store := sqlstore.New(db)
	rows := sqlmock.
		NewRows([]string{"accountId"})

	var elemID int64 = 1
	expect := []*model.Freelancer{
		{ ID: elemID },
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	u := testUser(t)
	u.ID = 1
	f := testFreelancer(t, u)

	// TODO: uncomment when validation will be implemented
	/*if err := f.Validate(); err != nil {
		t.Fatal()
	}*/

	//ok query
	mock.
		ExpectQuery(`INSERT INTO freelancers`).
		WithArgs(f.AccountId, f.RegistrationDate, f.Country, f.City, f.Address, f.Phone, f.TagLine,
			f.Overview, f.ExperienceLevelId, f.SpecialityId).
		WillReturnRows(rows)

	err = store.Freelancer().Create(f)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if f.ID != 1 {
		t.Errorf("bad id: want %v, have %v", u.ID, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectQuery(`INSERT INTO freelancers`).
		WithArgs(f.AccountId, f.RegistrationDate, f.Country, f.City, f.Address, f.Phone, f.TagLine,
			f.Overview, f.ExperienceLevelId, f.SpecialityId).
		WillReturnError(fmt.Errorf("bad query"))

	err = store.Freelancer().Create(f)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFreelancerRep_Find(t *testing.T) {
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

	var elemID int64 = 1

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "accountId", "registrationDate", "country", "city", "address", "phone", "tagLine",
		"overview", "experienceLevelId", "specialityId" })

	u := testUser(t)
	u.ID = elemID + 1
	expect := []*model.Freelancer{
		testFreelancer(t, u),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.AccountId, item.RegistrationDate, item.Country, item.City, item.Address,
			item.Phone, item.TagLine, item.Overview, item.ExperienceLevelId, item.SpecialityId)
	}

	mock.
		ExpectQuery("SELECT id, accountId, registrationDate, country, city, address, phone, tagLine, " +
		"overview, experienceLevelId, specialityId FROM freelancers WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	store := sqlstore.New(db)

	item, err := store.Freelancer().Find(elemID)
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
		ExpectQuery("SELECT id, accountId, registrationDate, country, city, address, phone, tagLine, " +
		"overview, experienceLevelId, specialityId FROM freelancers WHERE").
		WithArgs(elemID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = store.Freelancer().Find(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	expect = []*model.Freelancer{
		testFreelancer(t, u),
	}

	mock.
		ExpectQuery("SELECT id, accountId, registrationDate, country, city, address, phone, tagLine, " +
		"overview, experienceLevelId, specialityId FROM freelancers WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	_, err = store.Freelancer().Find(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestFreelancerRep_FindByUser(t *testing.T) {
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

	var elemID int64 = 1

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "accountId", "registrationDate", "country", "city", "address", "phone", "tagLine",
			"overview", "experienceLevelId", "specialityId" })

	u := testUser(t)
	u.ID = elemID + 1
	expect := []*model.Freelancer{
		testFreelancer(t, u),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.AccountId, item.RegistrationDate, item.Country, item.City, item.Address,
			item.Phone, item.TagLine, item.Overview, item.ExperienceLevelId, item.SpecialityId)
	}

	mock.
		ExpectQuery("SELECT id, accountId, registrationDate, country, city, address, phone, tagLine, " +
			"overview, experienceLevelId, specialityId FROM freelancers WHERE").
		WithArgs(u.ID).
		WillReturnRows(rows)

	store := sqlstore.New(db)

	item, err := store.Freelancer().FindByUser(u.ID)
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
		ExpectQuery("SELECT id, accountId, registrationDate, country, city, address, phone, tagLine, " +
			"overview, experienceLevelId, specialityId FROM freelancers WHERE").
		WithArgs(u.ID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = store.Freelancer().FindByUser(u.ID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	expect = []*model.Freelancer{
		testFreelancer(t, u),
	}

	mock.
		ExpectQuery("SELECT id, accountId, registrationDate, country, city, address, phone, tagLine, " +
			"overview, experienceLevelId, specialityId FROM freelancers WHERE").
		WithArgs(u.ID).
		WillReturnRows(rows)

	_, err = store.Freelancer().FindByUser(u.ID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestFreelancerRep_Edit(t *testing.T) {
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

	store := sqlstore.New(db)

	rows := sqlmock.
		NewRows([]string{"accountId"})

	var elemID int64 = 1
	expect := []*model.Freelancer{
		{ ID: elemID },
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	u := testUser(t)
	u.ID = 1
	f := testFreelancer(t, u)
	f.ID = 1

	// TODO: uncomment when validation will be implemented
	/*if err := f.Validate(); err != nil {
		t.Fatal()
	}*/

	//ok query
	f.Country = "England"
	f.City = "London"
	mock.
		ExpectQuery(`UPDATE freelancers SET`).
		WithArgs(f.Country, f.City, f.Address, f.Phone, f.TagLine,
		f.Overview, f.ExperienceLevelId, f.SpecialityId, f.ID).
		WillReturnRows(rows)

	err = store.Freelancer().Edit(f)
	if err != nil {
		t.Fatal(err)
	}

	if u.ID != 1 {
		t.Errorf("bad id: want %v, have %v", u.ID, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
