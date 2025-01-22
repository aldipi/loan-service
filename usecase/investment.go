package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/aldipi/loan-service/model"
)

func (u *LoanUsecase) GetInvestmentsByInvestorID(ctx context.Context, investorID int64, limit int, offset int) ([]*model.Investment, error) {
	investments, err := u.repo.GetInvestmentsByInvestorID(ctx, investorID, limit, offset)
	if err != nil {
		return nil, err
	}

	return investments, nil
}

func (u *LoanUsecase) CheckAvailableInvestmentByLoanID(ctx context.Context, loanID int64) (int, error) {
	loan, err := u.repo.GetLoanByID(ctx, loanID)
	if err != nil {
		return 0, model.ErrLoanNotFound
	}

	if loan.State != model.LoanStateApproved {
		return 0, model.ErrLoanNotApproved
	}

	investments, err := u.repo.GetInvestmentsByLoanID(ctx, loanID)
	if err != nil {
		return 0, err
	}

	var totalInvested int
	for _, investment := range investments {
		totalInvested += investment.Amount
	}

	return loan.PrincipalAmount - totalInvested, nil
}

// NOTE: CreateInvestment function is not thread safe because it doesn't use transaction.
//
//	It will need to implement proper transaction handling before it can be used in production.
func (u *LoanUsecase) CreateInvestment(ctx context.Context, investorID int64, loanID int64, amount int) (investment *model.Investment, err error) {
	investor, err := u.repo.GetInvestorByID(ctx, investorID)
	if err != nil {
		return nil, model.ErrInvestorNotFound
	}

	loan, err := u.repo.GetLoanByID(ctx, loanID)
	if err != nil {
		return nil, model.ErrLoanNotFound
	}

	if loan.State != model.LoanStateApproved {
		return nil, model.ErrLoanNotApproved
	}

	investments, err := u.repo.GetInvestmentsByLoanID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	// Check available investment amount
	var totalInvested int
	for _, investment := range investments {
		totalInvested += investment.Amount
	}

	if amount > (loan.PrincipalAmount - totalInvested) {
		return nil, model.ErrInvestmentInvalidAmount
	}

	investment = &model.Investment{
		Amount:          amount,
		InvestorID:      investor.ID,
		LoanID:          loan.ID,
		AgreementLetter: generateAgreementLetter(),
	}

	investmentID, err := u.repo.CreateInvestment(ctx, investment)
	if err != nil {
		return nil, err
	}
	investment.ID = investmentID

	// Postprocess loan investment
	// NOTE: Again this is not thread safe, if error occurs after this line, the loan will be in invalid state.
	//       It will need to implement proper transaction and retry mechanism.
	totalInvested += amount
	if totalInvested == loan.PrincipalAmount {
		loan.State = model.LoanStateInvested
		loan.InvestedAt = sql.NullTime{Time: time.Now(), Valid: true}
		loan.LastUpdatedAt = time.Now()

		err = u.repo.UpdateLoan(ctx, loan)
		if err != nil {
			return nil, err
		}
	}

	return investment, nil
}

func generateAgreementLetter() string {
	return "http://file.io/investor/agreement.pdf"
}
