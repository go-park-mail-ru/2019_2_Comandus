package responseHttp

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/respond"
	user_response "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ResponseHandler struct {
	ResponseUsecase	user_response.Usecase
	sanitizer		*bluemonday.Policy
	logger			*zap.SugaredLogger
	sessionStore	sessions.Store
}

func NewResponseHandler(m *mux.Router, rs user_response.Usecase, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &ResponseHandler{
		ResponseUsecase:	rs,
		sanitizer:			sanitizer,
		logger:				logger,
		sessionStore:		sessionStore,
	}

	m.HandleFunc("/jobs/proposal/{id:[0-9]+}", handler.HandleResponseJob).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/proposals", handler.HandleGetResponses).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/proposals/{id:[0-9]+}/accept", handler.HandleResponseAccept).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/proposals/{id:[0-9]+}/deny", handler.HandleResponseDeny).Methods(http.MethodPut, http.MethodOptions)
}

func (h *ResponseHandler) HandleResponseJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseJob<-strconv.Atoi: ")
		respond.Error(w, r, http.StatusBadRequest, err)
	}
	jobId := int64(id)

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleResponseJob: ")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := h.ResponseUsecase.CreateResponse(u, jobId); err != nil {
		err := errors.Wrapf(err,"HandleResponseJob<-ResponseUsecase.CreateResponse: ")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, struct {}{})
}


func (h *ResponseHandler) HandleGetResponses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleGetResponses: ")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	responses, err := h.ResponseUsecase.GetResponses(u)
	if err != nil {
		err := errors.Wrapf(err,"HandleGetResponses<-ResponseUsecase.GetResponses: ")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	for i, _ := range *responses{
		(*responses)[i].Sanitize(h.sanitizer)
	}
	respond.Respond(w, r, http.StatusOK, responses)
}

func (h * ResponseHandler) HandleResponseAccept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept<-strconv.Atoi: ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	responseId := int64(id)

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleResponseAccept: ")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := h.ResponseUsecase.AcceptResponse(u, responseId); err != nil {
		err := errors.Wrapf(err,"HandleResponseAccept<-ResponseUsecase.AcceptResponse: ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h * ResponseHandler) HandleResponseDeny(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept<-strconv.Atoi: ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	responseId := int64(id)

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleResponseAccept: ")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := h.ResponseUsecase.DenyResponse(u, responseId); err != nil {
		err := errors.Wrapf(err,"HandleResponseAccept<-ResponseUsecase.AcceptResponse: ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

