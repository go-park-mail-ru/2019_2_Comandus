package responseRepository

import (
	"database/sql"
	user_response "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/monitoring"
	"github.com/prometheus/client_golang/prometheus"
)

type ResponseRepository struct {
	db *sql.DB
}

func NewResponseRepository(db *sql.DB) user_response.Repository {
	return &ResponseRepository{db}
}

func (r *ResponseRepository) Create(response *model.Response) error {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep":"response", "method":"create"}))
	defer timer.ObserveDuration()

	return r.db.QueryRow(
		"INSERT INTO responses (freelancerId, jobId, files, date, statusManager, statusFreelancer, paymentAmount) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		response.FreelancerId,
		response.JobId,
		response.Files,
		response.Date,
		response.StatusManager,
		response.StatusFreelancer,
		response.PaymentAmount,
	).Scan(&response.ID)
}

func (r *ResponseRepository) Edit(response *model.Response) error {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep":"response", "method":"edit"}))
	defer timer.ObserveDuration()

	return r.db.QueryRow(
		"UPDATE responses SET files = $1, statusmanager = $2, statusFreelancer = $3, paymentAmount = $4 WHERE id = $5 "+
			"RETURNING id",
		response.Files,
		response.StatusManager,
		response.StatusFreelancer,
		response.PaymentAmount, /**/
		response.ID,
	).Scan(&response.ID)
}

func (r *ResponseRepository) ListForFreelancer(id int64) ([]model.Response, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep":"response", "method":"listForFreelancer"}))
	defer timer.ObserveDuration()

	var responses []model.Response
	rows, err := r.db.Query(
		"SELECT id, freelancerId, jobId, files, date, statusManager, statusFreelancer, paymentAmount "+
			"FROM responses WHERE freelancerId = $1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		r := model.Response{}
		err := rows.Scan(&r.ID, &r.FreelancerId, &r.JobId, &r.Files, &r.Date, &r.StatusManager,
			&r.StatusFreelancer, &r.PaymentAmount)
		if err != nil {
			return nil, err
		}
		responses = append(responses, r)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *ResponseRepository) ListForManager(id int64) ([]model.Response, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep":"response", "method":"listForManager"}))
	defer timer.ObserveDuration()

	var responses []model.Response
	rows, err := r.db.Query(
		"SELECT responses.id, responses.freelancerId, responses.jobId, responses.files, responses.date, "+
			"responses.statusManager, responses.statusFreelancer, responses.paymentAmount "+
			"FROM responses "+
			"INNER JOIN jobs "+
			"ON jobs.id = responses.jobId "+
			"WHERE jobs.managerId = $1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		r := model.Response{}
		err := rows.Scan(&r.ID, &r.FreelancerId, &r.JobId, &r.Files, &r.Date, &r.StatusManager,
			&r.StatusFreelancer, &r.PaymentAmount)
		if err != nil {
			return nil, err
		}
		responses = append(responses, r)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *ResponseRepository) Find(id int64) (*model.Response, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep":"response", "method":"find"}))
	defer timer.ObserveDuration()

	response := &model.Response{}
	if err := r.db.QueryRow(
		"SELECT id, freelancerId, jobId, files, date, statusManager, statusFreelancer, paymentAmount FROM responses WHERE id = $1",
		id,
	).Scan(
		&response.ID,
		&response.FreelancerId,
		&response.JobId,
		&response.Files,
		&response.Date,
		&response.StatusManager,
		&response.StatusFreelancer,
		&response.PaymentAmount,
	); err != nil {
		return nil, err
	}
	return response, nil
}


func (r *ResponseRepository) ListResponsesOnJobID(jobID int64) ([]model.ExtendResponse, error) {
	timer := prometheus.NewTimer(monitoring.DBQueryDuration.With(prometheus.
		Labels{"rep":"response", "method":"listResponsesOnJobID"}))
	defer timer.ObserveDuration()

	var responses []model.ExtendResponse
	rows, err := r.db.Query(
		"SELECT R.id, R.freelancerId, R.jobId, R.files, R.date, R.statusManager, R.statusFreelancer, " +
			"R.paymentAmount, U.firstname , U.secondname "+
			"FROM responses AS R " +
			"INNER JOIN freelancers AS F " +
			"ON R.freelancerid = F.id " +
			"INNER JOIN users AS U " +
			"ON U.accountid = F.accountid " +
			" WHERE jobid = $1", jobID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		r := model.Response{}
		exR := model.ExtendResponse{}
		err := rows.Scan(&r.ID, &r.FreelancerId, &r.JobId, &r.Files, &r.Date, &r.StatusManager,
			&r.StatusFreelancer, &r.PaymentAmount, &exR.FirstName, &exR.SecondName)
		if err != nil {
			return nil, err
		}
		exR.R = &r

		responses = append(responses, exR)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return responses, nil
}