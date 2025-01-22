package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aldipi/loan-service/model"
	"github.com/labstack/echo/v4"
)

type Usecase interface {
	GetLoans(ctx context.Context, limit int, offset int) ([]*model.Loan, error)
	GetLoansByBorrowerID(ctx context.Context, borrowerID int64, limit int, offset int) ([]*model.Loan, error)
	CreateLoan(ctx context.Context, userID int64, loanProductID int64, amount int) (loan *model.Loan, err error)
	ApproveLoan(ctx context.Context, loanID int64, employeeID int64, approvalProof string) error
	DisburseLoan(ctx context.Context, loanID int64, employeeID int64, agreementLetter string) error
	GetInvestmentsByInvestorID(ctx context.Context, investorID int64, limit int, offset int) ([]*model.Investment, error)
	CheckAvailableInvestmentByLoanID(ctx context.Context, loanID int64) (int, error)
	CreateInvestment(ctx context.Context, investorID int64, loanID int64, amount int) (investment *model.Investment, err error)
}

type HttpHanlder struct {
	uc Usecase
}

func NewHttpHandler(uc Usecase) *HttpHanlder {
	return &HttpHanlder{uc: uc}
}

func (h *HttpHanlder) GetAllLoans(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	loans, err := h.uc.GetLoans(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, loans)
}

func (h *HttpHanlder) GetLoans(c echo.Context) error {
	borrowerID, _ := strconv.ParseInt(c.Request().Header["X-User-Id"][0], 10, 64)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	loans, err := h.uc.GetLoansByBorrowerID(c.Request().Context(), borrowerID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, loans)
}

func (h *HttpHanlder) CreateLoan(c echo.Context) error {
	userID, _ := strconv.ParseInt(c.Request().Header["X-User-Id"][0], 10, 64)
	loanProductID, _ := strconv.ParseInt(c.FormValue("loanProductID"), 10, 64)
	amount, _ := strconv.Atoi(c.FormValue("amount"))
	loan, err := h.uc.CreateLoan(c.Request().Context(), userID, loanProductID, amount)
	if err != nil {
		if loanErr, ok := err.(model.LoanError); ok {
			return c.JSON(http.StatusBadRequest, loanErr.Error())
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, loan)
}

func (h *HttpHanlder) ApproveLoan(c echo.Context) error {
	loanID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	employeeID, _ := strconv.ParseInt(c.Request().Header["X-User-Id"][0], 10, 64)
	approvalProof := c.FormValue("approvalProof")
	err := h.uc.ApproveLoan(c.Request().Context(), loanID, employeeID, approvalProof)
	if err != nil {
		if loanErr, ok := err.(model.LoanError); ok {
			return c.JSON(http.StatusBadRequest, loanErr.Error())
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "Loan approved")
}

func (h *HttpHanlder) DisburseLoan(c echo.Context) error {
	loanID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	employeeID, _ := strconv.ParseInt(c.Request().Header["X-User-Id"][0], 10, 64)
	agreementLetter := c.FormValue("agreementLetter")
	err := h.uc.DisburseLoan(c.Request().Context(), loanID, employeeID, agreementLetter)
	if err != nil {
		if loanErr, ok := err.(model.LoanError); ok {
			return c.JSON(http.StatusBadRequest, loanErr.Error())
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "Loan disbursed")
}

func (h *HttpHanlder) GetInvestments(c echo.Context) error {
	investorID, _ := strconv.ParseInt(c.Request().Header["X-User-Id"][0], 10, 64)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	investments, err := h.uc.GetInvestmentsByInvestorID(c.Request().Context(), investorID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, investments)
}

func (h *HttpHanlder) CreateInvestment(c echo.Context) error {
	investorID, _ := strconv.ParseInt(c.Request().Header["X-User-Id"][0], 10, 64)
	loanID, _ := strconv.ParseInt(c.FormValue("loan_id"), 10, 64)
	amount, _ := strconv.Atoi(c.FormValue("amount"))
	investment, err := h.uc.CreateInvestment(c.Request().Context(), investorID, loanID, amount)
	if err != nil {
		if loanErr, ok := err.(model.LoanError); ok {
			return c.JSON(http.StatusBadRequest, loanErr.Error())
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, investment)
}

func (h *HttpHanlder) LoanAvailability(c echo.Context) error {
	loanID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	availableAmount, err := h.uc.CheckAvailableInvestmentByLoanID(c.Request().Context(), loanID)
	if err != nil {
		if loanErr, ok := err.(model.LoanError); ok {
			return c.JSON(http.StatusBadRequest, loanErr.Error())
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, availableAmount)
}
