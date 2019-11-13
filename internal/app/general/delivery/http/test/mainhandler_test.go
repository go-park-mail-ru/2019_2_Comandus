package test

import (
	"bytes"
	"encoding/json"
	apiserver "github.com/go-park-mail-ru/2019_2_Comandus/internal/app"
	mainHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/delivery/http"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/middleware"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServer(t *testing.T) (*apiserver.Server, *MockUserUsecase) {
	t.Helper()

	zapLogger, _ := zap.NewProduction()
	defer func() {
		if err := zapLogger.Sync(); err != nil {
			log.Println(err)
		}
	}()
	sugaredLogger := zapLogger.Sugar()

	s, err := apiserver.NewServer(apiserver.NewConfig(), sugaredLogger)
	if err != nil {
		t.Fatal()
	}

	userU := NewMockUserUsecase(gomock.NewController(t))
	mid := middleware.NewMiddleware(s.SessionStore, s.Logger, s.Token, userU, s.Config.ClientUrl)
	s.Mux.Use(mid.RequestIDMiddleware, mid.CORSMiddleware, mid.AccessLogMiddleware)
	mainHttp.NewMainHandler(s.Mux, userU, s.Sanitizer, s.Logger, s.SessionStore)
	return s, userU
}

func TestServer_HandleMain(t *testing.T) {
	s, _ := testServer(t)

	b := &bytes.Buffer{}
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", b)
	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestServer_HandleCreateUser(t *testing.T) {

	s, userU := testServer(t)
	testCases := []struct {
		name			string
		payload			interface{}
		user			*model.User
		expectRun		bool
		expectedCode	int
		expectError		error
	}{
		{
			name:			"valid",
			payload:		map[string]interface{}{
				"email":    "user@example.org",
				"password": "secret",
			},
			user:			&model.User{
				Email:           "user@example.org",
				Password:        "secret",
			},
			expectedCode:	http.StatusCreated,
			expectError:	nil,
			expectRun:		true,
		},
		{
			name:			"invalid payload",
			payload:		"invalid",
			user:			nil,
			expectedCode:	http.StatusBadRequest,
			expectError:	nil,
			expectRun:		false,
		},
		{
			name:			"invalid params",
			payload:		map[string]interface{}{
				"email":    "1",
				"password": "1",
			},
			user:			&model.User{
				Email:           "1",
				Password:        "1",
			},
			expectedCode:	http.StatusBadRequest,
			expectError:	errors.New("CreateUser: : email: must be a valid email address; password: the length must be between 6 and 100."),
			expectRun:		false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userU.
				EXPECT().
				CreateUser(tc.user).
				Return(tc.expectError)

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

func TestServer_HandleSessionCreate(t *testing.T) {
	s, userU := testServer(t)

	testCases := []struct {
		name			string
		payload			interface{}
		user			*model.User
		expectedCode	int
		expectFind		bool
	}{
		{
			name: 			"valid",
			payload:		map[string]interface{}{
				"email":    "user@example.org",
				"password": "password",
			},
			user:			&model.User{
				Email:           "user@example.org",
				Password:        "password",
			},
			expectedCode:	http.StatusOK,
			expectFind:		true,
		},
		{
			name:			"invalid payload",
			payload:		"invalid",
			expectedCode:	http.StatusBadRequest,
			expectFind:		false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectFind {
				userU.
					EXPECT().
					VerifyUser(tc.user).
					Return(int64(1), nil)
			}

			b := &bytes.Buffer{}
			if err := json.NewEncoder(b).Encode(tc.payload); err != nil {
				t.Fatal()
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/login", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
