package contractRepository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"reflect"
	"testing"
	"time"
)

const (
	responseId = 1
	companyId = 1
	freelancerId = 1
	)

func testContract(t * testing.T) *model.Contract {
	t.Helper()
	return &model.Contract{
		ID:            1,
		ResponseID:    responseId,
		CompanyID:     companyId,
		FreelancerID:  freelancerId,
		StartTime:     time.Now(),
		EndTime:       time.Time{},
		Status:        "review",
		Grade:         0,
		PaymentAmount: 100,
	}
}

func IsEqual(t *testing.T, c1 *model.Contract, c2 * model.Contract) bool {
	t.Helper()
	return c1.ID == c2.ID &&
		c1.ResponseID == c2.ResponseID &&
		c1.CompanyID == c2.CompanyID &&
		c1.FreelancerID == c2.FreelancerID &&
		c1.StartTime == c2.StartTime &&
		c1.EndTime == c2.EndTime &&
		c1.Status == c2.Status &&
		c1.Grade == c2.Grade &&
		c1.PaymentAmount == c2.PaymentAmount
}

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

	repo := NewContractRepository(db)
	rows := sqlmock.
		NewRows([]string{"id"})

	var elemID int64 = 1
	expect := []*model.Contract{
		{ID: elemID},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	contract := testContract(t)

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

	err = repo.Create(contract)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if contract.ID != 1 {
		t.Errorf("bad id: want %v, have %v", contract.ID, 1)
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

	err = repo.Create(contract)
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

	contract := testContract(t)
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

	repo := NewContractRepository(db)

	item, err := repo.Find(elemID)
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

	_, err = repo.Find(elemID)
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

	repo := NewContractRepository(db)

	rows := sqlmock.
		NewRows([]string{"id"})

	var elemID int64 = 1
	expect := []*model.Contract{
		{ ID: elemID },
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	contract := testContract(t)

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

	err = repo.Edit(contract)
	if err != nil {
		t.Fatal(err)
	}

	if contract.ID != 1 {
		t.Errorf("bad id: want %v, have %v", contract.ID, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestContractRepository_ListCompany(t *testing.T) {
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
		NewRows([]string{"id", "responseId", "companyId", "freelancerId", "startTime", "endTime", "status",
			"grade", "paymentAmount" })

	c1 := testContract(t)
	c2 := testContract(t)
	c3 := testContract(t)

	expect := []*model.Contract{
		c1,
		c2,
		c3,
	}

	// company mode
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.ResponseID, item.CompanyID, item.FreelancerID, item.StartTime,
			item.EndTime, item.Status, item.Grade, item.PaymentAmount)
	}

	mock.
		ExpectQuery("SELECT id, responseId, companyId, freelancerId, startTime, endTime, status, " +
			"grade, paymentAmount FROM contracts WHERE").
		WithArgs(companyId).
		WillReturnRows(rows)

	repo := NewContractRepository(db)

	contracts, err := repo.List(companyId, "company")
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	for i := 0; i < 3; i++ {
		if !IsEqual(t, &contracts[i], expect[i]) {
			t.Errorf("results not match, want %v, have %v", expect[i], contracts[i])
			return
		}
	}

	// query error
	mock.
		ExpectQuery("SELECT id, responseId, companyId, freelancerId, startTime, endTime, status, " +
			"grade, paymentAmount FROM contracts WHERE").
		WithArgs(companyId).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.List(companyId, "company")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestContractRepository_ListFreelancer(t *testing.T) {
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
		NewRows([]string{"id", "responseId", "companyId", "freelancerId", "startTime", "endTime", "status",
			"grade", "paymentAmount" })

	c1 := testContract(t)
	c2 := testContract(t)
	c3 := testContract(t)

	expect := []*model.Contract{
		c1,
		c2,
		c3,
	}

	repo := NewContractRepository(db)

	// freelancer mode
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.ResponseID, item.CompanyID, item.FreelancerID, item.StartTime,
			item.EndTime, item.Status, item.Grade, item.PaymentAmount)
	}

	mock.
		ExpectQuery("SELECT id, responseId, companyId, freelancerId, startTime, endTime, status, " +
			"grade, paymentAmount FROM contracts WHERE").
		WithArgs(freelancerId).
		WillReturnRows(rows)

	contracts, err := repo.List(freelancerId, "freelancer")
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	for i := 0; i < 3; i++ {
		if !IsEqual(t, &contracts[i], expect[i]) {
			t.Errorf("results not match, want %v, have %v", expect[i], contracts[i])
			return
		}
	}

	// query error
	mock.
		ExpectQuery("SELECT id, responseId, companyId, freelancerId, startTime, endTime, status, " +
			"grade, paymentAmount FROM contracts WHERE").
		WithArgs(freelancerId).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.List(freelancerId, "freelancer")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}