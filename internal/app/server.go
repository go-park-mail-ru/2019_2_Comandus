package apiserver

import (
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
)

type ctxKey int8

const (
	sessionName                    = "user-session"
	ctxKeyUser              ctxKey = iota
	userFreelancer                 = "freelancer"
	userCustomer                   = "client"
	userTypeCookieName             = "user_type"
	hireManagerIdCookieName        = "hire-manager-id"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type server struct {
	mux *mux.Router
	// store
	store        *store.Store
	usersdb      *model.UsersDB
	sessionStore sessions.Store
	config       *Config
	userType     string
}

func newServer(sessionStore sessions.Store) *server {
	s := &server{
		mux:          mux.NewRouter(),
		usersdb:      model.NewUsersDB(),
		sessionStore: sessionStore,
		//store: store,
	}
	s.ConfigureServer()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *server) ConfigureStore() error {
	cnfg := store.NewConfig()
	st := store.New(cnfg)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}

/*var uploadFormTmpl = []byte(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Document</title>
  </head>
  <body>
    <form
      enctype="multipart/form-data"
      action="http://localhost:8080/account/upload-avatar"
      method="post"
    >
      <input type="file" name="myFile" />
      <input type="submit" value="upload" />
    </form>
  </body>
</html>
`)
func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write(uploadFormTmpl)
}*/

// СЮДА СВОИ ХАНДЛЕРЫ
func (s *server) ConfigureServer() {
	s.mux.HandleFunc("/signup", s.HandleCreateUser).Methods(http.MethodPost)
	s.mux.HandleFunc("/login", s.HandleSessionCreate).Methods(http.MethodPost)
	// only for authenticated users
	private := s.mux.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/setusertype", s.HandleSetUserType).Methods(http.MethodPost)
	private.HandleFunc("/account", s.HandleShowProfile).Methods(http.MethodGet)
	private.HandleFunc("/account", s.HandleEditProfile).Methods(http.MethodPut)
	private.HandleFunc("/account/upload-avatar", s.HandleUploadAvatar).Methods(http.MethodPost)
	private.HandleFunc("/account/download-avatar", s.HandleDownloadAvatar).Methods(http.MethodGet)
	private.HandleFunc("/account/settings/password", s.HandleEditPassword).Methods(http.MethodPut)
	private.HandleFunc("/account/settings/notifications", s.HandleEditNotifications).Methods(http.MethodPut)
	private.HandleFunc("/account/settings/auth-history", s.HandleGetAuthHistory).Methods(http.MethodGet)
	private.HandleFunc("/account/settings/security-question", s.HandleGetSecQuestion).Methods(http.MethodGet)
	private.HandleFunc("/account/settings/security-question", s.HandleEditSecQuestion).Methods(http.MethodPut)
	private.HandleFunc("/account/check-security-question", s.HandleCheckSecQuestion).Methods(http.MethodPut)
	private.HandleFunc("/roles", s.HandleRoles).Methods(http.MethodGet)
	private.HandleFunc("/logout", s.HandleLogout).Methods(http.MethodPost)
	private.HandleFunc("/jobs", s.HandleCreateJob).Methods(http.MethodPost)
	private.HandleFunc("/jobs/{id:[0-9]+}", s.HandleGetJob).Methods(http.MethodGet)

	// TODO: fix wrong paths
	/*s.mux.HandleFunc("/signup", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/login", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/logout", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/setusertype", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/private/account", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/account/upload-avatar", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/private/account/settings/password", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/private/account/settings/notifications", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/account/settings/security-question", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/account/check-security-question", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/roles", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/private/jobs", s.HandleOptions).Methods(http.MethodOptions)*/
	s.mux.HandleFunc("/signup", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/login", s.HandleOptions).Methods(http.MethodOptions)

	private.HandleFunc("/logout", s.HandleOptions).Methods(http.MethodOptions)
	private.HandleFunc("/setusertype", s.HandleOptions).Methods(http.MethodOptions)
	private.HandleFunc("/account", s.HandleOptions).Methods(http.MethodOptions)
	private.HandleFunc("/account/upload-avatar", s.HandleOptions).Methods(http.MethodOptions)
	private.HandleFunc("/account/settings/password", s.HandleOptions).Methods(http.MethodOptions)
	private.HandleFunc("/account/settings/notifications", s.HandleOptions).Methods(http.MethodOptions)
	private.HandleFunc("/account/settings/security-question", s.HandleOptions).Methods(http.MethodOptions)
	private.HandleFunc("/account/check-security-question", s.HandleOptions).Methods(http.MethodOptions)
	private.HandleFunc("/roles", s.HandleOptions).Methods(http.MethodOptions)
	private.HandleFunc("/jobs", s.HandleOptions).Methods(http.MethodOptions)
}

// error handlers
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
