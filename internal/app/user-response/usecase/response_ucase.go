package responseUcase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients/interfaces"
	user_response "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"log"
	"time"
)

type ResponseUsecase struct {
	responseRep   user_response.Repository
	freelancerClient clients.ClientFreelancer
	managerClient clients.ManagerClient
	jobClient 	  clients.ClientJob
}

func NewResponseUsecase(r user_response.Repository, fClient clients.ClientFreelancer, mclient clients.ManagerClient,
	jClient clients.ClientJob) user_response.Usecase {
	return &ResponseUsecase{
		responseRep:   r,
		freelancerClient: fClient,
		managerClient: mclient,
		jobClient: jClient,
	}
}

func (u *ResponseUsecase) CreateResponse(user *model.User, response *model.Response, jobId int64) error {
	if user.IsManager() {
		return errors.New("to response user need to be freelancer")
	}

	currFreelancer, err := u.freelancerClient.GetFreelancerByUserFromServer(user.ID)
	if err != nil {
		return errors.Wrap(err, "getFreelancerByUserFromServer()")
	}

	// TODO: get files from request
	response.FreelancerId = currFreelancer.ID
	response.JobId = jobId
	response.Date = time.Now()
	response.StatusManager = model.ResponseStatusReview
	response.StatusFreelancer = model.ResponseStatusBlock

	log.Println(response)

	if err := response.Validate(0); err != nil {
		return errors.Wrapf(err, "Validate()")
	}

	if err := u.responseRep.Create(response); err != nil {
		return errors.Wrapf(err, "responseRep.Create()")
	}
	return nil
}

func (u *ResponseUsecase) GetResponses(user *model.User) ([]model.Response, error) {
	var responses []model.Response

	if user.IsManager() {
		currManager, err := u.managerClient.GetManagerByUserFromServer(user.ID)
		if err != nil {
			err = errors.Wrapf(err, "getManagerByUserFromServer()")
			return nil, err
		}

		responses, err = u.responseRep.ListForManager(currManager.ID)
		if err != nil {
			err = errors.Wrapf(err, "responseRep.ListForManager()")
			return nil, err
		}
	} else {
		currFreelancer, err := u.freelancerClient.GetFreelancerByUserFromServer(user.ID)
		if err != nil {
			err = errors.Wrapf(err, "getFreelancerByUserFromServer()")
			return nil, err
		}

		responses, err = u.responseRep.ListForFreelancer(currFreelancer.ID)
		if err != nil {
			err = errors.Wrapf(err, "responseRep.ListForFreelancer()")
			return nil, err
		}
	}

	return responses, nil
}

func (u *ResponseUsecase) AcceptResponse(user *model.User, responseId int64) error {
	response, err := u.responseRep.Find(responseId)
	if err != nil {
		return errors.Wrapf(err, "responseRep.Find()")
	}

	job, err := u.jobClient.GetJobFromServer(response.JobId)
	if err != nil {
		return errors.Wrapf(err, "clients.getJobFromServer()")
	}

	if user.IsManager() {
		currManager, err := u.managerClient.GetManagerByUserFromServer(user.ID)
		if err != nil {
			return errors.Wrapf(err, "clients.getManagerByUserFromServer()")
		}

		if job.HireManagerId != currManager.ID {
			return errors.New("current manager cant accept this response")
		}
		response.StatusManager = model.ResponseStatusAccepted
		response.StatusFreelancer = model.ResponseStatusReview
	} else {
		currFreelancer, err := u.freelancerClient.GetFreelancerByUserFromServer(user.ID)
		if err != nil {
			return errors.Wrapf(err, "clients.getFreelancerByUserFromServer()")
		}

		if currFreelancer.ID != response.FreelancerId {
			return errors.New("current freelancer can't accept this response")
		}

		if response.StatusFreelancer == model.ResponseStatusBlock {
			return errors.New("freelancer can't accept response before manager")
		}
	}
	response.StatusManager = model.ResponseStatusAccepted
	err = u.responseRep.Edit(response)
	if err != nil {
		err = errors.Wrapf(err, "responseRep.Edit()")
		return err
	}
	return nil
}

func (u *ResponseUsecase) DenyResponse(user *model.User, responseId int64) error {
	response, err := u.responseRep.Find(responseId)
	if err != nil {
		return errors.Wrapf(err, "responseRep.Find()")
	}

	job, err := u.jobClient.GetJobFromServer(response.JobId)
	if err != nil {
		return errors.Wrapf(err, "getJobFromServer()")
	}

	if user.IsManager() {
		currManager, err := u.managerClient.GetManagerByUserFromServer(user.ID)
		if err != nil {
			return errors.Wrapf(err, "getManagerByUserFromServer()")
		}

		if job.HireManagerId != currManager.ID {
			return errors.New("current manager cant accept this response")
		}

		response.StatusManager = model.ResponseStatusDenied
		response.StatusFreelancer = model.ResponseStatusBlock
	} else {
		currFreelancer, err := u.freelancerClient.GetFreelancerByUserFromServer(user.ID)
		if err != nil {
			return errors.Wrapf(err, "getFreelancerByUserFromServer()")
		}

		if currFreelancer.ID != response.FreelancerId {
			return errors.New("current freelancer can't accept this response")
		}

		if response.StatusFreelancer == model.ResponseStatusBlock {
			return errors.New("freelancer can't accept response before manager")
		}
	}
	response.StatusManager = model.ResponseStatusDenied
	err = u.responseRep.Edit(response)
	if err != nil {
		err = errors.Wrapf(err, "responseRep.Edit()")
		return err
	}
	return nil
}

func (u *ResponseUsecase) Find(id int64) (*model.Response, error) {
	response, err := u.responseRep.Find(id)
	if err != nil {
		return nil, errors.Wrap(err, "responseRep.Find()")
	}
	return response, nil
}


func (u *ResponseUsecase) GetResponsesOnJobID(jobID int64) ([]model.ExtendResponse, error) {
	responses, err := u.responseRep.ListResponsesOnJobID(jobID)
	if err != nil {
		return nil, errors.Wrap(err, "responseRep.GetResponsesOnJobID()")
	}
	return responses, nil
}