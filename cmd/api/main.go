package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/romanitalian/carch-go/config"
	"github.com/romanitalian/carch-go/internal/pkg/database"
	"github.com/romanitalian/carch-go/internal/pkg/logger"
	"github.com/romanitalian/carch-go/internal/repository"
	"github.com/romanitalian/carch-go/internal/service"
	"github.com/romanitalian/carch-go/internal/transport/grpc"
	httpTransport "github.com/romanitalian/carch-go/internal/transport/http"

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

	// Try to connect to the database with the application user
	db, err := repository.NewPostgresDB(repository.PostgresConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.SSLMode,
		Logger:   log,
	})

	// If connection fails, try with postgres user
	if err != nil {
		if strings.Contains(err.Error(), "password authentication failed") {
			log.Warn("Failed to connect with configured user, trying with postgres user",
				map[string]interface{}{"error": err.Error()})

			db, err = repository.NewPostgresDB(repository.PostgresConfig{
				Host:     cfg.DB.Host,
				Port:     cfg.DB.Port,
				User:     "postgres",
				Password: "postgres",
				DBName:   cfg.DB.DBName,
				SSLMode:  cfg.DB.SSLMode,
				Logger:   log,
			})

			if err != nil {
				log.Fatal("Failed to initialize database with postgres user", err, map[string]interface{}{"error": err.Error()})
			}
		} else {
			log.Fatal("Failed to initialize database", err, map[string]interface{}{"error": err.Error()})
		}
	}
	defer db.Close()

	// Run database migrations
	migrationManager := database.NewMigrationManager(db.SQLDb, log)
	if err := migrationManager.RunMigrations(context.Background(), "./migrations"); err != nil {
		log.Fatal("Failed to run migrations", err, map[string]interface{}{"error": err.Error()})
	}
	log.Info("Database migrations completed successfully", nil)

	// Setting up AMQP/RabbitMQ connection
	messageQueue, err := repository.NewRabbitMQ(repository.RabbitMQConfig{
		URL:    cfg.RabbitMQ.URL,
		Logger: log,
	})
	if err != nil {
		// Try with default guest credentials if custom credentials fail
		if strings.Contains(err.Error(), "authentication failure") {
			log.Warn("Failed to connect to RabbitMQ with configured credentials, trying with guest user",
				map[string]interface{}{"error": err.Error()})

			messageQueue, err = repository.NewRabbitMQ(repository.RabbitMQConfig{
				URL:    "amqp://guest:guest@localhost:5672/",
				Logger: log,
			})

			if err != nil {
				log.Fatal("Failed to connect to RabbitMQ with guest user", err, map[string]interface{}{"error": err.Error()})
			}
		} else {
			log.Fatal("Failed to connect to RabbitMQ", err, map[string]interface{}{"error": err.Error()})
		}
	}
	defer messageQueue.Close()

	// Initializing repositories
	repos := repository.NewRepositories(db, messageQueue)

	// Initializing services
	services := service.NewServices(service.Deps{
		Repos: &service.Repositories{
			User: repos.User,
		},
		MessageQueue: messageQueue,
		Logger:       log,
	})

	// HTTP server with REST and GraphQL
	httpServer := httpTransport.NewServer(&httpTransport.Config{
		Address: cfg.HTTP.Address,
		Port:    cfg.HTTP.Port,
	}, services, log)

	// gRPC server
	grpcServer := grpc.NewServer(cfg.GRPC.Address+":"+cfg.GRPC.Port, services, log)

	// Creating errgroup for goroutine management
	serverErrors := make(chan error, 2)

	// Starting HTTP server
	go func() {
		log.Info("Starting HTTP server", map[string]interface{}{
			"address": cfg.HTTP.Address,
			"port":    cfg.HTTP.Port,
		})
		if err := httpServer.Run(); err != nil && err != http.ErrServerClosed {
			serverErrors <- fmt.Errorf("HTTP server error: %v", err)
		}
	}()

	// Starting gRPC server
	go func() {
		log.Info("Starting gRPC server", map[string]interface{}{
			"address": cfg.GRPC.Address,
			"port":    cfg.GRPC.Port,
		})
		if err := grpcServer.Run(); err != nil {
			serverErrors <- fmt.Errorf("gRPC server error: %v", err)
		}
	}()

	// Signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Error("Server error", err, nil)
	case sig := <-quit:
		log.Info("Received signal", map[string]interface{}{"signal": sig.String()})
	}

	log.Info("Shutting down servers", nil)

	// Creating context with timeout for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Graceful shutdown of servers
	shutdownErrors := make(chan error, 2)
	shutdownDone := make(chan struct{}, 2)

	go func() {
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			shutdownErrors <- fmt.Errorf("HTTP server shutdown error: %v", err)
		}
		shutdownDone <- struct{}{}
	}()

	go func() {
		if err := grpcServer.Shutdown(shutdownCtx); err != nil {
			shutdownErrors <- fmt.Errorf("gRPC server shutdown error: %v", err)
		}
		shutdownDone <- struct{}{}
	}()

	// Waiting for shutdown completion or timeout
	for i := 0; i < 2; i++ {
		select {
		case err := <-shutdownErrors:
			log.Error("Shutdown error", err, nil)
		case <-shutdownDone:
			log.Info("Server shutdown completed successfully", nil)
		}
	}

	log.Info("Servers gracefully stopped", nil)
}
