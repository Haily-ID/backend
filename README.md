# Haily Backend

Enterprise backend service built with Clean Architecture supporting multiple modules: User Management, Company Management, and more.

## Tech Stack

- **Language**: Golang
- **Framework**: Echo
- **Database**: PostgreSQL + GORM
- **Cache**: Redis
- **Background Jobs**: Asynq
- **ID Generation**: Snowflake
- **Architecture**: Clean Architecture

## Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Air (for hot reload)

### Setup

1. Clone the repository:

```bash
git clone https://github.com/Haily-ID/backend.git
cd backend
```

2. Copy environment file:

```bash
cp .env.example .env
```

3. Start infrastructure services:

```bash
make docker-up
```

4. Install Air for hot reload:

```bash
go install github.com/cosmtrek/air@latest
```

5. Run API server:

```bash
make dev-api
```

6. Run worker (in another terminal):

```bash
make dev-worker
```

## Project Structure

```
haily-backend/
├── cmd/                    # Application entry points
│   ├── api/               # API server
│   └── worker/            # Background worker
├── internal/              # Private application code
│   ├── delivery/http/     # HTTP handlers & routes
│   ├── usecase/           # Business logic
│   ├── repository/        # Data access
│   ├── domain/            # Domain entities & DTOs
│   └── pkg/               # Internal packages
├── config/                # Configuration files
└── docs/                  # Documentation
```

## Available Commands

```bash
make help              # Show all commands
make run-api           # Run API server
make run-worker        # Run worker
make dev-api           # Run API with hot reload
make dev-worker        # Run worker with hot reload
make build             # Build binaries
make test              # Run tests
make test-coverage     # Run tests with coverage
make lint              # Run linter
make docker-up         # Start Docker services
make docker-down       # Stop Docker services
make clean             # Clean build artifacts
```

## Features

- ✅ User Management
- ✅ Company Management
- ✅ User can create multiple companies
- ✅ User can join other companies

## Architecture

This project follows Clean Architecture principles with 4 layers:

1. **Domain Layer**: Core entities, DTOs, interfaces
2. **Repository Layer**: Data persistence
3. **Use Case Layer**: Business logic
4. **Delivery Layer**: HTTP handlers

## Development

See [AGENTS.md](AGENTS.md) for detailed development guidelines and architecture documentation.

## License

MIT
