package test

import (
	apiserver "github.com/go-park-mail-ru/2019_2_Comandus/internal/app"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/userUcase"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initUserRep(t * testing.T) *MockRepository{
	t.Helper()
	ctrl := gomock.NewController(t)
	return NewMockRepository(ctrl)
}

func TestUcase_CreateUser(t *testing.T) {
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

	ucase := userUcase.NewUserUsecase(initUserRep(t))


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
