package userHttp

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
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

	u, ok := r.Context().Value(general.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleShowProfile: ")
		general.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	u.Sanitize(h.sanitizer)
	general.Respond(w, r, http.StatusOK, u)
}

func (h *UserHandler) HandleEditProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	currUser, ok := r.Context().Value(general.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleEditProfile: ")
		general.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditProfile<-rBodyClose: ")
			general.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	userInput := new(model.User)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(userInput)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditProfile<-Decode: ")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.UserUsecase.EditUser(userInput, currUser); err != nil {
		err = errors.Wrapf(err, "HandleEditProfile<-EditUser: ")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	log.Println("edited profile: ", userInput)
	general.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *UserHandler) HandleEditPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	currUser, ok := r.Context().Value(general.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleEditProfile: ")
		general.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditPassword<-rBodyClose:")
			general.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	bodyPassword := new(model.BodyPassword)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(bodyPassword)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditPassword<-Decode: ")
		general.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := h.UserUsecase.EditUserPassword(bodyPassword, currUser); err != nil {
		err = errors.Wrap(err, "HandleEditPassword<-UserUseCase.EditUserPassword: ")
		general.Respond(w, r, http.StatusBadRequest, err)
		return
	}
	general.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *UserHandler) HandleUploadAvatar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		err = errors.Wrapf(err, "HandleUploadAvatar<-ParseMultipartForm:")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		err = errors.Wrapf(err, "HandleUploadAvatar<-FormFile:")
		general.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			general.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	currUser, ok := r.Context().Value(general.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleEditProfile: ")
		general.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	image := bytes.NewBuffer(nil)
	_, err = io.Copy(image, file)
	if err != nil {
		err = errors.Wrapf(err, "HandleUploadAvatar<-ioCopy:")
		general.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	currUser.Avatar = image.Bytes()
	if err := h.UserUsecase.EditUser(currUser, currUser); err != nil {
		err = errors.Wrapf(err, "HandleUploadAvatar<-UserUsecase.EditUser(): ")
		general.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	general.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *UserHandler) HandleDownloadAvatar(w http.ResponseWriter, r *http.Request) {
	currUser, ok := r.Context().Value(general.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleDownloadAvatar: ")
		general.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	avatar, err := h.UserUsecase.GetAvatar(currUser)
	if err != nil {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleEditProfile<-UserUseCase.GetAvatar(): ")
		general.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	if _, err := w.Write(avatar); err != nil {
		err = errors.Wrapf(err, "HandleDownloadAvatar<-Write():")
		general.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=avatar")
	w.Header().Set("Content-Type", "multipart/form-data")
	w.Header().Set("Content-Length", strconv.Itoa(len(avatar)))

	general.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *UserHandler) HandleGetAvatar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetAvatar<-Atoi(wrong id): ")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	currUser, err := h.UserUsecase.Find(int64(id))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetAvatar<-UserUseCase.Find(): ")
		general.Error(w, r, http.StatusNotFound, err)
		return
	}

	avatar, err := h.UserUsecase.GetAvatar(currUser)
	if err != nil {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleEditProfile<-UserUseCase.GetAvatar(): ")
		general.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	if _, err := w.Write(avatar); err != nil {
		err = errors.Wrapf(err, "HandleDownloadAvatar<-Write():")
		general.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=avatar")
	w.Header().Set("Content-Type", "multipart/form-data")
	w.Header().Set("Content-Length", strconv.Itoa(len(avatar)))

	general.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *UserHandler) HandleSetUserType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type Input struct {
		UserType string `json:"type"`
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleSetUserType<-rBodyClose:")
			general.Error(w, r, http.StatusInternalServerError, err)
		}
	}()


	decoder := json.NewDecoder(r.Body)
	newInput := new(Input)
	err := decoder.Decode(newInput)
	if err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-Decode:")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	currUser, ok := r.Context().Value(general.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleSetUserType: ")
		general.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := h.UserUsecase.SetUserType(currUser, newInput.UserType); err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-UserUsecaseSetUserType:")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	general.Respond(w, r, http.StatusOK, currUser.UserType)
}

func (h * UserHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionStore.Get(r, sessionName)
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

func (h *UserHandler) HandleRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currUser, ok := r.Context().Value(general.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleRoles: ")
		general.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	roles, err := h.UserUsecase.GetRoles(currUser)
	if err != nil {
		err := errors.Wrapf(err, "HandleRoles<-UserUsecase.GetRoles(): ")
		general.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	general.Respond(w, r, http.StatusOK, roles)
}