package responseHttp

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/respond"
	user_response "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/proposal"
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

type ResponseHandler struct {
	ResponseUsecase user_response.Usecase
	sanitizer       *bluemonday.Policy
	logger          *zap.SugaredLogger
	sessionStore    sessions.Store
}

func NewResponseHandler(m *mux.Router, rs user_response.Usecase, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &ResponseHandler{
		ResponseUsecase: rs,
		sanitizer:       sanitizer,
		logger:          logger,
		sessionStore:    sessionStore,
	}

	m.HandleFunc("/jobs/proposal/{id:[0-9]+}", handler.HandleResponseJob).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/proposals", handler.HandleGetResponses).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/proposals/{id:[0-9]+}", handler.HandleGetResponse).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/proposals/{id:[0-9]+}/accept", handler.HandleResponseAccept).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/proposals/{id:[0-9]+}/deny", handler.HandleResponseDeny).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/proposals/{id:[0-9]+}/cancel", handler.HandleResponseCancel).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/job/{jobid:[0-9]+}/proposals", handler.HandleGetResponsesOnJobID).Methods(http.MethodGet, http.MethodOptions)
}

func (h *ResponseHandler) HandleResponseJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/jobs/proposal/id", "method": r.Method}))
	defer timer.ObserveDuration()

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseJob<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
	}

	jobId := int64(id)

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleResponseJob<-Close()")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseJob<-ioutil.ReadAll()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	response := new(model.Response)
	if err := response.UnmarshalJSON(body); err != nil {
		err = errors.Wrapf(err, "currCompany.UnmarshalJSON()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	log.Println(response)

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleResponseJob()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := h.ResponseUsecase.CreateResponse(u, response, jobId); err != nil {
		err := errors.Wrapf(err, "HandleResponseJob<-UCase.CreateResponse()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *ResponseHandler) HandleGetResponses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/proposals", "method": r.Method}))
	defer timer.ObserveDuration()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleGetResponses()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	responses, err := h.ResponseUsecase.GetResponses(u)
	if err != nil {
		err := errors.Wrapf(err, "HandleGetResponses<-ResponseUsecase.GetResponses()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	for i, _ := range responses {
		(responses)[i].R.Sanitize(h.sanitizer)
	}
	respond.Respond(w, r, http.StatusOK, responses)
}

func (h *ResponseHandler) HandleResponseCancel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/proposals/id/cancel", "method": r.Method}))
	defer timer.ObserveDuration()

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	responseId := int64(id)

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleResponseAccept()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := h.ResponseUsecase.CancelResponse(u, responseId); err != nil {
		err := errors.Wrapf(err, "HandleResponseCancel<-ResponseUsecase.CancelResponse()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *ResponseHandler) HandleResponseAccept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/proposals/id/accept", "method": r.Method}))
	defer timer.ObserveDuration()

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	responseId := int64(id)

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleResponseAccept()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := h.ResponseUsecase.AcceptResponse(u, responseId); err != nil {
		err := errors.Wrapf(err, "HandleResponseAccept<-ResponseUsecase.AcceptResponse()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *ResponseHandler) HandleResponseDeny(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/proposals/id/deny", "method": r.Method}))
	defer timer.ObserveDuration()

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseDeny<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	responseId := int64(id)

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleResponseDeny()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := h.ResponseUsecase.DenyResponse(u, responseId); err != nil {
		err := errors.Wrapf(err, "HandleResponseDeny<-ResponseUsecase.DenyResponse()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *ResponseHandler) HandleGetResponsesOnJobID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/job/id/proposals", "method": r.Method}))
	defer timer.ObserveDuration()

	vars := mux.Vars(r)
	ids := vars["jobid"]
	jobid, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetResponsesOnJobID<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	exResp, err := h.ResponseUsecase.GetResponsesOnJobID(int64(jobid))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetResponsesOnJobID<-GetResponsesOnJobID()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	respond.Respond(w, r, http.StatusOK, exResp)
}

func (h *ResponseHandler) HandleGetResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetResponse<-strconv.Atoi()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	response, err := h.ResponseUsecase.GetResponse(id)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetResponse<-Usecase.GetResponse()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, response)
}
