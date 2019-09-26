package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/frontend-park-mail-ru/2019_2_Comandus/server/internal/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
)

type ctxKey int8
const (
	sessionName = "user-session"
	ctxKeyUser  ctxKey = iota
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type server struct {
	mux *mux.Router
	// store
	usersdb *model.UsersDB
	sessionStore sessions.Store
}

func newServer(sessionStore sessions.Store) *server {
	s := &server{mux: mux.NewRouter(), usersdb: model.NewUsersDB(), sessionStore:sessionStore}
	s.ConfigureServer()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}


// СЮДА СВОИ ХАНДЛЕРЫ
func (s * server) ConfigureServer() {
	s.mux.HandleFunc("/users", s.HandleCreateUser).Methods("POST")
	s.mux.HandleFunc("/sessions", s.HandleSessionCreate).Methods("POST")

	// only for authenticated users
	private := s.mux.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoami).Methods("GET")
	private.HandleFunc("/logout", s.HandleLogout).Methods("GET")
	private.HandleFunc("/profile", s.HandleShowProfile)
	private.HandleFunc("/profile/edit", s.HandleEditProfile)
}

func (s *server) HandleShowProfile(w http.ResponseWriter, r *http.Request) {
	// TODO Dima
	fmt.Println("HELLO profile")
}

func (s *server) HandleEditProfile(w http.ResponseWriter, r *http.Request) {
	// TODO Dima
}

func (s *server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := s.sessionStore.Get(r, sessionName)
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if session == nil {
		s.error(w, r, http.StatusNotFound, errors.New("failed to delete session"))
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		s.error(w, r, http.StatusExpectationFailed, errors.New("failed to delete session"))
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		var u *model.User
		var found bool
		for i := 0; i < len(s.usersdb.Users); i++ {
			if id == s.usersdb.Users[i].ID {
				u = &s.usersdb.Users[i]
				found = true
			}
		}

		if !found {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}


func (s *server) HandleSessionCreate(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// TODO: handle err
		r.Body.Close()
	}()

	decoder := json.NewDecoder(r.Body)
	newUserInput := new(model.UserInput)
	err := decoder.Decode(newUserInput)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	fmt.Println(newUserInput)

	for i:=0; i < len(s.usersdb.Users); i++ {
		// TODO: secure passwoord
		if s.usersdb.Users[i].Name == newUserInput.Name &&
			s.usersdb.Users[i].Password == newUserInput.Password {

			u := s.usersdb.Users[i]
			session, err := s.sessionStore.Get(r, sessionName)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			session.Values["user_id"] = u.ID
			if err := s.sessionStore.Save(r, w, session); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			s.respond(w, r, http.StatusOK, nil)
			return
		}
	}
	s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
	return
}

func (s *server) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// TODO: handle err
		r.Body.Close()
	}()
	decoder := json.NewDecoder(r.Body)
	newUserInput := new(model.UserInput)
	err := decoder.Decode(newUserInput)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	fmt.Println(newUserInput)

	s.usersdb.Mu.Lock()
	var id int64
	if len(s.usersdb.Users) > 0 {
		id = s.usersdb.Users[len(s.usersdb.Users)-1].ID + 1
		// первый почему то не ставится
		//id = DataSignerMd5(newUserInput.Name)
	}

	s.usersdb.Users = append(s.usersdb.Users, model.User{
		ID:       id,
		Name:     newUserInput.Name,
		Password: newUserInput.Password,
	})
	s.usersdb.Mu.Unlock()

	// TODO
	s.respond(w, r, http.StatusCreated, newUserInput)
}

/*func (s * server) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	s.usersdb.mu.Lock()
	err := encoder.Encode(s.usersdb.users)
	s.usersdb.mu.Unlock()

	if err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
}*/

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