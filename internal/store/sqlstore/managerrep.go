package sqlstore

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"time"
)

type ManagerRepository struct {
	store *Store
}

type HireManager struct {
	ID					int		`json:"id"`
	AccountID 			int		`json:"accountId"`
	RegistrationDate	time.Time	`json:"registrationDate"`
	Location			string 		`json:"location"`
	CompanyID			int 		`json:"companyId"`
}

func (r *ManagerRepository) Create(m *model.HireManager) (int64, error) {
	result, err := r.store.db.Exec(
		"INSERT INTO managers (accountId, registrationDate, location, companyId) " +
			"VALUES ($1, $2, $3, $4) RETURNING id",
		m.AccountID,
		m.RegistrationDate,
		m.Location,
		m.CompanyID,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ManagerRepository) Find(id int64) (*model.HireManager, error) {
	m := &model.HireManager{}
	if err := r.store.db.QueryRow(
		"SELECT id, accountId, registrationDate, location, companyId FROM managers WHERE id = $1",
		id,
	).Scan(
		&m.ID,
		&m.AccountID,
		&m.RegistrationDate,
		&m.Location,
		&m.CompanyID,
	); err != nil {
		return nil, err
	}
	return m, nil
}

func (r *ManagerRepository) FindByUser(accountId int64) (*model.HireManager, error) {
	m := &model.HireManager{}
	if err := r.store.db.QueryRow(
		"SELECT id, accountId, registrationDate, location, companyId FROM managers WHERE accountId = $1",
		accountId,
	).Scan(
		&m.ID,
		&m.AccountID,
		&m.RegistrationDate,
		&m.Location,
		&m.CompanyID,
	); err != nil {
		return nil, err
	}
	return m, nil
}

func (r *ManagerRepository) Edit(m * model.HireManager) (int64, error) {
	result, err := r.store.db.Exec(
		"UPDATE managers SET location = $1, companyId = $2 WHERE id = $3 RETURNING id",
		m.Location,
		m.CompanyID,
		m.ID,
	)

	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

