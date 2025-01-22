package repository

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aldipi/loan-service/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetLoanProductByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	rate := decimal.NewFromFloat(6.5)
	roi := decimal.NewFromFloat(3.0)
	createdAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")
	lastUpdatedAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")

	rows := sqlmock.NewRows([]string{"id", "name", "rate", "roi", "created_at", "last_updated_at"}).
		AddRow(1, "Loan Product 1", rate, roi, createdAt, lastUpdatedAt)

	mock.ExpectQuery("SELECT id, name, rate, roi, created_at, last_updated_at FROM loan_products WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	loanProduct, err := repo.GetLoanProductByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, loanProduct)
	assert.True(t, reflect.DeepEqual(loanProduct, &model.LoanProduct{
		ID:            1,
		Name:          "Loan Product 1",
		Rate:          rate,
		ROI:           roi,
		CreatedAt:     createdAt,
		LastUpdatedAt: lastUpdatedAt,
	}))
}
