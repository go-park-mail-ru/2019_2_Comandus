package freelancerUcase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/mocks/repository_mocks"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testUcase(t *testing.T) (*repository_mocks.MockFreelancerRepository, freelancer.Usecase) {
	t.Helper()
	ctrl := gomock.NewController(t)
	freelancerRep := repository_mocks.NewMockFreelancerRepository(ctrl)
	freelancerUcase := NewFreelancerUsecase(freelancerRep)
	return freelancerRep, freelancerUcase
}

func TestNewFreelancerUsecase_Create(t *testing.T) {
	freelancerRep, freelancerUcase := testUcase(t)

	user := &model.User{
		ID:        1,
		FirstName: "ddd",
		Email:     "ddd@hj.cv",
	}

	freelancer := &model.Freelancer{
		ID:        1,
		AccountId: user.ID,
	}

	testCases := []struct {
		name          string
		newFreelancer *model.Freelancer
		expectError   error
	}{
		{
			name: "valid",
			newFreelancer: &model.Freelancer{
				ID:        1,
				AccountId: user.ID,
			},
			expectError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			freelancerRep.
				EXPECT().
				Create(&model.Freelancer{
					AccountId: user.ID,
				}).Do(func(arg *model.Freelancer) {
				arg.ID = 1
			}).
				Return(tc.expectError)
			fr, err := freelancerUcase.Create(user.ID)

			if tc.expectError == nil {
				assert.Equal(t, freelancer, fr)
				assert.Equal(t, nil, err)
				return
			}

			if err != nil {
				assert.Equal(t, tc.expectError.Error(), err.Error())
				return
			}

			t.Fatal()
		})
	}

}

func TestFreelancerUsecase_Edit(t *testing.T) {
	freelancerRep, freelancerUcase := testUcase(t)

	user := &model.User{
		ID:        1,
		FirstName: "ddd",
		Email:     "ddd@hj.cv",
	}

	freelancer := &model.Freelancer{
		ID:        1,
		AccountId: user.ID,
	}

	testCases := []struct {
		name          string
		newFreelancer *model.Freelancer
		expectError   error
	}{
		{
			name: "valid",
			newFreelancer: &model.Freelancer{
				ID:        freelancer.ID,
				AccountId: user.ID,
				Country:   "russia",
				City:      "moscow",
			},
			expectError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			freelancerRep.
				EXPECT().
				Edit(tc.newFreelancer).
				Return(tc.expectError)

			err := freelancerUcase.Edit(tc.newFreelancer, freelancer)

			if tc.expectError == nil {
				assert.Equal(t, nil, err)
				return
			}

			if err != nil {
				assert.Equal(t, tc.expectError.Error(), err.Error())
				return
			}

			t.Fatal()
		})
	}
}

func TestFreelancerUsecase_Find(t *testing.T) {
	freelancerRep, freelancerUcase := testUcase(t)

	user := &model.User{
		ID:        1,
		FirstName: "ddd",
		Email:     "ddd@hj.cv",
	}

	freelancer := &model.Freelancer{
		ID:        1,
		AccountId: user.ID,
	}

	testCases := []struct {
		name          string
		newFreelancer *model.Freelancer
		expectError   error
	}{
		{
			name: "valid",
			newFreelancer: &model.Freelancer{
				ID:        1,
				AccountId: user.ID,
			},
			expectError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			freelancerRep.
				EXPECT().
				Find(testCases[0].newFreelancer.ID).
				Return(freelancer, tc.expectError)
			fr, err := freelancerUcase.Find(testCases[0].newFreelancer.ID)

			if tc.expectError == nil {
				assert.Equal(t, freelancer, fr)
				assert.Equal(t, nil, err)
				return
			}

			if err != nil {
				assert.Equal(t, tc.expectError.Error(), err.Error())
				return
			}

			t.Fatal()
		})
	}
}

func TestFreelancerUsecase_FindByUser(t *testing.T) {
	freelancerRep, freelancerUcase := testUcase(t)

	user := &model.User{
		ID:        1,
		FirstName: "ddd",
		Email:     "ddd@hj.cv",
	}

	freelancer := &model.Freelancer{
		ID:        1,
		AccountId: user.ID,
	}

	testCases := []struct {
		name          string
		newFreelancer *model.Freelancer
		expectError   error
	}{
		{
			name: "valid",
			newFreelancer: &model.Freelancer{
				ID:        1,
				AccountId: user.ID,
			},
			expectError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			freelancerRep.
				EXPECT().
				FindByUser(testCases[0].newFreelancer.AccountId).
				Return(freelancer, tc.expectError)
			fr, err := freelancerUcase.FindByUser(testCases[0].newFreelancer.AccountId)

			if tc.expectError == nil {
				assert.Equal(t, freelancer, fr)
				assert.Equal(t, nil, err)
				return
			}

			if err != nil {
				assert.Equal(t, tc.expectError.Error(), err.Error())
				return
			}

			t.Fatal()
		})
	}
}
