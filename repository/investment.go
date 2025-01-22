package repository

import (
	"context"

	"github.com/aldipi/loan-service/model"
)

func (r *LoanRepository) GetInvestmentsByLoanID(ctx context.Context, loanID int64) ([]*model.Investment, error) {
	query := `
		SELECT
			id, loan_id, investor_id, amount, agreement_letter, created_at, last_updated_at
		FROM
			investments
		WHERE
			loan_id = ?
	`

	rows, err := r.DB.QueryContext(ctx, query, loanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	investments := []*model.Investment{}
	for rows.Next() {
		investment := &model.Investment{}
		err = rows.Scan(
			&investment.ID,
			&investment.LoanID,
			&investment.InvestorID,
			&investment.Amount,
			&investment.AgreementLetter,
			&investment.CreatedAt,
			&investment.LastUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		investments = append(investments, investment)
	}

	return investments, nil
}

func (r *LoanRepository) GetInvestmentsByInvestorID(ctx context.Context, investorID int64, limit int, offset int) ([]*model.Investment, error) {
	query := `
		SELECT
			id, loan_id, investor_id, amount, agreement_letter, created_at, last_updated_at
		FROM
			investments
		WHERE
			investor_id = ?
		LIMIT ? OFFSET ?
	`

	rows, err := r.DB.QueryContext(ctx, query, investorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	investments := []*model.Investment{}
	for rows.Next() {
		investment := &model.Investment{}
		err = rows.Scan(
			&investment.ID,
			&investment.LoanID,
			&investment.InvestorID,
			&investment.Amount,
			&investment.AgreementLetter,
			&investment.CreatedAt,
			&investment.LastUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		investments = append(investments, investment)
	}

	return investments, nil
}

func (r *LoanRepository) CreateInvestment(ctx context.Context, investment *model.Investment) (id int64, err error) {
	query := `
		INSERT INTO investments (loan_id, investor_id, amount, agreement_letter)
		VALUES (?, ?, ?, ?)
	`

	res, err := r.DB.ExecContext(ctx, query,
		investment.LoanID,
		investment.InvestorID,
		investment.Amount,
		investment.AgreementLetter,
	)

	if err != nil {
		return
	}

	return res.LastInsertId()
}
