package userHttp

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/respond"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/monitoring"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	sessionName                    = "user-session"
)

type ResponseError struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserUsecase		user.Usecase
	sanitizer		*bluemonday.Policy
	logger			*zap.SugaredLogger
	sessionStore	sessions.Store
}

func NewUserHandler(m *mux.Router, us user.Usecase, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &UserHandler{
		UserUsecase:	us,
		sanitizer:		sanitizer,
		logger:			logger,
		sessionStore:	sessionStore,
	}

	m.HandleFunc("/account", handler.HandleShowProfile).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/account", handler.HandleEditProfile).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/account/settings/password", handler.HandleEditPassword).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/account/upload-avatar", handler.HandleUploadAvatar).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/account/download-avatar", handler.HandleDownloadAvatar).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/account/avatar/{id:[0-9]+}", handler.HandleGetAvatar).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/setusertype", handler.HandleSetUserType).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/roles", handler.HandleRoles).Methods(http.MethodGet, http.MethodOptions)
}

func (h *UserHandler) HandleShowProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/account", "method":r.Method}))
	defer timer.ObserveDuration()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleShowProfile()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	u.Sanitize(h.sanitizer)
	respond.Respond(w, r, http.StatusOK, u)
}

func (h *UserHandler) HandleEditProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/account", "method":r.Method}))
	defer timer.ObserveDuration()

	currUser, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleEditProfile()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditProfile<-rBodyClose()")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditProfile<-ioutil.ReadAll()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	userInput := currUser
	if err := userInput.UnmarshalJSON(body); err != nil {
		err = errors.Wrapf(err, "HandleEditProfile<-currCompany.UnmarshalJSON()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.UserUsecase.EditUser(userInput, currUser); err != nil {
		err = errors.Wrapf(err, "HandleEditProfile<-Ucase.EditUser()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	log.Println("edited profile: ", userInput)
	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *UserHandler) HandleEditPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/account/settings/password", "method":r.Method}))
	defer timer.ObserveDuration()

	currUser, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleEditPassword()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditPassword<-rBodyClose()")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditPassword<-ioutil.ReadAll()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	bodyPassword := new(model.BodyPassword)
	if err := bodyPassword.UnmarshalJSON(body); err != nil {
		err = errors.Wrapf(err, "bodyPassword.UnmarshalJSON()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := h.UserUsecase.EditUserPassword(bodyPassword, currUser); err != nil {
		err = errors.Wrap(err, "HandleEditPassword<-UCase.EditUserPassword()")
		respond.Respond(w, r, http.StatusBadRequest, err)
		return
	}
	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *UserHandler) HandleUploadAvatar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/account/upload-avatar", "method":r.Method}))
	defer timer.ObserveDuration()

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		err = errors.Wrapf(err, "HandleUploadAvatar<-ParseMultipartForm()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		err = errors.Wrapf(err, "HandleUploadAvatar<-FormFile()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	currUser, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleUploadAvatar")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	image := bytes.NewBuffer(nil)
	_, err = io.Copy(image, file)
	if err != nil {
		err = errors.Wrapf(err, "HandleUploadAvatar<-ioCopy()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	newUser := currUser
	newUser.Avatar = image.Bytes()
	if err := h.UserUsecase.EditUser(newUser, currUser); err != nil {
		err = errors.Wrapf(err, "HandleUploadAvatar<-UCase.EditUser()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *UserHandler) HandleDownloadAvatar(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/account/download-avatar", "method":r.Method}))
	defer timer.ObserveDuration()

	currUser, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleDownloadAvatar()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	avatar, err := h.UserUsecase.GetAvatar(currUser)
	if err != nil {
		err := errors.Wrapf(err, "HandleDownloadAvatar<-UserUseCase.GetAvatar()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	if _, err := w.Write(avatar); err != nil {
		err = errors.Wrapf(err, "HandleDownloadAvatar<-Write()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	log.Println("IMAGE SIZE:", len(avatar))

	w.Header().Set("Content-Disposition", "attachment; filename=avatar")
	w.Header().Set("Content-Type", "multipart/form-data")
	w.Header().Set("Content-Length", strconv.Itoa(len(avatar)))

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *UserHandler) HandleGetAvatar(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/account/avatar/id", "method":r.Method}))
	defer timer.ObserveDuration()

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetAvatar<-Atoi(wrong id)")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	currUser, err := h.UserUsecase.Find(int64(id))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetAvatar<-UserUseCase.Find()")
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	avatar, err := h.UserUsecase.GetAvatar(currUser)
	if err != nil {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleGetAvatar<-UserUseCase.GetAvatar()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	if _, err := w.Write(avatar); err != nil {
		err = errors.Wrapf(err, "HandleGetAvatar<-Write()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=avatar")
	w.Header().Set("Content-Type", "multipart/form-data")
	w.Header().Set("Content-Length", strconv.Itoa(len(avatar)))

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *UserHandler) HandleSetUserType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/setusertype", "method":r.Method}))
	defer timer.ObserveDuration()

	type Input struct {
		UserType string `json:"type"`
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleSetUserType<-rBodyClose()")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	newInput := new(Input)
	err := decoder.Decode(newInput)
	if err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-Decode()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	currUser, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleSetUserType: ")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := h.UserUsecase.SetUserType(currUser, newInput.UserType); err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-UCase.SetUserType()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, currUser.UserType)
}

func (h *UserHandler) HandleRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/roles", "method":r.Method}))
	defer timer.ObserveDuration()

	currUser, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleRoles()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	roles, err := h.UserUsecase.GetRoles(currUser)
	if err != nil {
		err := errors.Wrapf(err, "HandleRoles<-UserUsecase.GetRoles()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, roles)
}