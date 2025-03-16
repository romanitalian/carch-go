package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/romanitalian/carch-go/config"
	"github.com/romanitalian/carch-go/internal/repository"
	"github.com/romanitalian/carch-go/internal/worker"
)

func main() {
	// Loading configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initializing context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setting up RabbitMQ connection
	messageQueue, err := repository.NewRabbitMQ(repository.RabbitMQConfig{
		URL: cfg.RabbitMQ.URL,
	})
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer messageQueue.Close()

	// Initializing and starting worker
	worker := worker.NewWorker(messageQueue)
	go worker.Run(ctx)

	// Waiting for signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
	cancel()
}
