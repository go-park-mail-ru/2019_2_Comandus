package jobHttp

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/respond"
	user_job "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job"
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

type JobHandler struct {
	jobUsecase   user_job.Usecase
	sanitizer    *bluemonday.Policy
	logger       *zap.SugaredLogger
	sessionStore sessions.Store
}

func NewJobHandler(m *mux.Router, private *mux.Router, js user_job.Usecase, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &JobHandler{
		jobUsecase:   js,
		sanitizer:    sanitizer,
		logger:       logger,
		sessionStore: sessionStore,
	}

	private.HandleFunc("/jobs", handler.HandleCreateJob).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/jobs", handler.HandleGetAllJobs).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/jobs/{id:[0-9]+}", handler.HandleGetJob).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/jobs/{id:[0-9]+}", handler.HandleUpdateJob).Methods(http.MethodPut, http.MethodOptions)
	private.HandleFunc("/jobs/{id:[0-9]+}", handler.HandleDeleteJob).Methods(http.MethodDelete, http.MethodOptions)
	m.HandleFunc("/search/jobs", handler.HandleSearchJob).Methods(http.MethodGet, http.MethodOptions)
}

func (h *JobHandler) HandleCreateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/jobs", "method": r.Method}))
	defer timer.ObserveDuration()

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateJob<-Close()")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateJob<-ioutil.ReadAll()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	job := new(model.Job)
	if err := job.UnmarshalJSON(body); err != nil {
		err = errors.Wrapf(err, "currCompany.UnmarshalJSON()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleCreateJob()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := h.jobUsecase.CreateJob(u, job); err != nil {
		err := errors.Wrap(err, "HandleCreateJob<-JobUseCase.CreateJob()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, job)
}

func (h *JobHandler) HandleGetJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/jobs", "method": r.Method}))
	defer timer.ObserveDuration()

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)

	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-Atoi(wrong id)")
		respond.Error(w, r, http.StatusBadRequest, err)
	}

	job, err := h.jobUsecase.FindJob(int64(id))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-jobUsecase.FindJob()")
		respond.Error(w, r, http.StatusNotFound, err)
	}

	job.Sanitize(h.sanitizer)
	respond.Respond(w, r, http.StatusOK, job)
}

func (h *JobHandler) HandleGetAllJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var jobs []model.Job
	var err error
	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/jobs/id", "method": r.Method}))
	defer timer.ObserveDuration()

	pattern, ok := r.URL.Query()["manid"]
	if !ok || len(pattern[0]) < 1 {

		jobs, err = h.jobUsecase.GetAllJobs()
		if err != nil {
			err = errors.Wrapf(err, "HandleGetAllJobs<-jobUsecase.GetAllJobs()")
			respond.Error(w, r, http.StatusNotFound, err)
			return
		}
	} else {
		manID, err := strconv.Atoi(pattern[0])
		if err != nil {
			err = errors.Wrapf(err, "HandleGetAllJobs<-")
			respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		jobs, err = h.jobUsecase.GetMyJobs(int64(manID))
		if err != nil {
			err = errors.Wrapf(err, "HandleGetAllJobs<-")
			respond.Error(w, r, http.StatusNotFound, err)
			return

		}
	}
	for i, _ := range jobs {
		jobs[i].Sanitize(h.sanitizer)
	}

	respond.Respond(w, r, http.StatusOK, &jobs)

}

func (h *JobHandler) HandleDeleteJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/jobs/id", "method": r.Method}))
	defer timer.ObserveDuration()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleDeleteJob()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleDeleteJob<-Atoi(wrong type id)")
		respond.Error(w, r, http.StatusBadRequest, err)
	}

	if err := h.jobUsecase.MarkAsDeleted(int64(id), u); err != nil {
		err = errors.Wrapf(err, "HandleDeleteJob<-jobUsecase.MarkAsDeleted()")
		respond.Error(w, r, http.StatusBadRequest, err)
	}
	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *JobHandler) HandleUpdateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/jobs/id", "method": r.Method}))
	defer timer.ObserveDuration()

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"), "HandleUpdateJob()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleUpdateJob<-Close()")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrapf(err, "HandleUpdateJob<-ioutil.ReadAll()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	inputJob := new(model.Job)
	if err := inputJob.UnmarshalJSON(body); err != nil {
		err = errors.Wrapf(err, "currCompany.UnmarshalJSON()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleUpdateJob<-Atoi(wrong type id)()")
		respond.Error(w, r, http.StatusBadRequest, err)
	}

	if err := h.jobUsecase.EditJob(u, inputJob, int64(id)); err != nil {
		err = errors.Wrapf(err, "HandleUpdateJob<-jobUsecase.EditJob()")
		respond.Error(w, r, http.StatusBadRequest, err)
	}
	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *JobHandler) HandleSearchJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path": "/search/jobs", "method": r.Method}))
	defer timer.ObserveDuration()

	pattern := r.URL.Query().Get("q")
	if pattern == "" {
		log.Println("HERE")
		err := errors.Wrapf(errors.New("No search pattern"), "HandleSearchJob()")
		respond.Error(w, r, http.StatusBadRequest, err)
	}


	var err error
	params := new(model.JobSearchParams)
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

	proposals := r.URL.Query().Get("proposals")
	if proposals != "" {
		params.Proposals, err = strconv.ParseInt(proposals, 10, 64)
		if err != nil {
			respond.Error(w, r, http.StatusBadRequest, errors.Wrap(err, "HandleSearchJob()"))
		}
	}

	expLevel := r.URL.Query().Get("proposals")
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

	jobs, err := h.jobUsecase.PatternSearch(pattern, *params)
	if err != nil {
		err = errors.Wrapf(err, "HandleSearchJob<-jobUsecase.PatternSearch()")
		respond.Error(w, r, http.StatusInternalServerError, err)
	}

	respond.Respond(w, r, http.StatusOK, jobs)
}
