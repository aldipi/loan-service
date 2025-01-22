package usecase

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/aldipi/loan-service/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetLoans(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	dummyTime, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")

	repo.On("GetLoans", mock.Anything, 10, 0).Return([]*model.Loan{
		{
			ID:              1,
			State:           model.LoanStateApproved,
			BorrowerID:      123,
			PrincipalAmount: 1000000,
			Rate:            decimal.NewFromFloat(10.0),
			ROI:             decimal.NewFromFloat(5.5),
			ApprovalProof:   sql.NullString{String: "https://file.io/123/proof.jpg", Valid: true},
			ApprovedBy:      sql.NullInt64{Int64: 456, Valid: true},
			CreatedAt:       dummyTime,
			ApprovedAt:      sql.NullTime{Time: dummyTime, Valid: true},
			LastUpdatedAt:   dummyTime,
		},
		{
			ID:              2,
			State:           model.LoanStateProposed,
			BorrowerID:      456,
			PrincipalAmount: 2000000,
			Rate:            decimal.NewFromFloat(12.0),
			ROI:             decimal.NewFromFloat(6.0),
			CreatedAt:       dummyTime,
			LastUpdatedAt:   dummyTime,
		},
	}, nil)

	loans, err := uc.GetLoans(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, loans, 2)
	assert.True(t, reflect.DeepEqual(loans, []*model.Loan{
		{
			ID:              1,
			State:           model.LoanStateApproved,
			BorrowerID:      123,
			PrincipalAmount: 1000000,
			Rate:            decimal.NewFromFloat(10.0),
			ROI:             decimal.NewFromFloat(5.5),
			ApprovalProof:   sql.NullString{String: "https://file.io/123/proof.jpg", Valid: true},
			ApprovedBy:      sql.NullInt64{Int64: 456, Valid: true},
			CreatedAt:       dummyTime,
			ApprovedAt:      sql.NullTime{Time: dummyTime, Valid: true},
			LastUpdatedAt:   dummyTime,
		},
		{
			ID:              2,
			State:           model.LoanStateProposed,
			BorrowerID:      456,
			PrincipalAmount: 2000000,
			Rate:            decimal.NewFromFloat(12.0),
			ROI:             decimal.NewFromFloat(6.0),
			CreatedAt:       dummyTime,
			LastUpdatedAt:   dummyTime,
		},
	}))
}

func TestGetLoansReturnEmpty(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	repo.On("GetLoans", mock.Anything, 10, 0).Return([]*model.Loan{}, nil)

	loans, err := uc.GetLoans(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, loans, 0)
}

func TestGetLoansByBorrowerID(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	dummyTime, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")

	repo.On("GetLoansByBorrowerID", mock.Anything, int64(123), 10, 0).Return([]*model.Loan{
		{
			ID:              1,
			State:           model.LoanStateApproved,
			BorrowerID:      123,
			PrincipalAmount: 1000000,
			Rate:            decimal.NewFromFloat(10.0),
			ROI:             decimal.NewFromFloat(5.5),
			ApprovalProof:   sql.NullString{String: "https://file.io/123/proof.jpg", Valid: true},
			ApprovedBy:      sql.NullInt64{Int64: 456, Valid: true},
			CreatedAt:       dummyTime,
			ApprovedAt:      sql.NullTime{Time: dummyTime, Valid: true},
			LastUpdatedAt:   dummyTime,
		},
		{
			ID:              2,
			State:           model.LoanStateProposed,
			BorrowerID:      123,
			PrincipalAmount: 2000000,
			Rate:            decimal.NewFromFloat(12.0),
			ROI:             decimal.NewFromFloat(6.0),
			CreatedAt:       dummyTime,
			LastUpdatedAt:   dummyTime,
		},
	}, nil)

	loans, err := uc.GetLoansByBorrowerID(context.Background(), 123, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, loans, 2)
	assert.True(t, reflect.DeepEqual(loans, []*model.Loan{
		{
			ID:              1,
			State:           model.LoanStateApproved,
			BorrowerID:      123,
			PrincipalAmount: 1000000,
			Rate:            decimal.NewFromFloat(10.0),
			ROI:             decimal.NewFromFloat(5.5),
			ApprovalProof:   sql.NullString{String: "https://file.io/123/proof.jpg", Valid: true},
			ApprovedBy:      sql.NullInt64{Int64: 456, Valid: true},
			CreatedAt:       dummyTime,
			ApprovedAt:      sql.NullTime{Time: dummyTime, Valid: true},
			LastUpdatedAt:   dummyTime,
		},
		{
			ID:              2,
			State:           model.LoanStateProposed,
			BorrowerID:      123,
			PrincipalAmount: 2000000,
			Rate:            decimal.NewFromFloat(12.0),
			ROI:             decimal.NewFromFloat(6.0),
			CreatedAt:       dummyTime,
			LastUpdatedAt:   dummyTime,
		},
	}))
}

