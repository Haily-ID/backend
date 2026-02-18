# AGENTS.md - Haily Backend

## Project Overview

**Haily Backend** is an enterprise backend service supporting multiple front-end modules:

- Finance Management
- ERP (Enterprise Resource Planning)
- HR Management
- Inventory Management

Built with Clean Architecture for maintainability and scalability in complex enterprise environments.

---

## Tech Stack

**Core Technologies**:

- **Language**: Golang (latest stable)
- **Framework**: Echo
- **Database**: PostgreSQL + GORM
- **Cache**: Redis
- **Background Jobs**: Asynq (Redis-based distributed queue)
- **ID Generation**: Snowflake (64-bit distributed IDs)
- **Architecture**: Clean Architecture (4 layers)
- **Structure**: Monorepo
- **Dev Tool**: Air (hot reload for API & Worker)
- **Testing**: Go testing + testify
- **Config**: YAML + Environment Variables

**Development Principles**:

1. **Clean Architecture**: Strict layer separation (Domain → Repository → Use Case → Delivery)
2. **Avoid DRY**: Prefer clarity over abstraction, explicit code > clever code
3. **Performance First**: Redis caching, Asynq async jobs, connection pooling, Snowflake IDs
4. **Modular Design**: Monorepo with logical module separation
5. **API Design**: Simple response format `{data}` or `{error}`, specific error codes, semantic HTTP codes
6. **Minimal Comments**: Code should be self-explanatory, comment only complex logic
7. **Best Practices**: AutoMigrate, Air hot reload, semantic commits, PR-first workflow

---

## Project Structure

```
haily-backend/
├── cmd/
│   ├── api/main.go              # API server entry
│   └── worker/main.go           # Worker entry (separate binary)
├── internal/
│   ├── delivery/http/           # Presentation Layer
│   │   ├── handler/{module}/    # HTTP handlers per module
│   │   ├── middleware/          # Auth, logging, CORS, rate limit
│   │   └── route/               # Route definitions
│   ├── usecase/{module}/        # Business Logic Layer
│   ├── repository/              # Data Access Layer
│   │   ├── postgres/{module}/   # GORM repositories
│   │   └── redis/               # Redis cache
│   ├── domain/                  # Domain Layer (Core)
│   │   ├── entity/{module}/     # GORM models
│   │   ├── dto/{module}/        # API response DTOs
│   │   └── repository/          # Repository interfaces
│   └── pkg/                     # Internal shared packages
│       ├── config/              # Config loading
│       ├── database/            # GORM connection
│       ├── snowflake/           # ID generator
│       ├── asynq/               # Job queue
│       ├── logger/              # Logging
│       ├── validator/           # Validation
│       └── response/            # Response formatter
├── pkg/                         # Public packages
├── docs/                        # API docs (Swagger)
├── config/config.yaml           # Default config
├── go.mod
├── Makefile
├── .air.toml                    # API hot reload config
├── .air.worker.toml             # Worker hot reload config
├── docker-compose.yml
├── Dockerfile.api
├── Dockerfile.worker
├── .env.example
└── README.md
```

**Monorepo Benefits**: Shared code, consistent dependencies, atomic changes, easier refactoring, simplified CI/CD

---

## Clean Architecture Layers

### 1. Domain Layer (`internal/domain/`)

**Responsibility**: Core entities, DTOs, interfaces

**Components**:

- **Entity**: GORM models (Snowflake ID, timestamps, soft delete, relations)
- **DTO**: API responses (string IDs, Unix timestamps, `ToDTO()` method)
- **Repository Interface**: Data access contracts

**Principles**:

- ✅ Zero dependencies on other layers
- ✅ Snowflake ID (int64) for primary keys
- ✅ Separate DTOs from Entities
- ✅ Interface definitions only

**Entity Structure**:

- Primary Key: `ID int64` (Snowflake)
- Timestamps: `CreatedAt`, `UpdatedAt` (auto)
- Soft Delete: `DeletedAt gorm.DeletedAt`
- Relations: Foreign keys, Preload support

