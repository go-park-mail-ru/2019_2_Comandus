package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"log"
	"net/http"
)
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
	var idp int = 1 //+
	if len(s.usersdb.Users) > 0 {
		id = s.usersdb.Users[len(s.usersdb.Users)-1].ID + 1
		idf = s.usersdb.Freelancers[len(s.usersdb.Freelancers)-1].ID + 1
		idc = s.usersdb.Customers[len(s.usersdb.Customers)-1].ID + 1
		idp = s.usersdb.Profiles[len(s.usersdb.Profiles) - 1].ID + 1 // +
	}

	user := model.User{
		ID:              id,
		FreelancerID:    idf,
		CustomerID:      idc,
		Name:            newUserInput.Name,
		Email:           newUserInput.Email,
		Password:        newUserInput.Password,
	}
	profile := model.Profile{
		ID:idp,
		ContactInformation: model.ContactInfo{},
		AdditionalInfo: model.AdditionalInformation{},
		InnerInformation: model.InnerInfo{},
	}
	err = user.BeforeCreate()
	if err != nil {
		s.respond(w, r, http.StatusInternalServerError, newUserInput)
	}
	s.usersdb.Users = append(s.usersdb.Users, user)

	s.usersdb.Freelancers = append(s.usersdb.Freelancers, model.Freelancer {
		ID:       idf,
		ProfileID: idp,
	})
	s.usersdb.Customers = append(s.usersdb.Customers, model.Customer {
		ID:       idc,
		ProfileID: idp,
	})
	s.usersdb.Profiles = append(s.usersdb.Profiles, profile)

	fmt.Println(s.usersdb.Users[id])
	s.usersdb.Mu.Unlock()

	// TODO
	s.respond(w, r, http.StatusCreated, newUserInput)
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
		if s.usersdb.Users[i].Email == newUserInput.Email &&
			s.usersdb.Users[i].ComparePassword(newUserInput.Password) {

			u := s.usersdb.Users[i]
			session, err := s.sessionStore.Get(r, sessionName)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			session.Values["user_id"] = u.ID
			session.Values["user_type"] = userFreelancer
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

func (s * server) HandleSetUserType(w http.ResponseWriter, r *http.Request) {
	// TODO check if input user type invalid
	type Input struct {
		UserType     string `json:"user_type"`
	}
	defer func() {
		// TODO: handle err
		r.Body.Close()
	}()

	decoder := json.NewDecoder(r.Body)
	newInput := new(Input)
	err := decoder.Decode(newInput)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	session, err := s.sessionStore.Get(r, sessionName)
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	session.Values["user_type"] = newInput.UserType
	session.Save(r,w)
}

func (s *server) HandleShowProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SHOW PROFILE")
	session, err := s.sessionStore.Get(r, sessionName)
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	uidInteface := session.Values["user_id"]
	uid := uidInteface.(int)
	utypeInteface := session.Values["user_type"]
	utype := utypeInteface.(string)

	user := s.usersdb.GetUserByID(uid)
	if utype == userFreelancer {
		fmt.Println(userFreelancer)
		freelancer := s.usersdb.GetFreelancerByID(user.FreelancerID)
		profile := s.usersdb.GetProfileByID(freelancer.ProfileID)
		fmt.Println(*profile)
		s.respond(w, r, http.StatusOK, *profile)//r.Context().Value(ctxKeyUser).(model.Freelancer))
	} else if utype == userCustomer {
		fmt.Println(userCustomer)
		customer := s.usersdb.GetCustomerByID(user.CustomerID)
		profile := s.usersdb.GetProfileByID(customer.ProfileID)
		fmt.Println(*profile)
		s.respond(w, r, http.StatusOK, profile)
	} else {
		s.error(w,r, http.StatusInternalServerError, errors.New("user type not set"))
	}
}

func (s *server) HandleEditProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Edit profile")
	var profile *model.Profile
	session, err := s.sessionStore.Get(r, sessionName)
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	uidInteface := session.Values["user_id"]
	uid := uidInteface.(int)
	utypeInteface := session.Values["user_type"]
	utype := utypeInteface.(string)
	user := s.usersdb.GetUserByID(uid)
	if utype == userFreelancer {
		fmt.Println(userFreelancer)
		freelancer := s.usersdb.GetFreelancerByID(user.FreelancerID)
		profile = s.usersdb.GetProfileByID(freelancer.ProfileID)
	} else if utype == userCustomer {
		fmt.Println(userCustomer)
		customer := s.usersdb.GetCustomerByID(user.CustomerID)
		profile = s.usersdb.GetProfileByID(customer.ProfileID)
	} else {
		s.error(w,r, http.StatusInternalServerError, errors.New("user type not set"))
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(profile)
	fmt.Println(*profile)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}

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