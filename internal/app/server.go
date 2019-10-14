package apiserver

import (
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
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
	store        *sqlstore.Store
	usersdb      *model.UsersDB
	sessionStore sessions.Store
	config       *Config
	userType     string
	clientUrl 	 string
}

func newServer(sessionStore sessions.Store) *server {
	s := &server{
		mux:          mux.NewRouter(),
		usersdb:      model.NewUsersDB(),
		sessionStore: sessionStore,
		clientUrl: 	"https://comandus.now.sh",
	}
	s.ConfigureServer()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *server) ConfigureStore() error {
	cnfg := sqlstore.NewConfig()
	st := sqlstore.New(cnfg)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}

func (s *server) ConfigureServer() {
	s.mux.HandleFunc("/", s.HandleMain)
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
	private.HandleFunc("/account/avatar/{id:[0-9]}", s.HandleGetAvatar).Methods(http.MethodGet)

	private.HandleFunc("/account/settings/password", s.HandleEditPassword).Methods(http.MethodPut)
	private.HandleFunc("/account/settings/notifications", s.HandleEditNotifications).Methods(http.MethodPut)
	private.HandleFunc("/account/settings/auth-history", s.HandleGetAuthHistory).Methods(http.MethodGet)
	private.HandleFunc("/account/settings/security-question", s.HandleGetSecQuestion).Methods(http.MethodGet)
	private.HandleFunc("/account/settings/security-question", s.HandleEditSecQuestion).Methods(http.MethodPut)
	private.HandleFunc("/account/check-security-question", s.HandleCheckSecQuestion).Methods(http.MethodPut)
	private.HandleFunc("/roles", s.HandleRoles).Methods(http.MethodGet)
	private.HandleFunc("/logout", s.HandleLogout).Methods(http.MethodDelete)
	private.HandleFunc("/jobs", s.HandleCreateJob).Methods(http.MethodPost)
	private.HandleFunc("/jobs/{id:[0-9]+}", s.HandleGetJob).Methods(http.MethodGet)
	private.HandleFunc("/freelancer", s.HandleEditFreelancer).Methods(http.MethodPut)
	private.HandleFunc("/freelancer/{freelancerId}", s.HandleGetFreelancer).Methods(http.MethodGet)


	// TODO: fix wrong paths
	s.mux.HandleFunc("/signup", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/login", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/private/freelancer/{freelancerId}", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/private/logout", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/setusertype", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/private/account", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/account/upload-avatar", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/private/account/settings/password", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/private/account/settings/notifications", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/account/settings/security-question", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/account/check-security-question", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/roles", s.HandleOptions).Methods(http.MethodOptions)
	s.mux.HandleFunc("/private/jobs", s.HandleOptions).Methods(http.MethodOptions)
}

func (s * server) HandleMain(w http.ResponseWriter, r *http.Request) {
	s.respond(w,r,http.StatusOK, "hello from server")
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
