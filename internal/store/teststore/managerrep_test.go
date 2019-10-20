package teststore

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func testManager(t *testing.T, user * model.User) *model.HireManager {
	t.Helper()
	return &model.HireManager{
		AccountID: 			user.ID,
		RegistrationDate:	time.Now(),
		Location:			"Moscow",
	}
}

func TestManagerRepository_Create(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users", "managers")

	store := sqlstore.New(db)

	u := testUser(t)
	u.Email = "user231@example.org"

	assert.NoError(t, store.User().Create(u))

	m := testManager(t, u)
	assert.NoError(t, store.Manager().Create(m))
	assert.NotNil(t, u)
}

func TestManagerRepository_Find(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users", "managers")

	store := sqlstore.New(db)

	u := testUser(t)
	u.Email = "user2@example.org"
	err := store.User().Create(u)
	if err != nil {
		t.Fatal(err)
	}

	m1 := testManager(t, u)
	err = store.Manager().Create(m1)
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
	u.Email = "user3@example.org"
	if err := store.User().Create(u); err != nil {
		t.Fatal(err)
	}

	m := testManager(t, u)

	m1, err := store.Manager().FindByUser(u.ID)
	assert.EqualError(t, err, "sql: no rows in result set")//store.ErrRecordNotFound.Error())
	assert.Nil(t, m1)

	if err := store.Manager().Create(m); err != nil {
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
	u.Email = "user4@example.org"
	if err := store.User().Create(u); err != nil {
		t.Fatal(err)
	}

	m := testManager(t, u)
	if err := store.Manager().Create(m); err != nil {
		t.Fatal(err)
	}

	m.Location = "London"

	if err := store.Manager().Edit(m); err != nil {
		t.Fatal(err)
	}
}


