package test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
	"reflect"
	"testing"
)

func TestContractRep_Create(t *testing.T) {
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
	expect := []*model.Contract{
		{ID: elemID},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	u := testUser(t)
	f := testFreelancer(t, u)
	m := testManager(t, u)
	j := testJob(t, m)
	r := testResponse(t, f, j)
	c := testCompany(t)

	contract := testContract(t, r, c, f)

	// TODO: uncomment when validation will be implemented
	/*if err := c.Validate(); err != nil {
		t.Fatal()
	}*/

	//ok query
	mock.
		ExpectQuery(`INSERT INTO contracts`).
		WithArgs(contract.ResponseID, contract.CompanyID, contract.FreelancerID, contract.StartTime,
			contract.EndTime, contract.Status, contract.Grade, contract.PaymentAmount).
		WillReturnRows(rows)

	err = store.Contract().Create(contract)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if c.ID != 1 {
		t.Errorf("bad id: want %v, have %v", c.ID, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectQuery(`INSERT INTO contracts`).
		WithArgs(contract.ResponseID, contract.CompanyID, contract.FreelancerID, contract.StartTime,
			contract.EndTime, contract.Status, contract.Grade, contract.PaymentAmount).
		WillReturnError(fmt.Errorf("bad query"))

	err = store.Contract().Create(contract)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestContractRep_Find(t *testing.T) {
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
		NewRows([]string{"id", "responseId", "companyId", "freelancerId", "startTime", "endTime", "status",
			"grade", "paymentAmount" })

	u := testUser(t)
	f := testFreelancer(t, u)
	m := testManager(t, u)
	j := testJob(t, m)
	r := testResponse(t, f, j)
	c := testCompany(t)
	contract := testContract(t, r, c, f)
	expect := []*model.Contract{
		contract,
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.ResponseID, item.CompanyID, item.FreelancerID, item.StartTime, item.EndTime,
			item.Status, item.Grade, item.PaymentAmount)
	}

	mock.
		ExpectQuery("SELECT id, responseId, companyId, freelancerId, startTime, endTime, status, grade, " +
			"paymentAmount FROM contracts WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	store := sqlstore.New(db)

	item, err := store.Contract().Find(elemID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(item, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], item)
		return
	}

	// query error
	mock.
		ExpectQuery("SELECT id, responseId, companyId, freelancerId, startTime, endTime, status, grade, " +
			"paymentAmount FROM contracts WHERE").
		WithArgs(elemID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = store.Contract().Find(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	expect = []*model.Contract{
		contract,
	}

	mock.
		ExpectQuery("SELECT id, responseId, companyId, freelancerId, startTime, endTime, status, grade, " +
			"paymentAmount FROM contracts WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	_, err = store.Contract().Find(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestContractRep_Edit(t *testing.T) {
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
	expect := []*model.Contract{
		{ ID: elemID },
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	u := testUser(t)
	f := testFreelancer(t, u)
	m := testManager(t, u)
	j := testJob(t, m)
	r := testResponse(t, f, j)
	c := testCompany(t)
	contract := testContract(t, r, c, f)

	// TODO: uncomment when validation will be implemented
	/*if err := f.Validate(); err != nil {
		t.Fatal()
	}*/

	//ok query
	contract.Grade = 5
	contract.Status = "done"
	mock.
		ExpectQuery(`UPDATE contracts SET`).
		WithArgs(contract.FreelancerID, contract.StartTime, contract.EndTime, contract.Status, contract.Grade,
			contract.PaymentAmount, contract.ID).
		WillReturnRows(rows)

	err = store.Contract().Edit(contract)
	if err != nil {
		t.Fatal(err)
	}

	if contract.ID != 1 {
		t.Errorf("bad id: want %v, have %v", c.ID, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