**DTO Structure**:

- `ID string` (converted from int64)
- Timestamps as Unix int64
- `ToDTO()` conversion method
- Flat structure for unnormalized response

### 2. Repository Layer (`internal/repository/`)

**Responsibility**: Data persistence and retrieval

**Postgres (GORM)**:

- Implement domain interfaces
- GORM query builder
- Redis caching (Cache-Aside)
- Context for cancellation
- Error mapping (GORM → domain errors)

**Redis**:

- Cache operations: Get, Set, Delete, Invalidate
- TTL per entity type
- Key format: `{entity}:{field}:{value}`
- JSON serialization

**Principles**:

- ✅ Implement domain interfaces
- ✅ NO business logic (pure data access)
- ✅ Cache hot data with proper TTL
- ✅ Handle GORM errors properly

**Cache Strategy**:

- Cache-Aside: Check cache → DB on miss → populate cache
- Write-Through: Update DB → invalidate cache
- TTL: Users (1h), static data (24h), reports (30m)

### 3. Use Case Layer (`internal/usecase/`)

**Responsibility**: Business logic and orchestration

**Components**:

- Use case struct with dependencies (repos, cache, validator, asynq)
- Execute method with request DTO
- Business validation
- Multi-repo orchestration
- Job enqueueing

**Principles**:

- ✅ ALL business logic here
- ✅ Orchestrate multiple repositories
- ✅ Validate business rules
- ✅ Enqueue background jobs
- ✅ Return entities (not DTOs)
- ✅ Independent & testable

**Flow Pattern**:

1. Validate input
2. Check preconditions
3. Execute business logic
4. Update database
5. Invalidate caches
6. Enqueue jobs (if needed)
7. Return result

### 4. Delivery Layer (`internal/delivery/http/`)

**Responsibility**: HTTP request/response handling

**Handler Components**:

- Handler struct with use case dependencies
- Endpoint methods
- Request binding
- Use case execution
- Entity → DTO conversion
- Unnormalized response formatting
- Error handling with HTTP codes

**Principles**:

- ✅ HTTP handling ONLY
- ✅ NO business logic
- ✅ Bind & validate format
- ✅ Call use cases
- ✅ Convert Entity → DTO
- ✅ Format unnormalized response
- ✅ Proper HTTP codes

**Request Flow**:

1. Bind request
2. Format validation
3. Call use case
4. Handle errors
5. Convert to DTO
6. Create response
7. Return JSON

**Middleware**: Auth (JWT), Authorization (RBAC), Logging, CORS, Rate Limiting, Recovery, Request ID

---

## Naming Conventions

**Files**: `snake_case` (`user_handler.go`)
**Packages**: `lowercase` singular (`package handler`)
**Structs**: `PascalCase` (`UserHandler`)
**Interfaces**: `PascalCase` (`UserRepository`)
**Methods**: `PascalCase` exported, `camelCase` unexported
**Constants**: `PascalCase` or `UPPER_SNAKE_CASE`

**Module Pattern** (consistent across all modules):

```
internal/
├── domain/entity/{module}/
├── domain/dto/{module}/
├── domain/repository/
├── repository/postgres/{module}/
├── usecase/{module}/
└── delivery/http/handler/{module}/
```

---

## Database Management

### PostgreSQL + GORM

**Connection**: DSN, connection pooling (MaxOpenConns, MaxIdleConns), PrepareStmt enabled

**AutoMigrate**:

- Run on startup: `db.AutoMigrate(&Entity1{}, &Entity2{}...)`
- Auto create/update tables from structs
- Safe (no delete columns/tables)
- Idempotent
- Limitations: Can't rename columns, no rollback

**Query Optimization**:

1. **Indexing**: `gorm:"index"` on frequently queried columns
2. **Avoid N+1**: Use `Preload()` or `Joins()`
3. **Pagination**: Always use `Offset()` + `Limit()`
4. **Select Fields**: Avoid loading unnecessary columns
5. **Raw SQL**: For complex queries

