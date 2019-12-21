package jobHttp

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

	ctrl := gomock.NewController(t)
	userU := ucase_mocks.NewMockUserUsecase(ctrl)
	jobU := ucase_mocks.NewMockJobUsecase(ctrl)

	mid := middleware.NewMiddleware(s.SessionStore, s.Logger, s.Token, userU, s.Config.ClientUrl)
	s.Mux.Use(mid.RequestIDMiddleware, mid.CORSMiddleware, mid.AccessLogMiddleware)
	private := s.Mux.PathPrefix("").Subrouter()
	private.Use(mid.AuthenticateUser, mid.CheckTokenMiddleware)
	NewJobHandler(private, jobU, s.Sanitizer, s.Logger, s.SessionStore)

	sc := securecookie.New([]byte(s.Config.SessionKey), nil)
	user := &model.User{
		ID:              1,
		Email:           "ddjhd@mail.com",
		UserType:        model.UserCustomer,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user.UserType = tc.userType



			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/jobs", b)

			cookieStr, _ := sc.Encode("user-session", map[interface{}]interface{}{
				"user_id": 1,
			})

			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", "user-session", cookieStr))

			sess, _ := s.SessionStore.Get(req, "user-session")
			token, _ := s.Token.Create(sess, time.Now().Add(24*time.Hour).Unix())
			req.Header.Set("csrf-token", token)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}*/
