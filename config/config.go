package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTP struct {
		Address string `yaml:"address" env:"HTTP_ADDRESS" env-default:"0.0.0.0"`
		Port    string `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	} `yaml:"http"`
	GRPC struct {
		Address string `yaml:"address" env:"GRPC_ADDRESS" env-default:"0.0.0.0"`
		Port    string `yaml:"port" env:"GRPC_PORT" env-default:"9090"`
	} `yaml:"grpc"`
	DB struct {
		Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
		Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
		User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
		Password string `yaml:"password" env:"DB_PASSWORD" env-default:"postgres"`
		DBName   string `yaml:"dbname" env:"DB_NAME" env-default:"Carch-go"`
		SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
	} `yaml:"db"`
	RabbitMQ struct {
		URL      string `yaml:"url" env:"RABBITMQ_URL" env-default:"amqp://guest:guest@localhost:5672/"`
		User     string `yaml:"user" env:"RABBITMQ_USER" env-default:"guest"`
		Password string `yaml:"password" env:"RABBITMQ_PASSWORD" env-default:"guest"`
		VHost    string `yaml:"vhost" env:"RABBITMQ_VHOST" env-default:"/"`
	} `yaml:"rabbitmq"`
}

// Load loads configuration from .env file and environment variables
func Load() (*Config, error) {
	// Try to load .env file, but continue if it doesn't exist
	_ = godotenv.Load()

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// New is an alias for Load for compatibility with the example
func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	var c Config
	if err := cleanenv.ReadEnv(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
