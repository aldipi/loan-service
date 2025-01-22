package repository

import (
	"context"

	"github.com/aldipi/loan-service/model"
)

func (r *LoanRepository) GetLoanProductByID(ctx context.Context, id int64) (*model.LoanProduct, error) {
	query := `
		SELECT
			id, name, rate, roi, created_at, last_updated_at
		FROM
			loan_products
		WHERE
			id = ?
	`

	row := r.DB.QueryRowContext(ctx, query, id)

	loanProduct := &model.LoanProduct{}
	err := row.Scan(
		&loanProduct.ID,
		&loanProduct.Name,
		&loanProduct.Rate,
		&loanProduct.ROI,
		&loanProduct.CreatedAt,
		&loanProduct.LastUpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return loanProduct, nil
}
