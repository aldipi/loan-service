package repository

import (
	"context"
	"database/sql"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aldipi/loan-service/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetLoanByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	rate := decimal.NewFromFloat(6.5)
	roi := decimal.NewFromFloat(3.0)
	createdAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")
	approvedAt := sql.NullTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	investedAt := sql.NullTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	disbursedAt := sql.NullTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	lastUpdatedAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")
	approvalProof := sql.NullString{String: "https://file.io/123/approval_proof.jpg", Valid: true}
	agreementLetter := sql.NullString{String: "https://file.io/123/agreement_letter.pdf", Valid: true}

	rows := sqlmock.NewRows([]string{"id", "state", "borrower_id", "principal_amount", "rate", "roi", "approval_proof", "approved_by", "agreement_letter", "disbursed_by", "created_at", "approved_at", "invested_at", "disbursed_at", "last_updated_at"}).
		AddRow(1, model.LoanStateDisbursed, 123, 1000000, rate, roi, approvalProof, 555, agreementLetter, 777, createdAt, approvedAt, investedAt, disbursedAt, lastUpdatedAt)

	mock.ExpectQuery("SELECT id, state, borrower_id, principal_amount, rate, roi, approval_proof, approved_by, agreement_letter, disbursed_by, created_at, approved_at, invested_at, disbursed_at, last_updated_at FROM loans WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	loan, err := repo.GetLoanByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, loan)
	assert.True(t, reflect.DeepEqual(loan, &model.Loan{
		ID:              1,
		State:           model.LoanStateDisbursed,
		BorrowerID:      123,
		PrincipalAmount: 1000000,
		Rate:            rate,
		ROI:             roi,
		ApprovalProof:   approvalProof,
		ApprovedBy:      sql.NullInt64{Int64: 555, Valid: true},
		AgreementLetter: agreementLetter,
		DisbursedBy:     sql.NullInt64{Int64: 777, Valid: true},
		CreatedAt:       createdAt,
		ApprovedAt:      approvedAt,
		InvestedAt:      investedAt,
		DisbursedAt:     disbursedAt,
		LastUpdatedAt:   lastUpdatedAt,
	}))
}

func TestGetLoanByIDReturnEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	rows := sqlmock.NewRows([]string{"id", "state", "borrower_id", "principal_amount", "rate", "roi", "approval_proof", "approved_by", "agreement_letter", "disbursed_by", "created_at", "approved_at", "invested_at", "disbursed_at", "last_updated_at"})

	mock.ExpectQuery("SELECT id, state, borrower_id, principal_amount, rate, roi, approval_proof, approved_by, agreement_letter, disbursed_by, created_at, approved_at, invested_at, disbursed_at, last_updated_at FROM loans WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	loan, err := repo.GetLoanByID(context.Background(), 1)

	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Nil(t, loan)
}

func TestGetLoans(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	rate := decimal.NewFromFloat(6.5)
	roi := decimal.NewFromFloat(3.0)
	createdAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")
	approvedAt := sql.NullTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	investedAt := sql.NullTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	disbursedAt := sql.NullTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	lastUpdatedAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")
	approvalProof := sql.NullString{String: "https://file.io/123/approval_proof.jpg", Valid: true}
	agreementLetter := sql.NullString{String: "https://file.io/123/agreement_letter.pdf", Valid: true}

	rows := sqlmock.NewRows([]string{"id", "state", "borrower_id", "principal_amount", "rate", "roi", "approval_proof", "approved_by", "agreement_letter", "disbursed_by", "created_at", "approved_at", "invested_at", "disbursed_at", "last_updated_at"}).
		AddRow(1, model.LoanStateDisbursed, 123, 1000000, rate, roi, approvalProof, 555, agreementLetter, 777, createdAt, approvedAt, investedAt, disbursedAt, lastUpdatedAt).
		AddRow(2, model.LoanStateApproved, 456, 2000000, rate, roi, approvalProof, 333, nil, nil, createdAt, approvedAt, nil, nil, lastUpdatedAt)

	query := regexp.QuoteMeta(`
		SELECT
			id, state, borrower_id, principal_amount, rate, roi,
			approval_proof, approved_by, agreement_letter, disbursed_by,
			created_at, approved_at, invested_at, disbursed_at, last_updated_at
		FROM
			loans
		LIMIT ? OFFSET ?
	`)

	mock.ExpectQuery(query).
		WithArgs(10, 0).
		WillReturnRows(rows)

	loans, err := repo.GetLoans(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.NotEmpty(t, loans)
	assert.Len(t, loans, 2)
	assert.True(t, reflect.DeepEqual(loans[0], &model.Loan{
		ID:              1,
		State:           model.LoanStateDisbursed,
		BorrowerID:      123,
		PrincipalAmount: 1000000,
		Rate:            rate,
		ROI:             roi,
		ApprovalProof:   approvalProof,
		ApprovedBy:      sql.NullInt64{Int64: 555, Valid: true},
		AgreementLetter: agreementLetter,
		DisbursedBy:     sql.NullInt64{Int64: 777, Valid: true},
		CreatedAt:       createdAt,
		ApprovedAt:      approvedAt,
		InvestedAt:      investedAt,
		DisbursedAt:     disbursedAt,
		LastUpdatedAt:   lastUpdatedAt,
	}))
	assert.True(t, reflect.DeepEqual(loans[1], &model.Loan{
		ID:              2,
		State:           model.LoanStateApproved,
		BorrowerID:      456,
		PrincipalAmount: 2000000,
		Rate:            rate,
		ROI:             roi,
		ApprovalProof:   approvalProof,
		ApprovedBy:      sql.NullInt64{Int64: 333, Valid: true},
		AgreementLetter: sql.NullString{Valid: false},
		DisbursedBy:     sql.NullInt64{Valid: false},
		CreatedAt:       createdAt,
		ApprovedAt:      approvedAt,
		InvestedAt:      sql.NullTime{Valid: false},
		DisbursedAt:     sql.NullTime{Valid: false},
		LastUpdatedAt:   lastUpdatedAt,
	}))
}

func TestGetLoansReturnEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	rows := sqlmock.NewRows([]string{"id", "state", "borrower_id", "principal_amount", "rate", "roi", "approval_proof", "approved_by", "agreement_letter", "disbursed_by", "created_at", "approved_at", "invested_at", "disbursed_at", "last_updated_at"})

	query := regexp.QuoteMeta(`
		SELECT
			id, state, borrower_id, principal_amount, rate, roi,
			approval_proof, approved_by, agreement_letter, disbursed_by,
			created_at, approved_at, invested_at, disbursed_at, last_updated_at
		FROM
			loans
		LIMIT ? OFFSET ?
	`)

	mock.ExpectQuery(query).
		WithArgs(10, 0).
		WillReturnRows(rows)

	loans, err := repo.GetLoans(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Empty(t, loans)
}

func TestGetLoansByBorrowerID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	rate := decimal.NewFromFloat(6.5)
	roi := decimal.NewFromFloat(3.0)
	createdAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")
	approvedAt := sql.NullTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	investedAt := sql.NullTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	disbursedAt := sql.NullTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	lastUpdatedAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")
	approvalProof := sql.NullString{String: "https://file.io/123/approval_proof.jpg", Valid: true}
	agreementLetter := sql.NullString{String: "https://file.io/123/agreement_letter.pdf", Valid: true}

	rows := sqlmock.NewRows([]string{"id", "state", "borrower_id", "principal_amount", "rate", "roi", "approval_proof", "approved_by", "agreement_letter", "disbursed_by", "created_at", "approved_at", "invested_at", "disbursed_at", "last_updated_at"}).
		AddRow(1, model.LoanStateDisbursed, 123, 1000000, rate, roi, approvalProof, 555, agreementLetter, 777, createdAt, approvedAt, investedAt, disbursedAt, lastUpdatedAt).
		AddRow(2, model.LoanStateApproved, 123, 2000000, rate, roi, approvalProof, 333, nil, nil, createdAt, approvedAt, nil, nil, lastUpdatedAt)

	query := regexp.QuoteMeta(`
		SELECT
			id, state, borrower_id, principal_amount, rate, roi,
			approval_proof, approved_by, agreement_letter, disbursed_by,
			created_at, approved_at, invested_at, disbursed_at, last_updated_at
		FROM
			loans
		WHERE
			borrower_id = ?
		LIMIT ? OFFSET ?
	`)

	mock.ExpectQuery(query).
		WithArgs(123, 10, 0).
		WillReturnRows(rows)

	loans, err := repo.GetLoansByBorrowerID(context.Background(), 123, 10, 0)

	assert.NoError(t, err)
	assert.NotEmpty(t, loans)
	assert.Len(t, loans, 2)
	assert.True(t, reflect.DeepEqual(loans[0], &model.Loan{
		ID:              1,
		State:           model.LoanStateDisbursed,
		BorrowerID:      123,
		PrincipalAmount: 1000000,
		Rate:            rate,
		ROI:             roi,
		ApprovalProof:   approvalProof,
		ApprovedBy:      sql.NullInt64{Int64: 555, Valid: true},
		AgreementLetter: agreementLetter,
		DisbursedBy:     sql.NullInt64{Int64: 777, Valid: true},
		CreatedAt:       createdAt,
		ApprovedAt:      approvedAt,
		InvestedAt:      investedAt,
		DisbursedAt:     disbursedAt,
		LastUpdatedAt:   lastUpdatedAt,
	}))
	assert.True(t, reflect.DeepEqual(loans[1], &model.Loan{
		ID:              2,
		State:           model.LoanStateApproved,
		BorrowerID:      123,
		PrincipalAmount: 2000000,
		Rate:            rate,
		ROI:             roi,
		ApprovalProof:   approvalProof,
		ApprovedBy:      sql.NullInt64{Int64: 333, Valid: true},
		AgreementLetter: sql.NullString{Valid: false},
		DisbursedBy:     sql.NullInt64{Valid: false},
		CreatedAt:       createdAt,
		ApprovedAt:      approvedAt,
		InvestedAt:      sql.NullTime{Valid: false},
		DisbursedAt:     sql.NullTime{Valid: false},
		LastUpdatedAt:   lastUpdatedAt,
	}))
}

func TestGetLoansByBorrowerIDReturnEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	rows := sqlmock.NewRows([]string{"id", "state", "borrower_id", "principal_amount", "rate", "roi", "approval_proof", "approved_by", "agreement_letter", "disbursed_by", "created_at", "approved_at", "invested_at", "disbursed_at", "last_updated_at"})

	query := regexp.QuoteMeta(`
		SELECT
			id, state, borrower_id, principal_amount, rate, roi,
			approval_proof, approved_by, agreement_letter, disbursed_by,
			created_at, approved_at, invested_at, disbursed_at, last_updated_at
		FROM
			loans
		WHERE
			borrower_id = ?
		LIMIT ? OFFSET ?
	`)

	mock.ExpectQuery(query).
		WithArgs(123, 10, 0).
		WillReturnRows(rows)

	loans, err := repo.GetLoansByBorrowerID(context.Background(), 123, 10, 0)

	assert.NoError(t, err)
	assert.Empty(t, loans)
}

