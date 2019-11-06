package apiserver

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func (s * server) HandleCreateContract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateContract:")
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()

	decoder := json.NewDecoder(r.Body)
	contract := new(model.Contract)
	err := decoder.Decode(contract)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateContract:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept-strconv.Atoi: ")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}
	responseId := int64(id)

	response, err := s.store.Response().Find(responseId)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateContract<-Response().Find: ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	job, err := s.store.Job().Find(response.JobId)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateContract<-Job().Find: ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	manager, err := s.store.Manager().Find(job.HireManagerId)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateContract<-Manager().Find: ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	contract.ResponseID = response.ID
	contract.CompanyID = manager.CompanyID
	contract.FreelancerID = response.FreelancerId
	contract.Status = model.ContractStatusUnderDevelopment
	contract.PaymentAmount = response.PaymentAmount
	contract.Grade = 0

	if err := s.store.Contract().Create(contract); err != nil {
		err = errors.Wrapf(err, "HandleCreateContract<-Contract().Create: ")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	// TODO: other responses should be denied
	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s * server) SetStatusContract(user * model.User, contract *model.Contract, status string) error {
	// TODO: fix if add new modes
	if !user.IsManager() && status != model.ContractStatusDone {
		err := errors.New("freelancer can change status only to done status")
		return errors.Wrapf(err, "SetStatusContract<-GetUserFromRequest:")
	}

	contract.Status = status
	if err := s.store.Contract().Edit(contract); err != nil {
		return errors.Wrapf(err, "SetStatusContract<-Contract().Edit:")
	}
	return nil
}

func (s * server) HandleTickContractAsDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user, err, status := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleTickContractAsDone<-GetUserFromRequest: ")
		s.error(w, r, status, err)
		return
	}

	if user.IsManager() {
		err = errors.New("user must be freelancer")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	freelancer, err := s.store.Freelancer().FindByUser(user.ID)
	if err != nil {
		err = errors.Wrapf(err, "HandleTickContractAsDone<-Freelancer().FindByUser: ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleTickContractAsDone<-strconv.Atoi: ")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	contractId := int64(id)
	contract, err := s.store.Contract().Find(contractId)
	if err != nil {
		err = errors.Wrapf(err, "HandleTickContractAsDone<-Contract().Find: ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}


	if contract.FreelancerID != freelancer.ID {
		err = errors.New("current freelancer can't manage this contract")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := s.SetStatusContract(user, contract, model.ContractStatusDone); err != nil {
		err = errors.Wrapf(err, "HandleTickContractAsDone<-SetStatusContract: ")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	s.respond(w,r, http.StatusOK, struct{}{})
}

func (s *server) HandleReviewContract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user, err, status := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleReviewContract<-GetUserFromRequest: ")
		s.error(w, r, status, err)
		return
	}

	if !user.IsManager() {
		err = errors.New("user must be manager")
		s.error(w, r, http.StatusInternalServerError, errors.Wrap(err, "HandleReviewContract: "))
		return
	}

	manager, err := s.store.Manager().FindByUser(user.ID)
	if err != nil {
		err = errors.Wrapf(err, "HandleReviewContract<-Manager().FindByUser: ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleReviewContract: ")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}()

	type Input struct {
		Grade int `json:"grade"`
	}

	decoder := json.NewDecoder(r.Body)
	input := new(Input)
	if err := decoder.Decode(input); err != nil {
		err = errors.Wrapf(err, "HandleReviewContract: ")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	// TODO: max and min grades

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleReviewContract<-strconv.Atoi: ")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	contractId := int64(id)
	contract, err := s.store.Contract().Find(contractId)
	if err != nil {
		err = errors.Wrapf(err, "HandleReviewContract<-Contract().Find: ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	if contract.CompanyID != manager.CompanyID {
		err = errors.New("current manager cant manage this contract")
		s.error(w, r, http.StatusBadRequest, errors.Wrap(err, "HandleReviewContract: "))
	}

	contract.Grade = input.Grade
	contract.Status = model.ContractStatusReviewed
	if err := s.store.Contract().Edit(contract); err != nil {
		s.error(w, r, http.StatusInternalServerError, errors.Wrap(err, "HandleReviewContract<-Contract().Edit: "))
		return
	}

	s.respond(w,r, http.StatusOK, struct{}{})
}

