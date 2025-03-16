package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/romanitalian/carch-go/internal/domain"
)

func TestPostgresUserRepository_Create(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	ctx := context.Background()
	user := &domain.User{
		ID:        "user-123",
		Email:     "test@example.com",
		Password:  "hashed_password",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Expected query setup
	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO users (id, email, password_hash, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`)).WithArgs(
		user.ID,
		user.Email,
		user.Password,
		user.Name,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))

	// Act
	err = repo.Create(ctx, user)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresUserRepository_GetByID(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	ctx := context.Background()
	userID := "user-123"
	expectedUser := &domain.User{
		ID:        userID,
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Expected query setup
	rows := sqlmock.NewRows([]string{"id", "email", "name", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Name, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, email, name, created_at, updated_at
		FROM users
		WHERE id = $1`)).
		WithArgs(userID).
		WillReturnRows(rows)

	// Act
	user, err := repo.GetByID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresUserRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	ctx := context.Background()
	userID := "non-existent-id"

	// Expected query setup
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, email, name, created_at, updated_at
		FROM users
		WHERE id = $1`)).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	// Act
	user, err := repo.GetByID(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresUserRepository_Update(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	ctx := context.Background()
	user := &domain.User{
		ID:        "user-123",
		Email:     "updated@example.com",
		Name:      "Updated User",
		UpdatedAt: time.Now(),
	}

	// Expected query setup
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE users
		SET email = $1, name = $2, updated_at = $3
		WHERE id = $4`)).
		WithArgs(user.Email, user.Name, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Act
	err = repo.Update(ctx, user)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresUserRepository_Update_NotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	ctx := context.Background()
	user := &domain.User{
		ID:        "non-existent-id",
		Email:     "updated@example.com",
		Name:      "Updated User",
		UpdatedAt: time.Now(),
	}

	// Expected query setup
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE users
		SET email = $1, name = $2, updated_at = $3
		WHERE id = $4`)).
		WithArgs(user.Email, user.Name, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Act
	err = repo.Update(ctx, user)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresUserRepository_Delete(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	ctx := context.Background()
	userID := "user-123"

	// Expected query setup
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Act
	err = repo.Delete(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresUserRepository_Delete_NotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	ctx := context.Background()
	userID := "non-existent-id"

	// Expected query setup
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Act
	err = repo.Delete(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresUserRepository_List(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	ctx := context.Background()
	expectedUsers := []*domain.User{
		{
			ID:        "user-1",
			Email:     "user1@example.com",
			Name:      "User 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "user-2",
			Email:     "user2@example.com",
			Name:      "User 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Expected query setup
	rows := sqlmock.NewRows([]string{"id", "email", "name", "created_at", "updated_at"})
	for _, user := range expectedUsers {
		rows.AddRow(user.ID, user.Email, user.Name, user.CreatedAt, user.UpdatedAt)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, email, name, created_at, updated_at
		FROM users
		ORDER BY created_at DESC`)).
		WillReturnRows(rows)

	// Act
	users, err := repo.List(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, users, len(expectedUsers))
	assert.NoError(t, mock.ExpectationsWereMet())
}
