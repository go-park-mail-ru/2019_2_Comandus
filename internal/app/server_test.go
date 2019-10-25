package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleCreateUser(t *testing.T) {
	config := NewConfig()
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))

	db, truncate := testStore(t, databaseURL)
	defer truncate("users", "managers", "freelancers")

	store := sqlstore.New(db)
	s := newServer(sessionStore, store)

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
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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
	config := NewConfig()
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))

	db, truncate := testStore(t, databaseURL)
	defer truncate("users", "managers", "freelancers")

	store := sqlstore.New(db)
	s := newServer(sessionStore, store)

	err := s.addUser2Server(t)
	if err != nil {
		t.Fatal()
	}

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"email":    "user@example.org",
				"password": "password",
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

func TestServer_AuthenticateUser(t *testing.T) {
	testCases := []struct {
		name         string
		cookieValue  map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "authenticated",
			cookieValue: map[interface{}]interface{}{
				"user_id": 3,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "not authenticated",
			cookieValue:  nil,
			expectedCode: http.StatusUnauthorized,
		},
	}

	config := NewConfig()
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))

	db, truncate := testStore(t, databaseURL)
	defer truncate("users", "managers", "freelancers")

	store := sqlstore.New(db)
	s := newServer(sessionStore, store)

	err := s.addUser2Server(t)
	if err != nil {
		t.Fail()
	}

	sc := securecookie.New([]byte(config.SessionKey), nil)
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

	config := NewConfig()
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))

	db, truncate := testStore(t, databaseURL)
	defer truncate("users", "managers", "freelancers")

	store := sqlstore.New(db)
	s := newServer(sessionStore, store)

	err := s.addUser2Server(t)
	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	sc := securecookie.New([]byte(config.SessionKey), nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/private/setusertype", b)
			cookieStr, _ := sc.Encode(sessionName, map[interface{}]interface{}{
				"user_id": 4,
			})
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			s.ServeHTTP(rec, req)
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
				"user_id":   5,
				"user_type": userCustomer,
			},
			expectedCode: http.StatusOK,
			userType:     userCustomer,
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
				"user_id": 5,
			},
			expectedCode: http.StatusInternalServerError,
			userType:     userFreelancer,
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

	config := NewConfig()
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))

	db, truncate := testStore(t, databaseURL)
	defer truncate("users", "freelancers", "jobs")

	store := sqlstore.New(db)
	s := newServer(sessionStore, store)

	err := s.addUser2Server(t)
	if err != nil {
		t.Fatal()
	}

	sc := securecookie.New([]byte(config.SessionKey), nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s.userType = tc.userType

			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/private/jobs", b)
			cookieStr, _ := sc.Encode(sessionName, tc.cookie)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

/*func TestServer_HandleLogout(t *testing.T) {
	testCases := []struct {
		name         string
		cookie      interface{}
		expectedCode int
	}{
		{
			name: "auth user",
			cookie : map[interface{}]interface{}{
				"user_id":   0,
				"user_type": userCustomer,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "unauth user",
			cookie : "invalid",
			expectedCode: http.StatusUnauthorized,
		},
	}

	secretKey := []byte("secret")
	s := newServer(sessions.NewCookieStore(secretKey))
	sc := securecookie.New(secretKey, nil)

	err := s.addUser2Server()
	if err != nil {
		t.Fail()
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/private/logout", b)
			cookieStr, _ := sc.Encode(sessionName, tc.cookie)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleGetJob(t *testing.T) {
	testCases := []struct {
		name         string
		cookie      interface{}
		request 	string
		expectedCode int
	}{
		{
			name: "auth user",
			cookie : map[interface{}]interface{}{
				"user_id":   0,
				"user_type": userCustomer,
			},
			request : "/private/jobs/0",
			expectedCode: http.StatusOK,
		},
		{
			name: "unauth user",
			cookie : "invalid",
			request : "/private/jobs/0",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "nonexistent job",
			cookie : map[interface{}]interface{}{
				"user_id":   0,
				"user_type": userCustomer,
			},
			request : "/private/jobs/1",
			expectedCode: http.StatusNotFound,
		},
	}

	secretKey := []byte("secret")
	s := newServer(sessions.NewCookieStore(secretKey))
	sc := securecookie.New(secretKey, nil)

	err := s.addUser2Server()
	if err != nil {
		t.Fail()
	}
	s.addJob2Server()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tc.request, nil)
			cookieStr, _ := sc.Encode(sessionName, tc.cookie)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleShowProfile(t *testing.T) {
	testCases := []struct {
		name         string
		cookie      interface{}
		expectedCode int
	}{
		{
			name: "auth user",
			cookie : map[interface{}]interface{}{
				"user_id":   0,
				"user_type": userCustomer,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "unauth user",
			cookie : "invalid",
			expectedCode: http.StatusUnauthorized,
		},
	}

	secretKey := []byte("secret")
	s := newServer(sessions.NewCookieStore(secretKey))
	sc := securecookie.New(secretKey, nil)

	err := s.addUser2Server()
	if err != nil {
		t.Fail()
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/private/account", nil)
			cookieStr, _ := sc.Encode(sessionName, tc.cookie)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleEditPassword(t *testing.T) {
	testCases := []struct {
		name         string
		cookie      interface{}
		passwords 	interface{}
		expectedCode int
	}{
		{
			name: "correct request",
			cookie : map[interface{}]interface{}{
				"user_id":   0,
				"user_type": userCustomer,
			},
			passwords : map[interface{}]interface{}{
				"password":   "secret",
				"newPassword": "1234",
				"newPasswordConfirmation": "1234",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "different new passwords",
			cookie : map[interface{}]interface{}{
				"user_id":   0,
				"user_type": userCustomer,
			},
			passwords : map[interface{}]interface{}{
				"password":   "secret",
				"newPassword": "1234",
				"newPasswordConfirmation": "12345",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "wrong old password",
			cookie : map[interface{}]interface{}{
				"user_id":   0,
				"user_type": userCustomer,
			},
			passwords : map[interface{}]interface{}{
				"password":   "secret1",
				"newPassword": "1234",
				"newPasswordConfirmation": "1234",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "some parameters not set",
			cookie : map[interface{}]interface{}{
				"user_id":   0,
				"user_type": userCustomer,
			},
			passwords : map[interface{}]interface{}{
				"password":   "secret",
				"newPassword": "1234",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "not auth user",
			cookie : "invalid",
			passwords : map[interface{}]interface{}{
				"password":   "secret1",
				"newPassword": "1234",
				"newPasswordConfirmation": "1234",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	secretKey := []byte("secret")
	s := newServer(sessions.NewCookieStore(secretKey))
	sc := securecookie.New(secretKey, nil)

	err := s.addUser2Server()
	if err != nil {
		t.Fail()
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPut, "/private/account/settings/password", b)
			json.NewEncoder(b).Encode(tc.passwords)
			cookieStr, _ := sc.Encode(sessionName, tc.cookie)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}*/
