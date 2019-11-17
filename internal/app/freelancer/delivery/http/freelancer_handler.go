package freelancerHttp

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general"
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

type FreelancerHandler struct {
	FreelancerUsecase		freelancer.Usecase
	sanitizer				*bluemonday.Policy
	logger					*zap.SugaredLogger
	sessionStore			sessions.Store
}

func NewFreelancerHandler(m *mux.Router, uf freelancer.Usecase, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &FreelancerHandler{
		FreelancerUsecase:	uf,
		sanitizer:			sanitizer,
		logger:				logger,
		sessionStore:		sessionStore,
	}

	m.HandleFunc("/freelancer", handler.HandleEditFreelancer).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/freelancer/{freelancerId:[0-9]+}", handler.HandleGetFreelancer).Methods(http.MethodGet, http.MethodOptions)
}

func (h *FreelancerHandler) HandleEditFreelancer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u, ok := r.Context().Value(general.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleEditFreelancer: ")
		general.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditFreelancer<-rBodyClose: ")
			general.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	freelancer, err := h.FreelancerUsecase.FindByUser(u)
	if err != nil {
		err = errors.Wrap(err, "FreelancerUsecase.FindByUser()")
		general.Error(w, r, http.StatusInternalServerError, err)
	}

	currFreelancer := freelancer
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(currFreelancer)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-Decode(): ")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.FreelancerUsecase.Edit(freelancer, currFreelancer); err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-FreelancerUsecase.Edit(): ")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	general.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *FreelancerHandler) HandleGetFreelancer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["freelancerId"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancer<-Atoi(wrong id): ")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	currFreelancer, err := h.FreelancerUsecase.Find(int64(id))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancer<-Find: ")
		general.Error(w, r, http.StatusNotFound, err)
		return
	}

	currFreelancer.Sanitize(h.sanitizer)
	general.Respond(w, r, http.StatusOK, currFreelancer)
}
