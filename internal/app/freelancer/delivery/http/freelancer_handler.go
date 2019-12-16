package freelancerHttp

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer"
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

type FreelancerHandler struct {
	FreelancerUsecase freelancer.Usecase
	sanitizer         *bluemonday.Policy
	logger            *zap.SugaredLogger
	sessionStore      sessions.Store
}

func NewFreelancerHandler(m *mux.Router, private *mux.Router, uf freelancer.Usecase, sanitizer *bluemonday.Policy,
	logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &FreelancerHandler{
		FreelancerUsecase: uf,
		sanitizer:         sanitizer,
		logger:            logger,
		sessionStore:      sessionStore,
	}

	private.HandleFunc("/freelancer", handler.HandleEditFreelancer).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/freelancers/{pageID}", handler.HandleGetFreelancers).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/freelancer/{freelancerId}", handler.HandleGetFreelancer).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/search/freelancers", handler.HandleSearchFreelancers).Methods(http.MethodGet, http.MethodOptions)
}

func (h *FreelancerHandler) HandleEditFreelancer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/freelancer", "method": r.Method}))
	defer timer.ObserveDuration()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleEditFreelancer()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditFreelancer<-rBodyClose()")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-ioutil.ReadAll()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	currFreelancer := new(model.Freelancer)
	if err := currFreelancer.UnmarshalJSON(body); err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-currFreelancer.UnmarshalJSON()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.FreelancerUsecase.Edit(u.ID, currFreelancer); err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-FreelancerUsecase.Edit()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, struct{}{})
}

type combined struct {
	freelancer *model.FreelancerOutput
	user       *model.User
}

func (h *FreelancerHandler) HandleGetFreelancer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/freelancer/id", "method": r.Method}))
	defer timer.ObserveDuration()

	vars := mux.Vars(r)
	ids := vars["freelancerId"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancer<-Atoi(wrong id)")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	currFreelancer, err := h.FreelancerUsecase.Find(int64(id))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancer<-clients.GetFreelancerFromServer()")
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	currUser, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no currUser in context"), "HandleEditProfile()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	combined := combined{
		freelancer: currFreelancer.OuFreel,
		user:       currUser,
	}
	respond.Respond(w, r, http.StatusOK, combined)
}

func (h *FreelancerHandler) HandleSearchFreelancers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "search/freelancer", "method": r.Method}))
	defer timer.ObserveDuration()


	var err error
	params := new(model.SearchParams)
	pattern := r.URL.Query().Get("q")
	minGrade := r.URL.Query().Get("minGrade")
	if minGrade != "" {
		params.MinGrade, err = strconv.ParseInt(minGrade, 10, 64)
		if err != nil {
			respond.Error(w, r, http.StatusBadRequest, errors.Wrap(err, "HandleSearchJob()"))
		}
	}

	maxGrade := r.URL.Query().Get("maxGrade")
	if maxGrade != "" {
		params.MaxGrade, err = strconv.ParseInt(maxGrade, 10, 64)
		if err != nil {
			respond.Error(w, r, http.StatusBadRequest, errors.Wrap(err, "HandleSearchJob()"))
		}
	}

	minPayment := r.URL.Query().Get("minPaymentAmount")
	if minPayment != "" {
		params.MinPaymentAmount, err = strconv.ParseFloat(minPayment, 64)
		if err != nil {
			respond.Error(w, r, http.StatusBadRequest, errors.Wrap(err, "HandleSearchJob()"))
		}
	}

	maxPayment := r.URL.Query().Get("maxPaymentAmount")
	if maxPayment != "" {
		params.MaxPaymentAmount, err = strconv.ParseFloat(maxPayment, 64)
		if err != nil {
			respond.Error(w, r, http.StatusBadRequest, errors.Wrap(err, "HandleSearchJob()"))
		}
	}

	country := r.URL.Query().Get("country")
	if country != "" {
		params.Country, err = strconv.ParseInt(country, 10, 64)
		if err != nil {
			respond.Error(w, r, http.StatusBadRequest, errors.Wrap(err, "HandleSearchJob()"))
		}
	} else {
		params.Country = -1
	}

	city := r.URL.Query().Get("city")
	if city != "" {
		params.City, err = strconv.ParseInt(city, 10, 64)
		if err != nil {
			respond.Error(w, r, http.StatusBadRequest, errors.Wrap(err, "HandleSearchJob()"))
		}
	} else {
		params.City = -1
	}

	expLevel := r.URL.Query().Get("experienceLevel")
	if expLevel != "" {
		levels, err := strconv.ParseInt(expLevel, 10, 64)
		if err != nil {
			respond.Error(w, r, http.StatusBadRequest, errors.Wrap(err, "HandleSearchJob()"))
		}

		if levels / 100 != 0 {
			params.ExperienceLevel[0] = true
		}

		if (levels % 100) / 10 != 0 {
			params.ExperienceLevel[1] = true
		}

		if (levels % 100) % 10 != 0 {
			params.ExperienceLevel[2] = true
		}
	} else {
		params.ExperienceLevel[0] = true
		params.ExperienceLevel[1] = true
		params.ExperienceLevel[2] = true
	}

	desc := r.URL.Query().Get("desc")
	if desc != "" {
		params.Desc, err = strconv.ParseBool(desc)
		if err != nil {
			respond.Error(w, r, http.StatusBadRequest, errors.Wrap(err, "HandleSearchJob()"))
		}
	}

	limit := r.URL.Query().Get("limit")
	if limit != "" {
		params.Limit, err = strconv.ParseInt(limit, 10, 64)
		if err != nil {
			respond.Error(w, r, http.StatusBadRequest, errors.Wrap(err, "HandleSearchJob()"))
		}
	}

	extendedFreelancers, err := h.FreelancerUsecase.PatternSearch(pattern, *params)
	if err != nil {
		err = errors.Wrapf(err, "HandleSearchFreelancers<-Ucase.PatternSearch()")
		respond.Error(w, r, http.StatusInternalServerError, err)
	}

	respond.Respond(w, r, http.StatusOK, extendedFreelancers)
}

func (h *FreelancerHandler) HandleGetFreelancers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/freelancers/{pageID}", "method": r.Method}))
	defer timer.ObserveDuration()

	vars := mux.Vars(r)
	pageIDIn := vars["pageID"]
	pageID, err := strconv.Atoi(pageIDIn)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancers<-Atoi(wrong id)")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	// TODO get limit from request / now default limit = 20
	var limit = 20
	offset := 20 * (pageID - 1)

	extendedFreelancers, err := h.FreelancerUsecase.FindPart(offset, limit)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancers<-Ucase.FindPart()")
		respond.Error(w, r, http.StatusInternalServerError, err)
	}

	for i , _ := range extendedFreelancers {
		extendedFreelancers[i].Sanitize(h.sanitizer)
	}

	respond.Respond(w, r, http.StatusOK, extendedFreelancers)
}
