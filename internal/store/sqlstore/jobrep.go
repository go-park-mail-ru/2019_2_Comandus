package sqlstore

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type JobRepository struct {
	store *Store
}

func (r *JobRepository) Create(j *model.Job, m *model.HireManager) error {
	return r.store.db.QueryRow(
		"INSERT INTO jobs (managerId, title, description, files, specialityId, experienceLevelId, paymentAmount, " +
			"country, city, jobTypeId) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		m.ID,
		j.Title,
		j.Description,
		j.Files,
		j.SpecialityId,
		j.ExperienceLevelId,
		j.PaymentAmount,
		j.Country,
		j.City,
		j.JobTypeId,
	).Scan(&j.ID)
}

func (r *JobRepository) Find(id int) (*model.Job, error) {
	j := &model.Job{}
	if err := r.store.db.QueryRow(
		"SELECT id, managerId, title, description, files, specialityId, experienceLevelId, paymentAmount, " +
			"country, city, jobTypeId FROM jobs WHERE id = $1",
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
	); err != nil {
		return nil, err
	}
	return j, nil
}


func (r *JobRepository) GetAllJobs() ([]model.Job, error) {
	var jobs []model.Job
	rows, err := r.store.db.Query(
		"SELECT id, managerId, title, description, files, specialityId, experienceLevelId, paymentAmount, " +
			"country, city, jobTypeId FROM jobs LIMIT 10")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		j := model.Job{}
		err := rows.Scan(&j.ID, &j.HireManagerId, &j.Title, &j.Description, &j.Files, &j.SpecialityId,
			&j.ExperienceLevelId, &j.PaymentAmount, &j.Country, &j.City, &j.JobTypeId)
		if err != nil {
			return nil , err
		}
		jobs = append(jobs , j)
	}
	return jobs, nil
}

func (r *JobRepository) Edit(j *model.Job) error {
	return r.store.db.QueryRow("UPDATE jobs SET title = $1, description = $2, files = $3, " +
		"specialityId = $4, experienceLevelId = $5, paymentAmount = $6, country = $7, city = $8, " +
		"jobTypeId = $9 WHERE id = $10 RETURNING id",
		j.Title,
		j.Description,
		j.Files,
		j.SpecialityId,
		j.ExperienceLevelId,
		j.PaymentAmount,
		j.Country,
		j.City,
		j.JobTypeId,
		j.ID,
	).Scan(&j.ID)
}

