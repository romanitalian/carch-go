#!/bin/bash

# Example of running the script:
# ./scripts/test_endpoints.sh
# Testing API endpoints...

# Creating user:
# {
#   "id": "b8e61f4c-ee3f-4def-a0d0-484f51020c24",
#   "email": "test@example.com",
#   "name": "Test User",
#   "created_at": "2025-03-16T21:40:09.831586+03:00",
#   "updated_at": "2025-03-16T21:40:09.831586+03:00"
# }

# Created user ID: b8e61f4c-ee3f-4def-a0d0-484f51020c24

# Getting all users:
# [
#   {
#     "id": "b8e61f4c-ee3f-4def-a0d0-484f51020c24",
#     "email": "test@example.com",
#     "name": "Test User",
#     "created_at": "2025-03-16T21:40:09.831586Z",
#     "updated_at": "2025-03-16T21:40:09.831586Z"
#   }
# ]

# Getting user by ID:
# {
#   "id": "b8e61f4c-ee3f-4def-a0d0-484f51020c24",
#   "email": "test@example.com",
#   "name": "Test User",
#   "created_at": "2025-03-16T21:40:09.831586Z",
#   "updated_at": "2025-03-16T21:40:09.831586Z"
# }

# Updating user:
# {
#   "id": "b8e61f4c-ee3f-4def-a0d0-484f51020c24",
#   "email": "updated@example.com",
#   "name": "Updated Name",
#   "created_at": "0001-01-01T00:00:00Z",
#   "updated_at": "2025-03-16T21:40:09.907922+03:00"
# }

# Getting updated user:
# {
#   "id": "b8e61f4c-ee3f-4def-a0d0-484f51020c24",
#   "email": "updated@example.com",
#   "name": "Updated Name",
#   "created_at": "2025-03-16T21:40:09.831586Z",
#   "updated_at": "2025-03-16T21:40:09.907922Z"
# }

# Deleting user:

# Done!

echo "Testing API endpoints..."

# Create user
echo -e "\nCreating user:"
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"secret123","name":"Test User"}')
echo $RESPONSE | jq '.'
USER_ID=$(echo $RESPONSE | jq -r '.id')

echo -e "\nCreated user ID: $USER_ID"

# Get all users
echo -e "\nGetting all users:"
curl -s http://localhost:8080/api/v1/users | jq '.'

# Get user by ID
echo -e "\nGetting user by ID:"
curl -s http://localhost:8080/api/v1/users/$USER_ID | jq '.'

# Update user
echo -e "\nUpdating user:"
curl -s -X PUT http://localhost:8080/api/v1/users/$USER_ID \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Name","email":"updated@example.com"}' | jq '.'

# Get updated user
echo -e "\nGetting updated user:"
curl -s http://localhost:8080/api/v1/users/$USER_ID | jq '.'

# Delete user
echo -e "\nDeleting user:"
curl -s -X DELETE http://localhost:8080/api/v1/users/$USER_ID
echo -e "\nDone!" 
