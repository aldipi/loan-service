package repository

import (
	"context"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aldipi/loan-service/model"
	"github.com/stretchr/testify/assert"
)

func TestGetInvestmentsByLoanID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	createdAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")
	lastUpdatedAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")

	rows := sqlmock.NewRows([]string{"id", "loan_id", "investor_id", "amount", "agreement_letter", "created_at", "last_updated_at"}).
		AddRow(1, 1, 123, 200000, "https://file.io/123/agreement_letter.pdf", createdAt, lastUpdatedAt).
		AddRow(2, 1, 456, 300000, "https://file.io/456/agreement_letter.pdf", createdAt, lastUpdatedAt)

	query := regexp.QuoteMeta(`
		SELECT
			id, loan_id, investor_id, amount, agreement_letter, created_at, last_updated_at
		FROM
			investments
		WHERE
			loan_id = ?
	`)

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	investments, err := repo.GetInvestmentsByLoanID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotEmpty(t, investments)
	assert.Len(t, investments, 2)
	assert.True(t, reflect.DeepEqual(investments[0], &model.Investment{
		ID:              1,
		LoanID:          1,
		InvestorID:      123,
		Amount:          200000,
		AgreementLetter: "https://file.io/123/agreement_letter.pdf",
		CreatedAt:       createdAt,
		LastUpdatedAt:   lastUpdatedAt,
	}))
	assert.True(t, reflect.DeepEqual(investments[1], &model.Investment{
		ID:              2,
		LoanID:          1,
		InvestorID:      456,
		Amount:          300000,
		AgreementLetter: "https://file.io/456/agreement_letter.pdf",
		CreatedAt:       createdAt,
		LastUpdatedAt:   lastUpdatedAt,
	}))
}

func TestGetInvestmentsByLoanIDReturnEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	rows := sqlmock.NewRows([]string{"id", "loan_id", "investor_id", "amount", "agreement_letter", "created_at", "last_updated_at"})

	query := regexp.QuoteMeta(`
		SELECT
			id, loan_id, investor_id, amount, agreement_letter, created_at, last_updated_at
		FROM
			investments
		WHERE
			loan_id = ?
	`)

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	investments, err := repo.GetInvestmentsByLoanID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Empty(t, investments)
}

func TestGetInvestmentsByInvestorID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	createdAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")
	lastUpdatedAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")

	rows := sqlmock.NewRows([]string{"id", "loan_id", "investor_id", "amount", "agreement_letter", "created_at", "last_updated_at"}).
		AddRow(1, 1, 123, 200000, "https://file.io/123/agreement_letter.pdf", createdAt, lastUpdatedAt).
		AddRow(5, 2, 123, 300000, "https://file.io/456/agreement_letter.pdf", createdAt, lastUpdatedAt)

	query := regexp.QuoteMeta(`
		SELECT
			id, loan_id, investor_id, amount, agreement_letter, created_at, last_updated_at
		FROM
			investments
		WHERE
			investor_id = ?
		LIMIT ? OFFSET ?
	`)

	mock.ExpectQuery(query).WithArgs(123, 10, 0).WillReturnRows(rows)

	investments, err := repo.GetInvestmentsByInvestorID(context.Background(), 123, 10, 0)

	assert.NoError(t, err)
	assert.NotEmpty(t, investments)
	assert.Len(t, investments, 2)
	assert.True(t, reflect.DeepEqual(investments[0], &model.Investment{
		ID:              1,
		LoanID:          1,
		InvestorID:      123,
		Amount:          200000,
		AgreementLetter: "https://file.io/123/agreement_letter.pdf",
		CreatedAt:       createdAt,
		LastUpdatedAt:   lastUpdatedAt,
	}))
	assert.True(t, reflect.DeepEqual(investments[1], &model.Investment{
		ID:              5,
		LoanID:          2,
		InvestorID:      123,
		Amount:          300000,
		AgreementLetter: "https://file.io/456/agreement_letter.pdf",
		CreatedAt:       createdAt,
		LastUpdatedAt:   lastUpdatedAt,
	}))
}

func TestGetInvestmentsByInvestorIDReturnEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	rows := sqlmock.NewRows([]string{"id", "loan_id", "investor_id", "amount", "agreement_letter", "created_at", "last_updated_at"})

	query := regexp.QuoteMeta(`
		SELECT
			id, loan_id, investor_id, amount, agreement_letter, created_at, last_updated_at
		FROM
			investments
		WHERE
			investor_id = ?
		LIMIT ? OFFSET ?
	`)

	mock.ExpectQuery(query).WithArgs(123, 10, 0).WillReturnRows(rows)

	investments, err := repo.GetInvestmentsByInvestorID(context.Background(), 123, 10, 0)

	assert.NoError(t, err)
	assert.Empty(t, investments)
}

func TestCreateInvestment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	query := regexp.QuoteMeta(`
		INSERT INTO investments
			(loan_id, investor_id, amount, agreement_letter)
		VALUES
			(?, ?, ?, ?)
	`)

	mock.ExpectExec(query).
		WithArgs(1, 123, 200000, "https://file.io/123/agreement_letter.pdf").
		WillReturnResult(sqlmock.NewResult(1, 1))

	investment := &model.Investment{
		LoanID:          1,
		InvestorID:      123,
		Amount:          200000,
		AgreementLetter: "https://file.io/123/agreement_letter.pdf",
	}

	id, err := repo.CreateInvestment(context.Background(), investment)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}
