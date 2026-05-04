# gRPC SQLC User Service

A production-ready gRPC user management service built with Go, sqlc, and PostgreSQL. This project demonstrates type-safe database operations using sqlc with automatic schema migrations.

## Features

- **gRPC API** - High-performance RPC framework for service communication
- **sqlc Integration** - Type-safe SQL code generation from SQL queries
- **PostgreSQL** - Robust relational database for data persistence
- **Schema Migrations** - Version-controlled database schema changes using golang-migrate
- **Docker Support** - Containerized deployment with Docker Compose
- **Protocol Buffers** - Strongly typed API contracts

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   gRPC      │────▶│   Handlers  │────▶│    sqlc     │
│   Client    │     │   (users.go)│     │   Queries   │
└─────────────┘     └─────────────┘     └──────┬──────┘
                                              │
                                       ┌──────▼──────┐
                                       │  PostgreSQL │
                                       └─────────────┘
```

## Project Structure

```
.
├── cmd/
│   └── server.go              # Application entry point
├── internal/
│   └── api/handlers/
│       └── users.go           # gRPC handler implementations
├── database/
│   └── database.go            # Database connection manager
├── db/
│   ├── db.go                  # sqlc generated database code
│   ├── models.go              # sqlc generated models
│   └── query.sql.go           # sqlc generated queries
├── migrations/
│   ├── 000001_init_schema.up.sql     # Initial schema
│   └── 000001_init_schema.down.sql   # Schema rollback
├── proto/
│   ├── main.proto             # Protocol buffer definitions
│   └── gen/                   # Generated Go protobuf code
├── query/
│   └── query.sql              # SQL queries for sqlc
├── schema.sql                 # Database schema (reference)
├── sqlc.yaml                  # sqlc configuration
├── Dockerfile                 # Container build instructions
├── docker-compose.yml         # Multi-container orchestration
└── Makefile                   # Build and migration commands
```

## API Endpoints

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `CreateUser` | `UserRequest` | `UserResponse` | Create a new user |
| `GetUser` | `UserIdRequest` | `UserResponse` | Retrieve user by ID |
| `UpdateUser` | `UserUpdateRequest` | `UserResponse` | Update existing user |
| `DeactivateUser` | `UserIdRequest` | `UserResponse` | Delete user by ID |

### Message Types

```protobuf
message UserRequest {
    string first_name = 1;
    string last_name = 2;
    string phone_number = 3;
}

message UserUpdateRequest {
    int32 user_id = 1;
    string first_name = 2;
    string last_name = 3;
    string phone_number = 4;
}

message UserIdRequest {
    int32 user_id = 1;
}

message UserResponse {
    string message = 1;
    bool status = 2;
}
```

## Quick Start

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- [Go 1.26+](https://golang.org/dl/) (for local development)
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) (for code generation)
- [protoc](https://grpc.io/docs/protoc-installation/) (for protobuf generation)

### Running with Docker Compose

```bash
# Create the Docker network
make create_local_network

# Start all services (postgres, migrations, grpc server)
make up

# View logs
make logs

# Stop all services
make down
```

The gRPC server will be available at `localhost:50052`.

### Local Development

```bash
# 1. Start PostgreSQL
make up

# 2. Run migrations
make migrate-up

# 3. Generate sqlc code (after schema changes)
make sqlc-generate

# 4. Run the server locally
go run cmd/server.go
```

## Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for schema versioning.

### Creating New Migrations

```bash
# Create up and down migration files
make migrate-create name=add_email_column

# Edit the generated files:
# - migrations/000002_add_email_column.up.sql
# - migrations/000002_add_email_column.down.sql
```

### Migration Commands

| Command | Description |
|---------|-------------|
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-version` | Check current migration version |
| `make migrate-create name=<name>` | Create new migration files |

### Migration Workflow

1. Create migration: `make migrate-create name=feature_name`
2. Write `*.up.sql` for schema changes
3. Write `*.down.sql` for rollback
4. Run `make migrate-up` to apply
5. Run `make sqlc-generate` to update Go code

## sqlc Code Generation

sqlc generates type-safe Go code from SQL queries.

```bash
# Regenerate code after query changes
make sqlc-generate

# Or manually:
sqlc generate
```

Configuration is in `sqlc.yaml`:
- Schema source: `migrations/` directory
- Queries: `query/query.sql`
- Output: `db/` package
- Driver: `pgx/v5`

## Testing with grpcurl

```bash
# Create a user
grpcurl -plaintext -d '{
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "+1234567890"
}' localhost:50052 main.UserService/CreateUser

# Get a user
grpcurl -plaintext -d '{"user_id": 1}' localhost:50052 main.UserService/GetUser

# Update a user
grpcurl -plaintext -d '{
    "user_id": 1,
    "first_name": "Jane",
    "last_name": "Smith",
    "phone_number": "+0987654321"
}' localhost:50052 main.UserService/UpdateUser

# Deactivate a user
grpcurl -plaintext -d '{"user_id": 1}' localhost:50052 main.UserService/DeactivateUser
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://grpc_sqlc:s3cr3t$password@postgres:5432/grpc_sqlc_db?sslmode=disable` | PostgreSQL connection string |
| `DB_USER` | `grpc_sqlc` | Database username |
| `DB_PASSWORD` | `s3cr3t$password` | Database password |
| `DB_NAME` | `grpc_sqlc_db` | Database name |
| `DB_HOST` | `postgres` | Database host |
| `DB_PORT` | `5432` | Database port |
| `SERVER_PORT` | `50052` | gRPC server port |

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make up` | Start all Docker services |
| `make down` | Stop all Docker services |
| `make logs` | View service logs |
| `make ps` | Show running containers |
| `make migrate-up` | Apply database migrations |
| `make migrate-down` | Rollback migrations |
| `make migrate-create name=<name>` | Create new migration |
| `make sqlc-generate` | Regenerate sqlc Go code |
| `make create_local_network` | Create Docker network |

## Database Schema

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL UNIQUE
);
```

## Dependencies

- [pgx/v5](https://github.com/jackc/pgx) - PostgreSQL driver
- [grpc-go](https://github.com/grpc/grpc-go) - gRPC implementation
- [sqlc](https://sqlc.dev/) - SQL code generator
- [golang-migrate](https://github.com/golang-migrate/migrate) - Database migrations
- [protobuf](https://developers.google.com/protocol-buffers) - Protocol buffers

## License

MIT
