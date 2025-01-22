package usecase

import (
	"context"

	"github.com/aldipi/loan-service/model"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetLoanByID(ctx context.Context, id int64) (*model.Loan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Loan), args.Error(1)
}

func (m *MockRepository) GetLoans(ctx context.Context, limit int, offset int) ([]*model.Loan, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*model.Loan), args.Error(1)
}

func (m *MockRepository) GetLoansByBorrowerID(ctx context.Context, borrowerID int64, limit int, offset int) ([]*model.Loan, error) {
	args := m.Called(ctx, borrowerID, limit, offset)
	return args.Get(0).([]*model.Loan), args.Error(1)
}

func (m *MockRepository) CreateLoan(ctx context.Context, loan *model.Loan) (int64, error) {
	args := m.Called(ctx, loan)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) UpdateLoan(ctx context.Context, loan *model.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockRepository) GetLoanProductByID(ctx context.Context, id int64) (*model.LoanProduct, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LoanProduct), args.Error(1)
}

func (m *MockRepository) GetInvestmentsByLoanID(ctx context.Context, loanID int64) ([]*model.Investment, error) {
	args := m.Called(ctx, loanID)
	return args.Get(0).([]*model.Investment), args.Error(1)
}

func (m *MockRepository) GetInvestmentsByInvestorID(ctx context.Context, investorID int64, limit int, offset int) ([]*model.Investment, error) {
	args := m.Called(ctx, investorID, limit, offset)
	return args.Get(0).([]*model.Investment), args.Error(1)
}

func (m *MockRepository) CreateInvestment(ctx context.Context, investment *model.Investment) (int64, error) {
	args := m.Called(ctx, investment)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepository) GetEmployeeByID(ctx context.Context, id int64) (*model.Employee, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Employee), args.Error(1)
}

func (m *MockRepository) GetInvestorByID(ctx context.Context, id int64) (*model.Investor, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Investor), args.Error(1)
}
