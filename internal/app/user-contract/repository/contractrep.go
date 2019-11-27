package contractRepository

import (
	"database/sql"
	user_contract "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/monitoring"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ContractListByCompany    = "company"
	ContractListByFreelancer = "freelancer"
)

type ContractRepository struct {
	db *sql.DB
}

func NewContractRepository(db *sql.DB) user_contract.Repository {
	return &ContractRepository{db}
}

func (r *ContractRepository) Create(contract *model.Contract) error {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep":"contract", "method":"create"}))
	defer timer.ObserveDuration()

	return r.db.QueryRow(
		"INSERT INTO contracts (responseId, companyId, freelancerId, startTime, endTime, status, "+
			"paymentAmount, clientgrade, freelancergrade, clientcomment, freelancercomment)" +
			" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
		contract.ResponseID,
		contract.CompanyID,
		contract.FreelancerID,
		contract.StartTime,
		contract.EndTime,
		contract.Status,
		contract.PaymentAmount,
		contract.ClientGrade,
		contract.FreelancerGrade,
		contract.ClientComment,
		contract.FreelancerComment,
	).Scan(&contract.ID)
}

func (r *ContractRepository) Find(id int64) (*model.Contract, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep":"contract", "method":"find"}))
	defer timer.ObserveDuration()

	c := &model.Contract{}
	if err := r.db.QueryRow(
		"SELECT id, responseId, companyId, freelancerId, startTime, endTime, status, clientGrade, "+
			"clientComment, freelancerGrade, freelancerComment, paymentAmount FROM contracts WHERE id = $1",
		id,
	).Scan(
		&c.ID,
		&c.ResponseID,
		&c.CompanyID,
		&c.FreelancerID,
		&c.StartTime,
		&c.EndTime,
		&c.Status,
		&c.ClientGrade,
		&c.ClientComment,
		&c.FreelancerGrade,
		&c.FreelancerComment,
		&c.PaymentAmount,
	); err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ContractRepository) Edit(c *model.Contract) error {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep":"contract", "method":"edit"}))
	defer timer.ObserveDuration()
	return r.db.QueryRow("UPDATE contracts SET freelancerId = $1, startTime = $2, "+
		"endTime = $3, status = $4, clientGrade = $5, clientComment = $6, freelancerGrade = $7, " +
		"freelancerComment = $8, paymentAmount = $9 WHERE id = $10 RETURNING id",
		c.FreelancerID,
		c.StartTime,
		c.EndTime,
		c.Status,
		c.ClientGrade,
		c.ClientComment,
		c.FreelancerGrade,
		c.FreelancerComment,
		c.PaymentAmount,
		c.ID,
	).Scan(&c.ID)
}

func (r *ContractRepository) List(id int64, mode string) ([]model.Contract, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep":"contract", "method":"list"}))
	defer timer.ObserveDuration()

	var contracts []model.Contract

	var query string
	if mode == ContractListByCompany {
		query = "SELECT id, responseId, companyId, freelancerId, startTime, endTime, status, clientGrade, "+
			"clientComment, freelancerGrade, freelancerComment, paymentAmount FROM contracts WHERE companyId = $1"
	} else if mode == ContractListByFreelancer {
		query = "SELECT id, responseId, companyId, freelancerId, startTime, endTime, status, clientGrade, "+
			"clientComment, freelancerGrade, freelancerComment, paymentAmount FROM contracts WHERE freelancerId = $1"
	}

	rows, err := r.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := model.Contract{}
		err := rows.Scan(&c.ID, &c.ResponseID, &c.CompanyID, &c.FreelancerID, &c.StartTime, &c.EndTime,
			&c.Status, &c.ClientGrade, &c.ClientComment, &c.FreelancerGrade, &c.FreelancerComment, &c.PaymentAmount)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, c)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return contracts, nil
}
