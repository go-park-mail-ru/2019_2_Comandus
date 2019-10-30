package teststore

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJobRepository_Create(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users", "managers", "jobs")

	store := sqlstore.New(db)

	u := testUser(t)
	u.Email = "jobrepository1@example.org"
	if err := store.User().Create(u); err != nil {
		t.Fatal()
	}

	m := testManager(t, u)
	if err := store.Manager().Create(m); err != nil {
		t.Fatal()
	}

	j := testJob(t, m)
	if err := store.Job().Create(j, m); err != nil {
		t.Fatal()
	}
}

func TestJobRepository_Find(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users", "managers", "jobs")

	store := sqlstore.New(db)
	u := testUser(t)
	u.Email = "jobrepository2@example.org"
	if err := store.User().Create(u); err != nil {
		t.Fatal()
	}

	m := testManager(t, u)
	if err := store.Manager().Create(m); err != nil {
		t.Fatal()
	}

	j := testJob(t, m)
	if err := store.Job().Create(j, m); err != nil {
		t.Fatal()
	}

	j1, err := store.Job().Find(j.ID)
	assert.NoError(t, err)
	assert.NotNil(t, j1)
}

func TestJobRepository_Edit(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users", "managers", "jobs")

	store := sqlstore.New(db)
	u := testUser(t)
	u.Email = "jobrepository3@example.org"
	if err := store.User().Create(u); err != nil {
		t.Fatal()
	}

	m := testManager(t, u)
	if err := store.Manager().Create(m); err != nil {
		t.Fatal()
	}

	j := testJob(t, m)
	if err := store.Job().Create(j, m); err != nil {
		t.Fatal()
	}

	j.City = "London"
	err := store.Job().Edit(j)
	assert.NoError(t, err)
}