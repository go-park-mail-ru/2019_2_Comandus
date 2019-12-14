package jobRepository

import (
	"database/sql"
	user_job "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/monitoring"
	"github.com/prometheus/client_golang/prometheus"
)

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) user_job.Repository {
	return &JobRepository{db}
}

// TODO: remove hire manager
func (r *JobRepository) Create(j *model.Job) error {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep":"job", "method":"create"}))
	defer timer.ObserveDuration()

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
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep":"job", "method":"find"}))
	defer timer.ObserveDuration()

	j := &model.Job{}
	if err := r.db.QueryRow(
		"SELECT id, managerId, title, description, files, specialityId, experienceLevelId, paymentAmount, " +
			"country, city, jobTypeId, date, status FROM jobs WHERE id = $1 AND status != $2",
		id,
		model.JobStateDeleted,
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
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep":"job", "method":"edit"}))
	defer timer.ObserveDuration()

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
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep":"job", "method":"list"}))
	defer timer.ObserveDuration()

	var jobs []model.Job
	rows, err := r.db.Query(
		"SELECT id, managerId, title, description, files, specialityId, experienceLevelId, paymentAmount, " +
			"country, city, jobTypeId, date, status " +
			"FROM jobs WHERE status != $1 " +
			"ORDER BY id END DESC LIMIT 20",
			model.JobStateDeleted)

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

func (r *JobRepository) ListOnPattern(pattern string, params model.JobSearchParams) ([]model.Job, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep":"job", "method":"listOnPattern"}))
	defer timer.ObserveDuration()

	var jobs []model.Job
	rows, err := r.db.Query(
		"SELECT id, managerId, title, description, files, specialityId, experienceLevelId, paymentAmount, "+
			"country, city, jobTypeId, date, status "+
			"FROM jobs " +
			"WHERE to_tsvector('russian' , title) @@ plainto_tsquery('russian', $1) AND " +
			"status != $1 AND " +
			"($2 = 0 OR paymentAmount <= $2 AND paymentAmount >= $3) AND " +
			"($4 = 0 OR grade <= $4 AND grade >= $5) AND " +
			"($6 = '' OR country = $6) AND " +
			"($7 = '' OR city = $7) AND " +
			"(($8 AND experienceLevelId = 1) OR ($9 AND experienceLevelId = 2) OR ($10 AND experienceLevelId = 3)) " +
			"ORDER BY " +
			"CASE WHEN $11 THEN id END DESC " +
			"CASE WHEN !$11 THEN id END ASC " +
			"LIMIT 10",
			model.JobStateDeleted,
			params.MaxPaymentAmount, params.MinPaymentAmount,
			params.MaxGrade, params.MinGrade,
			params.Country,
			params.City,
			params.ExperienceLevel[0], params.ExperienceLevel[1], params.ExperienceLevel[2],
			params.Desc)
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
