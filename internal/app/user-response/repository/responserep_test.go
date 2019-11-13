package responseRepository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"testing"
	"time"
)

const (
	freelancerId = 1
	jobId = 1
	managerId = 1
)

func testResponse(t *testing.T) *model.Response {
	t.Helper()
	return &model.Response{
		ID:            1,
		FreelancerId:  freelancerId,
		JobId:         jobId,
		Files:         "no files",
		Date:          time.Time{},
		StatusManager: model.ResponseStatusReview,
		PaymentAmount: 10000,
	}
}

func TestResponseRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	repo := NewResponseRepository(db)
	rows := sqlmock.
		NewRows([]string{"id"})

	var elemID int64 = 1
	expect := []*model.Response{
		{ ID: elemID },
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	r := testResponse(t)

	r.BeforeCreate()
	if err := r.Validate(1); err != nil {
		t.Fatal(err)
	}

	//ok query
	mock.
		ExpectQuery(`INSERT INTO responses`).
		WithArgs(r.FreelancerId, r.JobId, r.Files, r.Date, r.StatusManager, r.StatusFreelancer, r.PaymentAmount).
		WillReturnRows(rows)

	err = repo.Create(r)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if r.ID != 1 {
		t.Errorf("bad id: want %v, have %v", r.ID, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectQuery(`INSERT INTO responses`).
		WithArgs(r.FreelancerId, r.JobId, r.Files, r.Date, r.StatusManager, r.StatusFreelancer, r.PaymentAmount).
		WillReturnError(fmt.Errorf("bad query"))

	err = repo.Create(r)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestResponseRepository_Edit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	repo := NewResponseRepository(db)

	rows := sqlmock.
		NewRows([]string{"id"})

	var elemID int64 = 1
	expect := []*model.Response{
		{ ID: elemID },
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	r := testResponse(t)

	r.StatusManager = model.ResponseStatusAccepted
	r.StatusFreelancer = model.ResponseStatusReview
	r.PaymentAmount = 2000
	r.BeforeCreate()
	if err := r.Validate(1); err != nil {
		t.Fatal(err)
	}

	mock.
		ExpectQuery(`UPDATE responses SET`).
		WithArgs(r.Files, r.StatusManager, r.StatusFreelancer, r.PaymentAmount, r.ID).
		WillReturnRows(rows)

	err = repo.Edit(r)
	if err != nil {
		t.Fatal(err)
	}

	if r.ID != 1 {
		t.Errorf("bad id: want %v, have %v", r.ID, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestResponseRepository_ListForFreelancer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "freelancerId", "jobId", "files", "date", "statusManager", "statusFreelancer", "paymentAmount"})

	expect := []*model.Response{
		testResponse(t),
		testResponse(t),
		testResponse(t),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.FreelancerId, item.JobId, item.Files, item.Date, item.StatusManager,
			item.StatusFreelancer, item.PaymentAmount)
	}

	mock.
		ExpectQuery("SELECT id, freelancerId, jobId, files, date, statusManager, statusFreelancer, paymentAmount " +
			"FROM responses WHERE").
		WithArgs(freelancerId).
		WillReturnRows(rows)

	repo := NewResponseRepository(db)

	responses, err := repo.ListForFreelancer(freelancerId)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	for i := 0; i < 3; i++ {
		if !responses[i].IsEqual(expect[i]) {
			t.Errorf("results not match, want %v, have %v", expect[i], responses[i])
			return
		}
	}

	// query error
	mock.
		ExpectQuery("SELECT id, freelancerId, jobId, files, date, statusManager, statusFreelancer, paymentAmount " +
			"FROM responses WHERE").
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.ListForFreelancer(freelancerId)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestResponseRepository_ListForManager(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "freelancerId", "jobId", "files", "date", "statusManager", "statusFreelancer", "paymentAmount" })

	expect := []*model.Response{
		testResponse(t),
		testResponse(t),
		testResponse(t),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.FreelancerId, item.JobId, item.Files, item.Date, item.StatusManager,
			item.StatusFreelancer, item.PaymentAmount)
	}

	mock.
		ExpectQuery("SELECT responses.id, responses.freelancerId, responses.jobId, responses.files, responses.date, " +
		"responses.statusManager, responses.statusFreelancer, responses.paymentAmount " +
		"FROM responses " +
		"INNER JOIN jobs " +
		"ON jobs.id = responses.jobId " +
		"WHERE ").
		WithArgs(managerId).
		WillReturnRows(rows)

	repo := NewResponseRepository(db)

	responses, err := repo.ListForManager(managerId)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	for i := 0; i < 3; i++ {
		if !responses[i].IsEqual(expect[i]) {
			t.Errorf("results not match, want %v, have %v", expect[i], responses[i])
			return
		}
	}

	// query error
	mock.
		ExpectQuery("SELECT responses.id, responses.freelancerId, responses.jobId, responses.files, responses.date, " +
		"responses.statusManager, responses.statusFreelancer, responses.paymentAmount " +
		"FROM responses " +
		"INNER JOIN jobs " +
		"ON jobs.id = responses.jobId " +
		"WHERE ").
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.ListForManager(managerId)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestResponseRepository_Find(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	var elemID int64 = 1

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "freelancerId", "jobId", "files", "date", "statusManager", "statusFreelancer", "paymentAmount"})
	expect := []*model.Response{
		testResponse(t),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.FreelancerId, item.JobId, item.Files, item.Date, item.StatusManager,
			item.StatusFreelancer, item.PaymentAmount)
	}

	mock.
		ExpectQuery("SELECT id, freelancerId, jobId, files, date, " +
			"statusManager, statusFreelancer, paymentAmount FROM responses WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	repo := NewResponseRepository(db)

	item, err := repo.Find(elemID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !item.IsEqual(expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], item)
		return
	}

	// query error
	mock.
		ExpectQuery("SELECT id, freelancerId, jobId, files, date, " +
			"statusManager, statusFreelancer, paymentAmount FROM responses WHERE ").
		WithArgs(elemID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.Find(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}