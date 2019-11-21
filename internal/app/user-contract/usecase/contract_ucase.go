package contractUcase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients"
	user_contract "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"time"
)

type ContractUsecase struct {
	contractRep		user_contract.Repository
}

func NewContractUsecase(c user_contract.Repository) user_contract.Usecase {
	return &ContractUsecase{
		contractRep:	c,
	}
}

func (u *ContractUsecase) CreateContract(user *model.User, responseId int64) error {
	response, err := clients.GetResponseFromServer(responseId)
	if err != nil {
		return errors.Wrapf(err, "clients.GetResponseFromServer()")
	}

	job, err := clients.GetJobFromServer(response.JobId)
	if err != nil {
		return errors.Wrapf(err, "clients.GetJobFromServer()")
	}

	currManager, err := clients.GetManagerFromServer(job.HireManagerId)
	if err != nil {
		return errors.Wrapf(err, "clients.GetManagerFromServer()")
	}

	// TODO: write struct for start time and end time
	contract := &model.Contract{
		ID:            0,
		ResponseID:    response.ID,
		CompanyID:     currManager.CompanyId,
		FreelancerID:  response.FreelancerId,
		StartTime:     time.Time{},
		EndTime:       time.Time{},
		Status:        model.ContractStatusUnderDevelopment,
		Grade:         0,
		PaymentAmount: response.PaymentAmount,
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

	currFreelancer, err := clients.GetFreelancerByUserFromServer(user.ID)
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

	currManager, err := clients.GetManagerByUserFromServer(user.ID)
	if err != nil {
		return errors.Wrapf(err, "clients.GetManagerByUserFromServer()")
	}

	if contract.CompanyID != currManager.CompanyId {
		return errors.New("current manager cant manage this contract")
	}

	if err := u.contractRep.Edit(contract); err != nil {
		return errors.Wrap(err, "contractRep.Edit()")
	}

	return nil
}