package mainHttp

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

type ctxKey int8

const (
	CtxKeyUser              ctxKey = iota
	sessionName                    = "user-session"
)

type ResponseError struct {
	Message string `json:"message"`
}

type MainHandler struct {
	UserUsecase		user.Usecase
	sanitizer		*bluemonday.Policy
	logger			*zap.SugaredLogger
	sessionStore	sessions.Store
}

func NewMainHandler(m *mux.Router, us user.Usecase, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &MainHandler{
		UserUsecase:	us,
		sanitizer:		sanitizer,
		logger:			logger,
		sessionStore:	sessionStore,
	}

	m.HandleFunc("/", handler.HandleMain)
	m.HandleFunc("/signup", handler.HandleCreateUser).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/login", handler.HandleSessionCreate).Methods(http.MethodPost, http.MethodOptions)
}

func (h *MainHandler) HandleMain(w http.ResponseWriter, r *http.Request) {
	general.Respond(w, r, http.StatusOK, "hello from server")
}

func (h *MainHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser:<-Body.Close: ")
			general.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	newUser := new(model.User)
	err := decoder.Decode(newUser)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-Decode: ")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.UserUsecase.CreateUser(newUser); err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-UserUsecase.CreateUser(): ")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	session, err := h.sessionStore.Get(r, sessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-sessionGet: ")
		general.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	session.Values["user_id"] = newUser.ID
	if err := h.sessionStore.Save(r, w, session); err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-sessionSave: ")
		general.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	newUser.Sanitize(h.sanitizer)
	general.Respond(w, r, http.StatusCreated, newUser)
}

func (h * MainHandler) HandleSessionCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleSessionCreate<-rBodyClose:")
			general.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	currUser := new(model.User)
	err := decoder.Decode(currUser)
	if err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-DecodeUser:")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.UserUsecase.VerifyUser(currUser)
	if err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-UserUseCase.VerifyUser(): ")
		general.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	session, err := h.sessionStore.Get(r, sessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-sessionGet:")
		general.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	session.Values["user_id"] = id
	if err := h.sessionStore.Save(r, w, session); err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-sessionSave:")
		general.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	general.Respond(w, r, http.StatusOK, struct{}{})
}