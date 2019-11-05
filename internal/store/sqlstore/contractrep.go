package sqlstore

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

const (
	ContractListByCompany = "company"
	ContractListByFreelancer = "freelancer"
	)

type ContractRepository struct {
	store *Store
}

func (r *ContractRepository)  Create(contract *model.Contract) error {
	return r.store.db.QueryRow(
		"INSERT INTO contracts (responseId, companyId, freelancerId, startTime, endTime, status, grade" +
			"paymentAmount) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING accountId",
		contract.ResponseID,
		contract.CompanyID,
		contract.FreelancerID,
		contract.StartTime,
		contract.EndTime,
		contract.Status,
		contract.Grade,
		contract.PaymentAmount,
	).Scan(&contract.ID)
}

func (r *ContractRepository) Find(id int64) (*model.Contract, error) {
	c := &model.Contract{}
	if err := r.store.db.QueryRow(
		"SELECT id, responseId, companyId, freelancerId, startTime, endTime, status, grade, " +
			"paymentAmount FROM contracts WHERE id = $1",
		id,
	).Scan(
		&c.ID,
		&c.ResponseID,
		&c.CompanyID,
		&c.FreelancerID,
		&c.StartTime,
		&c.EndTime,
		&c.Status,
		&c.Grade,
		&c.PaymentAmount,
	); err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ContractRepository) Edit(c * model.Contract) error {
	return r.store.db.QueryRow("UPDATE contracts SET freelancerId = $1, startTime = $2, " +
		"endTime = $3, status = $4, grade = $5, paymentAmount = $6 WHERE id = $7 RETURNING id",
		c.FreelancerID,
		c.StartTime,
		c.EndTime,
		c.Status,
		c.Grade,
		c.PaymentAmount,
		c.ID,
	).Scan(&c.ID)
}

func (r *ContractRepository) List(id int64, mode string) ([]model.Contract, error) {
	var contracts []model.Contract

	var query string
	if mode == ContractListByCompany {
		query = "SELECT id, responseId, companyId, freelancerId, startTime, endTime, status, grade, " +
			"paymentAmount FROM contracts WHERE companyId = $1"
	} else if mode == ContractListByFreelancer {
		query = "SELECT id, responseId, companyId, freelancerId, startTime, endTime, status, grade, " +
			"paymentAmount FROM contracts WHERE freelancerId = $1"
	}

	rows, err := r.store.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := model.Contract{}
		err := rows.Scan(&c.ID, &c.ResponseID, &c.CompanyID, &c.FreelancerID, &c.StartTime, &c.EndTime,
			&c.Status, &c.Grade, &c.PaymentAmount)
		if err != nil {
			return nil , err
		}
		contracts = append(contracts , c)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return contracts, nil
}