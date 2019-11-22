package companyhttp

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/respond"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Message string `json:"message"`
}

type CompanyHandler struct {
	CompanyUsecase company.Usecase
	sanitizer      *bluemonday.Policy
	logger         *zap.SugaredLogger
	sessionStore   sessions.Store
}

func NewCompanyHandler(m *mux.Router, uc company.Usecase, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &CompanyHandler{
		CompanyUsecase: uc,
		sanitizer:      sanitizer,
		logger:         logger,
		sessionStore:   sessionStore,
	}

	m.HandleFunc("/company", handler.HandleEditCompany).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/company/{companyId}", handler.HandleGetCompany).Methods(http.MethodGet, http.MethodOptions)
}

func (h *CompanyHandler) HandleEditCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleEditCompany: ")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditCompany<-rBodyClose: ")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	currCompany := new(model.Company)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(currCompany)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditCompany<-Decode(): ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.CompanyUsecase.Edit(u, currCompany); err != nil {
		err = errors.Wrapf(err, "HandleEditCompany<-CompanyUsecase.Edit(): ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *CompanyHandler) HandleGetCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")


	vars := mux.Vars(r)
	ids := vars["companyId"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetCompany<-Atoi(wrong id): ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	currCompany, err := h.CompanyUsecase.Find(int64(id))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetCompany<-Find: ")
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	currCompany.Sanitize(h.sanitizer)
	respond.Respond(w, r, http.StatusOK, currCompany)
}
