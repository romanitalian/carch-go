# Carch-go - Golang Clean Architecture Example

## Project Purpose

This project is an example of a Go microservice implementation using Clean Architecture principles. The project demonstrates:

- Code separation into layers according to Clean Architecture principles
- REST API implementation using Go's native http package with ServeMux (Go 1.22+)
- gRPC server implementation
- Working with PostgreSQL for data storage
- Integration with RabbitMQ for asynchronous message processing
- Graceful shutdown for proper service termination
- Configuration through environment variables and .env files

## QUICK START

```bash
# Clone repository
git clone https://github.com/your-username/carch-go.git
cd carch-go

# Copy example environment file and adjust as needed
cp .env.example .env

# Setup database and RabbitMQ (creates user, database, and grants permissions)
make setup-local

# Run all services
make run-all

# Alternatively, run just the API server
make run-api
```

### Available Make Commands

```bash
# Setup database and RabbitMQ
make setup-local  # or shorthand: make sl

# Run all services
make run-all      # or shorthand: make r

# Run API server only
make run-api

# Run worker only
make run-worker

# Run scheduler only
make run-scheduler

# Setup and run all services
make setup-and-run  # or shorthand: make sr

# Setup and run API server only
make setup-and-run-api  # or shorthand: make sa
```

## Project Structure

```
carch-go/
├── cmd/                    # Application entry points
│   └── api/               # API server
│   └── worker/            # Background processors
│   └── scheduler/         # Task scheduler (cron)
│   └── seed/              # Database and RabbitMQ initialization
├── config/                # Configuration 
├── internal/              # Internal application code
│   ├── domain/           # Business models and interfaces
│   ├── service/          # Business logic
│   ├── repository/       # Data storage operations
│   ├── transport/        # Transport layer
│   │   ├── http/        # REST API handlers
│   │   ├── grpc/        # gRPC handlers
│   │   ├── graphql/     # GraphQL handlers
│   │   └── middleware/  # Middleware components
│   ├── worker/          # Background processors
│   │   └── tasks/       # Task definitions
│   └── pkg/             # Internal utilities
├── pkg/                  # Public libraries
├── api/                  # API definitions
│   ├── proto/           # Protobuf definitions
│   └── graphql/         # GraphQL schemas
├── build/                # Compiled binaries and scripts
├── deployments/          # Deployment configurations
│   ├── docker/          # Dockerfiles
│   └── k8s/             # Kubernetes manifests
└── scripts/             # Various scripts
```

## Requirements

- Go 1.21+
- PostgreSQL
- RabbitMQ

## Configuration

Service configuration is done through the `config/config.yaml` file:

```yaml
http:
  address: "localhost"
  port: "8080"

grpc:
  address: "localhost"
  port: "9090"

db:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "postgres"
  dbname: "carch-go"
  sslmode: "disable"

rabbitmq:
  url: "amqp://guest:guest@localhost:5672/"
```

## Operation

### Running the Service

```bash
# Run in development mode
./build/carch-go

# Run with a specific configuration path
CONFIG_PATH=/path/to/config.yaml ./build/carch-go
```

### API Testing

To test the API, you can use the script:

```bash
./build/test_endpoints.sh
```

### Load Testing

For API load testing, use the script:

```bash
./build/bench_endpoints.sh
```

The script uses the `wrk` utility to measure API endpoint performance.

### Monitoring

The service provides metrics in Prometheus format at the `/metrics` endpoint.

### Logging

Logs are output to standard output (stdout) and can be redirected to a file or logging system.

### Graceful Shutdown

The service properly terminates when receiving SIGINT or SIGTERM signals, closing all connections and completing current requests.

## API Endpoints

### REST API
- POST /api/v1/users/ - Create a user
- GET /api/v1/users/:id - Get user by ID
- PUT /api/v1/users/:id - Update user
- DELETE /api/v1/users/:id - Delete user
- GET /api/v1/users/ - Get list of users

### gRPC
- Port 9090 - gRPC server with similar methods for user operations

## Development

### Adding New Endpoints

