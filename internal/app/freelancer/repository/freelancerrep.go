package freelancerRepository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
)

type FreelancerRepository struct {
	db *sql.DB
}

func NewFreelancerRepository(db *sql.DB) freelancer.Repository {
	return &FreelancerRepository{db}
}

func (r *FreelancerRepository) Create(f *model.Freelancer) error {
	return r.db.QueryRow(
		"INSERT INTO freelancers (accountId, country, city, address, phone, tagLine, "+
			"overview, experienceLevelId, specialityId) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
		f.AccountId,
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

func (r *FreelancerRepository) Find(id int64) (*model.Freelancer, error) {
	f := &model.Freelancer{}
	if err := r.db.QueryRow(
		"SELECT id, accountId, country, city, address, phone, tagLine, "+
			"overview, experienceLevelId, specialityId FROM freelancers WHERE id = $1",
		id,
	).Scan(
		&f.ID,
		&f.AccountId,
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

func (r *FreelancerRepository) FindByUser(accountId int64) (*model.Freelancer, error) {
	f := &model.Freelancer{}
	if err := r.db.QueryRow(
		"SELECT id, accountId, country, city, address, phone, tagLine, "+
			"overview, experienceLevelId, specialityId FROM freelancers WHERE accountId = $1",
		accountId,
	).Scan(
		&f.ID,
		&f.AccountId,
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

func (r *FreelancerRepository) Edit(f *model.Freelancer) error {
	return r.db.QueryRow("UPDATE freelancers SET country = $1, city = $2, address = $3, "+
		"phone = $4, tagLine = $5, overview = $6, experienceLevelId = $7, specialityId = $8 WHERE id = $9 RETURNING id",
		f.Country,
		f.City,
		f.Address,
		f.Phone,
		f.TagLine,
		f.Overview,
		f.ExperienceLevelId,
		f.SpecialityId,
		f.ID,
	).Scan(&f.ID)
}

func (r *FreelancerRepository) ListOnPattern(pattern string) ([]model.Freelancer, error) {
	var freelancers []model.Freelancer
	rows, err := r.db.Query(
		"SELECT F.id, F.accountId, F.country, F.city, F.address, F.phone, F.tagLine, "+
			" F.overview, F.experienceLevelId, F.specialityId "+
			"FROM freelancers AS F "+
			"INNER JOIN users AS U ON (F.accountid = U.accountid) "+
			"WHERE to_tsvector('russian' , U.firstname) || to_tsvector('russian' , U.secondname) "+
			" @@ plainto_tsquery('russian', $1) LIMIT 10",
		pattern,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		f := model.Freelancer{}
		err := rows.Scan(&f.ID, &f.AccountId, &f.Country, &f.City, &f.Address, &f.Phone,
			&f.TagLine, &f.Overview, &f.ExperienceLevelId, &f.SpecialityId)
		if err != nil {
			return nil, err
		}
		freelancers = append(freelancers, f)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return freelancers, nil
}
