package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	userHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/delivery/http"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/middleware"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/securecookie"
	"log"
	"time"

	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"

	//"github.com/gorilla/sessions"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app"
	"go.uber.org/zap"
	"testing"
)

func TestServer(t *testing.T) {
	t.Helper()
}

func TestServer_HandleSetUserType(t *testing.T) {
	testCases := []struct {
		name			string
		payload			interface{}
		expectedCode	int
		userType		string
	}{
		{
			name: "valid customer",
			payload: map[string]interface{}{
				"type": "client",
			},
			expectedCode: http.StatusOK,
			userType: model.UserCustomer,
		},
		{
			name: "valid",
			payload: map[string]interface{}{
				"type": "freelancer",
			},
			expectedCode: http.StatusOK,
			userType:	model.UserFreelancer,
		},
		{
			name:         "invalid params",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
			userType:		"wrong type",
		},
	}

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
	private := s.Mux.PathPrefix("").Subrouter()
	private.Use(mid.AuthenticateUser, mid.CheckTokenMiddleware)
	userHttp.NewUserHandler(private, userU, s.Sanitizer, s.Logger, s.SessionStore)

	sc := securecookie.New([]byte(s.Config.SessionKey), nil)
	user := &model.User{
		ID:              1,
		Email:           "ddjhd@mail.com",
		UserType:        model.UserFreelancer,
	}

	for _, tc := range testCases {
		log.Println(1)
		t.Run(tc.name, func(t *testing.T) {


			userU.
				EXPECT().
				Find(user.ID).
				Return(user, nil)
			userU.
				EXPECT().
				SetUserType(user, tc.userType).
				Return(nil)

			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/setusertype", b)


			cookieStr, _ := sc.Encode("user-session", map[interface{}]interface{}{
				"user_id": int64(1),
			})

			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", "user-session", cookieStr))

			sess, _ := s.SessionStore.Get(req, "user-session")
			token, _ := s.Token.Create(sess, time.Now().Add(24*time.Hour).Unix())
			req.Header.Set("csrf-token", token)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

/*func TestServer_HandleCreateJob(t *testing.T) {
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
				"title":                    "mocks job",
				"paymentAmount":     		"100",
				"country":		            "russia",
				"city":		                "moscow",
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
				"title":                    "mocks job",
				"paymentAmount":     		"100",
				"country":           		"russia",
				"city":              		"moscow",
			},
			cookie: map[interface{}]interface{}{
				"user_id": 1,
			},
			expectedCode: http.StatusOK,
			userType:     "wrong type",
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

	sanitizer := bluemonday.UGCPolicy()
	s := apiserver.NewServer(sessionStore, store, sugaredLogger, token, sanitizer)
	sc := securecookie.New([]byte(config.SessionKey), nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			user := &model.User{
				ID:              1,
				Email:           "ddjhd@mail.com",
				UserType:        "client",
			}

			manager := &model.HireManager{
				ID:               1,
				AccountID:        1,
				Location:         "moscow",
				CompanyID:        0,
			}

			job := &model.Job{
				ID:                0,
				HireManagerId:     0,
				Title:             "mocks job",
				Country:		   "russia",
				City:		   	   "moscow",
				PaymentAmount:     100,
			}

			store.userRepository.EXPECT().
				Find(int64(1)).
				Return(user, nil).
				AnyTimes()

			store.managerRepository.EXPECT().
				FindByUser(int64(1)).
				Return(manager, nil).
				AnyTimes()

			store.jobRepository.EXPECT().
				Create(job, manager).
				Return(nil).
				AnyTimes()

			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/jobs", b)

			cookieStr, _ := sc.Encode("user-session", map[interface{}]interface{}{
				"user_id": 1,
			})

			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", "user-session", cookieStr))

			sess, _ := sessionStore.Get(req, "user-session")
			token, _ := token.Create(sess, time.Now().Add(24*time.Hour).Unix())
			req.Header.Set("csrf-token", token)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleLogout(t *testing.T) {
	testCases := []struct {
		name         string
		cookie      interface{}
		expectedCode int
	}{
		{
			name: "auth user",
			cookie : map[interface{}]interface{}{
				"user_id":   1,
			},
			expectedCode: http.StatusOK,
		},
		/*{
			name: "unauth user",
			cookie : "invalid",
			expectedCode: http.StatusUnauthorized,
		},*/
	/*}

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

	sanitizer := bluemonday.UGCPolicy()
	s := apiserver.NewServer(sessionStore, store, sugaredLogger, token, sanitizer)
	sc := securecookie.New([]byte(config.SessionKey), nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := &model.User{
				ID:              1,
				Email:           "ddjhd@mail.com",
				UserType:        "client",
			}

			if tc.expectedCode == http.StatusOK {
				store.userRepository.EXPECT().
					Find(int64(1)).
					Return(user, nil).
					AnyTimes()
			}

			b := &bytes.Buffer{}

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/logout", b)

			cookieStr, _ := sc.Encode("user-session", tc.cookie)

			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", "user-session", cookieStr))

			sess, _ := sessionStore.Get(req, "user-session")
			token, _ := token.Create(sess, time.Now().Add(24*time.Hour).Unix())
			req.Header.Set("csrf-token", token)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}*/

