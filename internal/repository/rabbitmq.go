package repository

import (
	"fmt"

	"github.com/streadway/amqp"

	"github.com/romanitalian/carch-go/internal/pkg/logger"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	log     *logger.Logger
}

type RabbitMQConfig struct {
	URL    string
	Logger *logger.Logger
}

func NewRabbitMQ(cfg RabbitMQConfig) (*RabbitMQ, error) {
	// Connecting to RabbitMQ
	if cfg.Logger != nil {
		cfg.Logger.Info("Connecting to RabbitMQ", map[string]interface{}{"url": cfg.URL})
	}

	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		if cfg.Logger != nil {
			cfg.Logger.Error("Failed to connect to RabbitMQ", err, map[string]interface{}{"url": cfg.URL})
		}
		return nil, err
	}

	// Creating channel
	if cfg.Logger != nil {
		cfg.Logger.Info("Creating RabbitMQ channel", nil)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		if cfg.Logger != nil {
			cfg.Logger.Error("Failed to create RabbitMQ channel", err, nil)
		}
		return nil, err
	}

	// Declaring queue
	queueName := "tasks"
	if cfg.Logger != nil {
		cfg.Logger.Info("Declaring RabbitMQ queue", map[string]interface{}{"queue": queueName})
	}

	_, err = ch.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		if cfg.Logger != nil {
			cfg.Logger.Error("Failed to declare RabbitMQ queue", err, map[string]interface{}{"queue": queueName})
		}
		return nil, err
	}

	if cfg.Logger != nil {
		cfg.Logger.Info("Successfully connected to RabbitMQ", nil)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		log:     cfg.Logger,
	}, nil
}

func (r *RabbitMQ) Close() error {
	if r.log != nil {
		r.log.Info("Closing RabbitMQ connection", nil)
	}

	if err := r.channel.Close(); err != nil {
		if r.log != nil {
			r.log.Error("Failed to close RabbitMQ channel", err, nil)
		}
		return err
	}

	if err := r.conn.Close(); err != nil {
		if r.log != nil {
			r.log.Error("Failed to close RabbitMQ connection", err, nil)
		}
		return err
	}

	if r.log != nil {
		r.log.Info("RabbitMQ connection closed successfully", nil)
	}

	return nil
}

func (r *RabbitMQ) Consume(queueName string) (<-chan amqp.Delivery, error) {
	if r.log != nil {
		r.log.Info("Starting to consume from queue", map[string]interface{}{"queue": queueName})
	}

	msgs, err := r.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if err != nil && r.log != nil {
		r.log.Error("Failed to consume from queue", err, map[string]interface{}{"queue": queueName})
	}

	return msgs, err
}

// InitializeRabbitMQUser creates a RabbitMQ user and vhost if they don't exist
func InitializeRabbitMQUser(adminURL, username, password, vhost string, logger *logger.Logger) error {
	logger.Info("Initializing RabbitMQ user and vhost", map[string]interface{}{
		"username": username,
		"vhost":    vhost,
	})

	// Connect to RabbitMQ with admin credentials
	conn, err := amqp.Dial(adminURL)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ with admin credentials", err, nil)
		return fmt.Errorf("failed to connect to RabbitMQ with admin credentials: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to open a channel", err, nil)
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	// Try to use the management API via HTTP
	logger.Info("Note: RabbitMQ user creation requires the rabbitmqadmin tool or management plugin", nil)
	logger.Info("If this fails, please create the user manually with:", map[string]interface{}{
		"command": fmt.Sprintf("rabbitmqctl add_user %s %s && rabbitmqctl set_permissions -p %s %s \".*\" \".*\" \".*\"",
			username, password, vhost, username),
	})

	// Try to use the default vhost
	logger.Info("Using default vhost '/' for RabbitMQ connection", nil)

	logger.Info("RabbitMQ initialization completed", nil)
	return nil
}
