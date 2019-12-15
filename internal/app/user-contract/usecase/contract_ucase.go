package contractUcase

import (
	clients "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients/interfaces"
	user_contract "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"time"
)

type ContractUsecase struct {
	contractRep      user_contract.Repository
	freelancerClient clients.ClientFreelancer
	managerClient    clients.ManagerClient
	companyClient    clients.CompanyClient
	jobClient        clients.ClientJob
	responseClient   clients.ClientResponse
}

func NewContractUsecase(c user_contract.Repository, fClient clients.ClientFreelancer, mClient clients.ManagerClient,
	cClient clients.CompanyClient, jClient clients.ClientJob, rClient clients.ClientResponse) user_contract.Usecase {
	return &ContractUsecase{
		contractRep:      c,
		freelancerClient: fClient,
		managerClient:    mClient,
		companyClient:    cClient,
		jobClient:        jClient,
		responseClient:   rClient,
	}
}

func (u *ContractUsecase) CreateContract(user *model.User, responseId int64) error {
	response, err := u.responseClient.GetResponseFromServer(responseId)
	if err != nil {
		return errors.Wrapf(err, "clients.GetResponseFromServer()")
	}

	if response.StatusManager != model.ResponseStatusSent {
		return errors.New("Manager isn't accept proposal yet")
	}

	job, err := u.jobClient.GetJobFromServer(response.JobId)
	if err != nil {
		return errors.Wrapf(err, "clients.GetJobFromServer()")
	}

	currManager, err := u.managerClient.GetManagerFromServer(user.HireManagerId)
	if err != nil {
		return errors.Wrapf(err, "clients.GetManagerFromServer()")
	}

	if currManager.ID != job.HireManagerId {
		return errors.New("this manager can't create contract on that proposal")
	}

	// TODO: write struct for start time and end time
	contract := &model.Contract{
		ID:                   0,
		ResponseID:           response.ID,
		CompanyID:            currManager.CompanyId,
		FreelancerID:         response.FreelancerId,
		StartTime:            time.Time{},
		EndTime:              time.Time{},
		Status:               model.ContractStatusExpected,
		StatusFreelancerWork: model.FreelacncerNotReady,
		PaymentAmount:        response.PaymentAmount,
		TimeEstimation:       int(response.TimeEstimation),
		ClientGrade:          0,
		FreelancerGrade:      0,
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

func (u *ContractUsecase) SetAsDone(user *model.User, contractId int64) error {

	contract, err := u.contractRep.Find(contractId)
	if err != nil {
		return errors.Wrapf(err, "contractRep.Find()")
	}

	if contract.CompanyID != user.CompanyId {
		return errors.New("This company can't set done this contract")
	}

	if contract.Status != model.ContractStatusUnderDevelopment {
		return errors.New("incorrect Status, must be active")
	}

	if contract.StatusFreelancerWork != model.FreelancerReady {
		return errors.New("StatusFreelancerWork must be ready")
	}

	err = u.contractRep.ChangeStatus(contractId, model.ContractStatusDone)
	if err != nil {
		return errors.Wrapf(err, "ChangeStatus()")
	}

	return nil
}

func (u *ContractUsecase) ReviewContract(user *model.User, contractId int64, review *model.ReviewInput) error {
	if review.Grade < model.ContractMinGrade || review.Grade > model.ContractMaxGrade {
		return errors.New("grade must be between 1 and 5")
	}

	contract, err := u.contractRep.Find(contractId)
	if err != nil {
		return errors.Wrapf(err, "contractRep.Find()")
	}

	if contract.CompanyID != user.CompanyId && contract.FreelancerID != user.FreelancerId {
		return errors.New("you can't review this contract, you are not freelancer and not company")
	}

	if contract.Status != model.ContractStatusDone {
		return errors.New("contract status must be closed")
	}

	if user.IsManager() && contract.ClientGrade == 0 {
		contract.ClientGrade = review.Grade
		contract.ClientComment = review.Comment

	} else if !user.IsManager() && contract.FreelancerGrade == 0 {
		contract.FreelancerGrade = review.Grade
		contract.FreelancerComment = review.Comment
	} else {
		return errors.New("error , you reviewed this contract yet")
	}

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
	var listType string
	if user.IsManager() {
		manager, err := u.managerClient.GetManagerFromServer(user.HireManagerId)
		if err != nil {
			return nil, err
		}
		userID = manager.CompanyId
		listType = "company"
	} else {
		userID = user.FreelancerId
		listType = "freelancer"
	}

	list, err := u.contractRep.List(userID, listType)
	if err != nil {
		return nil, errors.Wrap(err, "contractRep.List()")
	}

	var res []model.ContractOutput
	for _, contract := range list {

		response, err := u.responseClient.GetResponseFromServer(contract.ResponseID)
		if err != nil {
			return nil, errors.Wrap(err, "clients.GetResponseFromServer()")
		}

		job, err := u.jobClient.GetJobFromServer(response.JobId)
		if err != nil {
			return nil, errors.Wrap(err, "clients.GetJobFromServer()")
		}

		contractOutput := model.ContractOutput{
			Job:      *job,
			Contract: contract,
		}

		res = append(res, contractOutput)
	}
	return res, nil
}

func (u *ContractUsecase) Find(user *model.User, id int64) (*model.ContractOutput, error) {
	contract, err := u.contractRep.Find(id)
	if err != nil {
		return nil, err
	}

	company, err := u.companyClient.GetCompanyFromServer(contract.CompanyID)
	if err != nil {
		return nil, err
	}

	manager, err := u.managerClient.GetManagerByUserFromServer(user.ID)
	if err != nil {
		return nil, err
	}

	if contract.FreelancerID != user.FreelancerId && manager.CompanyId != company.ID {
		return nil, errors.New("no access for current manager")
	}

	freelancer, err := u.freelancerClient.GetFreelancerFromServer(contract.FreelancerID)
	if err != nil {
		return nil, err
	}

	response, err := u.responseClient.GetResponseFromServer(contract.ResponseID)
	if err != nil {
		return nil, err
	}

	job, err := u.jobClient.GetJobFromServer(response.JobId)
	if err != nil {
		return nil, err
	}

	res := &model.ContractOutput{
		Company:    *company,
		Freelancer: *freelancer,
		Job:        *job,
		Contract:   *contract,
	}

	return res, nil
}

func (u *ContractUsecase) ChangeStatus(user *model.User, contractID int64, newStatus string) error {
	contract, err := u.contractRep.Find(contractID)
	if err != nil {
		return err
	}
	if contract.FreelancerID != user.FreelancerId {
		return errors.New("That freelancer can't to accept this contract")
	}
	if contract.Status != model.ContractStatusExpected {
		return errors.New("Incorrect status of contract, must be Expected")
	}

	err = u.contractRep.ChangeStatus(contractID, newStatus)

	return err
}

func (u *ContractUsecase) TickWorkAsReady(user *model.User, contractID int64) error {
	contract, err := u.contractRep.Find(contractID)
	if err != nil {
		return err
	}
	if contract.FreelancerID != user.FreelancerId {
		return errors.New("That freelancer can't to accept this contract")
	}
	if contract.Status != model.ContractStatusUnderDevelopment {
		return errors.New("Incorrect status of contract, must be active")
	}

	err = u.contractRep.ChangeStatusWorkAsReady(contractID)

	return err
}
