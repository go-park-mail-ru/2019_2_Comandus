package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (s * server) ConfigureTestServer() {
	s.mux.HandleFunc("/signup", s.HandleCreateUser).Methods(http.MethodPost)
	s.mux.HandleFunc("/login", s.HandleSessionCreate).Methods(http.MethodPost)
	s.mux.HandleFunc("/setusertype", s.HandleSetUserType).Methods(http.MethodPost)
	s.mux.HandleFunc("/account", s.HandleShowProfile).Methods(http.MethodGet)
	s.mux.HandleFunc("/account", s.HandleEditProfile).Methods(http.MethodPut)
	s.mux.HandleFunc("/account/upload-avatar", s.HandleUploadAvatar).Methods(http.MethodPost)
	s.mux.HandleFunc("/account/download-avatar", s.HandleDownloadAvatar).Methods(http.MethodGet)
	s.mux.HandleFunc("/account/settings/password", s.HandleEditPassword).Methods(http.MethodPut)
	s.mux.HandleFunc("/account/settings/notifications", s.HandleEditNotifications).Methods(http.MethodPut)
	s.mux.HandleFunc("/account/settings/auth-history", s.HandleGetAuthHistory).Methods(http.MethodGet)
	s.mux.HandleFunc("/logout", s.HandleLogout).Methods(http.MethodPost)
}


func TestServer_HandleCreateUser(t *testing.T) {
	config := NewConfig()
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	s := newServer(sessionStore)
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"email":    "user@example.org",
				"password": "secret",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]interface{}{
				"email":    "invalid",
				"password": "short",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/signup", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleSessionCreate(t *testing.T) {
	config := NewConfig()
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	s := newServer(sessionStore)

	u := model.User{
		ID:              0,
		FirstName:            "name",
		Email:           "user@example.org",
		Password:        "secret",
		EncryptPassword: "",
	}

	err:= u.BeforeCreate()
	if err != nil {
		t.Fail()
	}

	s.usersdb.Users = append(s.usersdb.Users, u)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"email":    "user@example.org",
				"password": "secret",
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/login", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestAuthenticateUser(t *testing.T) {
	testCases := []struct {
		name         string
		cookieValue  map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "authenticated",
			cookieValue: map[interface{}]interface{}{
				"user_id": 0,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "not authenticated",
			cookieValue:  nil,
			expectedCode: http.StatusUnauthorized,
		},
	}

	secretKey := []byte("secret")
	s := newServer(sessions.NewCookieStore(secretKey))

	u := model.User{
		ID:              0,
		FirstName:            "name",
		Email:           "user@example.org",
		Password:        "secret",
		EncryptPassword: "",
	}

	err:= u.BeforeCreate()
	if err != nil {
		t.Fail()
	}

	s.usersdb.Users = append(s.usersdb.Users, u)

	sc := securecookie.New(secretKey, nil)
	mw := s.authenticateUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			cookieStr, _ := sc.Encode(sessionName, tc.cookieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			mw.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleSetUserType(t *testing.T) {
	/*config := NewConfig()
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	s := &server{
		mux: mux.NewRouter(),
		usersdb: model.NewUsersDB(),
		sessionStore:sessionStore,
	}
	s.ConfigureTestServer()

	u := model.User{
		ID:              0,
		FirstName:            "name",
		Email:           "user@example.org",
		Password:        "secret",
		EncryptPassword: "",
	}
	u.BeforeCreate()
	s.usersdb.Users = append(s.usersdb.Users, u)

	testCases := []struct {
		name         string
		payload  map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "freelancer",
			payload: map[interface{}]interface{}{
				"type": "freelancer",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "hiremanager",
			payload: map[interface{}]interface{}{
				"type": "customer",
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "wrong input",
			payload:  nil,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		rec := httptest.NewRecorder()
		b := &bytes.Buffer{}
		json.NewEncoder(b).Encode(tc.payload)
		req, _ := http.NewRequest(http.MethodPost, "/setusertype", b)
		s.ServeHTTP(rec, req)
		assert.Equal(t, tc.expectedCode, rec.Code)
	}*/
}
