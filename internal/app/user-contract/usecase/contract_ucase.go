package contractUcase

import (
	clients "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients/interfaces"
	user_contract "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"time"
)

type ContractUsecase struct {
	contractRep		user_contract.Repository
	freelancerClient clients.ClientFreelancer
	managerClient   clients.ManagerClient
	companyClient 	clients.CompanyClient
	jobClient   	clients.ClientJob
	responseClient  clients.ClientResponse
}

func NewContractUsecase(c user_contract.Repository, fClient clients.ClientFreelancer, mClient clients.ManagerClient,
	cClient clients.CompanyClient ,jClient clients.ClientJob, rClient clients.ClientResponse) user_contract.Usecase {
	return &ContractUsecase{
		contractRep:	c,
		freelancerClient: fClient,
		managerClient: mClient,
		companyClient: cClient,
		jobClient: jClient,
		responseClient: rClient,
	}
}

func (u *ContractUsecase) CreateContract(user *model.User, responseId int64) error {
	response, err := u.responseClient.GetResponseFromServer(responseId)
	if err != nil {
		return errors.Wrapf(err, "clients.GetResponseFromServer()")
	}

	//job, err := u.jobClient.GetJobFromServer(response.JobId)
	//if err != nil {
	//	return errors.Wrapf(err, "clients.GetJobFromServer()")
	//}

	currManager, err := u.managerClient.GetManagerFromServer(user.HireManagerId)
	if err != nil {
		return errors.Wrapf(err, "clients.GetManagerFromServer()")
	}

	// TODO: write struct for start time and end time
	contract := &model.Contract{
		ID:            0,
		ResponseID:    response.ID,
		CompanyID:     currManager.CompanyId,
		FreelancerID:  user.FreelancerId,
		StartTime:     time.Time{},
		EndTime:       time.Time{},
		Status:        model.ContractStatusUnderDevelopment,
		PaymentAmount: response.PaymentAmount,
		ClientGrade: 0,
		FreelancerGrade: 0,
	}

	response.StatusManager = model.ResponseStatusContractSent
	if err := u.responseClient.UpdateResponseOnServer(response); err != nil {
		return errors.Wrapf(err, "responseCreate.UpdateResponseOnServer()")
	}

	if err := u.contractRep.Create(contract); err != nil {
		return errors.Wrapf(err, "contractRep.Create()")
	}

	return nil
}

func (u * ContractUsecase) SetStatusContract(user * model.User, contract *model.Contract, status string) error {
	// TODO: fix if add new modes
	if !user.IsManager() && status != model.ContractStatusDone {
		return errors.New("freelancer can change status only to done status")
	}
	contract.Status = status
	if err := u.contractRep.Edit(contract); err != nil {
		return errors.Wrapf(err, "contractRep.Edit()")
	}
	return nil
}

func (u * ContractUsecase) SetAsDone(user *model.User, contractId int64) error {
	if user.IsManager() {
		return errors.New("user must be freelancer")
	}

	currFreelancer, err := u.freelancerClient.GetFreelancerByUserFromServer(user.ID)
	if err != nil {
		return errors.Wrapf(err, "clients.GetFreelancerByUserFromServer()")
	}

	contract, err := u.contractRep.Find(contractId)
	if err != nil {
		return errors.Wrapf(err, "contractRep.Find()")
	}

	if contract.FreelancerID != currFreelancer.ID {
		return errors.New("current freelancer can't manage this contract")
	}

	if err := u.SetStatusContract(user, contract, model.ContractStatusDone); err != nil {
		return errors.Wrapf(err, "SetStatusContract()")
	}

	return nil
}

func (u * ContractUsecase) ReviewContract(user *model.User, contractId int64, review *model.ReviewInput) error {
	if review.Grade < model.ContractMinGrade || review.Grade > model.ContractMaxGrade {
		return errors.New("grade must be between 0 and 5")
	}

	contract, err := u.contractRep.Find(contractId)
	if err != nil {
		return errors.Wrapf(err, "contractRep.Find()")
	}

	if user.IsManager() {
		contract.ClientGrade = review.Grade
		contract.ClientComment = review.Comment
	} else {
		contract.FreelancerGrade = review.Grade
		contract.FreelancerComment = review.Comment
	}

	contract.Status = model.ContractStatusReviewed

	//currManager, err := u.managerClient.GetManagerByUserFromServer(user.ID)
	//if err != nil {
	//	return errors.Wrapf(err, "clients.GetManagerByUserFromServer()")
	//}

	/*if contract.CompanyID != currManager.CompanyId {
		return errors.New("current manager cant manage this contract")
	}*/

	if err := u.contractRep.Edit(contract); err != nil {
		return errors.Wrap(err, "contractRep.Edit()")
	}

	return nil
}


func (u *ContractUsecase) ReviewList(user *model.User) ([]model.Review, error) {
	if user.IsManager() {
		return nil, errors.New("user must be freelancer")
	}

	list, err := u.contractRep.List(user.FreelancerId, "freelancer")
	if err != nil {
		return nil, errors.Wrap(err, "contractRep.List()")
	}

	var reviews []model.Review
	for _, contract := range list {

		if contract.ClientGrade == 0 && contract.ClientComment == "" {
			continue
		}

		company, err := u.companyClient.GetCompanyFromServer(contract.CompanyID)

		if err != nil {
			return nil, errors.Wrap(err, "clients.GetCompanyFromServer()")
		}

		response, err := u.responseClient.GetResponseFromServer(contract.ResponseID)
		if err != nil {
			return nil, errors.Wrap(err, "clients.GetResponseFromServer()")
		}

		job, err := u.jobClient.GetJobFromServer(response.JobId)
		if err != nil {
			return nil, errors.Wrap(err, "clients.GetJobFromServer()")
		}

		review := model.Review{
			CompanyName:   company.CompanyName,
			JobTitle:      job.Title,
			ClientGrade:   contract.ClientGrade,
			ClientComment: contract.ClientComment,
		}

		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (u *ContractUsecase) ContractList(user *model.User) ([]model.ContractOutput, error) {
	var userID int64
	if user.IsManager() {
		manager, err := u.managerClient.GetManagerFromServer(user.HireManagerId)
		if err != nil {
			return nil, err
		}
		userID = manager.CompanyId
	} else {
		userID = user.FreelancerId
	}

	list, err := u.contractRep.List(userID, user.UserType)
	if err != nil {
		return nil, errors.Wrap(err, "contractRep.List()")
	}

	var res []model.ContractOutput
	for _, contract := range list {

		response, err := u.responseClient.GetResponseFromServer(contract.ResponseID)
		if err != nil {
			return nil, errors.Wrap(err, "clients.GetResponseFromServer()")
		}

		grpcjob, err := u.jobClient.GetJobFromServer(response.JobId)
		if err != nil {
			return nil, errors.Wrap(err, "clients.GetJobFromServer()")
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

		contractOutput := model.ContractOutput{
			Job:      job,
			Contract: contract,
		}

		res = append(res, contractOutput)
	}
	return res, nil
}