package jobHttp

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general"
	user_job "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
}

func (h *JobHandler) HandleCreateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateJob<-Close: ")
			general.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	job := new(model.Job)
	err := decoder.Decode(job)

	if err != nil {
		err = errors.Wrapf(err, "HandleCreateJob<-Decode: ")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}

	u, ok := r.Context().Value(general.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleCreateJob: ")
		general.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	if err := h.jobUsecase.CreateJob(u, job); err != nil {
		err := errors.Wrap(err, "HandleCreateJob <- JobUseCase.CreateJob(): ")
		general.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	general.Respond(w, r, http.StatusOK, job)
}

func (h *JobHandler) HandleGetJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)

	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-Atoi(wrong id): ")
		general.Error(w, r, http.StatusBadRequest, err)
	}

	job, err := h.jobUsecase.FindJob(int64(id))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-jobUsecase.FindJob(): ")
		general.Error(w, r, http.StatusNotFound, err)
	}

	job.Sanitize(h.sanitizer)
	general.Respond(w, r, http.StatusOK, job)
}

func (h *JobHandler) HandleGetAllJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jobs, err := h.jobUsecase.GetAllJobs()
	if err != nil {
		err = errors.Wrapf(err, "HandleGetAllJobs<-jobUsecase.GetAllJobs: ")
		general.Error(w, r, http.StatusNotFound, err)
	}

	for i, _ := range jobs{
		jobs[i].Sanitize(h.sanitizer)
	}

	general.Respond(w, r, http.StatusOK, &jobs)
}

func (h *JobHandler) HandleUpdateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u, ok := r.Context().Value(general.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleCreateJob: ")
		general.Error(w, r, http.StatusUnauthorized, err)
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
		general.Error(w, r, http.StatusBadRequest, err)
	}

	if err := h.jobUsecase.EditJob(u, inputJob, int64(id)); err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-jobUsecase.EditJob: ")
		general.Error(w, r, http.StatusBadRequest, err)
	}
	general.Respond(w, r, http.StatusOK, struct {}{})
}

/*func (s *server) HandleEditFreelancer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-GetUserFromRequest: ")
		s.error(w, r, codeStatus, err)
		return
	}

	freelancer, err := s.store.Freelancer().FindByUser(user.ID)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-FindByUser: ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditFreelancer<-rBodyClose: ")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}()
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(freelancer)

	if err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-Decode: ")
		s.error(w, r, http.StatusBadRequest, errors.New("invalid format of data"))
		return
	}
	// TODO: validate freelancer

	err = s.store.Freelancer().Edit(freelancer)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-Edit: ")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) HandleGetFreelancer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancer<-Atoi(wrong id): ")
		s.error(w, r, http.StatusBadRequest, err)
	}

	freelancer, err := s.store.Freelancer().Find(int64(id))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancer<-Find: ")
		s.error(w, r, http.StatusNotFound, err)
	}
	freelancer.Sanitize(s.sanitizer)
	s.respond(w, r, http.StatusOK, &freelancer)
}*/
