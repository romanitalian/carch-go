package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/romanitalian/carch-go/config"
	"github.com/romanitalian/carch-go/internal/pkg/logger"
	"github.com/romanitalian/carch-go/internal/repository"
	"github.com/romanitalian/carch-go/migrations"

	"github.com/rs/zerolog"
)

func main() {
	// Initialize logger
	log := logger.New(
		logger.WithLevel(zerolog.InfoLevel),
		logger.WithPretty(),
	)

	// Loading configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config", err, map[string]interface{}{"error": err.Error()})
	}

	// Save original values
	originalUser := cfg.DB.User
	originalPassword := cfg.DB.Password

	// Connect to postgres without specifying a database to check/create the database and user
	adminDB, err := repository.NewPostgresDBWithoutDB(repository.PostgresConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     "postgres", // Use postgres superuser for initialization
		Password: "postgres", // Use postgres superuser password
		SSLMode:  cfg.DB.SSLMode,
		Logger:   log,
	})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL server", err, map[string]interface{}{"error": err.Error()})
	}
	defer adminDB.Close()

	// Initialize seed manager
	seedManager := migrations.NewSeedManager(adminDB, log)

	// Try to initialize database with user
	err = seedManager.EnsureUserExists(cfg.DB.User, cfg.DB.Password)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			log.Warn("No permission to create user, will try to continue with existing postgres user",
				map[string]interface{}{"error": err.Error()})
			log.Info(fmt.Sprintf("Please create the user manually with: CREATE USER \"%s\" WITH PASSWORD '%s';",
				cfg.DB.User, cfg.DB.Password), nil)

			// We'll continue with postgres user for now
			// Use postgres user for database operations
			cfg.DB.User = "postgres"
			cfg.DB.Password = "postgres"

			// But we'll restore the original values later
			defer func() {
				cfg.DB.User = originalUser
				cfg.DB.Password = originalPassword
			}()
		} else {
			log.Fatal("Failed to ensure user exists", err, map[string]interface{}{"error": err.Error()})
		}
	}

	if err := seedManager.EnsureDatabaseExists(cfg.DB.DBName); err != nil {
		log.Fatal("Failed to ensure database exists", err, map[string]interface{}{"error": err.Error()})
	}

	// Only try to grant permissions if we're not using the postgres user
	if cfg.DB.User != "postgres" {
		if err := seedManager.InitializeDatabase(cfg.DB.DBName, cfg.DB.User); err != nil {
			log.Fatal("Failed to initialize database", err, map[string]interface{}{"error": err.Error()})
		}
	} else {
		log.Info(fmt.Sprintf("Please grant permissions manually with: GRANT ALL PRIVILEGES ON DATABASE \"%s\" TO \"%s\";",
			cfg.DB.DBName, originalUser), nil)
	}

	// Initialize RabbitMQ user if needed
	// Only attempt to initialize if we're using a custom user (not guest)
	if cfg.RabbitMQ.User != "guest" {
		adminRabbitMQURL := "amqp://guest:guest@localhost:5672/"
		if err := seedManager.InitializeRabbitMQUser(
			adminRabbitMQURL,
			cfg.RabbitMQ.User,
			cfg.RabbitMQ.Password,
			cfg.RabbitMQ.VHost,
		); err != nil {
			log.Warn("Failed to initialize RabbitMQ user", map[string]interface{}{"error": err.Error()})
		}
	}

	fmt.Println("Database and RabbitMQ initialization completed successfully!")
	fmt.Println("Note: If you don't have superuser privileges, you may need to create users and grant permissions manually.")
	os.Exit(0)
}
