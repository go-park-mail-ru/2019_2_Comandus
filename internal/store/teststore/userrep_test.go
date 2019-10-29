package teststore

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	//"database/sql"
	//"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
	"github.com/stretchr/testify/assert"
	//"net/http"
	//"net/http/httptest"
	"reflect"
	"testing"
)

func testUser(t *testing.T) *model.User {
	t.Helper()
	return &model.User{
		ID: 1,
		FirstName: "masha",
		SecondName: "ivanova",
		UserName: "masha1996",
		Email: "masha@mail.ru",
		Password: "123456",
		EncryptPassword: "",
		Avatar: nil,
		UserType: "freelancer",
	}
}

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	store := sqlstore.New(db)

	u := testUser(t)
	if err := u.Validate(); err != nil {
		t.Fatal()
	}

	if err := u.BeforeCreate(); err != nil {
		t.Fatal()
	}

	//ok query
	mock.
		ExpectExec(`INSERT INTO users`).
		WithArgs(u.FirstName, u.SecondName, u.UserName, u.Email, u.EncryptPassword, u.UserType).
		WillReturnResult(sqlmock.NewResult(1,1))

	lastId, err := store.User().Create(u)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if lastId != 1 {
		t.Errorf("bad id: want %v, have %v", lastId, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectExec(`INSERT INTO users`).
		WithArgs(u.FirstName, u.SecondName, u.UserName, u.Email, u.EncryptPassword, u.UserType).
		WillReturnError(fmt.Errorf("bad query"))

	_, err = store.User().Create(u)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// result error
	mock.
		ExpectExec(`INSERT INTO users`).
		WithArgs(u.FirstName, u.SecondName, u.UserName, u.Email, u.EncryptPassword, u.UserType).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("bad_result")))

	_, err = store.User().Create(u)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Find(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemID int64 = 2

	// good query
	rows := sqlmock.
		NewRows([]string{"accountId", "firsName", "secondName", "username", "email", "password", "encryptPassword",
			"avatar", "userType"})
	expect := []*model.User{
		{ 	ID: elemID,
			FirstName: "masha",
			SecondName: "ivanova",
			UserName: "masha1996",
			Email: "masha1@mail.ru",
			Password: "123456",
			EncryptPassword: "",
			Avatar: nil,
			UserType: "freelancer"},
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

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users")

	store := sqlstore.New(db)


	u1 := testUser(t)
	u1.Email = "userrep3@example.org"
	_, err := store.User().FindByEmail(u1.Email)
	assert.EqualError(t, err, "sql: no rows in result set")//store.ErrRecordNotFound.Error())

	_, err = store.User().Create(u1)
	if err != nil {
		t.Fatal(err)
	}

	u2, err := store.User().FindByEmail(u1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestUserRepository_Edit(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users")

	store := sqlstore.New(db)

	u := testUser(t)
	u.Email = "userrep4@example.org"
	if _, err := store.User().Create(u); err != nil {
		t.Fatal(err)
	}

	u.SecondName = "Second name"

	if _, err := store.User().Edit(u); err != nil {
		t.Fatal(err)
	}
}