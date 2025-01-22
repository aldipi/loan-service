package repository

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aldipi/loan-service/model"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoanRepository(db)

	createdAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")
	lastUpdatedAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")
	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "last_updated_at"}).
		AddRow(1, "John Doe", createdAt, lastUpdatedAt)

	mock.ExpectQuery("SELECT id, name, created_at, last_updated_at FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	user, err := repo.GetUserByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.True(t, reflect.DeepEqual(user, &model.User{
		ID:            1,
		Name:          "John Doe",
		CreatedAt:     createdAt,
		LastUpdatedAt: lastUpdatedAt,
	}))
}
