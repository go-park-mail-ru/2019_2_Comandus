package test

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/mocks"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	userUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/usecase"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testUcase(t *testing.T) (*mocks.MockUserRepository, user.Usecase){
	t.Helper()
	ctrl := gomock.NewController(t)
	userRep := mocks.NewMockUserRepository(ctrl)
	freelancerRep := mocks.NewMockFreelancerRepository(ctrl)
	managerRep := mocks.NewMockManagerRepository(ctrl)
	companyRep := mocks.NewMockCompanyRepository(ctrl)
	userU := userUcase.NewUserUsecase(userRep, managerRep, freelancerRep, companyRep)
	return userRep, userU
}

func TestUcase_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRep := mocks.NewMockUserRepository(ctrl)
	freelancerRep := mocks.NewMockFreelancerRepository(ctrl)
	managerRep := mocks.NewMockManagerRepository(ctrl)
	companyRep := mocks.NewMockCompanyRepository(ctrl)
	userU := userUcase.NewUserUsecase(userRep, managerRep, freelancerRep, companyRep)

	testCases := []struct {
		name			string
		user			*model.User
		expectRun		bool
		expectError		error
	}{
		{
			name:			"valid",
			user:			&model.User{
				Email:           "user@example.org",
				Password:        "secret",
			},
			expectError:	nil,
			expectRun:		true,
		},
		{
			name:			"invalid params",
			user:			&model.User{
				Email:           "1",
				Password:        "1",
			},
			expectError:	errors.Wrap(errors.New("email: must be a valid email address; " +
				"password: the length must be between 6 and 100."), "CreateUser"),
			expectRun:		false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userRep.
				EXPECT().
				Create(tc.user).
				Return(tc.expectError)

			if tc.expectRun {
				companyRep.
					EXPECT().
					Create(gomock.Any()).
					Return(nil)

				freelancerRep.
					EXPECT().
					Create(gomock.Any()).
					Return(nil)

				managerRep.
					EXPECT().
					Create(gomock.Any()).
					Return(nil)
			}

			err := userU.CreateUser(tc.user)

			if tc.expectRun {
				assert.Equal(t, nil, err)
				return
			}

			if !tc.expectRun && err != nil {
				assert.Equal(t, tc.expectError.Error(), err.Error())
				return
			}

			t.Fatal()
		})
	}
}

