package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/aldipi/loan-service/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetInvestmentsByInvestorID(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	repo.On("GetInvestmentsByInvestorID", mock.Anything, int64(100), 10, 0).Return([]*model.Investment{
		{
			ID:         1,
			InvestorID: 100,
			LoanID:     1,
			Amount:     1000,
		},
		{
			ID:         2,
			InvestorID: 100,
			LoanID:     2,
			Amount:     2000,
		},
	}, nil)

	investments, err := uc.GetInvestmentsByInvestorID(context.Background(), 100, 10, 0)

	assert.NoError(t, err)
	assert.Len(t, investments, 2)
	assert.True(t, reflect.DeepEqual(investments[0], &model.Investment{
		ID:         1,
		InvestorID: 100,
		LoanID:     1,
		Amount:     1000,
	}))
	assert.True(t, reflect.DeepEqual(investments[1], &model.Investment{
		ID:         2,
		InvestorID: 100,
		LoanID:     2,
		Amount:     2000,
	}))
}

func TestGetInvestmentsByInvestorIDReturnEmpty(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	repo.On("GetInvestmentsByInvestorID", mock.Anything, int64(100), 10, 0).Return([]*model.Investment{}, nil)

	investments, err := uc.GetInvestmentsByInvestorID(context.Background(), 100, 10, 0)

	assert.NoError(t, err)
	assert.Len(t, investments, 0)
}

func TestCheckAvailableInvestmentByLoanID(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	loan := &model.Loan{
		ID:              1,
		PrincipalAmount: 1000,
		State:           model.LoanStateApproved,
	}

	investments := []*model.Investment{
		{
			ID:         1,
			InvestorID: 100,
			LoanID:     1,
			Amount:     100,
		},
		{
			ID:         2,
			InvestorID: 101,
			LoanID:     1,
			Amount:     200,
		},
	}

	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(loan, nil)
	repo.On("GetInvestmentsByLoanID", mock.Anything, int64(1)).Return(investments, nil)

	availableAmount, err := uc.CheckAvailableInvestmentByLoanID(context.Background(), 1)

	// Expected calculation
	totalInvested := 0
	for _, investment := range investments {
		totalInvested += investment.Amount
	}
	expectedAmount := loan.PrincipalAmount - totalInvested

	assert.NoError(t, err)
	assert.Equal(t, expectedAmount, availableAmount)
}

func TestCheckAvailableInvestmentByLoanIDLoanNotFound(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(nil, model.ErrLoanNotFound)

	_, err := uc.CheckAvailableInvestmentByLoanID(context.Background(), 1)

	assert.ErrorIs(t, err, model.ErrLoanNotFound)
}

func TestCheckAvailableInvestmentByLoanIDLoanNotApproved(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	loan := &model.Loan{
		ID:              1,
		PrincipalAmount: 100000,
		State:           model.LoanStateProposed,
	}

	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(loan, nil)

	_, err := uc.CheckAvailableInvestmentByLoanID(context.Background(), 1)

	assert.ErrorIs(t, err, model.ErrLoanNotApproved)
}

func TestCreateInvestment(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	investor := &model.Investor{
		ID: 100,
	}

	loan := &model.Loan{
		ID:              1,
		PrincipalAmount: 1000000,
		State:           model.LoanStateApproved,
	}

	investments := []*model.Investment{
		{
			ID:         1,
			InvestorID: 101,
			LoanID:     1,
			Amount:     200000,
		},
		{
			ID:         2,
			InvestorID: 102,
			LoanID:     1,
			Amount:     300000,
		},
	}

	repo.On("GetInvestorByID", mock.Anything, int64(100)).Return(investor, nil)
	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(loan, nil)
	repo.On("GetInvestmentsByLoanID", mock.Anything, int64(1)).Return(investments, nil)
	repo.On("CreateInvestment", mock.Anything, mock.Anything).Return(int64(3), nil)
	repo.On("UpdateLoan", mock.Anything, loan).Return(nil)

	investment, err := uc.CreateInvestment(context.Background(), 100, 1, 500000)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), investment.ID)
	assert.Equal(t, int64(100), investment.InvestorID)
	assert.Equal(t, int64(1), investment.LoanID)
	assert.Equal(t, 500000, investment.Amount)
	assert.Equal(t, model.LoanStateInvested, loan.State)
}

