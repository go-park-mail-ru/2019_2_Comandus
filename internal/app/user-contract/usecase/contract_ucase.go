package contractUcase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager"
	user_contract "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract"
	user_job "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job"
	user_response "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"time"
)

type ContractUsecase struct {
	managerRep		manager.Repository
	freelancerRep	freelancer.Repository
	jobRep			user_job.Repository
	responseRep		user_response.Repository
	contractRep		user_contract.Repository
}

func NewContractUsecase(m manager.Repository,
	f freelancer.Repository,
	j user_job.Repository,
	r user_response.Repository,
	c user_contract.Repository) user_contract.Usecase {
	return &ContractUsecase{
		managerRep:		m,
		freelancerRep:	f,
		jobRep:			j,
		responseRep:	r,
		contractRep:	c,
	}
}

func (u *ContractUsecase) CreateContract(user *model.User, responseId int64) error {
	response, err := u.responseRep.Find(responseId)
	if err != nil {
		return errors.Wrapf(err, "responseRep.Find(): ")
	}

	job, err := u.jobRep.Find(response.JobId)
	if err != nil {
		return errors.Wrapf(err, "jobRep.Find(): ")
	}

	currManager, err := u.managerRep.Find(job.HireManagerId)
	if err != nil {
		return errors.Wrapf(err, "managerRep.Find(): ")
	}

	// TODO: write struct for start time and end time
	contract := &model.Contract{
		ID:            0,
		ResponseID:    response.ID,
		CompanyID:     currManager.CompanyID,
		FreelancerID:  response.FreelancerId,
		StartTime:     time.Time{},
		EndTime:       time.Time{},
		Status:        model.ContractStatusUnderDevelopment,
		Grade:         0,
		PaymentAmount: response.PaymentAmount,
	}

	if err := u.contractRep.Create(contract); err != nil {
		return errors.Wrapf(err, "contractRep.Create(): ")
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
		return errors.Wrapf(err, "contractRep.Edit(): ")
	}
	return nil
}

func (u * ContractUsecase) SetAsDone(user *model.User, contractId int64) error {
	if user.IsManager() {
		return errors.New("user must be freelancer")
	}

	currFreelancer, err := u.freelancerRep.FindByUser(user.ID)
	if err != nil {
		return errors.Wrapf(err, "freelancerRep.FindByUser(): ")
	}

	contract, err := u.contractRep.Find(contractId)
	if err != nil {
		return errors.Wrapf(err, "contractRep.Find(): ")
	}

	if contract.FreelancerID != currFreelancer.ID {
		return errors.New("current freelancer can't manage this contract")
	}

	if err := u.SetStatusContract(user, contract, model.ContractStatusDone); err != nil {
		return errors.Wrapf(err, "SetStatusContract(): ")
	}

	return nil
}

func (u * ContractUsecase) ReviewContract(user *model.User, contractId int64, grade int) error {
	if !user.IsManager() {
		return errors.New("user must be manager")
	}

	if grade < model.ContractMinGrade || grade > model.ContractMaxGrade {
		return errors.New("grade must be between 0 and 5")
	}

	contract, err := u.contractRep.Find(contractId)
	if err != nil {
		return errors.Wrapf(err, "contractRep.Find(): ")
	}

	contract.Grade = grade
	contract.Status = model.ContractStatusReviewed

	currManager, err := u.managerRep.FindByUser(user.ID)
	if err != nil {
		return errors.Wrapf(err, "managerRep.FindByUser(): ")
	}

	if contract.CompanyID != currManager.CompanyID {
		return errors.New("current manager cant manage this contract")
	}

	if err := u.contractRep.Edit(contract); err != nil {
		return errors.Wrap(err, "contractRep.Edit(): ")
	}

	return nil
}