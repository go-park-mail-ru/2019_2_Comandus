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
	sessionName = "user-session"
	ctxKeyUser  ctxKey = iota
	userFreelancer = "freelancer"
	userCustomer = "customer"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type server struct {
	mux *mux.Router
	// store
	store *store.Store
	usersdb *model.UsersDB
	sessionStore sessions.Store
	config *Config
}

func newServer(sessionStore sessions.Store) *server {
	s := &server{
		mux: mux.NewRouter(),
		usersdb: model.NewUsersDB(),
		sessionStore:sessionStore,
		//store: store,
	}
	s.ConfigureServer()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s * server) ConfigureStore() error {
	cnfg := store.NewConfig()
	st := store.New(cnfg)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}

// СЮДА СВОИ ХАНДЛЕРЫ
func (s * server) ConfigureServer() {
	s.mux.HandleFunc("/users", s.HandleCreateUser)//.Methods("POST")
	s.mux.HandleFunc("/sessions", s.HandleSessionCreate).Methods("POST")
	// only for authenticated users
	private := s.mux.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/setusertype", s.HandleSetUserType).Methods("POST")
	private.HandleFunc("/profile", s.HandleShowProfile)
	//private.HandleFunc("/whoami", s.handleWhoami).Methods("GET")
	private.HandleFunc("/logout", s.HandleLogout).Methods("GET")
	private.HandleFunc("/profile/edit", s.HandleEditProfile).Methods(http.MethodPost)
}


func (s *server) handleWhoami(w http.ResponseWriter, r *http.Request)  {
	s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
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