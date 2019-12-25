package jobRepository

import (
	"database/sql"
	"errors"
	user_job "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/job"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/monitoring"
	"github.com/prometheus/client_golang/prometheus"
	"log"
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
		Labels{"rep": "job", "method": "create"}))
	defer timer.ObserveDuration()

	return r.db.QueryRow(
		"INSERT INTO jobs (managerId, title, description, files, specialityId, experienceLevelId, paymentAmount, "+
			"country, city, jobTypeId, date, status, tagLine) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, "+
			"$11, $12, $13) RETURNING id",
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
		j.TagLine,
	).Scan(&j.ID)
}

func (r *JobRepository) Find(id int64) (*model.Job, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "job", "method": "find"}))
	defer timer.ObserveDuration()

	j := &model.Job{}
	if err := r.db.QueryRow(
		"SELECT id, managerId, title, description, files, specialityId, experienceLevelId, paymentAmount, "+
			"country, city, jobTypeId, date, status, tagLine FROM jobs WHERE id = $1 AND status != $2",
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
		&j.TagLine,
	); err != nil {
		return nil, err
	}
	return j, nil
}

func (r *JobRepository) Edit(j *model.Job) error {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "job", "method": "edit"}))
	defer timer.ObserveDuration()

	return r.db.QueryRow("UPDATE jobs SET title = $1, description = $2, files = $3, "+
		"specialityId = $4, experienceLevelId = $5, paymentAmount = $6, country = $7, city = $8, "+
		"jobTypeId = $9, status = $10, tagLine = $11 WHERE id = $12 RETURNING id",
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
		j.TagLine,
		j.ID,
	).Scan(&j.ID)
}

func (r *JobRepository) List() ([]model.Job, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "job", "method": "list"}))
	defer timer.ObserveDuration()

	var jobs []model.Job
	rows, err := r.db.Query(
		"SELECT id, managerId, title, description, files, specialityId, experienceLevelId, paymentAmount, "+
			"country, city, jobTypeId, date, status, tagLine "+
			"FROM jobs WHERE status != $1 "+
			"ORDER BY id DESC LIMIT 20",
		model.JobStateDeleted)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		j := model.Job{}
		err := rows.Scan(&j.ID, &j.HireManagerId, &j.Title, &j.Description, &j.Files, &j.SpecialityId,
			&j.ExperienceLevelId, &j.PaymentAmount, &j.Country, &j.City, &j.JobTypeId, &j.Date, &j.Status, &j.TagLine)
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

func (r *JobRepository) ListOnPattern(pattern string, params model.SearchParams) ([]model.Job, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "job", "method": "listOnPattern"}))
	defer timer.ObserveDuration()

	log.Println(params.JobType)
	var jobs []model.Job
	rows, err := r.db.Query(
		"SELECT J.id, J.managerId, J.title, J.description, J.files, J.specialityId, J.experienceLevelId, J.paymentAmount, "+
			"J.country, J.city, J.jobTypeId, J.date, J.status, J.tagLine "+
			"FROM jobs AS J "+
			"LEFT OUTER JOIN responses AS R "+
			"ON J.id = R.jobId "+
			"WHERE (LOWER(J.title) LIKE '%' || LOWER($1) || '%' OR LOWER(J.tagLine) LIKE '%'|| LOWER($1) ||'%') AND "+
			"J.status != $2 AND J.status != $15 AND "+
			"($3 = 0 OR J.paymentAmount <= $3) AND "+
			"($4 = 0 OR J.paymentAmount >= $4) AND "+
			"($5 = -1 OR J.country = $5) AND "+
			"($6 = -1 OR J.city = $6) AND "+
			"($14 = -1 OR jobTypeId = $14) AND "+
			"($16 = -1 OR specialityId = $16) AND "+
			"(($7 AND J.experienceLevelId = 1) OR ($8 AND J.experienceLevelId = 2) OR ($9 AND J.experienceLevelId = 3)) "+
			"GROUP BY J.id "+
			"HAVING ($12 = 0 OR COUNT(*) >= $12) AND "+
			"($13 = 0 OR COUNT(*) <= $13) "+
			"ORDER BY "+
			"CASE WHEN $10 THEN J.id END DESC, "+
			"CASE WHEN NOT $10 THEN J.id END ASC "+
			"LIMIT CASE WHEN $11 > 0 THEN $11 END;",
		pattern, model.JobStateDeleted,
		params.MaxPaymentAmount,
		params.MinPaymentAmount,
		params.Country,
		params.City,
		params.ExperienceLevel[0], params.ExperienceLevel[1], params.ExperienceLevel[2],
		params.Desc,
		params.Limit,
		params.MinProposals,
		params.MaxProposals,
		params.JobType,
		model.JobStateClosed,
		params.SpecialityId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		j := model.Job{}
		err := rows.Scan(&j.ID, &j.HireManagerId, &j.Title, &j.Description, &j.Files, &j.SpecialityId,
			&j.ExperienceLevelId, &j.PaymentAmount, &j.Country, &j.City, &j.JobTypeId, &j.Date, &j.Status, &j.TagLine)
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

func (r *JobRepository) ListMyJobs(managerID int64) ([]model.Job, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "job", "method": "list"}))
	defer timer.ObserveDuration()

	var jobs []model.Job
	rows, err := r.db.Query(
		"SELECT id, managerId, title, description, files, specialityId, experienceLevelId, paymentAmount, "+
			"country, city, jobTypeId, date, status, tagLine FROM jobs AS j "+
			"WHERE status != $1 AND managerid = $2 ORDER BY id DESC",
		model.JobStateDeleted,
		managerID,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		j := model.Job{}
		err := rows.Scan(&j.ID, &j.HireManagerId, &j.Title, &j.Description, &j.Files, &j.SpecialityId,
			&j.ExperienceLevelId, &j.PaymentAmount, &j.Country, &j.City, &j.JobTypeId, &j.Date, &j.Status, &j.TagLine)
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

// ?
func (r *JobRepository) GetUserIDByJobID(jobID int64) (int64, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "job", "method": "GetUserIDByJobID"}))
	defer timer.ObserveDuration()

	var accID int64
	ok := r.db.QueryRow(
		"SELECT u.accountid "+
			"FROM jobs AS j "+
			"INNER JOIN managers as m "+
			"ON (m.id = j.managerid) "+
			"INNER JOIN users as u "+
			"ON (u.accountid = m.accountid) "+
			"WHERE j.id = $1 ", jobID)
	if ok == nil {
		return -1, errors.New("no job with this id")
	}
	err := ok.Scan(&accID)

	return accID, err
}

func (r *JobRepository) ChangeStatus(jobID int64, status string) error {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep": "job", "method": "changeStatus"}))
	defer timer.ObserveDuration()

	r.db.QueryRow(
		"UPDATE jobs SET status = $1 WHERE id = $2",
		status, jobID)
	return nil
}
