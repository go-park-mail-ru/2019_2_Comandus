package companyhttp

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/respond"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/monitoring"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"io/ioutil"
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

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/company", "method":r.Method}))
	defer timer.ObserveDuration()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleEditCompany()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditCompany<-rBodyClose()")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()


	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditPassword<-ioutil.ReadAll()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	currCompany := new(model.Company)
	if err := currCompany.UnmarshalJSON(body); err != nil {
		err = errors.Wrapf(err, "currCompany.UnmarshalJSON()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	if _, err := h.CompanyUsecase.Edit(u.ID, currCompany); err != nil {
		err = errors.Wrapf(err, "HandleGetCountries<-CompanyUsecase.Edit()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *CompanyHandler) HandleGetCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/company/id", "method":r.Method}))
	defer timer.ObserveDuration()

	vars := mux.Vars(r)
	ids := vars["companyId"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetCompany<-Atoi(wrong id)")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	currCompany, err := h.CompanyUsecase.Find(int64(id))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetCompany<-Find()")
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, currCompany)
}