### Redis Caching

**Cache-Aside Pattern**:

1. Check cache (GET)
2. On hit: return cached
3. On miss: query DB → cache (SET with TTL) → return

**Cache Keys**: `{entity}:{field}:{value}`
Examples: `user:id:123`, `invoice:id:INV-001`

**TTL**: Users (1h), static (24h), reports (30m), session (configurable)

**Invalidation**: On update/delete, pattern matching with SCAN + DEL

---

## Snowflake ID Generation

**64-bit Distributed IDs**:

- ✅ K-sortable (by timestamp)
- ✅ Distributed (1-1024 machines)
- ✅ High performance (millions/sec)
- ✅ Contains timestamp
- ✅ Unique across instances

**Structure**: 1 bit unused + 41 bits timestamp + 10 bits machine ID + 12 bits sequence

**Setup**:

- Init with machine ID (1-1024)
- Set via ENV: `SNOWFLAKE_MACHINE_ID`
- Kubernetes: Use pod ordinal

**Usage**:

- Entity: `ID int64 gorm:"primaryKey;autoIncrement:false"`
- DTO: `ID string json:"id"` (converted)
- Conversion: `strconv.FormatInt()` / `strconv.ParseInt()`

---

## API Response Format

**Success Response** (200, 201):

```json
{
  "data": {
    /* entity atau array */
  }
}
```

**Error Response** (4xx, 5xx):

```json
{
  "error": "SPECIFIC_ERROR_CODE"
}
```

**Principles**:

- Simple, consistent structure
- `data` field only for success responses
- `error` field only for error responses
- Error codes are SPECIFIC to the problem (e.g., `USER_NOT_FOUND`, not just `NOT_FOUND`)
- No additional fields (no message, no meta)
- Client determines UI message based on error code

**Error Code Pattern**:

- Generic: `VALIDATION_ERROR`, `UNAUTHORIZED`, `FORBIDDEN`, `INTERNAL_SERVER_ERROR`
- Entity-specific: `{ENTITY}_{ACTION}_FAILED` - e.g., `USER_UPDATE_FAILED`, `COMPANY_DELETE_FAILED`
- Not found: `{ENTITY}_NOT_FOUND` - e.g., `USER_NOT_FOUND`, `COMPANY_NOT_FOUND`
- Already exists: `{FIELD}_ALREADY_EXISTS` - e.g., `EMAIL_ALREADY_EXISTS`
- Invalid input: `INVALID_{FIELD}` - e.g., `INVALID_USER_ID`, `INVALID_COMPANY_CODE`
- See `internal/pkg/response/response.go` for all available error codes

**HTTP Status Codes**:

- `200` OK, `201` Created, `204` No Content
- `400` Bad Request, `401` Unauthorized, `403` Forbidden, `404` Not Found
- `409` Conflict, `422` Unprocessable Entity, `429` Too Many Requests
- `500` Internal Error, `503` Service Unavailable

---

## Background Jobs (Asynq)

**Features**: Redis-based, automatic retries, priority queues, scheduled tasks, monitoring UI, graceful shutdown, dead letter queue

**Architecture**: Separate Worker binary (`cmd/worker/main.go`)

**Queue Priorities**:

- `critical` (60%): Payment, urgent notifications
- `default` (30%): Emails, regular tasks
- `low` (10%): Reports, cleanup

**Use Cases**: Email delivery, report generation, data processing, external API calls, cache warming, cleanup

**Usage Pattern**:

1. Execute critical DB ops
2. Create task payload
3. Enqueue with Asynq client
4. Return immediately
5. Worker processes async

**Options**: Queue selection, retry attempts, timeout, scheduled execution

---

## Configuration Management

**Single Source of Truth**:

- `config/config.yaml`: Default values + docs
- Environment Variables: Override for prod/staging
- Load order: YAML → ENV override

**Sections**: app, database, redis, jwt, snowflake, asynq

**ENV Variables**:

