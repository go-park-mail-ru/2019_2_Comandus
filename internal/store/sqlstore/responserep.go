package sqlstore

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
)

type ResponseRepository struct {
	store *Store
}

func (r *ResponseRepository) Create(response *model.Response) error {
	return r.store.db.QueryRow(
		"INSERT INTO responses (freelancerId, jobId, files, date, statusManager, statusFreelancer) " +
			"VALUES ($1, $2, $3, $4, $5, $6) RETURNING accountId",
		response.FreelancerId,
		response.JobId,
		response.Files,
		response.Date,
		response.StatusManager,
		response.StatusFreelancer,
	).Scan(&response.ID)
}

func (r *ResponseRepository) Edit(response *model.Response) error {
	return r.store.db.QueryRow(
		"UPDATE responses SET files = $1, statusManager = $2, statusFreelancer = $3 WHERE id = $4 " +
			"RETURNING accountId",
		response.Files,
		response.StatusManager,
		response.StatusFreelancer,
		response.ID,
	).Scan(&response.ID)
}

func (r *ResponseRepository) ListForFreelancer(id int64) ([]model.Response, error) {
	var responses []model.Response
	rows, err := r.store.db.Query(
		"SELECT id, freelancerId, jobId, files, date, statusManager, statusFreelancer " +
			"FROM responses WHERE freelancerId = $1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		r := model.Response{}
		err := rows.Scan(&r.ID, &r.FreelancerId, &r.JobId, &r.Files, &r.Date, &r.StatusManager,
			&r.StatusFreelancer)
		if err != nil {
			return nil , err
		}
		responses = append(responses , r)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *ResponseRepository) ListForManager(id int64) ([]model.Response, error) {
	var responses []model.Response
	rows, err := r.store.db.Query(
		"SELECT responses.id, responses.freelancerId, responses.jobId, responses.files, responses.date, " +
			"responses.statusManager, responses.statusFreelancer " +
			"FROM responses " +
			"INNER JOIN jobs" +
			"ON jobs.id = responses.jobId" +
			"WHERE jobs.managerId = $1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		r := model.Response{}
		err := rows.Scan(&r.ID, &r.FreelancerId, &r.JobId, &r.Files, &r.Date, &r.StatusManager,
			&r.StatusFreelancer)
		if err != nil {
			return nil , err
		}
		responses = append(responses , r)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return responses, nil
}

func (r * ResponseRepository) Find(id int64) (*model.Response, error) {
	response := &model.Response{}
	if err := r.store.db.QueryRow(
		"SELECT id, freelancerId, jobId, files, date, statusManager, statusFreelancer FROM responses WHERE id = $1",
		id,
	).Scan(
		&response.ID,
		&response.FreelancerId,
		&response.JobId,
		&response.Files,
		&response.Date,
		&response.StatusManager,
		&response.StatusFreelancer,
	); err != nil {
		return nil, err
	}
	return response, nil
}