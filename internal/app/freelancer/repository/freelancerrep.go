package freelancerRepository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/monitoring"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

const (
	PATH2AVATAR = "https://89.208.211.100:8080/account/avatar/"
)

type FreelancerRepository struct {
	db *sql.DB
}

func NewFreelancerRepository(db *sql.DB) freelancer.Repository {
	return &FreelancerRepository{db}
}

func (r *FreelancerRepository) Create(f *model.Freelancer) error {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "freelancer", "method": "create"}))
	defer timer.ObserveDuration()

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

func (r *FreelancerRepository) Find(id int64) (*model.ExtendFreelancer, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "freelancer", "method": "find"}))
	defer timer.ObserveDuration()

	f := &model.Freelancer{}
	exF := &model.ExtendFreelancer{}
	if err := r.db.QueryRow(
		"SELECT f.id, f.accountId, f.country, f.city, f.address, f.phone, f.tagLine, "+
			"f.overview, f.experienceLevelId, f.specialityId, u.firstName, u.secondName " +
			"FROM freelancers AS f " +
			"INNER JOIN users AS u " +
			"ON (f.accountid = u.accountid) " +
			"WHERE id = $1",
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
		&exF.FirstName,
		&exF.SecondName,
	); err != nil {
		return nil, err
	}
	f.Avatar = PATH2AVATAR + strconv.Itoa(int(f.AccountId))
	exF.F = f
	return exF, nil
}

func (r *FreelancerRepository) FindByUser(accountId int64) (*model.Freelancer, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "freelancer", "method": "findByUser"}))
	defer timer.ObserveDuration()

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
	f.Avatar = PATH2AVATAR + strconv.Itoa(int(f.AccountId))
	return f, nil
}

func (r *FreelancerRepository) Edit(f *model.Freelancer) error {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "freelancer", "method": "edit"}))
	defer timer.ObserveDuration()

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

func (r *FreelancerRepository) ListOnPattern(pattern string) ([]model.ExtendFreelancer, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "freelancer", "method": "listInPattern"}))
	defer timer.ObserveDuration()

	var exFreelancers []model.ExtendFreelancer
	rows, err := r.db.Query(
		"SELECT F.id, F.accountId, F.country, F.city, F.address, F.phone, F.tagLine, "+
			" F.overview, F.experienceLevelId, F.specialityId, U.firstname , U.secondname "+
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
		exFreelancer := model.ExtendFreelancer{}
		err := rows.Scan(&f.ID, &f.AccountId, &f.Country, &f.City, &f.Address, &f.Phone,
			&f.TagLine, &f.Overview, &f.ExperienceLevelId, &f.SpecialityId, &exFreelancer.FirstName, &exFreelancer.SecondName)

		if err != nil {
			return nil, err
		}

		f.Avatar = PATH2AVATAR + strconv.Itoa(int(f.AccountId))

		exFreelancer.F = &f
		exFreelancers = append(exFreelancers, exFreelancer)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return exFreelancers, nil
}

func (r *FreelancerRepository) FindPartByTime(offset int, limit int) ([]model.ExtendFreelancer, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "freelancer", "method": "listInPattern"}))
	defer timer.ObserveDuration()

	var exFreelancers []model.ExtendFreelancer

	rows, err := r.db.Query(
		"SELECT F.id, F.accountId, F.country, F.city, F.address, F.phone, F.tagLine, "+
			" F.overview, F.experienceLevelId, F.specialityId, U.firstname , U.secondname "+
			"FROM freelancers AS F "+
			"INNER JOIN users AS U ON (F.accountid = U.accountid) "+
			"OFFSET $1 LIMIT $2 ",
		offset,
		limit,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		f := model.Freelancer{}
		exFreelancer := model.ExtendFreelancer{}
		err := rows.Scan(&f.ID, &f.AccountId, &f.Country, &f.City, &f.Address, &f.Phone,
			&f.TagLine, &f.Overview, &f.ExperienceLevelId, &f.SpecialityId, &exFreelancer.FirstName, &exFreelancer.SecondName)

		if err != nil {
			return nil, err
		}
		f.Avatar = PATH2AVATAR + strconv.Itoa(int(f.AccountId))

		exFreelancer.F = &f
		exFreelancers = append(exFreelancers, exFreelancer)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return exFreelancers, nil
}