func TestGetLoansByBorrowerIDReturnEmpty(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	repo.On("GetLoansByBorrowerID", mock.Anything, int64(123), 10, 0).Return([]*model.Loan{}, nil)

	loans, err := uc.GetLoansByBorrowerID(context.Background(), 123, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, loans, 0)
}

func TestCreateLoan(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	repo.On("GetUserByID", mock.Anything, int64(123)).Return(&model.User{ID: 123}, nil)
	repo.On("GetLoanProductByID", mock.Anything, int64(100)).Return(&model.LoanProduct{ID: 1, Rate: decimal.NewFromFloat(10.0), ROI: decimal.NewFromFloat(5.5)}, nil)
	repo.On("CreateLoan", mock.Anything, mock.Anything).Return(int64(1), nil)

	loan, err := uc.CreateLoan(context.Background(), int64(123), int64(100), 1000000)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), loan.ID)
	assert.Equal(t, int64(123), loan.BorrowerID)
	assert.Equal(t, 1000000, loan.PrincipalAmount)
	assert.Equal(t, decimal.NewFromFloat(10.0), loan.Rate)
	assert.Equal(t, decimal.NewFromFloat(5.5), loan.ROI)
}

func TestCreateLoanUserNotFound(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	repo.On("GetUserByID", mock.Anything, int64(123)).Return(nil, model.ErrUserNotFound)

	loan, err := uc.CreateLoan(context.Background(), int64(123), int64(100), 1000000)

	assert.Error(t, err)
	assert.Nil(t, loan)
}

func TestCreateLoanProductNotFound(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	repo.On("GetUserByID", mock.Anything, int64(123)).Return(&model.User{ID: 123}, nil)
	repo.On("GetLoanProductByID", mock.Anything, int64(100)).Return(nil, model.ErrLoanProductNotFound)

	loan, err := uc.CreateLoan(context.Background(), int64(123), int64(100), 1000000)

	assert.Error(t, err)
	assert.Nil(t, loan)
}

func TestApproveLoan(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	loan := &model.Loan{ID: 1, State: model.LoanStateProposed}

	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(loan, nil)
	repo.On("GetEmployeeByID", mock.Anything, int64(555)).Return(&model.Employee{ID: 555}, nil)
	repo.On("UpdateLoan", mock.Anything, mock.Anything).Return(nil)

	err := uc.ApproveLoan(context.Background(), int64(1), int64(555), "https://file.io/123/proof.jpg")

	assert.NoError(t, err)
	repo.AssertCalled(t, "GetLoanByID", mock.Anything, int64(1))
	repo.AssertCalled(t, "GetEmployeeByID", mock.Anything, int64(555))
	repo.AssertCalled(t, "UpdateLoan", mock.Anything, loan)
	assert.Equal(t, model.LoanStateApproved, loan.State)
	assert.Equal(t, int64(555), loan.ApprovedBy.Int64)
	assert.Equal(t, "https://file.io/123/proof.jpg", loan.ApprovalProof.String)
}

