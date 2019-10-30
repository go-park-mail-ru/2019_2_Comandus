package test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	//"database/sql"
	//"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
	//"net/http"
	//"net/http/httptest"
	"reflect"
	"testing"
)

func TestUserRepository_Find(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemID int64 = 1

	// good query
	rows := sqlmock.
		NewRows([]string{"accountId", "firsName", "secondName", "username", "email", "password", "encryptPassword",
			"avatar", "userType"})
	expect := []*model.User{
		testUser(t),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.FirstName, item.SecondName, item.UserName, item.Email, item.Password,
			item.EncryptPassword, item.Avatar, item.UserType)
	}

	mock.
		ExpectQuery("SELECT accountId, firstName, secondName, username, email, " +
			"'' as password, encryptPassword, avatar, userType  FROM users WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	store := sqlstore.New(db)
	repo := store.User()

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
		ExpectQuery("SELECT accountId, firstName, secondName, username, email, " +
			"'' as password, encryptPassword, avatar, userType  FROM users WHERE").
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
	rows = sqlmock.NewRows([]string{"id", "firstName", "secondName", "username", "email" }).
		AddRow(1, "masha", "ivanova", "masha1996", "masha@mail.ru")

	mock.
		ExpectQuery("SELECT accountId, firstName, secondName, username, email, " +
			"'' as password, encryptPassword, avatar, userType  FROM users WHERE").
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

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	store := sqlstore.New(db)

	rows := sqlmock.
		NewRows([]string{"accountId"})

	var elemID int64 = 1
	expect := []*model.User{
		{ ID: elemID },
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	u := testUser(t)
	if err := u.Validate(); err != nil {
		t.Fatal()
	}
	if err := u.BeforeCreate(); err != nil {
		t.Fatal()
	}

	//ok query
	mock.
		ExpectQuery(`INSERT INTO users`).
		WithArgs(u.FirstName, u.SecondName, u.UserName, u.Email, u.EncryptPassword, u.UserType).
		WillReturnRows(rows)

	err = store.User().Create(u)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if u.ID != 1 {
		t.Errorf("bad id: want %v, have %v", u.ID, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectQuery(`INSERT INTO users`).
		WithArgs(u.FirstName, u.SecondName, u.UserName, u.Email, u.EncryptPassword, u.UserType).
		WillReturnError(fmt.Errorf("bad query"))

	err = store.User().Create(u)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// good query
	rows := sqlmock.
		NewRows([]string{"accountId", "firsName", "secondName", "username", "email", "password", "encryptPassword",
			"avatar", "userType"})
	expect := []*model.User{
		testUser(t),
	}

	u := testUser(t)

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.FirstName, item.SecondName, item.UserName, item.Email, item.Password,
			item.EncryptPassword, item.Avatar, item.UserType)
	}

	mock.
		ExpectQuery("SELECT accountId, firstName, secondName, username, email, " +
			"'' as password, encryptPassword, avatar, userType  FROM users WHERE").
		WithArgs(u.Email).
		WillReturnRows(rows)

	store := sqlstore.New(db)
	repo := store.User()

	item, err := repo.FindByEmail(u.Email)
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
		ExpectQuery("SELECT accountId, firstName, secondName, username, email, " +
			"'' as password, encryptPassword, avatar, userType  FROM users WHERE").
		WithArgs(u.Email).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.FindByEmail(u.Email)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	rows = sqlmock.NewRows([]string{"id", "firstName", "secondName", "username", "email" }).
		AddRow(1, "masha", "ivanova", "masha1996", "masha@mail.ru")

	mock.
		ExpectQuery("SELECT accountId, firstName, secondName, username, email, " +
			"'' as password, encryptPassword, avatar, userType  FROM users WHERE").
		WithArgs(u.Email).
		WillReturnRows(rows)

	_, err = repo.FindByEmail(u.Email)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestUserRepository_Edit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	store := sqlstore.New(db)

	rows := sqlmock.
		NewRows([]string{"accountId"})

	var elemID int64 = 1
	expect := []*model.User{
		{ ID: elemID },
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	u := testUser(t)
	if err := u.Validate(); err != nil {
		t.Fatal()
	}
	if err := u.BeforeCreate(); err != nil {
		t.Fatal()
	}

	//ok query
	u.UserName = "dasha"
	mock.
		ExpectQuery(`UPDATE users SET`).
		WithArgs(u.FirstName, u.SecondName, u.UserName, u.EncryptPassword, u.Avatar, u.UserType, u.ID).
		WillReturnRows(rows)

	err = store.User().Edit(u)
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