1. Define a model in `internal/domain`
2. Create a repository interface in `internal/repository`
3. Implement business logic in `internal/service`
4. Add handlers in `internal/transport/http` and/or `internal/transport/grpc`

### Running Tests

```bash
go test ./...
```

## Unit Tests

The project contains a complete set of unit tests for all components:

### Service Layer Tests (`internal/service/user_test.go`):
- Tests for all UserService methods (Create, GetByID, Update, Delete, List)
- Error handling verification, including cases when a user is not found
- Using mocks to simulate the repository

### HTTP Handler Tests (`internal/transport/http/handler_test.go`):
- Tests for all REST API endpoints
- Verification of correct request and response handling
- Input data validation testing
- Error handling testing

### HTTP Server Tests (`internal/transport/http/server_test.go`):
- Server initialization tests
- Server startup tests
- Graceful shutdown tests

### gRPC Server Tests (`internal/transport/grpc/server_test.go`):
- gRPC server startup tests
- Graceful shutdown tests
- Using buffered connections for gRPC testing

### PostgreSQL Repository Tests (`internal/repository/postgres_user_test.go`):
- Tests for all repository methods
- Using sqlmock to simulate the database
- SQL query and result handling verification

### Integration Tests (`cmd/api/main_test.go`):
- Test for verifying startup and graceful shutdown of the entire application
- Test for verifying the graceful shutdown mechanism

All tests follow the AAA pattern (Arrange-Act-Assert) and use mocks to isolate the components being tested. The following libraries are used for testing:
- `github.com/stretchr/testify/assert` - for assertion verification
- `github.com/stretchr/testify/mock` - for creating mocks
- `github.com/DATA-DOG/go-sqlmock` - for database simulation
- `net/http/httptest` - for HTTP handler testing
- `google.golang.org/grpc/test/bufconn` - for gRPC server testing

These tests provide good code coverage and help identify potential issues early in development.

## License

Creative Commons Zero (CC0)

This project is released into the public domain using CC0. You can copy, modify, distribute, and use the code without restrictions, even for commercial purposes, without the need to ask for permission.

# Environment Configuration

The application uses environment variables for configuration. You can set these variables in a `.env` file in the root directory of the project.

## Available Environment Variables

```
# HTTP Server
HTTP_ADDRESS=0.0.0.0
HTTP_PORT=8080

# gRPC Server
GRPC_ADDRESS=0.0.0.0
GRPC_PORT=9090

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=Carch-go
DB_SSLMODE=disable

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

## Database Migrations

The application automatically runs database migrations at startup. Migrations are stored in the `migrations` directory and are applied using the [golang-migrate](https://github.com/golang-migrate/migrate) library.

### Migration Files

Migration files follow the naming convention:

```
{version}_{description}.{up|down}.sql
```

For example:
- `000001_create_users_table.up.sql` - Creates the users table
- `000001_create_users_table.down.sql` - Drops the users table

### Creating New Migrations

To create a new migration, add two files to the `migrations` directory:
1. `{version}_{description}.up.sql` - SQL commands to apply the migration
2. `{version}_{description}.down.sql` - SQL commands to revert the migration

Example of creating a new migration:

```bash
# Create migration files
touch migrations/000002_add_user_roles.up.sql
touch migrations/000002_add_user_roles.down.sql

# Edit up migration
echo "ALTER TABLE users ADD COLUMN role VARCHAR(50) NOT NULL DEFAULT 'user';" > migrations/000002_add_user_roles.up.sql

# Edit down migration
echo "ALTER TABLE users DROP COLUMN role;" > migrations/000002_add_user_roles.down.sql
```

### Migration Process

At startup, the application:
1. Checks if the database exists and creates it if needed
2. Connects to the database
3. Runs all pending migrations in order

This ensures that the database schema is always up-to-date with the application code.

## Logging

The application uses `zerolog` for structured logging. The logger is initialized in `main.go` and passed through all layers of the application using the Option func pattern.

Example log output:
```
2023-03-16T12:34:56Z INF Starting HTTP server address=0.0.0.0 port=8080
```