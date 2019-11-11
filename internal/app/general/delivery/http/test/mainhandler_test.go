package test

import (
	"bytes"
	"encoding/json"
	apiserver "github.com/go-park-mail-ru/2019_2_Comandus/internal/app"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleCreateUser(t *testing.T) {
	config := apiserver.NewConfig()
	zapLogger, _ := zap.NewProduction()
	sugaredLogger := zapLogger.Sugar()

	token, err := apiserver.NewHMACHashToken(config.TokenSecret)
	if err != nil {
	}

	defer func() {
		if err := zapLogger.Sync(); err != nil {
		}
	}()

	sanitizer := bluemonday.UGCPolicy()
	sessionStore := sessions.NewCookieStore([]byte("config.SessionKey"))
	s := apiserver.NewServer(sessionStore, sugaredLogger, token, sanitizer, nil)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
		expectCreate bool
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"email":    "user@example.org",
				"password": "secret",
			},
			expectedCode: http.StatusCreated,
			expectCreate: true,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
			expectCreate: false,
		},
		{
			name: "invalid params",
			payload: map[string]interface{}{
				"email":    "invalid",
				"password": "short",
			},
			expectedCode: http.StatusBadRequest,
			expectCreate: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectCreate {
				store.userRepository.
					EXPECT().
					Create(gomock.Any()).
					Return(nil)

				store.freelancerRepository.
					EXPECT().
					Create(gomock.Any()).
					Return(nil)

				store.managerRepository.
					EXPECT().
					Create(gomock.Any()).
					Return(nil)
			}

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
	config := apiserver.NewConfig()
	zapLogger, _ := zap.NewProduction()
	sugaredLogger := zapLogger.Sugar()

	token, err := apiserver.NewHMACHashToken(config.TokenSecret)
	if err != nil {
	}
	store := New(t)

	defer func() {
		if err := zapLogger.Sync(); err != nil {
		}
	}()


	sessionStore := sessions.NewCookieStore([]byte("config.SessionKey"))
	sanitizer := bluemonday.UGCPolicy()
	s := apiserver.NewServer(sessionStore, store, sugaredLogger, token, sanitizer)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
		expectFind bool
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"email":    "user@example.org",
				"password": "password",
			},
			expectedCode: http.StatusOK,
			expectFind: true,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
			expectFind: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := &model.User{
				ID:              1,
				Email:           "user@example.org",
				Password:        "password",
			}

			_ = user.BeforeCreate()

			if tc.expectFind {
				store.userRepository.
					EXPECT().
					FindByEmail("user@example.org").
					Return(&model.User{}, nil)
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
