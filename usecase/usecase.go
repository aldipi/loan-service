package usecase

import (
	"context"

	"github.com/aldipi/loan-service/model"
)

type Repository interface {
	GetLoanByID(ctx context.Context, id int64) (*model.Loan, error)
	GetLoans(ctx context.Context, limit int, offset int) ([]*model.Loan, error)
	GetLoansByBorrowerID(ctx context.Context, borrowerID int64, limit int, offset int) ([]*model.Loan, error)
	CreateLoan(ctx context.Context, loan *model.Loan) (id int64, err error)
	UpdateLoan(ctx context.Context, loan *model.Loan) error

	GetLoanProductByID(ctx context.Context, id int64) (*model.LoanProduct, error)

	GetInvestmentsByLoanID(ctx context.Context, loanID int64) ([]*model.Investment, error)
	GetInvestmentsByInvestorID(ctx context.Context, investorID int64, limit int, offset int) ([]*model.Investment, error)
	CreateInvestment(ctx context.Context, investment *model.Investment) (id int64, err error)

	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	GetEmployeeByID(ctx context.Context, id int64) (*model.Employee, error)
	GetInvestorByID(ctx context.Context, id int64) (*model.Investor, error)
}

type LoanUsecase struct {
	repo Repository
}

func NewLoanUsecase(repo Repository) *LoanUsecase {
	return &LoanUsecase{repo: repo}
}
