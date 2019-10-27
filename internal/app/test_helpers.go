package apiserver

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"log"
	"strings"
	"testing"
	"time"
)

const (
	databaseURL = "host=localhost dbname=restapi_dev sslmode=disable port=5432 password=1234 user=d"
)

func testStore(t *testing.T, databaseURL string) (*sql.DB, func(...string)) {
	t.Helper()
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Println("fail open sql")
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Println("fail ping sql")
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			_, err = db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE;", strings.Join(tables, ", ")))
			if err != nil {
				log.Println("fail truncate")
				t.Fatal(err)
			}
		}

		if err := db.Close(); err != nil {
			log.Println("fail close sql")
			t.Fatal(err)
		}
	}
}

func testUser(t *testing.T) *model.User {
	t.Helper()
	return &model.User{
		Email:    "user@example.org",
		Password: "password",
	}
}

func testFreelancer(t *testing.T, user *model.User) *model.Freelancer {
	t.Helper()
	return &model.Freelancer{
		AccountId:        user.ID,
		RegistrationDate: time.Now(),
		Country:          "Russia",
		City:             "Moscow",
		Address:          "my address",
		Phone:            "my phone",
		Overview:         "overview",
	}
}

func testManager(t *testing.T, user *model.User) *model.HireManager {
	t.Helper()
	return &model.HireManager{
		AccountID:        user.ID,
		RegistrationDate: time.Now(),
		Location:         "Moscow",
	}
}

func (s *server) addUser2Server(t *testing.T) error {
	t.Helper()

	u := testUser(t)

	err := s.store.User().Create(u)
	log.Println(u.ID)
	if err != nil {
		return err
	}
	log.Println("ID: ", u.ID)

	m := testManager(t, u)
	err = s.store.Manager().Create(m)
	if err != nil {
		return err
	}

	f := testFreelancer(t, u)
	err = s.store.Freelancer().Create(f)
	if err != nil {
		return err
	}
	return nil
}

func (s *server) addJob2Server(t *testing.T) {
	t.Helper()

	/*j := model.Job{
		ID:                0,
		HireManagerId:     0,
		Title:             "first job",
		Description:       "work hard",
		Files:             "",
		SpecialityId:      0,
		ExperienceLevelId: 0,
		PaymentAmout:      0,
		Country:           "Russia",
		City:              "Moscow",
		JobTypeId:         0,
	}*/

	// TODO: add to db when jobs create func is impl
}