func TestCreateLoan(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	rate := decimal.NewFromFloat(6.5)
	roi := decimal.NewFromFloat(3.0)
	createdAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")
	lastUpdatedAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")

	query := regexp.QuoteMeta(`
        INSERT INTO loans (
            state, borrower_id, principal_amount, rate, roi
        ) VALUES (
            ?, ?, ?, ?, ?
        )
    `)

	mock.ExpectExec(query).
		WithArgs(model.LoanStateProposed, 123, 1000000, rate, roi).
		WillReturnResult(sqlmock.NewResult(1, 1))

	loan := &model.Loan{
		State:           model.LoanStateProposed,
		BorrowerID:      123,
		PrincipalAmount: 1000000,
		Rate:            rate,
		ROI:             roi,
		CreatedAt:       createdAt,
		LastUpdatedAt:   lastUpdatedAt,
	}

	id, err := repo.CreateLoan(context.Background(), loan)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func TestUpdateLoan(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	approvedAt := sql.NullTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	investedAt := sql.NullTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	disbursedAt := sql.NullTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	approvalProof := sql.NullString{String: "https://file.io/123/approval_proof.jpg", Valid: true}
	agreementLetter := sql.NullString{String: "https://file.io/123/agreement_letter.pdf", Valid: true}

	query := regexp.QuoteMeta(`
		UPDATE loans
		SET state = ?,
			approval_proof = ?,
			approved_by = ?,
			agreement_letter = ?,
			disbursed_by = ?,
			approved_at = ?,
			invested_at = ?,
			disbursed_at = ?,
			last_updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`)

	mock.ExpectExec(query).
		WithArgs(model.LoanStateDisbursed, approvalProof, 333, agreementLetter, 555, approvedAt, investedAt, disbursedAt, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	loan := &model.Loan{
		ID:              1,
		State:           model.LoanStateDisbursed,
		ApprovalProof:   approvalProof,
		ApprovedBy:      sql.NullInt64{Int64: 333, Valid: true},
		AgreementLetter: agreementLetter,
		DisbursedBy:     sql.NullInt64{Int64: 555, Valid: true},
		ApprovedAt:      approvedAt,
		InvestedAt:      investedAt,
		DisbursedAt:     disbursedAt,
	}

	err = repo.UpdateLoan(context.Background(), loan)

	assert.NoError(t, err)
}
