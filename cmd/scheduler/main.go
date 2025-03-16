package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/romanitalian/carch-go/config"
	"github.com/romanitalian/carch-go/internal/scheduler"
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

	// Initializing scheduler
	scheduler := scheduler.NewScheduler(cfg)

	// Registering tasks
	scheduler.RegisterTasks()

	// Starting scheduler
	go scheduler.Run(ctx)

	// Waiting for signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down scheduler...")
	cancel()
}
