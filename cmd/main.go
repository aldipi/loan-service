package main

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/aldipi/loan-service/handler"
	"github.com/aldipi/loan-service/repository"
	"github.com/aldipi/loan-service/usecase"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	connStr := os.Getenv("DB_CONN_STRING")
	if connStr == "" {
		panic("DB_CONN_STRING environment variable not set")
	}
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := repository.NewLoanRepository(db)
	uc := usecase.NewLoanUsecase(repo)
	h := handler.NewHttpHandler(uc)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/loans/all", h.GetAllLoans)
	e.GET("/loans", h.GetLoans)
	e.POST("/loans", h.CreateLoan)
	e.PATCH("/loans/:id/approval", h.ApproveLoan)
	e.PATCH("/loans/:id/disbursement", h.DisburseLoan)
	e.GET("/loans/:id/availability", h.LoanAvailability)

	e.GET("/investments", h.GetInvestments)
	e.POST("/investments", h.CreateInvestment)

	e.Logger.Fatal(e.Start(":8080"))
}
