#!/bin/bash

# =====================================================================
# Script for load testing API endpoints of the service
# =====================================================================
#
# Purpose:
# - Measuring API endpoint performance under load
# - Determining maximum server throughput
# - Identifying performance bottlenecks
# - Comparing performance of different endpoints
#
# Usage:
#   ./build/bench_endpoints.sh
#
# Requirements:
# - Running service on localhost:8080
# - Installed wrk utility (https://github.com/wg/wrk)
# - Installed jq utility for JSON processing
#
# Example output:
# =============================================================
# Benchmarking API endpoints...
# Ensure your server is running before continuing...
#
# Creating test user for benchmarks:
# {
#   "id": "0d3389e2-022e-4155-865a-0cbce44959d3",
#   "email": "benchmark_1742145200@example.com",
#   "name": "Benchmark User"
# }
# Created test user ID: 0d3389e2-022e-4155-865a-0cbce44959d3
#
# === Benchmarking GET http://localhost:8080/api/v1/users/ ===
#
# --- Concurrency: 10 connections ---
# Running 3s test @ http://localhost:8080/api/v1/users/
#   2 threads and 10 connections
#   Thread Stats   Avg      Stdev     Max   +/- Stdev
#     Latency     5.93ms   13.29ms 154.72ms   94.15%
#     Req/Sec     2.44k   752.12     3.36k    70.00%
#   14601 requests in 3.01s, 10.03MB read
# Requests/sec:   4849.85
# Transfer/sec:      3.33MB
# =============================================================
#
# Interpreting results:
# - Requests/sec: number of requests per second (throughput)
# - Latency: response time (average, standard deviation, maximum)
# - Transfer/sec: volume of data transferred per second
#
# The script tests the following endpoints:
# - GET /api/v1/users/ - get list of users
# - GET /api/v1/users/:id - get user by ID
# - POST /api/v1/users/ - create user
# - PUT /api/v1/users/:id - update user
#
# =====================================================================

echo "Benchmarking API endpoints..."

# Make sure the server is running
echo "Ensure your server is running before continuing..."

# Generate a unique email for each benchmark run
TIMESTAMP=$(date +%s)
UNIQUE_EMAIL="benchmark_${TIMESTAMP}@example.com"

# Create a test user first to have a valid ID for benchmarking
echo -e "\nCreating test user for benchmarks:"
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users/ \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$UNIQUE_EMAIL\",\"password\":\"benchmark123\",\"name\":\"Benchmark User\"}")
echo $RESPONSE | jq '.'
USER_ID=$(echo $RESPONSE | jq -r '.id')

# Check if user was created successfully
if [ "$USER_ID" = "null" ] || [ -z "$USER_ID" ]; then
  echo "Failed to create test user. Exiting."
  exit 1
fi

echo -e "Created test user ID: $USER_ID"

# Create Lua scripts for POST and PUT requests
cat > /tmp/post_request.lua << EOF
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.body = '{"email":"test_' .. math.random(1000000) .. '@example.com","password":"secret123","name":"Test User"}'
request = function()
  local random_email = 'test_' .. math.random(1000000) .. '@example.com'
  local body = '{"email":"' .. random_email .. '","password":"secret123","name":"Test User"}'
  return wrk.format(nil, nil, nil, body)
end
EOF

cat > /tmp/put_request.lua << EOF
wrk.method = "PUT"
wrk.headers["Content-Type"] = "application/json"
wrk.body = '{"name":"Updated Name","email":"updated@example.com"}'
EOF

# Function to run benchmarks with different concurrency levels
run_benchmark() {
  local endpoint=$1
  local method=$2
  local script=$3
  
  echo -e "\n=== Benchmarking $method $endpoint ==="
  
  for c in 10 25; do
    echo -e "\n--- Concurrency: $c connections ---"
    if [ -z "$script" ]; then
      wrk -t2 -c$c -d3s $endpoint
    else
      wrk -t2 -c$c -d3s -s $script $endpoint
    fi
  done
}

# Run benchmarks for different endpoints
run_benchmark "http://localhost:8080/api/v1/users/" "GET"
run_benchmark "http://localhost:8080/api/v1/users/$USER_ID" "GET"
run_benchmark "http://localhost:8080/api/v1/users/" "POST" "/tmp/post_request.lua"
run_benchmark "http://localhost:8080/api/v1/users/$USER_ID" "PUT" "/tmp/put_request.lua"

# Clean up - delete the test user
echo -e "\nCleaning up - deleting test user:"
curl -s -X DELETE http://localhost:8080/api/v1/users/$USER_ID

# Remove temporary files
rm -f /tmp/post_request.lua /tmp/put_request.lua

echo -e "\nBenchmarking completed!"
