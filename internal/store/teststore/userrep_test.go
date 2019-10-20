package teststore

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testUser(t *testing.T) *model.User {
	t.Helper()
	return &model.User{
		Email:    "user@example.org",
		Password: "password",
	}
}

func TestUserRepository_Create(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users")

	store := sqlstore.New(db)

	u := testUser(t)
	u.Email = "userrep1@example.org"
	assert.NoError(t, store.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_Find(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users")

	store := sqlstore.New(db)

	u1 := testUser(t)
	u1.Email = "userrep2@example.org"

	err := store.User().Create(u1)
	if err != nil {
		t.Fatal(err)
	}

	u2, err := store.User().Find(u1.ID)

	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users")

	store := sqlstore.New(db)


	u1 := testUser(t)
	u1.Email = "userrep3@example.org"
	_, err := store.User().FindByEmail(u1.Email)
	assert.EqualError(t, err, "sql: no rows in result set")//store.ErrRecordNotFound.Error())

	err = store.User().Create(u1)
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
	if err := store.User().Create(u); err != nil {
		t.Fatal(err)
	}

	u.SecondName = "Second name"

	if err := store.User().Edit(u); err != nil {
		t.Fatal(err)
	}
}