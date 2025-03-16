package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/romanitalian/carch-go/config"
)

type Scheduler struct {
	cron *cron.Cron
	cfg  *config.Config
}

func NewScheduler(cfg *config.Config) *Scheduler {
	return &Scheduler{
		cron: cron.New(cron.WithSeconds()),
		cfg:  cfg,
	}
}

func (s *Scheduler) RegisterTasks() {
	// Registration of periodic tasks
	s.cron.AddFunc("0 * * * * *", func() { // Every minute
		if err := s.exampleTask(); err != nil {
			log.Printf("Error running example task: %v", err)
		}
	})

	s.cron.AddFunc("0 0 * * * *", func() { // Every hour
		if err := s.hourlyTask(); err != nil {
			log.Printf("Error running hourly task: %v", err)
		}
	})
}

func (s *Scheduler) Run(ctx context.Context) {
	s.cron.Start()
	defer s.cron.Stop()

	// Waiting for termination signal
	<-ctx.Done()
}

func (s *Scheduler) exampleTask() error {
	log.Printf("Running example task at %v", time.Now())
	return nil
}

func (s *Scheduler) hourlyTask() error {
	log.Printf("Running hourly task at %v", time.Now())
	return nil
}
