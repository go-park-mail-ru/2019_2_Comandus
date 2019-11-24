package jobRepository

import (
	"database/sql"
	user_job "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
)

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) user_job.Repository {
	return &JobRepository{db}
}

// TODO: remove hire manager
func (r *JobRepository) Create(j *model.Job) error {
	return r.db.QueryRow(
		"INSERT INTO jobs (managerId, title, description, files, specialityId, experienceLevelId, paymentAmount, "+
			"country, city, jobTypeId, date, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id",
		j.HireManagerId,
		j.Title,
		j.Description,
		j.Files,
		j.SpecialityId,
		j.ExperienceLevelId,
		j.PaymentAmount,
		j.Country,
		j.City,
		j.JobTypeId,
		j.Date,
		j.Status,
	).Scan(&j.ID)
}

func (r *JobRepository) Find(id int64) (*model.Job, error) {
	j := &model.Job{}
	if err := r.db.QueryRow(
		"SELECT id, managerId, title, description, files, specialityId, experienceLevelId, paymentAmount, "+
			"country, city, jobTypeId, date, status FROM jobs WHERE id = $1",
		id,
	).Scan(
		&j.ID,
		&j.HireManagerId,
		&j.Title,
		&j.Description,
		&j.Files,
		&j.SpecialityId,
		&j.ExperienceLevelId,
		&j.PaymentAmount,
		&j.Country,
		&j.City,
		&j.JobTypeId,
		&j.Date,
		&j.Status,
	); err != nil {
		return nil, err
	}
	return j, nil
}

func (r *JobRepository) Edit(j *model.Job) error {
	return r.db.QueryRow("UPDATE jobs SET title = $1, description = $2, files = $3, "+
		"specialityId = $4, experienceLevelId = $5, paymentAmount = $6, country = $7, city = $8, "+
		"jobTypeId = $9, status = $10 WHERE id = $11 RETURNING id",
		j.Title,
		j.Description,
		j.Files,
		j.SpecialityId,
		j.ExperienceLevelId,
		j.PaymentAmount,
		j.Country,
		j.City,
		j.JobTypeId,
		j.Status,
		j.ID,
	).Scan(&j.ID)
}

func (r *JobRepository) List() ([]model.Job, error) {
	var jobs []model.Job
	rows, err := r.db.Query(
		"SELECT id, managerId, title, description, files, specialityId, experienceLevelId, paymentAmount, " +
			"country, city, jobTypeId, date, status FROM jobs ORDER BY id DESC LIMIT 10")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		j := model.Job{}
		err := rows.Scan(&j.ID, &j.HireManagerId, &j.Title, &j.Description, &j.Files, &j.SpecialityId,
			&j.ExperienceLevelId, &j.PaymentAmount, &j.Country, &j.City, &j.JobTypeId, &j.Date, &j.Status)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *JobRepository) ListOnPattern(pattern string) ([]model.Job, error) {
	var jobs []model.Job
	rows, err := r.db.Query(
		"SELECT id, managerId, title, description, files, specialityId, experienceLevelId, paymentAmount, "+
			"country, city, jobTypeId, date, status "+
			"FROM jobs WHERE to_tsvector('russian' , title) @@ plainto_tsquery('russian', $1) LIMIT 10",
		pattern,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		j := model.Job{}
		err := rows.Scan(&j.ID, &j.HireManagerId, &j.Title, &j.Description, &j.Files, &j.SpecialityId,
			&j.ExperienceLevelId, &j.PaymentAmount, &j.Country, &j.City, &j.JobTypeId, &j.Date, &j.Status)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return jobs, nil
}
