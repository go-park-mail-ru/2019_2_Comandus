package apiserver

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

func (s *server) HandleResponseJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseJob<-strconv.Atoi: ")
		s.error(w, r, http.StatusBadRequest, err)
	}
	jobId := int64(id)

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseJob<-GetUserFromRequest: ")
		s.error(w, r, codeStatus, err)
		return
	}

	if user.IsManager() {
		err = errors.Wrapf(errors.New("to response user need to be freelancer"),
			"HandleResponseJob<-IsManager: ")
		s.error(w, r, codeStatus, err)
		return
	}

	freelancer, err := s.store.Freelancer().FindByUser(user.ID)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseJob<-Freelancer().FindByUser: ")
		s.error(w, r, codeStatus, err)
		return
	}

	// TODO: get files from request
	response := model.Response{
		ID:               0,
		FreelancerId:     freelancer.ID,
		JobId:            jobId,
		Files:            "",
		Date:             time.Now(),
		StatusManager:    model.ResponseStatusReview,
		StatusFreelancer: model.ResponseStatusBlock,
	}

	if err := response.Validate(0); err != nil {
		err = errors.Wrapf(err, "HandleResponseJob<-Validate: ")
		s.error(w, r, http.StatusBadRequest, err)
	}

	if err := s.store.Response().Create(&response); err != nil {
		err = errors.Wrapf(err, "HandleResponseJob<-Response().Create")
		s.error(w, r, http.StatusInternalServerError, err)
	}

	s.respond(w, r, http.StatusOK, struct {}{})
}

func (s * server) getManagerResponses(userId int64) (*[]model.Response, error){
	manager, err := s.store.Manager().FindByUser(userId)
	if err != nil {
		err = errors.Wrapf(err, " GetManagerResponses<-Manager().FindByUser: ")
		return nil, err
	}

	responses, err := s.store.Response().ListForManager(manager.ID)
	if err != nil {
		err = errors.Wrapf(err, "GetManagerResponses<-Responses().ListForManager: ")
		return nil, err
	}
	return &responses, nil
}

func (s * server) getFreelancerResponses(userId int64) (*[]model.Response, error){
	freelancer, err := s.store.Freelancer().FindByUser(userId)
	if err != nil {
		err = errors.Wrapf(err, " GetManagerResponses<-Manager().FindByUser: ")
		return nil, err
	}

	responses, err := s.store.Response().ListForManager(freelancer.ID)
	if err != nil {
		err = errors.Wrapf(err, "GetManagerResponses<-Responses().ListForManager: ")
		return nil, err
	}
	return &responses, nil
}

func (s *server) HandleGetResponses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetResponses<-GetUserFromRequest: ")
		s.error(w, r, codeStatus, err)
		return
	}

	var responses *[]model.Response
	if user.IsManager() {
		responses, err = s.getManagerResponses(user.ID)
		if err != nil {
			err = errors.Wrapf(err, "HandleGetResponses<-GetManagerResponses: ")
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	} else {
		responses, err = s.getFreelancerResponses(user.ID)
		if err != nil {
			err = errors.Wrapf(err, "HandleGetResponses<-GetFreelancerResponses: ")
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
	for i, _ := range *responses{
		(*responses)[i].Sanitize(s.sanitizer)
	}
	s.respond(w, r, http.StatusOK, responses)
}

func (s * server) HandleResponseAccept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept-strconv.Atoi: ")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}
	responseId := int64(id)

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept<-GetUserFromRequest: ")
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	response, err := s.store.Response().Find(responseId)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept<-Response().Find(): ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	if user.IsManager() {
		manager, err := s.store.Manager().FindByUser(user.ID)
		if err != nil {
			err = errors.Wrapf(err, "HandleResponseAccept<-Manager().FindByUser: ")
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		job, err := s.store.Job().Find(response.JobId)
		if err != nil {
			err = errors.Wrapf(err, "HandleResponseAccept<-Job.Find: ")
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if job.HireManagerId != manager.ID {
			err = errors.New("current manager cant accept this response")
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		response.StatusManager = model.ResponseStatusAccepted
		response.StatusFreelancer = model.ResponseStatusReview
	} else {
		freelancer, err := s.store.Freelancer().FindByUser(user.ID)
		if err != nil {
			err = errors.Wrapf(err, "HandleResponseAccept<-Freelancer().FindByUser: ")
			s.error(w, r, codeStatus, err)
			return
		}

		if freelancer.ID != response.FreelancerId {
			err = errors.New("current freelancer can't accept this response")
			s.error(w, r, codeStatus, err)
			return
		}

		if response.StatusFreelancer == model.ResponseStatusBlock {
			err = errors.New("freelancer can't accept response before manager")
			s.error(w, r, codeStatus, err)
			return
		}

		response.StatusManager = model.ResponseStatusAccepted
	}

	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s * server) HandleResponseDeny(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept-strconv.Atoi: ")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}
	responseId := int64(id)

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept<-GetUserFromRequest: ")
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	response, err := s.store.Response().Find(responseId)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept<-Response().Find(): ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	if user.IsManager() {
		manager, err := s.store.Manager().FindByUser(user.ID)
		if err != nil {
			err = errors.Wrapf(err, "HandleResponseAccept<-Manager().FindByUser: ")
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		job, err := s.store.Job().Find(response.JobId)
		if err != nil {
			err = errors.Wrapf(err, "HandleResponseAccept<-Job.Find: ")
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if job.HireManagerId != manager.ID {
			err = errors.New("current manager cant accept this response")
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		response.StatusManager = model.ResponseStatusDenied
		response.StatusFreelancer = model.ResponseStatusBlock
	} else {
		freelancer, err := s.store.Freelancer().FindByUser(user.ID)
		if err != nil {
			err = errors.Wrapf(err, "HandleResponseAccept<-Freelancer().FindByUser: ")
			s.error(w, r, codeStatus, err)
			return
		}

		if freelancer.ID != response.FreelancerId {
			err = errors.New("current freelancer can't accept this response")
			s.error(w, r, codeStatus, err)
			return
		}

		if response.StatusFreelancer == model.ResponseStatusBlock {
			err = errors.New("freelancer can't accept response before manager")
			s.error(w, r, codeStatus, err)
			return
		}

		response.StatusManager = model.ResponseStatusDenied
	}

	s.respond(w, r, http.StatusOK, struct{}{})
}
