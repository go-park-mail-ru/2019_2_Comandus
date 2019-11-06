package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (s *server) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser:<-Body.Close")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	user := new(model.User)
	err := decoder.Decode(user)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-Decode:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	fmt.Println(user)

	if err := user.Validate(); err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-Validate:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := user.BeforeCreate(); err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-BeforeCreate:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	err = s.store.User().Create(user)

	log.Println("USER ID: ", user.ID)

	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-CreateUser:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	f := model.Freelancer{
		AccountId:        user.ID,
		RegistrationDate: time.Now(),
	}

	err = s.store.Freelancer().Create(&f)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-CreateFreelancer:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	m := model.HireManager{
		AccountID:        user.ID,
		RegistrationDate: time.Now(),
	}

	err = s.store.Manager().Create(&m)

	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-CreateManager:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-sessionGet:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	session.Values["user_id"] = user.ID

	if err := s.sessionStore.Save(r, w, session); err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-sessionSave:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	user.Sanitize(s.sanitizer)
	// TODO: why we need manager cookie?
	cookie := http.Cookie{Name: hireManagerIdCookieName, Value: strconv.Itoa(1)} // m.Id
	http.SetCookie(w, &cookie)
	s.respond(w, r, http.StatusCreated, user)
}

func (s *server) AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(int64(id.(int)))

		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, &u)))
	})
}

func (s *server) HandleSessionCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleSessionCreate<-rBodyClose:")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}()
	decoder := json.NewDecoder(r.Body)
	user := new(model.User)
	err := decoder.Decode(user)
	if err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-DecodeUser:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	log.Println("input user:", user)

	u, err := s.store.User().FindByEmail(user.Email)
	if err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-FindByEmail:")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	if !u.ComparePassword(user.Password) {
		err = errors.Wrapf(errors.New("wrong password"), "HandleSessionCreate<-ComparePassword:")
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-sessionGet:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	session.Values["user_id"] = u.ID
	if err := s.sessionStore.Save(r, w, session); err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-sessionSave:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleLogout<-sessionGet:")
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		err = errors.Wrapf(err, "HandleLogout<-sessionSave:")
		s.error(w, r, http.StatusExpectationFailed, err)
	}
	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) HandleSetUserType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type Input struct {
		UserType string `json:"type"`
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleSetUserType<-rBodyClose:")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}()
	decoder := json.NewDecoder(r.Body)
	newInput := new(Input)
	err := decoder.Decode(newInput)
	if err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-Decode:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-GetUserFromRequest:")
		s.error(w, r, codeStatus, err)
		return
	}

	err = user.SetUserType(newInput.UserType)
	if err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-SetUserType:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	err = s.store.User().Edit(user)

	if err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-sessionEdit:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-sessionGet:")
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	err = session.Save(r, w)
	if err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-sessionSave:")
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	s.respond(w, r, http.StatusOK, user.UserType)
}

func (s *server) HandleRoles(w http.ResponseWriter, r *http.Request) {
	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleRoles<-GetUserFromRequest:")
		s.error(w, r, codeStatus, err)
		return
	}

	_, err = s.store.Manager().FindByUser(user.ID)
	if err != nil {
		err = errors.Wrapf(err, "HandleRoles<-ManagerFind:")
		s.error(w, r, http.StatusNotFound, err)
	}

	// TODO: rewrite after Roles and Companies db interfaces realization

	//company := s.usersdb.Companies[hireManager.ID]

	company := model.Company{
		ID:          0,
		CompanyName: "default company",
		Site:        "company.ru",
		Description: "default company",
		Country:     "Russia",
		City:        "Moscow",
		Address:     "Red square street",
		Phone:       "88888888888",
	}

	var roles []*model.Role
	clientRole := &model.Role{
		Role:   "client",
		Label:  company.CompanyName,
		Avatar: "../store//default.png",
	}
	freelanceRole := &model.Role{
		Role:   "freelancer",
		Label:  user.FirstName + " " + user.SecondName,
		Avatar: "/default.png",
	}
	roles = append(roles, clientRole)
	roles = append(roles, freelanceRole)
	s.respond(w, r, http.StatusOK, roles)
}

func (s * server) HandleGetToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sess , err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetToken<-sessionStore.Get :")
		s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
		return
	}

	token, err := s.token.Create(sess, time.Now().Add(24*time.Hour).Unix())
	s.respond(w, r, http.StatusOK, map[string]string{"csrf-token" : token})
}