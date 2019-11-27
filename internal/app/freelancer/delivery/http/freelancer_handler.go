package freelancerHttp

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/delivery/grpc/freelancer_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/respond"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/delivery/grpc/user_grpc"
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

type FreelancerHandler struct {
	FreelancerUsecase		freelancer.Usecase
	sanitizer				*bluemonday.Policy
	logger					*zap.SugaredLogger
	sessionStore			sessions.Store
}

func NewFreelancerHandler(m *mux.Router, uf freelancer.Usecase, sanitizer *bluemonday.Policy,
	logger *zap.SugaredLogger, sessionStore sessions.Store) {
		handler := &FreelancerHandler{
			FreelancerUsecase:	uf,
			sanitizer:			sanitizer,
			logger:				logger,
			sessionStore:		sessionStore,
		}

	m.HandleFunc("/freelancer", handler.HandleEditFreelancer).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/freelancers/{freelancerId}", handler.HandleGetFreelancer).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/search/freelancers", handler.HandleSearchFreelancers).Methods(http.MethodGet, http.MethodOptions)
}

func (h *FreelancerHandler) HandleEditFreelancer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/freelancer", "method":r.Method}))
	defer timer.ObserveDuration()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleEditFreelancer: ")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditFreelancer<-rBodyClose: ")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	freelancer, err := h.FreelancerUsecase.FindByUser(u.ID)
	if err != nil {
		err = errors.Wrap(err, "FreelancerUsecase.FindByUser()")
		respond.Error(w, r, http.StatusInternalServerError, err)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditPassword<-ioutil.ReadAll()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	currFreelancer := freelancer
	if err := currFreelancer.UnmarshalJSON(body); err != nil {
		err = errors.Wrapf(err, "currCompany.UnmarshalJSON()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.FreelancerUsecase.Edit(freelancer, currFreelancer); err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-FreelancerUsecase.Edit(): ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

type combined struct {
	*freelancer_grpc.Freelancer
	*user_grpc.User
}

func (h *FreelancerHandler) HandleGetFreelancer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/freelancer/id", "method":r.Method}))
	defer timer.ObserveDuration()

	vars := mux.Vars(r)
	ids := vars["freelancerId"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancer<-Atoi(wrong id): ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	currFreelancer, err := clients.GetFreelancerFromServer(int64(id))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancer<-FindFreelancer: ")
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	req := &user_grpc.UserID{
		ID:		currFreelancer.AccountId,
	}

	currUser, err := clients.GetUserFromServer(req)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancer<-FindUser: ")
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	combined  := combined {
		Freelancer: currFreelancer,
		User:       currUser,
	}
	respond.Respond(w, r, http.StatusOK, combined)
}

func (h *FreelancerHandler) HandleSearchFreelancers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"search/freelancer", "method":r.Method}))
	defer timer.ObserveDuration()

	pattern, ok := r.URL.Query()["q"]
	if !ok || len(pattern[0]) < 1 {
		err := errors.Wrapf(errors.New("No search pattern"),"HandleSearchFreelancers: ")
		respond.Error(w, r, http.StatusBadRequest, err)
	}

	log.Println(pattern[0])
	extendedFreelancers, err := h.FreelancerUsecase.PatternSearch(pattern[0])
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-FreelancerUsecase.PatternSearch: ")
		respond.Error(w, r, http.StatusInternalServerError, err)
	}

	respond.Respond(w, r, http.StatusOK, extendedFreelancers)
}