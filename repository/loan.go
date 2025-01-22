package repository

import (
	"context"

	"github.com/aldipi/loan-service/model"
)

func (r *LoanRepository) GetLoanByID(ctx context.Context, id int64) (*model.Loan, error) {
	query := `
		SELECT
			id, state, borrower_id, principal_amount, rate, roi,
			approval_proof, approved_by, agreement_letter, disbursed_by,
			created_at, approved_at, invested_at, disbursed_at, last_updated_at
		FROM
			loans
		WHERE
			id = ?
	`

	row := r.DB.QueryRowContext(ctx, query, id)

	loan := &model.Loan{}
	err := row.Scan(
		&loan.ID,
		&loan.State,
		&loan.BorrowerID,
		&loan.PrincipalAmount,
		&loan.Rate,
		&loan.ROI,
		&loan.ApprovalProof,
		&loan.ApprovedBy,
		&loan.AgreementLetter,
		&loan.DisbursedBy,
		&loan.CreatedAt,
		&loan.ApprovedAt,
		&loan.InvestedAt,
		&loan.DisbursedAt,
		&loan.LastUpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return loan, nil
}

func (r *LoanRepository) GetLoans(ctx context.Context, limit int, offset int) ([]*model.Loan, error) {
	query := `
		SELECT
			id, state, borrower_id, principal_amount, rate, roi,
			approval_proof, approved_by, agreement_letter, disbursed_by,
			created_at, approved_at, invested_at, disbursed_at, last_updated_at
		FROM
			loans
		LIMIT ? OFFSET ?
	`

	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loans := []*model.Loan{}
	for rows.Next() {
		loan := &model.Loan{}
		err := rows.Scan(
			&loan.ID,
			&loan.State,
			&loan.BorrowerID,
			&loan.PrincipalAmount,
			&loan.Rate,
			&loan.ROI,
			&loan.ApprovalProof,
			&loan.ApprovedBy,
			&loan.AgreementLetter,
			&loan.DisbursedBy,
			&loan.CreatedAt,
			&loan.ApprovedAt,
			&loan.InvestedAt,
			&loan.DisbursedAt,
			&loan.LastUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		loans = append(loans, loan)
	}

	return loans, nil
}

func (r *LoanRepository) GetLoansByBorrowerID(ctx context.Context, borrowerID int64, limit int, offset int) ([]*model.Loan, error) {
	query := `
		SELECT
			id, state, borrower_id, principal_amount, rate, roi,
			approval_proof, approved_by, agreement_letter, disbursed_by,
			created_at, approved_at, invested_at, disbursed_at, last_updated_at
		FROM
			loans
		WHERE
			borrower_id = ?
		LIMIT ? OFFSET ?
	`

	rows, err := r.DB.QueryContext(ctx, query, borrowerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loans := []*model.Loan{}
	for rows.Next() {
		loan := &model.Loan{}
		err := rows.Scan(
			&loan.ID,
			&loan.State,
			&loan.BorrowerID,
			&loan.PrincipalAmount,
			&loan.Rate,
			&loan.ROI,
			&loan.ApprovalProof,
			&loan.ApprovedBy,
			&loan.AgreementLetter,
			&loan.DisbursedBy,
			&loan.CreatedAt,
			&loan.ApprovedAt,
			&loan.InvestedAt,
			&loan.DisbursedAt,
			&loan.LastUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		loans = append(loans, loan)
	}

	return loans, nil
}

func (r *LoanRepository) CreateLoan(ctx context.Context, loan *model.Loan) (id int64, err error) {
	query := `
        INSERT INTO loans (
            state, borrower_id, principal_amount, rate, roi
        ) VALUES (
            ?, ?, ?, ?, ?
        )
    `

	res, err := r.DB.ExecContext(
		ctx,
		query,
		loan.State,
		loan.BorrowerID,
		loan.PrincipalAmount,
		loan.Rate,
		loan.ROI,
	)

	if err != nil {
		return
	}

	return res.LastInsertId()
}

func (r *LoanRepository) UpdateLoan(ctx context.Context, loan *model.Loan) error {
	query := `
		UPDATE loans
		SET state = ?,
			approval_proof = ?,
			approved_by = ?,
			agreement_letter = ?,
			disbursed_by = ?,
			approved_at = ?,
			invested_at = ?,
			disbursed_at = ?,
			last_updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		loan.State,
		loan.ApprovalProof,
		loan.ApprovedBy,
		loan.AgreementLetter,
		loan.DisbursedBy,
		loan.ApprovedAt,
		loan.InvestedAt,
		loan.DisbursedAt,
		loan.ID,
	)

	return err
}
