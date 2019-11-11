package test

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/mocks"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	userUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/usecase"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initUserCase(t * testing.T) user.Usecase {
	t.Helper()
	ctrl := gomock.NewController(t)
	userRep := mocks.NewMockUserRepository(ctrl)
	freelancerRep := mocks.NewMockFreelancerRepository(ctrl)
	managerRep := mocks.NewMockManagerRepository(ctrl)
	return userUcase.NewUserUsecase(userRep, managerRep, freelancerRep)
}

func TestUcase_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRep := mocks.NewMockUserRepository(ctrl)
	freelancerRep := mocks.NewMockFreelancerRepository(ctrl)
	managerRep := mocks.NewMockManagerRepository(ctrl)
	userU := userUcase.NewUserUsecase(userRep, managerRep, freelancerRep)

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
			name:			"invalid payload",
			user:			nil,
			expectError:	nil,
			expectRun:		false,
		},
		{
			name:			"invalid params",
			user:			&model.User{
				Email:           "1",
				Password:        "1",
			},
			expectError:	errors.New("CreateUser: : email: must be a valid email address; password: the length must be between 6 and 100."),
			expectRun:		false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userRep.
				EXPECT().
				Create(tc.user).
				Return(tc.expectError)

			freelancerRep.
				EXPECT().
				Create().
				Return(tc.expectError)

			managerRep.
				EXPECT().
				Create()

			b := &bytes.Buffer{}
			if err := json.NewEncoder(b).Encode(tc.payload); err != nil {
				t.Fatal()
			}

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/signup", b)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
