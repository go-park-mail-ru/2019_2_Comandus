package freelancerRepository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"reflect"
	"testing"
	"time"
)

func testFreelancer(t *testing.T) *model.Freelancer {
	t.Helper()
	return &model.Freelancer{
		ID:				   1,
		AccountId:         1,
		RegistrationDate:  time.Now(),
		Country:           "Russia",
		City:              "Moscow",
		Address:           "moscow",
		Phone:             "111111111",
	}
}

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

	repo := NewFreelancerRepository(db)

	rows := sqlmock.
		NewRows([]string{"accountId"})

	var elemID int64 = 1
	expect := []*model.Freelancer{
		{ ID: elemID },
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	f := testFreelancer(t)

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

	err = repo.Create(f)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if f.ID != 1 {
		t.Errorf("bad id: want %v, have %v", 1, 1)
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

	err = repo.Create(f)
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

	expect := []*model.Freelancer{
		testFreelancer(t),
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

	repo := NewFreelancerRepository(db)

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
		ExpectQuery("SELECT id, accountId, registrationDate, country, city, address, phone, tagLine, " +
		"overview, experienceLevelId, specialityId FROM freelancers WHERE").
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
	expect = []*model.Freelancer{
		testFreelancer(t),
	}

	mock.
		ExpectQuery("SELECT id, accountId, registrationDate, country, city, address, phone, tagLine, " +
		"overview, experienceLevelId, specialityId FROM freelancers WHERE").
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

	//var elemID int64 = 1

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "accountId", "registrationDate", "country", "city", "address", "phone", "tagLine",
			"overview", "experienceLevelId", "specialityId" })

	expect := []*model.Freelancer{
		testFreelancer(t),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.AccountId, item.RegistrationDate, item.Country, item.City, item.Address,
			item.Phone, item.TagLine, item.Overview, item.ExperienceLevelId, item.SpecialityId)
	}

	mock.
		ExpectQuery("SELECT id, accountId, registrationDate, country, city, address, phone, tagLine, " +
			"overview, experienceLevelId, specialityId FROM freelancers WHERE").
		WithArgs(1).
		WillReturnRows(rows)

	repo := NewFreelancerRepository(db)

	item, err := repo.FindByUser(1)
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
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.FindByUser(1)
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
		testFreelancer(t),
	}

	mock.
		ExpectQuery("SELECT id, accountId, registrationDate, country, city, address, phone, tagLine, " +
			"overview, experienceLevelId, specialityId FROM freelancers WHERE").
		WithArgs(1).
		WillReturnRows(rows)

	_, err = repo.FindByUser(1)
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

	repo := NewFreelancerRepository(db)

	rows := sqlmock.
		NewRows([]string{"accountId"})

	var elemID int64 = 1
	expect := []*model.Freelancer{
		{ ID: elemID },
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	f := testFreelancer(t)
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

	err = repo.Edit(f)
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
