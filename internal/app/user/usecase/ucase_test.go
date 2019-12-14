package userUcase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/delivery/grpc/company_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/delivery/grpc/freelancer_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/delivery/grpc/manager_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/mocks/client_mocks"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/mocks/repository_mocks"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	//"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testUcase(t *testing.T) (*repository_mocks.MockUserRepository, user.Usecase, *client_mocks.MockClientFreelancer,
	*client_mocks.MockManagerClient, *client_mocks.MockCompanyClient) {
	t.Helper()
	ctrl := gomock.NewController(t)
	userRep := repository_mocks.NewMockUserRepository(ctrl)
	freelancerClient := client_mocks.NewMockClientFreelancer(ctrl)
	managerClient := client_mocks.NewMockManagerClient(ctrl)
	companyClient := client_mocks.NewMockCompanyClient(ctrl)
	userU := NewUserUsecase(userRep, freelancerClient, managerClient, companyClient)
	return userRep, userU, freelancerClient, managerClient, companyClient
}

func TestUcase_CreateUser(t *testing.T) {
	userRep, userU, _, _, _ := testUcase(t)
	//userRep , userU, freelancerClient, managerClient, companyClient := testUcase(t)
	testCases := []struct {
		name        string
		user        *model.User
		expectRun   bool
		expectError error
	}{
		{
			name: "valid",
			user: &model.User{
				Email:    "user@example.org",
				Password: "secret",
			},
			expectError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userRep.
				EXPECT().
				Create(tc.user).
				Return(tc.expectError)

			err := userU.CreateUser(tc.user)

			if err == nil {
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

func TestUcase_EditUser(t *testing.T) {
	userRep, userU, _, _, _ := testUcase(t)
	user := &model.User{
		ID:       1,
		Email:    "user@example.org",
		Password: "secret",
		UserType: model.UserFreelancer,
	}
	if err := user.BeforeCreate(); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name        string
		newUser     *model.User
		expectError error
	}{
		{
			name: "valid: edit user info",
			newUser: &model.User{
				ID:               user.ID,
				FirstName:        "ivan",
				SecondName:       "ivanov",
				UserName:         "ivan1970",
				Email:            "user@example.org",
				Password:         "",
				EncryptPassword:  user.EncryptPassword,
				UserType:         model.UserFreelancer,
				RegistrationDate: user.RegistrationDate,
			},
			expectError: nil,
		},
		//{
		//	name:			"invalid: edit user type",
		//	newUser:		&model.User{
		//		ID:			user.ID,
		//		Email:		"user@example.org",
		//		Password:	"secret",
		//		UserType:	model.UserCustomer,
		//		RegistrationDate:	user.RegistrationDate,
		//	},
		//	expectError:	errors.New("userRep.Edit(): can't change user type by edit"),
		//},
		//{
		//	name:			"invalid params: edit password",
		//	newUser:		&model.User{
		//		ID:			user.ID,
		//		Email:           "user@example.org",
		//		Password:        "1",
		//		UserType:	model.UserFreelancer,
		//		RegistrationDate:	user.RegistrationDate,
		//	},
		//	expectError:	errors.New("userRep.Edit(): ComparePassword: can't change password without validation"),
		//},
		{
			name: "invalid params: edit email",
			newUser: &model.User{
				ID:               user.ID,
				Email:            "user1@example.org",
				Password:         "secret",
				UserType:         model.UserFreelancer,
				RegistrationDate: user.RegistrationDate,
			},
			expectError: errors.Wrap(errors.New("can't change email"), "EditUser"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userRep.
				EXPECT().
				Edit(tc.newUser).
				Return(tc.expectError)

			err := userU.EditUser(tc.newUser, user)

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

//func TestUcase_EditPassword(t *testing.T) {
//	userRep, userU := testUcase(t)
//
//	user := &model.User{
//		ID:			1,
//		Email:		"user@example.org",
//		Password:	"secret",
//		UserType:	model.UserFreelancer,
//	}
//	if err := user.BeforeCreate(); err != nil {
//		t.Fatal(err)
//	}
//
//	testCases := []struct {
//		name			string
//		bodyPassword	*model.BodyPassword
//		expectError		error
//	}{
//		{
//			name:			"valid",
//			bodyPassword:	&model.BodyPassword{
//				Password:                "secret",
//				NewPassword:             "new.secret",
//				NewPasswordConfirmation: "new.secret",
//			},
//			expectError:	nil,
//		},
//		{
//			name:			"invalid: wrong old password",
//			bodyPassword:	&model.BodyPassword{
//				Password:                "wrong",
//				NewPassword:             "new.secret",
//				NewPasswordConfirmation: "new.secret",
//			},
//			expectError:	errors.Wrap(errors.New("wrong old password"),"model.user.ComparePassword"),
//		},
//		{
//			name:			"invalid: new password and confirmation dont match",
//			bodyPassword:	&model.BodyPassword{
//				Password:                "secret",
//				NewPassword:             "new.secret1",
//				NewPasswordConfirmation: "new.secret2",
//			},
//			expectError:	errors.New("new passwords are different"),
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			if tc.expectError == nil {
//				userRep.
//					EXPECT().
//					Edit(user).
//					Return(tc.expectError)
//			}
//
//			err := userU.EditUserPassword(tc.bodyPassword, user)
//
//			if tc.expectError == nil {
//				assert.Equal(t, nil, err)
//				return
//			}
//
//			if err != nil {
//				assert.Equal(t, tc.expectError.Error(), err.Error())
//				return
//			}
//
//			t.Fatal()
//		})
//	}
//}
//
//func TestUcase_VerifyUser(t *testing.T) {
//	userRep, userU := testUcase(t)
//
//	user := &model.User{
//		ID:			1,
//		Email:		"user@example.org",
//		Password:	"secret",
//		UserType:	model.UserFreelancer,
//	}
//	if err := user.BeforeCreate(); err != nil {
//		t.Fatal(err)
//	}
//
//	testCases := []struct {
//		name        string
//		user        *model.User
//		expectError error
//	}{
//		{
//			name:			"valid",
//			user:			&model.User{
//				Email:		"user@example.org",
//				Password:	"secret",
//			},
//			expectError:	nil,
//		},
//		{
//			name:			"invalid password",
//			user:			&model.User{
//				Email:		"user@example.org",
//				Password:	"secret1",
//			},
//			expectError:	errors.New("wrong password"),
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var err error
//			userRep.
//				EXPECT().
//				FindByEmail(tc.user.Email).
//				Return(user, err)
//
//			id, err := userU.VerifyUser(tc.user)
//
//			if tc.expectError == nil {
//				assert.Equal(t, user.ID, id)
//				assert.Equal(t, nil, err)
//				return
//			}
//
//			if err != nil {
//				assert.Equal(t, tc.expectError.Error(), err.Error())
//				return
//			}
//
//			t.Fatal()
//		})
//	}
//}
//
//func TestUcase_Find(t *testing.T) {
//	userRep, userU := testUcase(t)
//
//	user := &model.User{
//		ID:			1,
//		Email:		"user@example.org",
//		Password:	"secret",
//		UserType:	model.UserFreelancer,
//	}
//	if err := user.BeforeCreate(); err != nil {
//		t.Fatal(err)
//	}
//
//	testCases := []struct {
//		name        string
//		user        *model.User
//		expectError error
//	}{
//		{
//			name:			"valid",
//			user:			&model.User{
//				ID:			1,
//				Email:		"user@example.org",
//				Password:	"secret",
//			},
//			expectError:	nil,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var err error
//			userRep.
//				EXPECT().
//				Find(tc.user.ID).
//				Return(user, err)
//
//			_, err = userU.Find(tc.user.ID)
//
//			assert.Equal(t, nil, err)
//		})
//	}
//}
//
//func TestUcase_SetUserType(t *testing.T) {
//	userRep, userU := testUcase(t)
//
//	user := &model.User{
//		ID:			1,
//		Email:		"user@example.org",
//		Password:	"secret",
//		UserType:	model.UserFreelancer,
//	}
//	if err := user.BeforeCreate(); err != nil {
//		t.Fatal(err)
//	}
//
//	testCases := []struct {
//		name        string
//		userType	string
//		user        *model.User
//		expectError error
//	}{
//		{
//			name:			"valid: freelancer",
//			userType:		model.UserFreelancer,
//			expectError:	nil,
//		},
//		{
//			name:			"valid: manager",
//			userType:		model.UserCustomer,
//			expectError:	nil,
//		},
//		{
//			name:			"invalid",
//			userType:		"wrong type",
//			expectError:	errors.Wrap(errors.New("wrong user type"), "SetUserType()"),
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var err error
//			userRep.
//				EXPECT().
//				Edit(user).
//				Return(err)
//
//			err = userU.SetUserType(user, tc.userType)
//
//			if err == nil {
//				assert.Equal(t, tc.expectError, err)
//				return
//			}
//
//			assert.Equal(t, tc.expectError.Error(), err.Error())
//		})
//	}
//}
//
//func TestUcase_GetRoles(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	userRep := repository_mocks.NewMockUserRepository(ctrl)
//	freelancerRep := repository_mocks.NewMockFreelancerRepository(ctrl)
//	managerRep := repository_mocks.NewMockManagerRepository(ctrl)
//	companyRep := repository_mocks.NewMockCompanyRepository(ctrl)
//	userU := NewUserUsecase(userRep, managerRep, freelancerRep, companyRep)
//
//	user := &model.User{
//		ID:			1,
//		Email:		"user@example.org",
//		Password:	"secret",
//		UserType:	model.UserFreelancer,
//	}
//
//	manager := &model.HireManager{
//		ID:               1,
//		AccountID:        1,
//		Location:         "moscow",
//		CompanyID:        1,
//	}
//
//	company := &model.Company{
//		ID:          1,
//		CompanyName: "test company",
//	}
//
//	if err := user.BeforeCreate(); err != nil {
//		t.Fatal(err)
//	}
//
//	managerRep.
//		EXPECT().
//		FindByUser(user.ID).
//		Return(manager, nil)
//
//	companyRep.
//		EXPECT().
//		Find(manager.ID).
//		Return(company, nil)
//
//	_, err := userU.GetRoles(user)
//
//	assert.Equal(t, nil, err)
//}
//
//func TestUcase_GetAvatar(t *testing.T) {
//	_, userU := testUcase(t)
//
//	user := &model.User{
//		ID:			1,
//		Email:		"user@example.org",
//		Password:	"secret",
//		UserType:	model.UserFreelancer,
//	}
//
//	if err := user.BeforeCreate(); err != nil {
//		t.Fatal(err)
//	}
//
//	_, err := userU.GetAvatar(user)
//
//	assert.Equal(t, nil, err)
//}

func TestUcase_Find(t *testing.T) {
	userRep, userU, freelancerClient, managerClient, _ := testUcase(t)
	testCases := []struct {
		name        string
		user        *model.User
		expectRun   bool
		expectError error
	}{
		{
			name: "valid",
			user: &model.User{
				ID:            1,
				Email:         "user@example.org",
				FreelancerId:  1,
				HireManagerId: 1,
				CompanyId:     1,
			},
			expectError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userRep.
				EXPECT().
				Find(tc.user.ID).
				Return(tc.user, nil)

			freelancerClient.EXPECT().
				GetFreelancerByUserFromServer(tc.user.ID).
				Return(&freelancer_grpc.Freelancer{
					ID: tc.user.FreelancerId,
				}, nil)

			managerClient.EXPECT().
				GetManagerByUserFromServer(tc.user.ID).
				Return(&manager_grpc.Manager{
					ID:        tc.user.HireManagerId,
					CompanyId: tc.user.CompanyId,
				}, nil)

			u, err := userU.Find(tc.user.ID)

			if tc.expectError == nil {
				assert.Equal(t, tc.user, u)
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

func TestUcase_SetUserType(t *testing.T) {
	userRep, userU, _, _, _ := testUcase(t)
	user := &model.User{
		ID:       1,
		UserType: model.UserCustomer,
	}

	testCases := []struct {
		name        string
		newUser     *model.User
		expectError error
	}{
		{
			name: "valid",
			newUser: &model.User{
				ID:       user.ID,
				UserType: model.UserFreelancer,
			},
			expectError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userRep.
				EXPECT().
				Edit(tc.newUser).
				Do(func(arg *model.User) {
					arg.UserType = user.UserType
				}).
				Return(tc.expectError)

			err := userU.EditUser(tc.newUser, user)

			if tc.expectError == nil {
				assert.Equal(t, tc.newUser, user)
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

func TestUcase_VerifyUser(t *testing.T) {
	userRep, userU, _, _, _ := testUcase(t)
	enc, _ := model.EncryptString("Helloo")
	testCases := []struct {
		name        string
		user        *model.User
		expectRun   bool
		expectError error
	}{
		{
			name: "valid",
			user: &model.User{
				ID:              1,
				Email:           "user@example.org",
				Password:        "Helloo",
				EncryptPassword: enc,
			},
			expectError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userRep.
				EXPECT().
				FindByEmail(tc.user.Email).
				Return(tc.user, nil)

			uID, err := userU.VerifyUser(tc.user)

			if tc.expectError == nil {
				assert.Equal(t, tc.user.ID, uID)
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

func TestUcase_GetRoles(t *testing.T) {
	_, userU, _, managerClient, companyClient := testUcase(t)
	testCases := []struct {
		name        string
		user        *model.User
		expectRun   bool
		expectError error
	}{
		{
			name: "valid",
			user: &model.User{
				ID:         1,
				FirstName:  "Hi",
				SecondName: "World",
			},
			expectError: nil,
		},
	}

	exm := &manager_grpc.Manager{
		ID:        1,
		CompanyId: 1,
	}

	expC := &company_grpc.Company{
		ID:          1,
		CompanyName: "BroBro",
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			managerClient.
				EXPECT().
				GetManagerByUserFromServer(tc.user.ID).
				Return(exm, nil)

			companyClient.
				EXPECT().
				GetCompanyFromServer(exm.CompanyId).
				Return(expC, nil)

			role, err := userU.GetRoles(tc.user)

			if tc.expectError == nil {
				assert.Equal(t, expC.CompanyName, role[0].Label)
				assert.Equal(t, tc.user.FirstName+" "+tc.user.SecondName, role[1].Label)
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
