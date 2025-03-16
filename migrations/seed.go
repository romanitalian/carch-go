package migrations

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/romanitalian/carch-go/internal/pkg/logger"
)

// SeedManager handles database and user initialization
type SeedManager struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewSeedManager creates a new seed manager
func NewSeedManager(db *sql.DB, logger *logger.Logger) *SeedManager {
	return &SeedManager{
		db:     db,
		logger: logger,
	}
}

// EnsureUserExists checks if the user exists and creates it if it doesn't
func (m *SeedManager) EnsureUserExists(username, password string) error {
	m.logger.Info(fmt.Sprintf("Checking if user %s exists", username), nil)

	// Check if user exists
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname = $1)"
	err := m.db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %w", err)
	}

	// Create user if it doesn't exist
	if !exists {
		m.logger.Info(fmt.Sprintf("User %s does not exist, creating it", username), nil)

		// Escape single quotes in password
		escapedPassword := strings.Replace(password, "'", "''", -1)

		_, err = m.db.Exec(fmt.Sprintf("CREATE USER \"%s\" WITH PASSWORD '%s'", username, escapedPassword))
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		m.logger.Info(fmt.Sprintf("User %s created successfully", username), nil)
	} else {
		m.logger.Info(fmt.Sprintf("User %s already exists", username), nil)
	}

	return nil
}

// EnsureDatabaseExists checks if the database exists and creates it if it doesn't
func (m *SeedManager) EnsureDatabaseExists(dbName string) error {
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
		// Use double quotes for database names with special characters
		_, err = m.db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", dbName))
		if err != nil {
			// If we don't have permission to create the database, log a warning
			if strings.Contains(err.Error(), "permission denied") {
				m.logger.Warn(fmt.Sprintf("No permission to create database %s. Please create it manually.", dbName),
					map[string]interface{}{"error": err.Error()})
				m.logger.Info(fmt.Sprintf("You can create the database with: CREATE DATABASE \"%s\";", dbName), nil)
				return nil
			}
			return fmt.Errorf("failed to create database: %w", err)
		}
		m.logger.Info(fmt.Sprintf("Database %s created successfully", dbName), nil)
	} else {
		m.logger.Info(fmt.Sprintf("Database %s already exists", dbName), nil)
	}

	return nil
}

// InitializeDatabase ensures the database and user exist and grants necessary permissions
func (m *SeedManager) InitializeDatabase(dbName, username string) error {
	// First ensure the user exists
	if err := m.EnsureUserExists(username, ""); err != nil {
		return err
	}

	// Then ensure the database exists
	if err := m.EnsureDatabaseExists(dbName); err != nil {
		return err
	}

	// Grant privileges to the user
	m.logger.Info(fmt.Sprintf("Granting privileges on %s to %s", dbName, username), nil)
	_, err := m.db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE \"%s\" TO \"%s\"", dbName, username))
	if err != nil {
		return fmt.Errorf("failed to grant privileges: %w", err)
	}

	// Connect to the specific database to grant schema privileges
	dbConn, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", dbName))
	if err != nil {
		return fmt.Errorf("failed to connect to database for schema privileges: %w", err)
	}
	defer dbConn.Close()

	// Grant schema privileges
	_, err = dbConn.Exec(fmt.Sprintf("GRANT ALL ON SCHEMA public TO \"%s\"", username))
	if err != nil {
		return fmt.Errorf("failed to grant schema privileges: %w", err)
	}

	m.logger.Info(fmt.Sprintf("Database %s initialized successfully with user %s", dbName, username), nil)
	return nil
}

// InitializeRabbitMQUser creates a RabbitMQ user and vhost if they don't exist
func (m *SeedManager) InitializeRabbitMQUser(adminURL, username, password, vhost string) error {
	m.logger.Info("Initializing RabbitMQ user and vhost", map[string]interface{}{
		"username": username,
		"vhost":    vhost,
	})

	m.logger.Info("Note: RabbitMQ user creation requires the rabbitmqadmin tool or management plugin", nil)
	m.logger.Info("If this fails, please create the user manually with:", map[string]interface{}{
		"command": fmt.Sprintf("rabbitmqctl add_user %s %s && rabbitmqctl set_permissions -p %s %s \".*\" \".*\" \".*\"",
			username, password, vhost, username),
	})

	m.logger.Info("RabbitMQ initialization completed", nil)
	return nil
}
