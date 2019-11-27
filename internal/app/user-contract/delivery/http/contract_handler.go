package contractHttp

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/respond"
	user_contract "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract"
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

type ContractHandler struct {
	ContractUsecase	user_contract.Usecase
	sanitizer		*bluemonday.Policy
	logger			*zap.SugaredLogger
	sessionStore	sessions.Store
}

func NewContractHandler(m *mux.Router, cs user_contract.Usecase, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &ContractHandler{
		ContractUsecase:	cs,
		sanitizer:			sanitizer,
		logger:				logger,
		sessionStore:		sessionStore,
	}

	m.HandleFunc("/responses/{id:[0-9]+}/contract", handler.HandleCreateContract).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/contract/{id:[0-9]+/done}", handler.HandleTickContractAsDone).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/contract/{id:[0-9]+}", handler.HandleReviewContract).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/grades", handler.HandleGetContractsGrades).Methods(http.MethodGet, http.MethodOptions)
}

func (h * ContractHandler) HandleCreateContract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/responses/id/contract", "method":r.Method}))
	defer timer.ObserveDuration()

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateContract<-Body.Close()")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	//TODO: parse start end time here

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleCreateContract()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateContract<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	responseId := int64(id)


	if err := h.ContractUsecase.CreateContract(u, responseId); err != nil {
		err = errors.Wrapf(err, "HandleCreateContract<-UСase.CreateContract()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	// TODO: other responses should be denied
	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h * ContractHandler) HandleTickContractAsDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/contract/id/done", "method":r.Method}))
	defer timer.ObserveDuration()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleTickContractAsDone()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleTickContractAsDone<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	contractId := int64(id)
	if err := h.ContractUsecase.SetAsDone(u, contractId); err != nil {
		err = errors.Wrapf(err, "HandleTickContractAsDone<-ContractUsecase.SetAsDone()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	respond.Respond(w,r, http.StatusOK, struct{}{})
}

func (h *ContractHandler) HandleReviewContract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/contract/id", "method":r.Method}))
	defer timer.ObserveDuration()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleReviewContract()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}


	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleReviewContract()")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrapf(err, "HandleReviewContract<-ioutil.ReadAll()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	input := new(model.ReviewInput)
	if err := input.UnmarshalJSON(body); err != nil {
		err = errors.Wrapf(err, "currCompany.UnmarshalJSON()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleReviewContract<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	contractId := int64(id)

	if err := h.ContractUsecase.ReviewContract(u, contractId, input); err != nil {
		err = errors.Wrapf(err, "HandleReviewContract<-contractUsecase.ReviewContract()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.Respond(w,r, http.StatusOK, struct{}{})
}

func (h *ContractHandler) HandleGetContractsGrades(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/grades", "method":r.Method}))
	defer timer.ObserveDuration()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleGetContractsGrades()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	list, err := h.ContractUsecase.ReviewList(u)
	if err != nil {
		err = errors.Wrap(err, "HandleGetContractsGrades<-ContractUsecase.ReviewList()")
		respond.Error(w, r, http.StatusInternalServerError, err)
	}

	respond.Respond(w, r, http.StatusOK, list)
}
