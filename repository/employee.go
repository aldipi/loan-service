package repository

import (
	"context"

	"github.com/aldipi/loan-service/model"
)

func (r *LoanRepository) GetEmployeeByID(ctx context.Context, id int64) (*model.Employee, error) {
	query := `
		SELECT
			id, name, created_at, last_updated_at
		FROM
			employees
		WHERE
			id = ?
	`

	row := r.DB.QueryRowContext(ctx, query, id)

	employee := &model.Employee{}
	err := row.Scan(
		&employee.ID,
		&employee.Name,
		&employee.CreatedAt,
		&employee.LastUpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return employee, nil
}
