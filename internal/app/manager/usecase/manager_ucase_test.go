package managerUcase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/mocks/repository_mocks"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testUcase(t *testing.T) (*repository_mocks.MockManagerRepository, manager.Usecase) {
	t.Helper()
	ctrl := gomock.NewController(t)
	managerRep := repository_mocks.NewMockManagerRepository(ctrl)
	managerUcase := NewManagerUsecase(managerRep)
	return managerRep, managerUcase
}

func TestManagerUsecase_Create(t *testing.T) {
	managerRep, managerUcase := testUcase(t)

	testCases := []struct {
		name        string
		newCompany  *model.Company
		userType    string
		userID      int64
		compID      int64
		expectError error
	}{
		{
			name:        "valid",
			newCompany:  &model.Company{},
			userType:    model.UserCustomer,
			expectError: nil,
			userID:      1,
			compID:      1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			managerRep.
				EXPECT().
				Create(&model.HireManager{
					AccountID: tc.userID,
					CompanyID: tc.compID,
				}).
				Do(func(arg *model.HireManager) {
					arg.ID = 1
				}).
				Return(tc.expectError)

			c, err := managerUcase.Create(tc.userID, tc.compID)

			if tc.expectError == nil {
				assert.Equal(t, int64(1), c.ID)
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

func TestCompanyUsecase_Find(t *testing.T) {
	managerRep, managerUcase := testUcase(t)

	testCases := []struct {
		name        string
		newCompany  *model.Company
		userType    string
		userID      int64
		compID      int64
		expectError error
	}{
		{
			name:        "valid",
			newCompany:  &model.Company{},
			userType:    model.UserCustomer,
			expectError: nil,
			userID:      1,
			compID:      1,
		},
	}

	hID := int64(1)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			managerRep.
				EXPECT().
				Find(hID).
				Return(&model.HireManager{
					ID:        hID,
					AccountID: tc.userID,
					CompanyID: tc.compID,
				},
					tc.expectError)

			c, err := managerUcase.Find(hID)

			expectModel := &model.HireManager{
				ID:        1,
				AccountID: 1,
				CompanyID: 1,
			}

			if tc.expectError == nil {
				assert.Equal(t, expectModel, c)
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

func TestCompanyUsecase_FindByUser(t *testing.T) {
	managerRep, managerUcase := testUcase(t)

	testCases := []struct {
		name        string
		newCompany  *model.Company
		userType    string
		userID      int64
		compID      int64
		hID         int64
		expectError error
	}{
		{
			name:        "valid",
			newCompany:  &model.Company{},
			userType:    model.UserCustomer,
			expectError: nil,
			userID:      1,
			compID:      1,
			hID:         1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			managerRep.
				EXPECT().
				FindByUser(tc.userID).
				Return(&model.HireManager{
					ID:        tc.hID,
					AccountID: tc.userID,
					CompanyID: tc.compID,
				},
					tc.expectError)

			c, err := managerUcase.FindByUser(tc.userID)

			expectModel := &model.HireManager{
				ID:        1,
				AccountID: 1,
				CompanyID: 1,
			}

			if tc.expectError == nil {
				assert.Equal(t, expectModel, c)
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
