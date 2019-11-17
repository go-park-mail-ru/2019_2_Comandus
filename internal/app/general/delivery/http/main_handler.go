package mainHttp

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/token"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type MainHandler struct {
	UserUsecase		user.Usecase
	sanitizer		*bluemonday.Policy
	logger			*zap.SugaredLogger
	sessionStore	sessions.Store
	token 			*token.HashToken
}

func NewMainHandler(m *mux.Router,private *mux.Router, us user.Usecase, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger,
	sessionStore sessions.Store, thisToken *token.HashToken) {
		handler := &MainHandler{
		UserUsecase:	us,
		sanitizer:		sanitizer,
		logger:			logger,
		sessionStore:	sessionStore,
		token:			thisToken,
	}

	m.HandleFunc("/", handler.HandleMain)
	m.HandleFunc("/signup", handler.HandleCreateUser).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/login", handler.HandleSessionCreate).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/token", handler.HandleGetToken).Methods(http.MethodGet)
	private.HandleFunc("/logout", handler.HandleLogout).Methods(http.MethodDelete, http.MethodOptions)
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

	session, err := h.sessionStore.Get(r, general.SessionName)
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

	session, err := h.sessionStore.Get(r, general.SessionName)
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

func (h *MainHandler) HandleGetToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sess, err := h.sessionStore.Get(r, general.SessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetToken<-sessionStore.Get :")
		general.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	currToken, err := h.token.Create(sess, time.Now().Add(24*time.Hour).Unix())
	general.Respond(w, r, http.StatusOK, map[string]string{"csrf-token": currToken})
}

func (h * MainHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionStore.Get(r, general.SessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleLogout<-sessionGet:")
		general.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		err = errors.Wrapf(err, "HandleLogout<-sessionSave:")
		general.Error(w, r, http.StatusExpectationFailed, err)
		return
	}
	general.Respond(w, r, http.StatusOK, struct{}{})
}
