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
	ctxKeyUser ctxKey = iota
	sessionName                    = "user-session"
	userTypeCookieName             = "user_type"
	hireManagerIdCookieName        = "hire-manager-id"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type server struct {
	mux *mux.Router
	store        *sqlstore.Store
	usersdb      *model.UsersDB
	sessionStore sessions.Store
	config       *Config
	clientUrl 	 string
}

func newServer(sessionStore sessions.Store, store *sqlstore.Store) *server {
	s := &server{
		mux:          mux.NewRouter(),
		usersdb:      model.NewUsersDB(),
		sessionStore: sessionStore,
		clientUrl: 	"https://comandus.now.sh",
		store: store,
	}
	s.ConfigureServer()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *server) ConfigureServer() {
	s.mux.HandleFunc("/", s.HandleMain)
	s.mux.HandleFunc("/signup", s.HandleCreateUser).Methods(http.MethodPost, http.MethodOptions)
	s.mux.HandleFunc("/login", s.HandleSessionCreate).Methods(http.MethodPost , http.MethodOptions)
	s.mux.Use(s.CORSMiddleware)

	// only for authenticated users
	private := s.mux.PathPrefix("").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/setusertype", s.HandleSetUserType).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/account", s.HandleShowProfile).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/account", s.HandleEditProfile).Methods(http.MethodPut, http.MethodOptions)
	private.HandleFunc("/account/upload-avatar", s.HandleUploadAvatar).Methods(http.MethodPost, http.MethodOptions)

	private.HandleFunc("/account/download-avatar", s.HandleDownloadAvatar).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/account/avatar/{id:[0-9]}", s.HandleGetAvatar).Methods(http.MethodGet, http.MethodOptions)

	private.HandleFunc("/account/settings/password", s.HandleEditPassword).Methods(http.MethodPut, http.MethodOptions)
	private.HandleFunc("/account/settings/notifications", s.HandleEditNotifications).Methods(http.MethodPut, http.MethodOptions)
	private.HandleFunc("/account/settings/auth-history", s.HandleGetAuthHistory).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/account/settings/security-question", s.HandleGetSecQuestion).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/account/settings/security-question", s.HandleEditSecQuestion).Methods(http.MethodPut, http.MethodOptions)
	private.HandleFunc("/account/check-security-question", s.HandleCheckSecQuestion).Methods(http.MethodPut, http.MethodOptions)
	private.HandleFunc("/roles", s.HandleRoles).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/logout", s.HandleLogout).Methods(http.MethodDelete, http.MethodOptions)
	private.HandleFunc("/jobs", s.HandleCreateJob).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/jobs/{id:[0-9]+}", s.HandleGetJob).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/freelancer", s.HandleEditFreelancer).Methods(http.MethodPut, http.MethodOptions)
	private.HandleFunc("/freelancer/{freelancerId}", s.HandleGetFreelancer).Methods(http.MethodGet, http.MethodOptions)
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
