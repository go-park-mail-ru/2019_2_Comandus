package managerRepository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"reflect"
	"testing"
	"time"
)

func testManager(t *testing.T) *model.HireManager {
	t.Helper()
	return &model.HireManager{
		ID:					1,
		AccountID: 			1,
		RegistrationDate:	time.Now(),
		Location:			"Moscow",
	}
}

func TestManagerRepository_Create(t *testing.T) {
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

	repo := NewManagerRepository(db)
	rows := sqlmock.
		NewRows([]string{"accountId"})

	var elemID int64 = 1
	expect := []*model.HireManager{
		{ ID: elemID },
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	m := testManager(t)

	// TODO: uncomment when validation will be implemented
	/*if err := m.Validate(); err != nil {
		t.Fatal()
	}*/

	//ok query
	// id, accountId, registrationDate, location, companyId FROM managers WHERE accountId = $1
	mock.
		ExpectQuery(`INSERT INTO managers`).
		WithArgs(m.AccountID, m.RegistrationDate, m.Location, m.CompanyID).
		WillReturnRows(rows)

	err = repo.Create(m)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if m.ID != 1 {
		t.Errorf("bad id: want %v, have %v", 1, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectQuery(`INSERT INTO managers`).
		WithArgs(m.AccountID, m.RegistrationDate, m.Location, m.CompanyID).
		WillReturnError(fmt.Errorf("bad query"))

	err = repo.Create(m)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestManagerRepository_Find(t *testing.T) {
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
		NewRows([]string{"id", "accountId", "registrationDate", "location", "companyId" })

	expect := []*model.HireManager{
		testManager(t),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.AccountID, item.RegistrationDate, item.Location, item.CompanyID)
	}

	mock.
		ExpectQuery("SELECT id, accountId, registrationDate, location, companyId FROM managers WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	repo := NewManagerRepository(db)

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
		ExpectQuery("SELECT id, accountId, registrationDate, location, companyId FROM managers WHERE").
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
	expect = []*model.HireManager{
		testManager(t),
	}

	mock.
		ExpectQuery("SELECT id, accountId, registrationDate, location, companyId FROM managers WHERE").
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

func TestManagerRepository_FindByUser(t *testing.T) {
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
		NewRows([]string{"id", "accountId", "registrationDate", "location", "companyId" })

	expect := []*model.HireManager{
		testManager(t),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.AccountID, item.RegistrationDate, item.Location, item.CompanyID)
	}

	mock.
		ExpectQuery("SELECT id, accountId, registrationDate, location, companyId FROM managers WHERE").
		WithArgs(1).
		WillReturnRows(rows)

	repo := NewManagerRepository(db)

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
		ExpectQuery("SELECT id, accountId, registrationDate, location, companyId FROM managers WHERE").
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
	expect = []*model.HireManager{
		testManager(t),
	}

	mock.
		ExpectQuery("SELECT id, accountId, registrationDate, location, companyId FROM managers WHERE").
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

func TestManagerRepository_Edit(t *testing.T) {
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

	repo := NewManagerRepository(db)

	rows := sqlmock.
		NewRows([]string{"accountId"})

	var elemID int64 = 1
	expect := []*model.HireManager{
		{ ID: elemID },
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	m := testManager(t)
	m.ID = 1

	// TODO: uncomment when validation will be implemented
	/*if err := f.Validate(); err != nil {
		t.Fatal()
	}*/

	//ok query
	m.Location = "underwater"

	mock.
		ExpectQuery(`UPDATE managers SET`).
		WithArgs(m.Location, m.CompanyID, m.ID).
		WillReturnRows(rows)

	err = repo.Edit(m)
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
