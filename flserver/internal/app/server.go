package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"strconv"
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
	s.mux.HandleFunc("/profile", s.HandleShowProfile)
	// only for authenticated users
	private := s.mux.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoami).Methods("GET")
	private.HandleFunc("/logout", s.HandleLogout).Methods("GET")
	private.HandleFunc("/profile/edit", s.HandleEditProfile).Methods(http.MethodPost)
}

func (s *server) HandleShowProfile(w http.ResponseWriter, r *http.Request) {
	UserIDStr := r.FormValue("id")
	if UserIDStr != "" {
		http.Redirect(w , r , "/" , http.StatusBadRequest)
		return
	}
	UserID, err := strconv.Atoi(UserIDStr)
	if err != nil {
		log.Printf("error id isn't int: %s", err)
		http.Redirect(w, r , "/" , http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	s.usersdb.Mu.Lock()
	err = encoder.Encode(s.usersdb.GetUserByID(UserID))
	s.usersdb.Mu.Unlock()
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		http.Redirect(w, r , "/" , http.StatusBadRequest)
		return
	}
}

func (s *server) HandleEditProfile(w http.ResponseWriter, r *http.Request) {
	session, err := s.sessionStore.Get(r, sessionName)
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	uidInteface := session.Values["user_id"]
	uid := uidInteface.(int)
	user := s.usersdb.GetUserByID(uid)
	profile := s.usersdb.GetProfileByID((*user).ProfileID)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(profile)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}

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
	fmt.Println("logout")
	http.Redirect(w, r, "/", http.StatusUnauthorized)
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
			//fmt.Printf("id=%T ID=%T", id , s.usersdb.Users[i].ID)//, id == s.usersdb.Users[i].ID)
			if id == s.usersdb.Users[i].ID {
				u = &s.usersdb.Users[i]
				found = true
			}
		}
		//fmt.Println("id=", id)

		if !found {
			fmt.Println("HERE3s")
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
		if s.usersdb.Users[i].Email == newUserInput.Email &&
			s.usersdb.Users[i].ComparePassword(newUserInput.Password) {

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

	err = newUserInput.CheckEmail()
	if err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	s.usersdb.Mu.Lock()
	var id int
	var idf int
	var idc int
	var idp int //+
	if len(s.usersdb.Users) > 0 {
		id = s.usersdb.Users[len(s.usersdb.Users)-1].ID + 1
		idf = s.usersdb.Freelancers[len(s.usersdb.Freelancers)-1].ID + 1
		idc = s.usersdb.Customers[len(s.usersdb.Customers)-1].ID + 1
		idp = s.usersdb.Profiles[len(s.usersdb.Profiles) - 1].ID + 1 // +
	}

	// Может по указателю надо &model.User
	user := model.User{
		ID:              id,
		FreelancerID:    idf,
		CustomerID:      idc,
		ProfileID:		 idp,
		Name:            newUserInput.Name,
		Email:           newUserInput.Email,
		Password:        newUserInput.Password,
	}

	err = user.BeforeCreate()
	if err != nil {
		s.respond(w, r, http.StatusInternalServerError, newUserInput)
	}
	s.usersdb.Users = append(s.usersdb.Users, user)

	s.usersdb.Freelancers = append(s.usersdb.Freelancers, model.Freelancer {
		ID:       idf,
	})
	s.usersdb.Customers = append(s.usersdb.Customers, model.Customer {
		ID:       idc,
	})
	s.usersdb.Profiles = append(s.usersdb.Profiles, model.Profile { // +
		ID: idp,
	})

	fmt.Println(s.usersdb.Users[id])
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