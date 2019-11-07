package apiserver

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
)

func (s *server) HandleCreateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateJob<-Close: ")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	job := new(model.Job)
	err := decoder.Decode(job)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateJob<-Decode: ")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	if 	_, err = govalidator.ValidateStruct(job); err != nil {
		err = errors.Wrapf(err, "HandleCreateJob<-ValidateStruct:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateJob<-GetUserFromRequest: ")
		s.error(w, r, codeStatus, err)
		return
	}

	if !user.IsManager() {
		err = errors.New("HandleCreateJob:current user is not a manager : ")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	manager, err := s.store.Manager().FindByUser(user.ID)

	if err != nil {
		log.Println("fail find manager", err)
		err = errors.Wrapf(err, "HandleCreateJob<-FindByUser: ")
		s.error(w, r, http.StatusNotFound, err)
	}

	err = s.store.Job().Create(job, manager)
	if err != nil {
		log.Println("fail create job", err)
		err = errors.Wrapf(err, "HandleCreateJob<-Create: ")
		s.error(w, r, http.StatusInternalServerError, err)
	}

	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) HandleGetJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-Atoi(wrong id): ")
		s.error(w, r, http.StatusBadRequest, err)
	}

	job, err := s.store.Job().Find(int64(id))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-Find: ")
		s.error(w, r, http.StatusNotFound, err)
	}
	job.Sanitize(s.sanitizer)
	s.respond(w, r, http.StatusOK, &job)
}

func (s *server) HandleGetAllJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jobs, err := s.store.Job().List()
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-Find: ")
		s.error(w, r, http.StatusNotFound, err)
	}
	for i, _ := range jobs{
		jobs[i].Sanitize(s.sanitizer)
	}
	s.respond(w, r, http.StatusOK, &jobs)
}

func (s *server) HandleUpdateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	inputJob := new(model.Job)
	err := decoder.Decode(inputJob)

	// Validate Job
	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-Atoi(wrong type id): ")
		s.error(w, r, http.StatusBadRequest, err)
	}

	job, err := s.store.Job().Find(int64(id))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-Find: ")
		s.error(w, r, http.StatusNotFound, err)
	}
	inputJob.ID = job.ID
	inputJob.HireManagerId = job.HireManagerId
	err = s.store.Job().Edit(inputJob)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditProfile<-JobEdit")
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	s.respond(w, r, http.StatusOK, struct {}{})
}

func (s *server) HandleEditFreelancer(w http.ResponseWriter, r *http.Request) {
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

	if 	_, err = govalidator.ValidateStruct(freelancer); err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-ValidateStruct:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

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
}