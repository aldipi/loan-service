package model

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID            int64     `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	LastUpdatedAt time.Time `json:"last_updated_at" db:"last_updated_at"`
}

type Employee struct {
	ID            int64     `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	LastUpdatedAt time.Time `json:"last_updated_at" db:"last_updated_at"`
}

type Investor struct {
	ID            int64     `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	LastUpdatedAt time.Time `json:"last_updated_at" db:"last_updated_at"`
}

type LoanState int16

const (
	LoanStateProposed LoanState = iota
	LoanStateApproved
	LoanStateInvested
	LoanStateDisbursed
)

type Loan struct {
	ID              int64           `json:"id" db:"id"`
	State           LoanState       `json:"state" db:"state"`
	BorrowerID      int64           `json:"borrower_id" db:"borrower_id"`
	PrincipalAmount int             `json:"principal_amount" db:"principal_amount"`
	Rate            decimal.Decimal `json:"rate" db:"rate"`
	ROI             decimal.Decimal `json:"roi" db:"roi"`
	ApprovalProof   sql.NullString  `json:"approval_proof" db:"approval_proof"`
	ApprovedBy      sql.NullInt64   `json:"approved_by" db:"approved_by"`
	AgreementLetter sql.NullString  `json:"agreement_letter" db:"agreement_letter"`
	DisbursedBy     sql.NullInt64   `json:"disbursed_by" db:"disbursed_by"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	ApprovedAt      sql.NullTime    `json:"approved_at" db:"approved_at"`
	InvestedAt      sql.NullTime    `json:"invested_at" db:"invested_at"`
	DisbursedAt     sql.NullTime    `json:"disbursed_at" db:"disbursed_at"`
	LastUpdatedAt   time.Time       `json:"last_updated_at" db:"last_updated_at"`
}

type Investment struct {
	ID              int64     `json:"id" db:"id"`
	Amount          int       `json:"amount" db:"amount"`
	InvestorID      int64     `json:"investor_id" db:"investor_id"`
	LoanID          int64     `json:"loan_id" db:"loan_id"`
	AgreementLetter string    `json:"agreement_letter" db:"agreement_letter"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	LastUpdatedAt   time.Time `json:"last_updated_at" db:"last_updated_at"`
}

type LoanProduct struct {
	ID            int64           `json:"id" db:"id"`
	Name          string          `json:"name" db:"name"`
	Rate          decimal.Decimal `json:"rate" db:"rate"`
	ROI           decimal.Decimal `json:"roi" db:"roi"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	LastUpdatedAt time.Time       `json:"last_updated_at" db:"last_updated_at"`
}

const (
	ErrLoanNotProposed         = LoanError("loan not proposed")
	ErrLoanNotApproved         = LoanError("loan not approved")
	ErrLoanNotInvested         = LoanError("loan not invested")
	ErrLoanNotFound            = LoanError("loan not found")
	ErrLoanProductNotFound     = LoanError("loan product not found")
	ErrInvestmentNotFound      = LoanError("investment not found")
	ErrInvestmentInvalidAmount = LoanError("investment amount is invalid")
	ErrUserNotFound            = LoanError("user not found")
	ErrEmployeeNotFound        = LoanError("employee not found")
	ErrInvestorNotFound        = LoanError("investor not found")
)

type LoanError string

func (e LoanError) Error() string {
	return string(e)
}