func TestCreateInvestmentInvestorNotFound(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	repo.On("GetInvestorByID", mock.Anything, int64(100)).Return(nil, model.ErrInvestorNotFound)

	_, err := uc.CreateInvestment(context.Background(), 100, 1, 500000)

	assert.ErrorIs(t, err, model.ErrInvestorNotFound)
}

func TestCreateInvestmentLoanNotFound(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	investor := &model.Investor{
		ID: 100,
	}

	repo.On("GetInvestorByID", mock.Anything, int64(100)).Return(investor, nil)
	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(nil, model.ErrLoanNotFound)

	_, err := uc.CreateInvestment(context.Background(), 100, 1, 500000)

	assert.ErrorIs(t, err, model.ErrLoanNotFound)
}

func TestCreateInvestmentLoanNotApproved(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	investor := &model.Investor{
		ID: 100,
	}

	loan := &model.Loan{
		ID:              1,
		PrincipalAmount: 1000000,
		State:           model.LoanStateProposed,
	}

	repo.On("GetInvestorByID", mock.Anything, int64(100)).Return(investor, nil)
	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(loan, nil)

	_, err := uc.CreateInvestment(context.Background(), 100, 1, 500000)

	assert.ErrorIs(t, err, model.ErrLoanNotApproved)
}

func TestCreateInvestmentInvalidAmount(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	investor := &model.Investor{
		ID: 100,
	}

	loan := &model.Loan{
		ID:              1,
		PrincipalAmount: 1000000,
		State:           model.LoanStateApproved,
	}

	investments := []*model.Investment{
		{
			ID:         1,
			InvestorID: 101,
			LoanID:     1,
			Amount:     200000,
		},
		{
			ID:         2,
			InvestorID: 102,
			LoanID:     1,
			Amount:     300000,
		},
	}

	repo.On("GetInvestorByID", mock.Anything, int64(100)).Return(investor, nil)
	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(loan, nil)
	repo.On("GetInvestmentsByLoanID", mock.Anything, int64(1)).Return(investments, nil)

	_, err := uc.CreateInvestment(context.Background(), 100, 1, 700000)

	assert.ErrorIs(t, err, model.ErrInvestmentInvalidAmount)
}

func TestCreateInvestmentDoNotUpdateLoanState(t *testing.T) {
	repo := new(MockRepository)
	uc := NewLoanUsecase(repo)

	investor := &model.Investor{
		ID: 100,
	}

	loan := &model.Loan{
		ID:              1,
		PrincipalAmount: 1000000,
		State:           model.LoanStateApproved,
	}

	investments := []*model.Investment{
		{
			ID:         1,
			InvestorID: 101,
			LoanID:     1,
			Amount:     200000,
		},
		{
			ID:         2,
			InvestorID: 102,
			LoanID:     1,
			Amount:     300000,
		},
	}

	repo.On("GetInvestorByID", mock.Anything, int64(100)).Return(investor, nil)
	repo.On("GetLoanByID", mock.Anything, int64(1)).Return(loan, nil)
	repo.On("GetInvestmentsByLoanID", mock.Anything, int64(1)).Return(investments, nil)
	repo.On("CreateInvestment", mock.Anything, mock.Anything).Return(int64(3), nil)
	repo.On("UpdateLoan", mock.Anything, loan).Return(nil)

	investment, err := uc.CreateInvestment(context.Background(), 100, 1, 100000)

	assert.NoError(t, err)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), investment.ID)
	assert.Equal(t, int64(100), investment.InvestorID)
	assert.Equal(t, int64(1), investment.LoanID)
	assert.Equal(t, 100000, investment.Amount)
	assert.Equal(t, model.LoanStateApproved, loan.State)
}
