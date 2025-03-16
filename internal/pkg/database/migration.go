package database

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/romanitalian/carch-go/internal/pkg/logger"
)

// MigrationManager handles database migrations
type MigrationManager struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *sql.DB, logger *logger.Logger) *MigrationManager {
	return &MigrationManager{
		db:     db,
		logger: logger,
	}
}

// EnsureDatabaseExists checks if the database exists and creates it if it doesn't
func (m *MigrationManager) EnsureDatabaseExists(dbName string) error {
	m.logger.Info("Checking database existence", nil)

	// Check if database exists
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err := m.db.QueryRow(query, dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	// Create database if it doesn't exist
	if !exists {
		m.logger.Info(fmt.Sprintf("Database %s does not exist, creating it", dbName), nil)
		_, err = m.db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		m.logger.Info(fmt.Sprintf("Database %s created successfully", dbName), nil)
	} else {
		m.logger.Info(fmt.Sprintf("Database %s already exists", dbName), nil)
	}

	return nil
}

// RunMigrations runs all migrations from the migrations directory
func (m *MigrationManager) RunMigrations(ctx context.Context, migrationsPath string) error {
	m.logger.Info("Running database migrations", map[string]interface{}{
		"path": migrationsPath,
	})

	// Ensure migrations path is absolute
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for migrations: %w", err)
	}

	// Create postgres driver for migrations
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver for migrations: %w", err)
	}

	// Create migrate instance
	sourceURL := fmt.Sprintf("file://%s", absPath)
	m.logger.Info("Using migrations source", map[string]interface{}{
		"source": sourceURL,
	})

	migrator, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Run migrations
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	m.logger.Info("Database migrations completed successfully", nil)
	return nil
}