```bash
APP_ENV=production
APP_PORT=8080
DB_HOST=db.prod.com
DB_PASSWORD=secret
REDIS_HOST=redis.prod.com
JWT_SECRET=secret-key
SNOWFLAKE_MACHINE_ID=1
```

**Best Practices**: Never commit secrets, use `.env.example` template, document all options, validate on startup

---

## Authentication & Authorization

### JWT Authentication

**Flow**: Login → validate → generate JWT → return token → client includes in header → middleware validates

**Token Contents**: User ID, role, email, iat, exp

**Middleware**: Check header → parse JWT → validate signature → check expiration → extract claims → inject to context

### RBAC

**Roles**: admin, manager, user, guest

**Middleware**: Check role from context → compare with allowed roles → 403 if unauthorized

---

## Error Handling

**Domain Errors**: `ErrNotFound`, `ErrAlreadyExists`, `ErrUnauthorized`, `ErrForbidden`, `ErrBadRequest`, `ErrValidation`

**Response Pattern**:

- Handler: Map use case errors → specific error codes
- Use `response.Error(c, statusCode, errorCode)` with appropriate HTTP status and error constant
- Error codes defined in `internal/pkg/response/response.go`
- Always use specific error codes (e.g., `ErrUserNotFound`, not generic `ErrNotFound`)

**Pattern**:

- Use Case: Return domain errors with context
- Handler: Check error type/message → return specific error code
- Repository: GORM errors → domain errors

**Example**:

```go
if err != nil {
    if err.Error() == "email already registered" {
        return response.Error(c, http.StatusConflict, response.ErrEmailAlreadyExists)
    }
    return response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
}
```

---

## Testing

**Unit Tests**:

- Mock dependencies (testify/mock)
- Test use cases in isolation
- Table-driven tests for multiple scenarios

**Integration Tests**:

- Test database (in-memory or Docker)
- Test actual DB interactions
- Clean up after each test

**Commands**:

```bash
go test ./...                    # All tests
go test -v ./...                 # Verbose
go test -cover ./...             # Coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Development Workflow

**Setup**:

1. Clone repo
2. Copy `.env.example` to `.env`
3. `go mod download`
4. `docker-compose up -d`
5. `go install github.com/cosmtrek/air@latest`
6. `make dev-api` (API with hot reload)
7. `make dev-worker` (Worker with hot reload, separate terminal)

**Makefile Commands**:

```makefile
make run-api         # API without hot reload
make run-worker      # Worker without hot reload
make dev-api         # API with Air hot reload
make dev-worker      # Worker with Air hot reload
make build           # Build binaries
make test            # Run tests
make test-coverage   # Tests with coverage
make lint            # Run linter
make docker-up       # Start Docker services
make docker-down     # Stop Docker services
make clean           # Clean build artifacts
```

**Air Hot Reload**:

- `.air.toml` for API
- `.air.worker.toml` for Worker
- Auto-rebuild on file changes
- Fast iteration

---

## Git Workflow & Version Control

### Branch Strategy

**Main Branches**: `main` (production), `develop` (integration), `staging` (optional)

**Feature Branches** (Semantic):

```
feature/<desc>     # New features
bugfix/<desc>      # Bug fixes
hotfix/<desc>      # Urgent production fixes
refactor/<desc>    # Code refactoring
chore/<desc>       # Maintenance
docs/<desc>        # Documentation
test/<desc>        # Tests
perf/<desc>        # Performance
```

Examples: `feature/user-auth`, `bugfix/invoice-duplicate`, `hotfix/payment-error`

### Semantic Commits (Conventional Commits)

**Format**: `<type>(<scope>): <subject>`

**Types**: `feat`, `fix`, `refactor`, `perf`, `style`, `test`, `docs`, `chore`, `build`, `ci`, `revert`

**Examples**:

```bash
feat(auth): add JWT token refresh mechanism
fix(invoice): resolve duplicate number generation
refactor(user): separate DTO from entity model
perf(cache): implement Redis caching
docs(readme): add setup instructions
chore(deps): upgrade GORM to v1.25.0
```

**Scopes**: Module/area affected (auth, invoice, user, cache, database, api)

### PR Workflow (7 Steps)

1. **Create Branch**: `git checkout -b feature/name`
2. **Develop**: Make changes, commit with semantic messages
3. **Keep Updated**: `git rebase origin/develop`
4. **Push**: `git push origin feature/name`
5. **Create PR**: Semantic title, complete description, request reviewers
6. **Review**: Address comments, push updates
7. **Merge**: After approval + CI passes, squash and merge, delete branch

**PR Template**:

```markdown
## Description

