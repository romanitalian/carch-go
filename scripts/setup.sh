#!/bin/bash

# Make scripts executable
chmod +x scripts/setup_db.sh
chmod +x scripts/setup_rabbitmq.sh

# Run setup scripts
echo "Setting up PostgreSQL..."
./scripts/setup_db.sh

echo "Setting up RabbitMQ..."
./scripts/setup_rabbitmq.sh

echo "All services setup completed successfully!"
echo "You can now run 'make run-all' to start all services." 