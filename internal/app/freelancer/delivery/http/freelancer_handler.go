package freelancerHttp

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/respond"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Message string `json:"message"`
}

type FreelancerHandler struct {
	FreelancerUsecase		freelancer.Usecase
	UserUsecase				user.Usecase
	sanitizer				*bluemonday.Policy
	logger					*zap.SugaredLogger
	sessionStore			sessions.Store
}

func NewFreelancerHandler(m *mux.Router, uf freelancer.Usecase, uc user.Usecase, sanitizer *bluemonday.Policy,
	logger *zap.SugaredLogger, sessionStore sessions.Store) {
		handler := &FreelancerHandler{
			FreelancerUsecase:	uf,
			UserUsecase:		uc,
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

	currFreelancer := freelancer
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(currFreelancer)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-Decode(): ")
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
	*model.Freelancer
	*model.User
}

func (h *FreelancerHandler) HandleGetFreelancer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["freelancerId"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancer<-Atoi(wrong id): ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	currFreelancer, err := h.FreelancerUsecase.Find(int64(id))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancer<-FindFreelancer: ")
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	currUser, err := h.UserUsecase.Find(currFreelancer.AccountId)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancer<-FindUser: ")
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	currFreelancer.Sanitize(h.sanitizer)
	currUser.Sanitize(h.sanitizer)

	combined  := combined {
		Freelancer: currFreelancer,
		User:       currUser,
	}
	respond.Respond(w, r, http.StatusOK, combined)
}

func (h *FreelancerHandler) HandleSearchFreelancers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pattern, ok := r.URL.Query()["q"]
	if !ok || len(pattern[0]) < 1 {
		err := errors.Wrapf(errors.New("No search pattern"),"HandleSearchFreelancers: ")
		respond.Error(w, r, http.StatusBadRequest, err)
	}

	log.Println(pattern[0])
	freelancers, err := h.FreelancerUsecase.PatternSearch(pattern[0])
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-FreelancerUsecase.PatternSearch: ")
		respond.Error(w, r, http.StatusInternalServerError, err)
	}

	respond.Respond(w, r, http.StatusOK, freelancers)
}