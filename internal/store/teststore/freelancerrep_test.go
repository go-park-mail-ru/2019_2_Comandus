package teststore

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFreelancerRep_Create(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users", "freelancers")

	store := sqlstore.New(db)

	u := testUser(t)
	u.Email = "freelancerrepository1@example.org"
	if _, err := store.User().Create(u); err != nil {
		t.Fatal()
	}

	f := testFreelancer(t, u)
	if _, err := store.Freelancer().Create(f); err != nil {
		t.Fatal()
	}
}

// TODO: add not found case
func TestFreelancerRep_Find(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("freelancers")

	store := sqlstore.New(db)

	u := testUser(t)
	u.Email = "freelancerrepository2@example.org"
	if _, err := store.User().Create(u); err != nil {
		t.Fatal()
	}

	f := testFreelancer(t, u)
	if _, err := store.Freelancer().Create(f); err != nil {
		t.Fatal()
	}

	f1, err := store.Freelancer().Find(f.ID)
	assert.NotNil(t, f1)
	assert.NoError(t, err)
}

func TestFreelancerRep_FindByUser(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users", "freelancers")

	store := sqlstore.New(db)

	u := testUser(t)
	u.Email = "freelancerrepository3@example.org"
	if _, err := store.User().Create(u); err != nil {
		t.Fatal()
	}

	f := testFreelancer(t, u)
	if _, err := store.Freelancer().Create(f); err != nil {
		t.Fatal()
	}

	f1, err := store.Freelancer().FindByUser(u.ID)
	assert.NotNil(t, f1)
	assert.NoError(t, err)
}

func TestFreelancerRep_Edit(t *testing.T) {
	db, teardown := testStore(t, databaseURL)
	defer teardown("users", "freelancers")

	store := sqlstore.New(db)

	u := testUser(t)
	u.Email = "freelancerrepository4@example.org"
	if _, err := store.User().Create(u); err != nil {
		t.Fatal()
	}

	f := testFreelancer(t, u)
	if _, err := store.Freelancer().Create(f); err != nil {
		t.Fatal()
	}
	f.City = "London"

	if _, err := store.Freelancer().Edit(f); err != nil {
		t.Fatal()
	}
}