[Brief description]

## Type of Change

- [ ] feat / fix / refactor / perf / docs

## Changes Made

- [List changes]

## Testing

- [ ] Unit tests
- [ ] Integration tests
- [ ] Manual testing

## Checklist

- [ ] Clean Architecture
- [ ] Self-review
- [ ] Tests pass
- [ ] No lint errors
```

**Code Review Checklist** (10 points):

1. Clean Architecture compliance
2. No business logic in handlers/repos
3. Proper error handling
4. Security vulnerabilities
5. Test coverage
6. Performance implications
7. Code readability
8. Proper logging
9. No memory/goroutine leaks
10. Follows "Avoid DRY" philosophy

**Review Comments**:

- 🔴 Critical (must fix)
- 🟡 Suggestion (should fix)
- 🔵 Nit (minor)
- ✅ Great (positive)

**Protected Branches** (`main`, `develop`):

- Require 1-2 PR reviews
- Require CI checks pass
- Require up-to-date branch
- No direct commits
- Optional: Linear history, signed commits

**Best Practices**:

- ✅ Commit often (small, focused)
- ✅ Descriptive messages
- ✅ Test before commit
- ✅ Self-review before PR
- ✅ Atomic commits
- ✅ Regular rebase
- ✅ Delete merged branches

**Don'ts**:

- ❌ Commit secrets
- ❌ Commit commented code
- ❌ "WIP" messages (use git stash)
- ❌ Force push to shared branches
- ❌ Large files (use Git LFS)
- ❌ Mixed changes
- ❌ Ignore failing tests

**CI/CD**: Auto-run tests, linter, coverage check, build API/Worker, security scan. Block PR if fails.

---

## AI Agent Guidelines

### Core Principles (13 Rules)

1. ✅ **Always Clean Architecture**: Strict layer separation
2. ✅ **No business logic in handlers**: HTTP only
3. ✅ **No business logic in repos**: Data access only
4. ✅ **All business logic in use cases**: Core logic here
5. ✅ **Prefer duplication**: Clarity > abstraction
6. ✅ **Test use cases thoroughly**: Critical reliability
7. ✅ **Always use context**: Cancellation & tracing
8. ✅ **Cache wisely**: Hot data, proper TTL
9. ✅ **Validate at use case**: Business rules
10. ✅ **Log important events**: Avoid over-logging
11. ✅ **Semantic commits**: Conventional format
12. ✅ **Feature branches**: Never commit to main/develop
13. ✅ **Always PR**: Human review required

### Git Workflow for AI

**Branch**: `git checkout -b <type>/<description>`

**Commit**: `git commit -m "feat(scope): clear description"`

**PR Process**:

1. Create semantic branch
2. Make changes (Clean Architecture)
3. Commit semantically
4. Push branch
5. Create PR with description
6. Wait for human review
7. Address comments
8. Wait for approval

**Never**:

- ❌ Commit to main/develop
- ❌ Force push to shared branches
- ❌ Merge without review
- ❌ Generic commit messages
- ❌ Skip tests

### Code Generation Guidelines

**Code Quality**:

- Use explicit types (avoid `interface{}` where possible)
- Comprehensive error handling
- Follow naming conventions
- Validate at use case level
- Convert Entity → DTO in handlers
- Proper HTTP status codes with specific error codes

**Comments**:

- ✅ **Minimal comments** - code should be self-explanatory
- ✅ Comment ONLY for complex algorithms or non-obvious business logic
- ✅ NO comments for obvious operations (e.g., "// Create user", "// Validate input")
- ✅ Package documentation at top of main packages
- ❌ NO commented-out code (delete it, Git keeps history)
- ❌ NO TODO comments (use issue tracker)

**Logging**:

- Log important operations (create, update, delete)
- Log errors with context
- Avoid over-logging (no verbose debug logs in production)

**Performance**:

- Consider caching for frequent reads
- GORM transactions for multi-step operations
- Use connection pooling (already configured)

### Common Scenarios

**New Endpoint**:

1. Check if use case exists
2. Create use case (if needed)
3. Create handler method
4. Register route (+ middleware)
5. Add tests

**New Entity**:

1. Create domain entity (GORM, Snowflake ID)
2. Create DTO (ToDTO method)
3. Create repository interface
4. Implement repository (GORM + cache)
5. Create use cases
6. Create handlers
7. Register routes
8. AutoMigrate handles schema

**Performance Issue**:

1. Check caching
2. Verify indexes
3. Check N+1 queries
4. Optimize queries
5. Add monitoring
6. Consider background jobs

---

## Additional Topics

### API Documentation (Swagger)

- Use Swaggo for auto-generation
- Add annotations to handlers
- Generate: `swag init -g cmd/api/main.go`
- Access: `http://localhost:8080/swagger/index.html`

