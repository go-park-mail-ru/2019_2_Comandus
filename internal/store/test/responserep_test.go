package test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
	"testing"
)

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

	store := sqlstore.New(db)
	rows := sqlmock.
		NewRows([]string{"id"})

	var elemID int64 = 1
	expect := []*model.Response{
		{ ID: elemID },
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	u := testUser(t)
	m := testManager(t, u)
	f := testFreelancer(t, u)
	j := testJob(t, m)
	r := testResponse(t, f, j)

	r.BeforeCreate()
	if err := r.Validate(1); err != nil {
		t.Fatal(err)
	}

	//ok query
	mock.
		ExpectQuery(`INSERT INTO responses`).
		WithArgs(r.FreelancerId, r.JobId, r.Files, r.Date, r.StatusManager, r.StatusFreelancer, r.PaymentAmount).
		WillReturnRows(rows)

	err = store.Response().Create(r)
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

	err = store.Response().Create(r)
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
	store := sqlstore.New(db)

	rows := sqlmock.
		NewRows([]string{"id"})

	var elemID int64 = 1
	expect := []*model.Response{
		{ ID: elemID },
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	u := testUser(t)
	m := testManager(t, u)
	f := testFreelancer(t, u)
	j := testJob(t, m)
	r := testResponse(t, f, j)

	r.ID = 1
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

	err = store.Response().Edit(r)
	if err != nil {
		t.Fatal(err)
	}

	if u.ID != 1 {
		t.Errorf("bad id: want %v, have %v", u.ID, 1)
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

	u := testUser(t)
	m := testManager(t, u)
	j := testJob(t,m)

	f := testFreelancer(t, u)

	expect := []*model.Response{
		testResponse(t, f, j),
		testResponse(t, f, j),
		testResponse(t, f, j),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.FreelancerId, item.JobId, item.Files, item.Date, item.StatusManager,
			item.StatusFreelancer, item.PaymentAmount)
	}

	mock.
		ExpectQuery("SELECT id, freelancerId, jobId, files, date, statusManager, statusFreelancer, paymentAmount " +
			"FROM responses WHERE").
		WithArgs(f.ID).
		WillReturnRows(rows)

	store := sqlstore.New(db)

	responses, err := store.Response().ListForFreelancer(f.ID)
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

	_, err = store.Response().ListForFreelancer(f.ID)
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

	u := testUser(t)
	m := testManager(t, u)
	j := testJob(t, m)

	f1 := testFreelancer(t, u)
	f2 := testFreelancer(t, u)
	f2.ID = 2
	f3 := testFreelancer(t, u)
	f3.ID = 3

	expect := []*model.Response{
		testResponse(t, f1, j),
		testResponse(t, f2, j),
		testResponse(t, f3, j),
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
		WithArgs(m.ID).
		WillReturnRows(rows)

	store := sqlstore.New(db)

	responses, err := store.Response().ListForManager(m.ID)
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

	_, err = store.Response().ListForManager(m.ID)
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

	u := testUser(t)
	f := testFreelancer(t, u)
	m := testManager(t, u)
	j := testJob(t, m)

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "freelancerId", "jobId", "files", "date", "statusManager", "statusFreelancer", "paymentAmount"})
	expect := []*model.Response{
		testResponse(t, f, j),
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

	store := sqlstore.New(db)

	item, err := store.Response().Find(elemID)
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

	_, err = store.Response().Find(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}