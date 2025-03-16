#!/bin/bash

# Load environment variables from .env file
source .env

echo "Setting up RabbitMQ user $RABBITMQ_USER..."

# Check if rabbitmqctl is available
if command -v rabbitmqctl &> /dev/null; then
    # Create RabbitMQ user and set permissions
    rabbitmqctl add_user "$RABBITMQ_USER" "$RABBITMQ_PASSWORD" || true
    rabbitmqctl set_permissions -p "$RABBITMQ_VHOST" "$RABBITMQ_USER" ".*" ".*" ".*"
    rabbitmqctl set_user_tags "$RABBITMQ_USER" administrator
    
    echo "RabbitMQ setup completed successfully!"
else
    echo "rabbitmqctl command not found. Please install RabbitMQ or ensure it's in your PATH."
    echo "Alternatively, you can manually create the user with these commands:"
    echo "rabbitmqctl add_user $RABBITMQ_USER $RABBITMQ_PASSWORD"
    echo "rabbitmqctl set_permissions -p $RABBITMQ_VHOST $RABBITMQ_USER \".*\" \".*\" \".*\""
    echo "rabbitmqctl set_user_tags $RABBITMQ_USER administrator"
    exit 1
fi 