func TestApproveLoanNotFound(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(nil, model.ErrLoanNotFound)

	err := uc.ApproveLoan(context.Background(), int64(1), int64(555), "https://file.io/123/proof.jpg")

	assert.Error(t, err)
	assert.ErrorIs(t, err, model.ErrLoanNotFound)
}

func TestApproveLoanEmployeeNotFound(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	loan := &model.Loan{ID: 1, State: model.LoanStateProposed}

	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(loan, nil)
	repo.On("GetEmployeeByID", mock.Anything, int64(555)).Return(nil, model.ErrEmployeeNotFound)

	err := uc.ApproveLoan(context.Background(), int64(1), int64(555), "https://file.io/123/proof.jpg")

	assert.Error(t, err)
	assert.ErrorIs(t, err, model.ErrEmployeeNotFound)
}

func TestApproveLoanNotProposed(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	loan := &model.Loan{ID: 1, State: model.LoanStateApproved}

	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(loan, nil)
	repo.On("GetEmployeeByID", mock.Anything, int64(555)).Return(&model.Employee{ID: 555}, nil)
	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(loan, nil)

	err := uc.ApproveLoan(context.Background(), int64(1), int64(555), "https://file.io/123/proof.jpg")

	assert.Error(t, err)
	assert.ErrorIs(t, err, model.ErrLoanNotProposed)
}

func TestDisburseLoan(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	loan := &model.Loan{ID: 1, State: model.LoanStateInvested}

	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(loan, nil)
	repo.On("GetEmployeeByID", mock.Anything, int64(555)).Return(&model.Employee{ID: 555}, nil)
	repo.On("UpdateLoan", mock.Anything, mock.Anything).Return(nil)

	err := uc.DisburseLoan(context.Background(), int64(1), int64(555), "https://file.io/123/agreement.pdf")

	assert.NoError(t, err)
	repo.AssertCalled(t, "GetLoanByID", mock.Anything, int64(1))
	repo.AssertCalled(t, "GetEmployeeByID", mock.Anything, int64(555))
	repo.AssertCalled(t, "UpdateLoan", mock.Anything, loan)
	assert.Equal(t, model.LoanStateDisbursed, loan.State)
	assert.Equal(t, int64(555), loan.DisbursedBy.Int64)
	assert.Equal(t, "https://file.io/123/agreement.pdf", loan.AgreementLetter.String)
}

func TestDisburseLoanNotFound(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(nil, model.ErrLoanNotFound)

	err := uc.DisburseLoan(context.Background(), int64(1), int64(555), "https://file.io/123/agreement.pdf")

	assert.Error(t, err)
	assert.ErrorIs(t, err, model.ErrLoanNotFound)
}

func TestDisburseLoanEmployeeNotFound(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	loan := &model.Loan{ID: 1, State: model.LoanStateInvested}

	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(loan, nil)
	repo.On("GetEmployeeByID", mock.Anything, int64(555)).Return(nil, model.ErrEmployeeNotFound)

	err := uc.DisburseLoan(context.Background(), int64(1), int64(555), "https://file.io/123/agreement.pdf")

	assert.Error(t, err)
	assert.ErrorIs(t, err, model.ErrEmployeeNotFound)
}

func TestDisburseLoanNotInvested(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	loan := &model.Loan{ID: 1, State: model.LoanStateApproved}

	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(loan, nil)
	repo.On("GetEmployeeByID", mock.Anything, int64(555)).Return(&model.Employee{ID: 555}, nil)
	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(loan, nil)

	err := uc.DisburseLoan(context.Background(), int64(1), int64(555), "https://file.io/123/agreement.pdf")

	assert.Error(t, err)
	assert.ErrorIs(t, err, model.ErrLoanNotInvested)
}
