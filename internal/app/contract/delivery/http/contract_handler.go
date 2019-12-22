package contractHttp

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/respond"
	user_contract "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/contract"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/monitoring"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Message string `json:"message"`
}

type ContractHandler struct {
	ContractUsecase user_contract.Usecase
	sanitizer       *bluemonday.Policy
	logger          *zap.SugaredLogger
	sessionStore    sessions.Store
}

func NewContractHandler(m *mux.Router, cs user_contract.Usecase, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &ContractHandler{
		ContractUsecase: cs,
		sanitizer:       sanitizer,
		logger:          logger,
		sessionStore:    sessionStore,
	}

	m.HandleFunc("/proposals/{id:[0-9]+}/contract", handler.HandleCreateContract).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/contract/{id:[0-9]+}/freelancer/accept", handler.HandleFreelancerAccept).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/contract/{id:[0-9]+}/freelancer/deny", handler.HandleFreelancerDeny).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/contract/{id:[0-9]+}/freelancer/ready", handler.HandleFreelancerReady).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/contract/{id:[0-9]+}/done", handler.HandleTickContractAsDone).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/contract/{id:[0-9]+}/review", handler.HandleReviewContract).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/contracts", handler.HandleGetContracts).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/contracts/{id:[0-9]+}", handler.HandleGetContract).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/contracts/archive/{freelancerID:[0-9]+}", handler.HandleGetClosedContracts).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/grades", handler.HandleGetContractsGrades).Methods(http.MethodGet, http.MethodOptions)
}

func (h *ContractHandler) HandleCreateContract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/responses/id/contract", "method": r.Method}))
	defer timer.ObserveDuration()

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateContract<-Body.Close()")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleCreateContract()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateContract<-ioutil.ReadAll()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	input := new(model.ContractInput)
	if err := input.UnmarshalJSON(body); err != nil {
		err = errors.Wrapf(err, "CreateContract.UnmarshalJSON()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	log.Println(input)

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateContract<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	responseId := int64(id)

	if err := h.ContractUsecase.CreateContract(u, responseId, input); err != nil {
		err = errors.Wrapf(err, "HandleCreateContract<-UСase.CreateContract()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	// TODO: other responses should be denied
	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *ContractHandler) HandleTickContractAsDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/contract/id/done", "method": r.Method}))
	defer timer.ObserveDuration()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleTickContractAsDone()")
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

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *ContractHandler) HandleReviewContract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/contract/id", "method": r.Method}))
	defer timer.ObserveDuration()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleReviewContract()")
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
		err = errors.Wrapf(err, "Review.UnmarshalJSON()")
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

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *ContractHandler) HandleGetContractsGrades(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/grades", "method": r.Method}))
	defer timer.ObserveDuration()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleGetContractsGrades()")
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

func (h *ContractHandler) HandleGetContracts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/contracts", "method": r.Method}))
	defer timer.ObserveDuration()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleGetContracts()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	list, err := h.ContractUsecase.ContractList(u)
	if err != nil {
		err = errors.Wrap(err, "HandleGetContracts<-ContractUsecase.ContractList()")
		respond.Error(w, r, http.StatusInternalServerError, err)
	}

	respond.Respond(w, r, http.StatusOK, list)
}

func (h *ContractHandler) HandleGetContract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/contracts/{id}", "method": r.Method}))
	defer timer.ObserveDuration()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleGetContracts()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetContract<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	contract, err := h.ContractUsecase.Find(u, id)
	if err != nil {
		err = errors.Wrap(err, "HandleGetContract<-ContractUsecase.Find()")
		respond.Error(w, r, http.StatusInternalServerError, err)
	}

	respond.Respond(w, r, http.StatusOK, contract)
}

func (h *ContractHandler) HandleFreelancerAccept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/contract/{id}/freelancer/accept", "method": r.Method}))
	defer timer.ObserveDuration()

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleFreelancerAccept<-Body.Close()")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleFreelancerAccept()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleFreelancerAccept<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	contractID := int64(id)

	if err := h.ContractUsecase.ChangeStatus(u, contractID, model.ContractStatusUnderDevelopment); err != nil {
		err = errors.Wrapf(err, "HandleFreelancerAccept<-UСase.ChangeStatus()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *ContractHandler) HandleFreelancerDeny(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/contract/{id}/freelancer/deny", "method": r.Method}))
	defer timer.ObserveDuration()

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleFreelancerDeny<-Body.Close()")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleFreelancerDeny()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleFreelancerDeny<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	contractID := int64(id)

	if err := h.ContractUsecase.ChangeStatus(u, contractID, model.ContractStatusDenied); err != nil {
		err = errors.Wrapf(err, "HandleFreelancerDeny<-UСase.ChangeStatus()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *ContractHandler) HandleFreelancerReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/contract/{id}/freelancer/ready", "method": r.Method}))
	defer timer.ObserveDuration()

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleFreelancerReady<-Body.Close()")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleFreelancerReady()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleFreelancerReady<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	contractID := int64(id)

	if err := h.ContractUsecase.TickWorkAsReady(u, contractID); err != nil {
		err = errors.Wrapf(err, "HandleFreelancerReady<-UСase.ChangeStatus()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, struct{}{})
}



func (h *ContractHandler) HandleGetClosedContracts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
	Labels{"path": "/contracts/archive/{freelancerID}", "method": r.Method}))
	defer timer.ObserveDuration()

	vars := mux.Vars(r)
	ids := vars["freelancerID"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetContracts<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	freelancerID := int64(id)

	publicContracts , err := h.ContractUsecase.GetClosedContracts(freelancerID)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetContracts<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	for i, _ := range publicContracts {
		publicContracts[i].Sanitize(h.sanitizer)
	}

	respond.Respond(w, r, http.StatusOK, publicContracts)
}
