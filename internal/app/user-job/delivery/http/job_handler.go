package jobHttp

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/respond"
	user_job "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job"
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

type JobHandler struct {
	jobUsecase   user_job.Usecase
	sanitizer    *bluemonday.Policy
	logger       *zap.SugaredLogger
	sessionStore sessions.Store
}

func NewJobHandler(m *mux.Router, js user_job.Usecase, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &JobHandler{
		jobUsecase:   js,
		sanitizer:    sanitizer,
		logger:       logger,
		sessionStore: sessionStore,
	}

	m.HandleFunc("/jobs", handler.HandleCreateJob).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/jobs", handler.HandleGetAllJobs).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/jobs/{id:[0-9]+}", handler.HandleGetJob).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/jobs/{id:[0-9]+}", handler.HandleUpdateJob).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/search/jobs", handler.HandleSearchJob).Methods(http.MethodGet, http.MethodOptions)
}

func (h *JobHandler) HandleCreateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateJob<-Close: ")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	job := new(model.Job)
	err := decoder.Decode(job)

	if err != nil {
		err = errors.Wrapf(err, "HandleCreateJob<-Decode: ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleCreateJob: ")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := h.jobUsecase.CreateJob(u, job); err != nil {
		err := errors.Wrap(err, "HandleCreateJob <- JobUseCase.CreateJob(): ")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, job)
}

func (h *JobHandler) HandleGetJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)

	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-Atoi(wrong id): ")
		respond.Error(w, r, http.StatusBadRequest, err)
	}

	job, err := h.jobUsecase.FindJob(int64(id))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-jobUsecase.FindJob(): ")
		respond.Error(w, r, http.StatusNotFound, err)
	}

	job.Sanitize(h.sanitizer)
	respond.Respond(w, r, http.StatusOK, job)
}

func (h *JobHandler) HandleGetAllJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jobs, err := h.jobUsecase.GetAllJobs()
	if err != nil {
		err = errors.Wrapf(err, "HandleGetAllJobs<-jobUsecase.GetAllJobs: ")
		respond.Error(w, r, http.StatusNotFound, err)
	}

	for i, _ := range jobs{
		jobs[i].Sanitize(h.sanitizer)
	}

	respond.Respond(w, r, http.StatusOK, &jobs)
}

func (h *JobHandler) HandleUpdateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleCreateJob: ")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	inputJob := new(model.Job)
	err := decoder.Decode(inputJob)

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-Atoi(wrong type id): ")
		respond.Error(w, r, http.StatusBadRequest, err)
	}

	if err := h.jobUsecase.EditJob(u, inputJob, int64(id)); err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-jobUsecase.EditJob: ")
		respond.Error(w, r, http.StatusBadRequest, err)
	}
	respond.Respond(w, r, http.StatusOK, struct {}{})
}


func (h *JobHandler) HandleSearchJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pattern, ok := r.URL.Query()["q"]
	if !ok || len(pattern[0]) < 1 {
		err := errors.Wrapf(errors.New("No search pattern"),"HandleSearchJob: ")
		respond.Error(w, r, http.StatusBadRequest, err)
	}

	log.Println(pattern[0])
	jobs, err := h.jobUsecase.PatternSearch(pattern[0])
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-jobUsecase.PatternSearch: ")
		respond.Error(w, r, http.StatusInternalServerError, err)
	}

	respond.Respond(w, r, http.StatusOK, jobs)
}