### Performance Optimization

**Database**: Indexing, avoid N+1, connection pooling, prepared statements, pagination, select specific fields

**API**: Caching, compression, rate limiting, background jobs, HTTP/2, unnormalized responses

### Security Best Practices

1. ✅ Never commit secrets
2. ✅ Parameterized queries (GORM)
3. ✅ Validate all inputs
4. ✅ JWT with expiration
5. ✅ RBAC authorization
6. ✅ HTTPS only in prod
7. ✅ Rate limiting
8. ✅ CORS config
9. ✅ Keep dependencies updated
10. ✅ Don't leak info in errors

### Monitoring & Health Checks

**Metrics**: Request rate, response time (p50/p95/p99), error rate, pool usage, cache hit/miss, job queue size

**Health Endpoint** (`GET /health`):

```json
{
  "status": "healthy",
  "database": "up",
  "redis": "up",
  "version": "1.0.0"
}
```

### Deployment

**Build**:

```bash
# API
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api cmd/api/main.go

# Worker
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/worker cmd/worker/main.go

# Docker
docker build -f Dockerfile.api -t haily-api:latest .
docker build -f Dockerfile.worker -t haily-worker:latest .
```

**Services**: API (HTTP), Worker (jobs), PostgreSQL (DB), Redis (cache/queue)

**Scaling**: API horizontally (load balancer), Worker by queue size, unique Snowflake machine ID per instance

### Module Development

**Adding New Module** (e.g., Inventory):

1. Domain entities (GORM, Snowflake, relations)
2. DTOs (ToDTO method)
3. Repository interface
4. Implement repository (GORM + cache)
5. Use cases (business logic)
6. Handlers (endpoints)
7. Register routes
8. AutoMigrate

**Consistency**: Follow same pattern for all modules (finance, ERP, HR, inventory)

### "Avoid DRY" Philosophy

**Prefer clarity over cleverness**:

- Explicit code > clever abstractions
- Duplication > wrong abstraction
- Easy debugging > reduced lines
- Independent modules > tight coupling

**Allow Duplication**: Handler per entity (UserHandler, InvoiceHandler), use case per action, repository per entity

**Apply DRY**: Truly generic utilities (string/date utils, middleware, config loading, common packages, DB setup, error types)

**Rule**: Generic & stable → extract. Domain-specific or likely to diverge → duplicate.

### Common Patterns

**Repository**: Return entities, handle DB ops, implement caching, no business logic, context for cancellation, error mapping

**Use Case**: All business logic, orchestrate repos, validate rules, handle transactions, enqueue jobs, manage caching, testable

**Handler**: HTTP only, bind requests, call use cases, convert to DTO, format response, map errors to HTTP codes, no business logic
