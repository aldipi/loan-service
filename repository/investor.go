package repository

import (
	"context"

	"github.com/aldipi/loan-service/model"
)

func (r *LoanRepository) GetInvestorByID(ctx context.Context, id int64) (*model.Investor, error) {
	query := `
		SELECT
			id, name, created_at, last_updated_at
		FROM
			investors
		WHERE
			id = ?
	`

	row := r.DB.QueryRowContext(ctx, query, id)

	investor := &model.Investor{}
	err := row.Scan(
		&investor.ID,
		&investor.Name,
		&investor.CreatedAt,
		&investor.LastUpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return investor, nil
}
