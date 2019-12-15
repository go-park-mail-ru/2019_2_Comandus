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
		Labels{"rep": "contract", "method": "create"}))
	defer timer.ObserveDuration()

	return r.db.QueryRow(
		"INSERT INTO contracts (responseId, companyId, freelancerId, startTime, endTime, status, statusFreelancerWork, "+
			"paymentAmount, clientgrade, freelancergrade, clientcomment, freelancercomment, timeEstimation)"+
			" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id",
		contract.ResponseID,
		contract.CompanyID,
		contract.FreelancerID,
		contract.StartTime,
		contract.EndTime,
		contract.Status,
		contract.StatusFreelancerWork,
		contract.PaymentAmount,
		contract.ClientGrade,
		contract.FreelancerGrade,
		contract.ClientComment,
		contract.FreelancerComment,
		contract.TimeEstimation,
	).Scan(&contract.ID)
}

func (r *ContractRepository) Find(id int64) (*model.Contract, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "contract", "method": "find"}))
	defer timer.ObserveDuration()

	c := &model.Contract{}
	if err := r.db.QueryRow(
		"SELECT id, responseId, companyId, freelancerId, startTime, endTime, status, statusFreelancerWork, " +
			"clientGrade, clientComment, freelancerGrade, freelancerComment, paymentAmount, timeestimation " +
			"FROM contracts WHERE id = $1",
		id,
	).Scan(
		&c.ID,
		&c.ResponseID,
		&c.CompanyID,
		&c.FreelancerID,
		&c.StartTime,
		&c.EndTime,
		&c.Status,
		&c.StatusFreelancerWork,
		&c.ClientGrade,
		&c.ClientComment,
		&c.FreelancerGrade,
		&c.FreelancerComment,
		&c.PaymentAmount,
		&c.TimeEstimation,
	); err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ContractRepository) Edit(c *model.Contract) error {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "contract", "method": "edit"}))
	defer timer.ObserveDuration()

	return r.db.QueryRow("UPDATE contracts SET freelancerId = $1, startTime = $2, "+
		"endTime = $3, clientGrade = $4, clientComment = $5, freelancerGrade = $6, "+
		"freelancerComment = $7, paymentAmount = $8 WHERE id = $9 RETURNING id",
		c.FreelancerID,
		c.StartTime,
		c.EndTime,
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
		Labels{"rep": "contract", "method": "list"}))
	defer timer.ObserveDuration()

	var contracts []model.Contract

	var rows *sql.Rows
	var err error

	if mode == ContractListByCompany {
		rows, err = r.db.Query("SELECT id, responseId, companyId, freelancerId, startTime, endTime, status, " +
			" statusFreelancerWork, clientGrade, clientComment, freelancerGrade, freelancerComment, paymentAmount," +
			"timeestimation " +
			"FROM contracts WHERE companyId = $1", id)
	} else if mode == ContractListByFreelancer {
		rows, err = r.db.Query("SELECT id, responseId, companyId, freelancerId, startTime, endTime, status," +
			" statusFreelancerWork, clientGrade, clientComment, freelancerGrade, freelancerComment, paymentAmount, " +
			"timeestimation " +
			"FROM contracts WHERE freelancerId = $1", id)
	}
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := model.Contract{}
		err := rows.Scan(&c.ID, &c.ResponseID, &c.CompanyID, &c.FreelancerID, &c.StartTime, &c.EndTime,
			&c.Status, &c.StatusFreelancerWork, &c.ClientGrade, &c.ClientComment, &c.FreelancerGrade,
			&c.FreelancerComment, &c.PaymentAmount, &c.TimeEstimation)
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

func (r *ContractRepository) ChangeStatus(contractID int64, newStatus string) error {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "contract", "method": "ChangeStatus"}))
	defer timer.ObserveDuration()

	_, err := r.db.Exec("UPDATE contracts"+
		" SET status = $1 "+
		" WHERE id = $2",
		newStatus,
		contractID,
	)

	if err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) ChangeStatusWorkAsReady(contractID int64) error {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
	Labels{"rep": "contract", "method": "ChangeStatusWorkAsReady"}))
	defer timer.ObserveDuration()

	_, err := r.db.Exec("UPDATE contracts"+
		" SET statusFreelancerWork = $1 "+
		" WHERE id = $2",
		model.FreelancerReady,
		contractID,
	)

	if err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) GetClosedContracts(id int64) ([]model.PublicContractVersion, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
	Labels{"rep": "contract", "method": "GetClosedContracts"}))
	defer timer.ObserveDuration()

	var contracts []model.PublicContractVersion

	var rows *sql.Rows
	var err error

	rows, err = r.db.Query("SELECT con.clientGrade , con.clientComment, con.freelancerGrade, " +
		" con.freelancerComment, con.status, u.firstName , u.SecondName, cp.companyName, j.Title " +
		"FROM contracts AS con " +
		"INNER JOIN freelancers AS f ON (f.id = con.freelancerid) " +
		"INNER JOIN users AS u ON(u.accountid = f.accountid) " +
		"INNER JOIN companies AS cp ON(cp.id = con.companyid) " +
		"INNER JOIN responses AS r ON(r.id = con.responseid) " +
		"INNER JOIN jobs AS j ON (j.id = r.jobID) " +
		"WHERE con.freelancerId = $1 AND con.status = 'closed'",
		id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := model.PublicContractVersion{}
		err := rows.Scan(&c.ClientGrade, &c.ClientComment, &c.FreelancerGrade, &c.FreelancerComment, &c.Status,
			&c.FirstName, &c.SecondName, &c.CompanyName, &c.JobTitle)
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