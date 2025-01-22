package repository

import (
	"context"

	"github.com/aldipi/loan-service/model"
)

func (r *LoanRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT
			id, name, created_at, last_updated_at
		FROM
			users
		WHERE
			id = ?
	`

	row := r.DB.QueryRowContext(ctx, query, id)

	user := &model.User{}
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.CreatedAt,
		&user.LastUpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
