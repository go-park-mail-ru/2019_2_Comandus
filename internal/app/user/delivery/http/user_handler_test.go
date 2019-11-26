package userHttp

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/mocks/ucase_mocks"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/token"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/middleware"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"log"
	"testing"
)

type testConfig struct {
	sessionKey		string
	tokenSecret		string
	clientUrl		string
}

type testServer struct {
	config	testConfig
	logger	*zap.SugaredLogger
	token	*token.HashToken
	sanitizer	*bluemonday.Policy
	mux			*mux.Router
	ss			sessions.Store
}

func testConf(t *testing.T) (*testServer, *ucase_mocks.MockUserUsecase) {
	t.Helper()
	config := testConfig{
		sessionKey:  	"jdfhdfdj",
		tokenSecret: 	"golangsecpark",
		clientUrl:   	"https://comandus.now.sh",
	}

	zapLogger, _ := zap.NewProduction()
	defer func() {
		if err := zapLogger.Sync(); err != nil {
			log.Println(err)
		}
	}()
	sugaredLogger := zapLogger.Sugar()

	token, err := token.NewHMACHashToken(config.tokenSecret)
	if err != nil {
		t.Fatal()
	}

	sanitizer := bluemonday.UGCPolicy()
	m := mux.NewRouter()
	ss := sessions.NewCookieStore([]byte(config.sessionKey))

	server := &testServer{
		config:    config,
		token:     token,
		sanitizer: sanitizer,
		mux:       m,
		ss:			ss,
	}

	userU := ucase_mocks.NewMockUserUsecase(gomock.NewController(t))
	mid := middleware.NewMiddleware(ss, sugaredLogger, token, userU, config.clientUrl)
	m.Use(mid.RequestIDMiddleware, mid.CORSMiddleware, mid.AccessLogMiddleware)
	NewUserHandler(m, userU, sanitizer, sugaredLogger, ss)
	return server, userU
}

/*func TestServer_HandleSetUserType(t *testing.T) {
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

	server, userU := testConf(t)
	sc := securecookie.New([]byte(server.config.sessionKey), nil)


	user := &model.User{
		ID:              1,
		Email:           "ddjhd@mail.com",
		UserType:        model.UserFreelancer,
	}

	for _, tc := range testCases {
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

			sess, _ := server.ss.Get(req, "user-session")
			token, _ := server.token.Create(sess, time.Now().Add(24*time.Hour).Unix())
			req.Header.Set("csrf-token", token)

			server.mux.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}*/


/*func TestServer_HandleLogout(t *testing.T) {
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

