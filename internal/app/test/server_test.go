package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/securecookie"
	"log"
	//"github.com/golang/mock/gomock"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"

	//"github.com/gorilla/sessions"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app"
	"go.uber.org/zap"
	"testing"
)

func TestServer_HandleCreateUser(t *testing.T) {
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
	s := apiserver.NewServer(sessionStore, store, sugaredLogger, token)

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
	s := apiserver.NewServer(sessionStore, store, sugaredLogger, token)

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

func TestServer_HandleSetUserType(t *testing.T) {
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid customer",
			payload: map[string]interface{}{
				"type": "client",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "valid",
			payload: map[string]interface{}{
				"type": "freelancer",
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid params",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
	}

	config := apiserver.NewConfig()
	zapLogger, _ := zap.NewProduction()
	sugaredLogger := zapLogger.Sugar()
	defer func() {
		if err := zapLogger.Sync(); err != nil {
		}
	}()

	token, err := apiserver.NewHMACHashToken(config.TokenSecret)
	if err != nil {
		log.Println("TOKEN ERROR")
	}

	store := New(t)

	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))

	s := apiserver.NewServer(sessionStore, store, sugaredLogger, token)
	sc := securecookie.New([]byte(config.SessionKey), nil)

	for _, tc := range testCases {
		log.Println(1)
		t.Run(tc.name, func(t *testing.T) {



			store.userRepository.EXPECT().
				Find(int64(1)).
				Return(&model.User{}, nil)

			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/setusertype", b)
			cookieStr, _ := sc.Encode("user-session", map[interface{}]interface{}{
				"user_id": 1,
			})
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", "user-session", cookieStr))
			s.ServeHTTP(rec, req)
			log.Println(rec.Body)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_AuthenticateUser(t *testing.T) {
	testCases := []struct {
		name         string
		cookieValue  map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "authenticated",
			cookieValue: map[interface{}]interface{}{
				"user_id": 1,
			},
			expectedCode: http.StatusOK,
		},
		/*{
			name:         "not authenticated",
			cookieValue:  nil,
			expectedCode: http.StatusUnauthorized,
		},*/
	}

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

	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	s := apiserver.NewServer(sessionStore, store, sugaredLogger, token)

	sc := securecookie.New([]byte(config.SessionKey), nil)

	mw := s.AuthenticateUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedCode == http.StatusOK {
				store.userRepository.
					EXPECT().
					Find(int64(1)).
					Return(&model.User{}, nil)
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			cookieStr, _ := sc.Encode("user-session", tc.cookieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", "user-session", cookieStr))
			mw.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleCreateJob(t *testing.T) {
	testCases := []struct {
		name         string
		payload      interface{}
		cookie       interface{}
		expectedCode int
		userType     string
	}{
		{
			name: "correct user",
			payload: map[string]interface{}{
				"title":                    "golang server writing",
				"description":              "write server for fl.ru",
				"files":                    "",
				"specialityId,string":      "1",
				"experienceLevelId,string": "1",
				"paymentAmount,string":     "23000.34",
				"country":                  "Russia",
				"city":                     "Moscow",
				"jobTypeId,string":         "1",
			},
			cookie: map[interface{}]interface{}{
				"user_id":   1,
			},
			expectedCode: http.StatusOK,
			userType:     model.UserFreelancer,
		},
		{
			name: "user without user type",
			payload: map[string]interface{}{
				"title":                    "golang server writing",
				"description":              "write server for fl.ru",
				"files":                    "",
				"specialityId,string":      "1",
				"experienceLevelId,string": "1",
				"paymentAmount,string":     "23000.34",
				"country":                  "Russia",
				"city":                     "Moscow",
				"jobTypeId,string":         "1",
			},
			cookie: map[interface{}]interface{}{
				"user_id": 1,
			},
			expectedCode: http.StatusInternalServerError,
			userType:     "wrong type",
		},
		{
			name: "not auth user",
			payload: map[string]interface{}{
				"title":                    "golang server writing",
				"description":              "write server for fl.ru",
				"files":                    "",
				"specialityId,string":      "1",
				"experienceLevelId,string": "1",
				"paymentAmount,string":     "23000.34",
				"country":                  "Russia",
				"city":                     "Moscow",
				"jobTypeId,string":         "1",
			},
			cookie:       "nil",
			expectedCode: http.StatusUnauthorized,
			userType:     "",
		},
	}

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

	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	s := apiserver.NewServer(sessionStore, store, sugaredLogger, token)

	sc := securecookie.New([]byte(config.SessionKey), nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/jobs", b)
			cookieStr, _ := sc.Encode("user-session", tc.cookie)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", "user-session", cookieStr))
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
