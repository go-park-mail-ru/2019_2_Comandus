package teststore

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"testing"
	"time"
)

// TODO: all helpers in another file test_helpers.go
func testManager(t *testing.T, user * model.User) *model.HireManager {
	t.Helper()
	return &model.HireManager{
		AccountID: 			user.ID,
		RegistrationDate:	time.Now(),
		Location:			"Moscow",
	}
}

func testFreelancer(t *testing.T, user *model.User) *model.Freelancer {
	t.Helper()
	return &model.Freelancer{
		AccountId:         user.ID,
		RegistrationDate:  time.Now(),
		Country:           "Russia",
		City:              "Moscow",
		Address:           "moscow",
		Phone:             "111111111",
	}
}

func testJob(t *testing.T, manager * model.HireManager) *model.Job {
	t.Helper()
	return &model.Job{
		HireManagerId:	manager.ID,
		Title:          "title",
		Description:    "description",
		PaymentAmount:   11222,
		Country: 	"russia",
		City:	"moscow",
	}
}