func TestUcase_EditUser(t *testing.T) {
	userRep, userU := testUcase(t)
	user := &model.User{
		ID:			1,
		Email:		"user@example.org",
		Password:	"secret",
		UserType:	model.UserFreelancer,
	}
	if err := user.BeforeCreate(); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name			string
		newUser			*model.User
		expectError		error
	}{
		{
			name:			"valid: edit user info",
			newUser:		&model.User{
				ID:			user.ID,
				FirstName:	"ivan",
				SecondName:	"ivanov",
				UserName:	"ivan1970",
				Email:		"user@example.org",
				Password:	"",
				EncryptPassword: user.EncryptPassword,
				UserType:	model.UserFreelancer,
				RegistrationDate:	user.RegistrationDate,
			},
			expectError:	nil,
		},
		{
			name:			"invalid: edit user type",
			newUser:		&model.User{
				ID:			user.ID,
				Email:		"user@example.org",
				Password:	"secret",
				UserType:	model.UserCustomer,
				RegistrationDate:	user.RegistrationDate,
			},
			expectError:	errors.New("can't change user type by edit"),
		},
		{
			name:			"invalid params: edit password",
			newUser:		&model.User{
				ID:			user.ID,
				Email:           "user@example.org",
				Password:        "1",
				UserType:	model.UserFreelancer,
				RegistrationDate:	user.RegistrationDate,
			},
			expectError:	errors.Wrap(errors.New("can't change password without validation"), "ComparePassword"),
		},
		{
			name:			"invalid params: edit email",
			newUser:		&model.User{
				ID:			user.ID,
				Email:           "user1@example.org",
				Password:        "secret",
				UserType:	model.UserFreelancer,
				RegistrationDate:	user.RegistrationDate,
			},
			expectError:	errors.Wrap(errors.New("can't change email"), "EditUser"),
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

func TestUcase_EditPassword(t *testing.T) {
	userRep, userU := testUcase(t)

	user := &model.User{
		ID:			1,
		Email:		"user@example.org",
		Password:	"secret",
		UserType:	model.UserFreelancer,
	}
	if err := user.BeforeCreate(); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name			string
		bodyPassword	*model.BodyPassword
		expectError		error
	}{
		{
			name:			"valid",
			bodyPassword:	&model.BodyPassword{
				Password:                "secret",
				NewPassword:             "new.secret",
				NewPasswordConfirmation: "new.secret",
			},
			expectError:	nil,
		},
		{
			name:			"invalid: wrong old password",
			bodyPassword:	&model.BodyPassword{
				Password:                "wrong",
				NewPassword:             "new.secret",
				NewPasswordConfirmation: "new.secret",
			},
			expectError:	errors.Wrap(errors.New("wrong old password"),"model.user.ComparePassword"),
		},
		{
			name:			"invalid: new password and confirmation dont match",
			bodyPassword:	&model.BodyPassword{
				Password:                "secret",
				NewPassword:             "new.secret1",
				NewPasswordConfirmation: "new.secret2",
			},
			expectError:	errors.New("new passwords are different"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectError == nil {
				userRep.
					EXPECT().
					Edit(user).
					Return(tc.expectError)
			}

			err := userU.EditUserPassword(tc.bodyPassword, user)

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

func TestUcase_VerifyUser(t *testing.T) {
	userRep, userU := testUcase(t)

	user := &model.User{
		ID:			1,
		Email:		"user@example.org",
		Password:	"secret",
		UserType:	model.UserFreelancer,
	}
	if err := user.BeforeCreate(); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name        string
		user        *model.User
		expectError error
	}{
		{
			name:			"valid",
			user:			&model.User{
				Email:		"user@example.org",
				Password:	"secret",
			},
			expectError:	nil,
		},
		{
			name:			"invalid password",
			user:			&model.User{
				Email:		"user@example.org",
				Password:	"secret1",
			},
			expectError:	errors.New("wrong password"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			userRep.
				EXPECT().
				FindByEmail(tc.user.Email).
				Return(user, err)

			id, err := userU.VerifyUser(tc.user)

			if tc.expectError == nil {
				assert.Equal(t, user.ID, id)
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

func TestUcase_Find(t *testing.T) {
	userRep, userU := testUcase(t)

	user := &model.User{
		ID:			1,
		Email:		"user@example.org",
		Password:	"secret",
		UserType:	model.UserFreelancer,
	}
	if err := user.BeforeCreate(); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name        string
		user        *model.User
		expectError error
	}{
		{
			name:			"valid",
			user:			&model.User{
				ID:			1,
				Email:		"user@example.org",
				Password:	"secret",
			},
			expectError:	nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			userRep.
				EXPECT().
				Find(tc.user.ID).
				Return(user, err)

			_, err = userU.Find(tc.user.ID)

			assert.Equal(t, nil, err)
		})
	}
}

func TestUcase_SetUserType(t *testing.T) {
	userRep, userU := testUcase(t)

	user := &model.User{
		ID:			1,
		Email:		"user@example.org",
		Password:	"secret",
		UserType:	model.UserFreelancer,
	}
	if err := user.BeforeCreate(); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name        string
		userType	string
		user        *model.User
		expectError error
	}{
		{
			name:			"valid: freelancer",
			userType:		model.UserFreelancer,
			expectError:	nil,
		},
		{
			name:			"valid: manager",
			userType:		model.UserCustomer,
			expectError:	nil,
		},
		{
			name:			"invalid",
			userType:		"wrong type",
			expectError:	errors.Wrap(errors.New("wrong user type"), "SetUserType()"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			userRep.
				EXPECT().
				Edit(user).
				Return(err)

			err = userU.SetUserType(user, tc.userType)

			if err == nil {
				assert.Equal(t, tc.expectError, err)
				return
			}

			assert.Equal(t, tc.expectError.Error(), err.Error())
		})
	}
}

func TestUcase_GetRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRep := mocks.NewMockUserRepository(ctrl)
	freelancerRep := mocks.NewMockFreelancerRepository(ctrl)
	managerRep := mocks.NewMockManagerRepository(ctrl)
	companyRep := mocks.NewMockCompanyRepository(ctrl)
	userU := userUcase.NewUserUsecase(userRep, managerRep, freelancerRep, companyRep)

	user := &model.User{
		ID:			1,
		Email:		"user@example.org",
		Password:	"secret",
		UserType:	model.UserFreelancer,
	}

	manager := &model.HireManager{
		ID:               1,
		AccountID:        1,
		Location:         "moscow",
		CompanyID:        1,
	}

	company := &model.Company{
		ID:          1,
		CompanyName: "test company",
	}

	if err := user.BeforeCreate(); err != nil {
		t.Fatal(err)
	}

	managerRep.
		EXPECT().
		FindByUser(user.ID).
		Return(manager, nil)

	companyRep.
		EXPECT().
		Find(manager.ID).
		Return(company, nil)

	_, err := userU.GetRoles(user)

	assert.Equal(t, nil, err)
}

func TestUcase_GetAvatar(t *testing.T) {
	_, userU := testUcase(t)

	user := &model.User{
		ID:			1,
		Email:		"user@example.org",
		Password:	"secret",
		UserType:	model.UserFreelancer,
	}

	if err := user.BeforeCreate(); err != nil {
		t.Fatal(err)
	}

	_, err := userU.GetAvatar(user)

	assert.Equal(t, nil, err)
}