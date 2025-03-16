#!/bin/bash

# Install RabbitMQ server
brew install rabbitmq

# Start RabbitMQ server
brew services start rabbitmq

# Enable RabbitMQ management plugin
rabbitmq-plugins enable rabbitmq_management

# Start RabbitMQ management UI
open http://localhost:15672

# Create a new user
rabbitmqctl add_user admin admin

# Set user permissions
rabbitmqctl set_user_tags admin administrator

