package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/aldipi/loan-service/model"
)

func (u *LoanUsecase) GetLoans(ctx context.Context, limit int, offset int) ([]*model.Loan, error) {
	loans, err := u.repo.GetLoans(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return loans, nil
}

func (u *LoanUsecase) GetLoansByBorrowerID(ctx context.Context, borrowerID int64, limit int, offset int) ([]*model.Loan, error) {
	loans, err := u.repo.GetLoansByBorrowerID(ctx, borrowerID, limit, offset)
	if err != nil {
		return nil, err
	}

	return loans, nil
}

func (u *LoanUsecase) CreateLoan(ctx context.Context, userID int64, loanProductID int64, amount int) (loan *model.Loan, err error) {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	loanProduct, err := u.repo.GetLoanProductByID(ctx, loanProductID)
	if err != nil {
		return nil, model.ErrLoanProductNotFound
	}

	loan = &model.Loan{
		State:           0,
		BorrowerID:      user.ID,
		PrincipalAmount: amount,
		Rate:            loanProduct.Rate,
		ROI:             loanProduct.ROI,
	}

	loanID, err := u.repo.CreateLoan(ctx, loan)
	if err != nil {
		return nil, err
	}

	loan.ID = loanID

	return loan, nil
}

func (u *LoanUsecase) ApproveLoan(ctx context.Context, loanID int64, employeeID int64, approvalProof string) error {
	loan, err := u.repo.GetLoanByID(ctx, loanID)
	if err != nil {
		return model.ErrLoanNotFound
	}

	employee, err := u.repo.GetEmployeeByID(ctx, employeeID)
	if err != nil {
		return model.ErrEmployeeNotFound
	}

	if loan.State != model.LoanStateProposed {
		return model.ErrLoanNotProposed
	}

	loan.State = model.LoanStateApproved
	loan.ApprovedBy = sql.NullInt64{Int64: employee.ID, Valid: true}
	loan.ApprovalProof = sql.NullString{String: approvalProof, Valid: true}
	loan.ApprovedAt = sql.NullTime{Time: time.Now(), Valid: true}
	loan.LastUpdatedAt = time.Now()

	err = u.repo.UpdateLoan(ctx, loan)
	if err != nil {
		return err
	}

	return nil
}

func (u *LoanUsecase) DisburseLoan(ctx context.Context, loanID int64, employeeID int64, agreementLetter string) error {
	loan, err := u.repo.GetLoanByID(ctx, loanID)
	if err != nil {
		return model.ErrLoanNotFound
	}

	employee, err := u.repo.GetEmployeeByID(ctx, employeeID)
	if err != nil {
		return model.ErrEmployeeNotFound
	}

	if loan.State != model.LoanStateInvested {
		return model.ErrLoanNotInvested
	}

	loan.State = model.LoanStateDisbursed
	loan.DisbursedBy = sql.NullInt64{Int64: employee.ID, Valid: true}
	loan.AgreementLetter = sql.NullString{String: agreementLetter, Valid: true}
	loan.DisbursedAt = sql.NullTime{Time: time.Now(), Valid: true}
	loan.LastUpdatedAt = time.Now()

	err = u.repo.UpdateLoan(ctx, loan)
	if err != nil {
		return err
	}

	return nil
}
