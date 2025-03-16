package worker

import (
	"context"
	"log"

	"github.com/streadway/amqp"
)

type MessageQueue interface {
	Consume(queueName string) (<-chan amqp.Delivery, error)
	Close() error
}

type Worker struct {
	queue MessageQueue
}

func NewWorker(queue MessageQueue) *Worker {
	return &Worker{
		queue: queue,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	// Subscribing to the task queue
	messages, err := w.queue.Consume("tasks")
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-messages:
			if err := w.processMessage(msg); err != nil {
				log.Printf("Error processing message: %v", err)
			}
		}
	}
}

func (w *Worker) processMessage(msg amqp.Delivery) error {
	// Processing message
	log.Printf("Processing message: %s", string(msg.Body))

	// Acknowledging processing
	return msg.Ack(false)
}
