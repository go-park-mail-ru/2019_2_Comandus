package sqlstore

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
)

type FreelancerRepository struct {
	store *Store
}

func (r *FreelancerRepository) Create(f *model.Freelancer) error {
	return r.store.db.QueryRow(
		"INSERT INTO freelancers (accountId, registrationDate, country, city, address, phone, tagLine, " +
			"overview, experienceLevelId, specialityId) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		f.AccountId,
		f.RegistrationDate,
		f.Country,
		f.City,
		f.Address,
		f.Phone,
		f.TagLine,
		f.Overview,
		f.ExperienceLevelId,
		f.SpecialityId,
	).Scan(&f.ID)
}

func (r *FreelancerRepository) Find(id int) (*model.Freelancer, error) {
	f := &model.Freelancer{}
	if err := r.store.db.QueryRow(
		"SELECT id, accountId, registrationDate, country, city, address, phone, tagLine, " +
			"overview, experienceLevelId, specialityId FROM freelancers WHERE id = $1",
		id,
	).Scan(
		&f.AccountId,
		&f.RegistrationDate,
		&f.Country,
		&f.City,
		&f.Address,
		&f.Phone,
		&f.TagLine,
		&f.Overview,
		&f.ExperienceLevelId,
		&f.SpecialityId,
	); err != nil {
		return nil, err
	}
	return f, nil
}
