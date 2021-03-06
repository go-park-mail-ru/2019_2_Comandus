package responseUcase

import (
	clients "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients/interfaces"
	user_response "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/proposal"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type ResponseUsecase struct {
	responseRep      user_response.Repository
	freelancerClient clients.ClientFreelancer
	managerClient    clients.ManagerClient
	jobClient        clients.ClientJob
	chatClient		 clients.ChatClient
}

func NewResponseUsecase(r user_response.Repository, fClient clients.ClientFreelancer, mclient clients.ManagerClient,
	jClient clients.ClientJob, chClient clients.ChatClient) user_response.Usecase {
	return &ResponseUsecase{
		responseRep:      r,
		freelancerClient: fClient,
		managerClient:    mclient,
		jobClient:        jClient,
		chatClient:		  chClient,
	}
}

func (u *ResponseUsecase) CreateResponse(user *model.User, response *model.Response, jobId int64) error {

	currFreelancer, err := u.freelancerClient.GetFreelancerByUserFromServer(user.ID)
	if err != nil {
		return errors.Wrap(err, "getFreelancerByUserFromServer()")
	}

	uIDFromJob, err := u.jobClient.GetUserIDByJobID(jobId)
	if err != nil {
		return errors.Wrap(err, "GetUserIDByJobID()")
	}

	if uIDFromJob == user.ID {
		return errors.New("Вы не можете откликнуться на свой заказ")
	}

	IsResponseYet, err := u.responseRep.CheckForHavingResponse(jobId, currFreelancer.ID)
	if err != nil {
		return errors.Wrap(err, "GetUserIDByJobID()")
	}

	if IsResponseYet {
		return errors.New("Вы уже откликались, дождитесь ответа от заказчика")
	}

	response.FreelancerId = currFreelancer.ID
	response.JobId = jobId
	response.Date = time.Now()
	response.StatusManager = model.ResponseStatusReview
	response.StatusFreelancer = model.ResponseStatusSent

	if err := response.Validate(0); err != nil {
		return errors.Wrapf(err, "Validate()")
	}

	if err := u.responseRep.Create(response); err != nil {
		return errors.Wrapf(err, "responseRep.Create()")
	}
	return nil
}

func (u *ResponseUsecase) GetResponses(user *model.User) ([]model.ExtendResponse, error) {
	var responses []model.ExtendResponse

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

func (u *ResponseUsecase) CancelResponse(user *model.User, responseId int64) error {
	if user.IsManager() {
		return errors.New("user must be freelancer")
	}

	response, err := u.responseRep.Find(responseId)
	if err != nil {
		return errors.Wrapf(err, "responseRep.Find()")
	}

	currFreelancer, err := u.freelancerClient.GetFreelancerByUserFromServer(user.ID)
	if err != nil {
		return errors.Wrapf(err, "clients.getFreelancerByUserFromServer()")
	}

	if response.FreelancerId != currFreelancer.ID {
		return errors.New("no access")
	}

	response.StatusFreelancer = model.ResponseStatusCancel
	err = u.responseRep.Edit(response)
	if err != nil {
		err = errors.Wrapf(err, "responseRep.Edit()")
		return err
	}
	return nil
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

		if response.StatusManager != model.ResponseStatusReview {
			return errors.New("wrong current status")
		}

		chat := &model.Chat{
			Freelancer: response.FreelancerId,
			Manager:    currManager.ID,
			Name:       "Чат №" + strconv.Itoa(int(response.ID)),
			ProposalId: response.ID,
		}

		if err := u.chatClient.CreateChatOnServer(chat); err != nil {
			return errors.Wrap(err, "CreateChatOnServer")
		}

		response.StatusManager = model.ResponseStatusAccepted
	} else {
		currFreelancer, err := u.freelancerClient.GetFreelancerByUserFromServer(user.ID)
		if err != nil {
			return errors.Wrapf(err, "clients.getFreelancerByUserFromServer()")
		}

		if currFreelancer.ID != response.FreelancerId {
			return errors.New("current freelancer can't accept this response")
		}

		if response.StatusManager != model.ResponseStatusContractSent {
			return errors.New("wrong current status")
		}
		response.StatusFreelancer = model.ResponseStatusAccepted
	}

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

		if response.StatusManager != model.ResponseStatusReview {
			return errors.New("wrong current status")
		}
		response.StatusManager = model.ResponseStatusDenied
	} else {
		currFreelancer, err := u.freelancerClient.GetFreelancerByUserFromServer(user.ID)
		if err != nil {
			return errors.Wrapf(err, "getFreelancerByUserFromServer()")
		}

		if currFreelancer.ID != response.FreelancerId {
			return errors.New("current freelancer can't accept this response")
		}

		if response.StatusManager != model.ResponseStatusContractSent {
			return errors.New("wrong current status")
		}
		response.StatusFreelancer = model.ResponseStatusDenied
	}
	err = u.responseRep.Edit(response)
	if err != nil {
		err = errors.Wrapf(err, "responseRep.Edit()")
		return err
	}
	return nil
}

func (u *ResponseUsecase) GetResponse(id int64) (*model.ResponseOutputWithFreel, error) {
	response, err := u.responseRep.Find(id)
	if err != nil {
		return nil, errors.Wrap(err, "responseRep.Find()")
	}

	grpcjob, err := u.jobClient.GetJobFromServer(response.JobId)
	if err != nil {
		return nil, err
	}

	job := model.Job{
		ID:                grpcjob.ID,
		HireManagerId:     grpcjob.HireManagerId,
		Title:             grpcjob.Title,
		Description:       grpcjob.Description,
		Files:             grpcjob.Files,
		SpecialityId:      grpcjob.SpecialityId,
		ExperienceLevelId: grpcjob.ExperienceLevelId,
		PaymentAmount:     grpcjob.PaymentAmount,
		Country:           grpcjob.Country,
		City:              grpcjob.City,
		JobTypeId:         grpcjob.JobTypeId,
		Date:              time.Unix(grpcjob.Date.Seconds, int64(grpcjob.Date.Nanos)),
		Status:            grpcjob.Status,
	}

	grpcFreelancer, err := u.freelancerClient.GetFreelancerFromServer(response.FreelancerId)
	if err != nil {
		return nil, err
	}

	freelancer := &model.Freelancer{
		ID:                grpcFreelancer.Fr.ID,
		AccountId:         grpcFreelancer.Fr.AccountId,
		Country:           grpcFreelancer.Fr.Country,
		City:              grpcFreelancer.Fr.City,
		Address:           grpcFreelancer.Fr.Address,
		Phone:             grpcFreelancer.Fr.Phone,
		TagLine:           grpcFreelancer.Fr.TagLine,
		Overview:          grpcFreelancer.Fr.Overview,
		ExperienceLevelId: grpcFreelancer.Fr.ExperienceLevelId,
		SpecialityId:      grpcFreelancer.Fr.SpecialityId,
		Avatar:            "https://fwork.live/api/account/avatar/" + strconv.FormatInt(grpcFreelancer.Fr.AccountId, 10),
	}
	exFr := model.ExtendFreelancer{
		F:          freelancer,
		FirstName:  grpcFreelancer.FirstName,
		SecondName: grpcFreelancer.SecondName,
	}

	res := new(model.ResponseOutputWithFreel)
	res.Response = *response
	res.Job = job
	res.Freelancer = exFr

	return res, nil
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

func (u *ResponseUsecase) Update(response *model.Response) error {
	if err := u.responseRep.Edit(response); err != nil {
		return err
	}
	return nil
}
