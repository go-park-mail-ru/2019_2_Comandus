package contractUcase

import (
	server_clients "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients/server-clients"
	user_contract "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"time"
)

type ContractUsecase struct {
	contractRep		user_contract.Repository
	grpcClients		*server_clients.ServerClients
}

func NewContractUsecase(c user_contract.Repository, clients *server_clients.ServerClients) user_contract.Usecase {
	return &ContractUsecase{
		contractRep:	c,
		grpcClients:	clients,
	}
}

func (u *ContractUsecase) CreateContract(user *model.User, responseId int64) error {
	response, err := u.grpcClients.ResponseClient.GetResponseFromServer(responseId)
	if err != nil {
		return errors.Wrapf(err, "clients.GetResponseFromServer()")
	}

	job, err := u.grpcClients.JobClient.GetJobFromServer(response.JobId)
	if err != nil {
		return errors.Wrapf(err, "clients.GetJobFromServer()")
	}

	currManager, err := u.grpcClients.ManagerClient.GetManagerFromServer(job.HireManagerId)
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

	currFreelancer, err := u.grpcClients.FreelancerClient.GetFreelancerByUserFromServer(user.ID)
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

	currManager, err := u.grpcClients.ManagerClient.GetManagerByUserFromServer(user.ID)
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


func (u *ContractUsecase) ReviewList(user *model.User) ([]model.Review, error) {
	if user.IsManager() {
		return nil, errors.New("user must be freelancer")
	}

	list, err := u.contractRep.List(user.ID, "freelancer")
	if err != nil {
		return nil, errors.Wrap(err, "contractRep.List()")
	}

	var reviews []model.Review
	for _, contract := range list {

		if contract.ClientGrade == 0 && contract.ClientComment == "" {
			continue
		}

		company, err := u.grpcClients.CompanyClient.GetCompanyFromServer(contract.CompanyID)
		if err != nil {
			return nil, errors.Wrap(err, "clients.GetCompanyFromServer()")
		}

		response, err := u.grpcClients.ResponseClient.GetResponseFromServer(contract.ResponseID)
		if err != nil {
			return nil, errors.Wrap(err, "clients.GetResponseFromServer()")
		}

		job, err := u.grpcClients.JobClient.GetJobFromServer(response.JobId)
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