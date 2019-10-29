package teststore

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManagerRepository_Create(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users", "managers")

	store := sqlstore.New(db)

	u := testUser(t)
	u.Email = "managerrep1@example.org"

	_, err := store.User().Create(u)
	assert.NoError(t, err)

	m := testManager(t, u)
	_, err = store.Manager().Create(m)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestManagerRepository_Find(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users", "managers")

	store := sqlstore.New(db)

	u := testUser(t)
	u.Email = "managerrep2@example.org"
	_, err := store.User().Create(u)
	if err != nil {
		t.Fatal(err)
	}

	m1 := testManager(t, u)
	_, err = store.Manager().Create(m1)
	if err != nil {
		t.Fatal(err)
	}

	m2, err := store.Manager().Find(m1.ID)

	assert.NoError(t, err)
	assert.NotNil(t, m2)
}

func TestManagerRepository_FindByUser(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users", "managers")

	store := sqlstore.New(db)

	u := testUser(t)
	u.Email = "managerrep3@example.org"
	if _, err := store.User().Create(u); err != nil {
		t.Fatal(err)
	}

	m := testManager(t, u)

	m1, err := store.Manager().FindByUser(u.ID)
	assert.EqualError(t, err, "sql: no rows in result set")//store.ErrRecordNotFound.Error())
	assert.Nil(t, m1)

	if _, err := store.Manager().Create(m); err != nil {
		t.Fatal(err)
	}

	m2, err := store.Manager().FindByUser(u.ID)
	assert.NoError(t, err)
	assert.NotNil(t, m2)
}

func TestManagerRepository_Edit(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users", "managers")

	store := sqlstore.New(db)

	u := testUser(t)
	u.Email = "managerrep4@example.org"
	if _, err := store.User().Create(u); err != nil {
		t.Fatal(err)
	}

	m := testManager(t, u)
	if _, err := store.Manager().Create(m); err != nil {
		t.Fatal(err)
	}

	m.Location = "London"

	if _, err := store.Manager().Edit(m); err != nil {
		t.Fatal(err)
	}
}


