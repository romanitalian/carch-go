package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/romanitalian/carch-go/internal/pkg/logger"
)

// PostgresConfig holds configuration for PostgreSQL connection
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Logger   *logger.Logger
}

// DB is a wrapper around sqlx.DB that exposes the underlying sql.DB
type DB struct {
	*sqlx.DB
	SQLDb *sql.DB // Expose the underlying sql.DB for migrations
}

// NewPostgresDB creates a new PostgreSQL connection
func NewPostgresDB(cfg PostgresConfig) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	if cfg.Logger != nil {
		cfg.Logger.Info("Connecting to PostgreSQL", map[string]interface{}{
			"host": cfg.Host,
			"port": cfg.Port,
			"user": cfg.User,
			"db":   cfg.DBName,
		})
	}

	sqlxDB, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		if cfg.Logger != nil {
			cfg.Logger.Error("Failed to connect to PostgreSQL", err, map[string]interface{}{
				"host": cfg.Host,
				"port": cfg.Port,
			})
		}
		return nil, err
	}

	err = sqlxDB.Ping()
	if err != nil {
		if cfg.Logger != nil {
			cfg.Logger.Error("Failed to ping PostgreSQL", err, nil)
		}
		return nil, err
	}

	if cfg.Logger != nil {
		cfg.Logger.Info("Successfully connected to PostgreSQL", nil)
	}

	return &DB{
		DB:    sqlxDB,
		SQLDb: sqlxDB.DB,
	}, nil
}

// NewPostgresDBWithoutDB creates a new PostgreSQL connection without specifying a database
// This is useful for administrative tasks like creating a database
func NewPostgresDBWithoutDB(cfg PostgresConfig) (*sql.DB, error) {
	cfg.Logger.Info("Connecting to PostgreSQL server", map[string]interface{}{
		"host": cfg.Host,
		"port": cfg.Port,
		"user": cfg.User,
	})

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	cfg.Logger.Info("Connected to PostgreSQL server", nil)
	return db, nil
}
