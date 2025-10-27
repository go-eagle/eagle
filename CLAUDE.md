# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Eagle is a Go microservice framework suitable for rapid business development. It can quickly build API services or Web sites with support for HTTP and gRPC protocols.

## Essential Commands

### Development
```bash
# Run the main application
make run

# Build the binary (with version info)
make build

# Run tests
make test

# Run a single test
go test -v -run TestFunctionName ./path/to/package

# Run tests with coverage
make cover

# View test coverage in browser
make view-cover
```

### Code Quality
```bash
# Run linter (uses golangci-lint)
make lint

# Format code
make fmt

# Run vet
make vet
```

### Code Generation
```bash
# Generate Swagger documentation
make docs
# Access docs at: http://localhost:8080/swagger/index.html

# Generate mock files
make mockgen
```

### Other
```bash
# Clean build artifacts
make clean

# Build Docker image
make docker

# Generate CA certificates
make ca

# Generate call graph visualization
make graph
```

### Eagle CLI Tool
The project provides a CLI tool for code generation:
```bash
# Install the CLI
go install github.com/go-eagle/eagle/cmd/eagle@latest

# Create a new project
eagle new project-name

# Generate code components
eagle handler <name>    # Generate handler
eagle service <name>    # Generate service
eagle repo <name>       # Generate repository
eagle cache <name>      # Generate cache
eagle proto <name>      # Generate proto files
eagle task <name>       # Generate task
eagle model <name>      # Generate model
```

## Architecture

### Layered Structure

Eagle follows a classic layered architecture with dependency injection:

```
Handler (API Layer)
    ↓
Service (Business Logic)
    ↓
Repository (Data Access Interface)
    ↓
Model/Cache (Data Layer)
```

### Key Directories

- **internal/handler**: HTTP/gRPC request handlers, receives user requests
- **internal/service**: Business logic layer, contains core business rules
- **internal/repository**: Repository layer, wraps data access interfaces (database, cache, RPC)
- **internal/model**: Database models and GORM definitions
- **internal/cache**: Cache layer implementations
- **internal/routers**: Route registration and middleware setup
- **internal/server**: Server initialization (HTTP/gRPC)
- **internal/middleware**: Custom middleware implementations
- **internal/ecode**: Error code definitions
- **pkg**: Reusable packages (can be imported by external projects)

### Application Initialization Flow

1. Parse flags and load configuration from `config/` directory using Viper
2. Initialize logger (zap)
3. Initialize database connections (GORM)
4. Initialize Redis (optional)
5. Build service dependencies: `service.New(repository.New(db))`
6. Create HTTP server with Gin router
7. Start app with graceful shutdown support

### Router and Middleware

Routes are defined in `internal/routers/router.go`. The framework applies middleware in this order:
1. Recovery (panic handling)
2. NoCache, Options, Secure (security headers)
3. Logging
4. RequestID
5. Metrics (Prometheus)
6. Tracing (OpenTelemetry)
7. Timeout (default 3s)
8. Translations (i18n)
9. Auth (JWT - applied per route group)

### Configuration

- Configuration files are in `config/` directory
- Supports multiple environments via `-e` flag
- Uses Viper for config parsing
- Main config loaded into `eagle.Config` struct

### Dependency Injection

The framework uses manual dependency injection:
- Services are initialized with repository dependencies
- Repositories are initialized with database/cache connections
- Global service instance: `service.Svc`

## Testing

- Uses GoConvey for testing framework
- Test files should be named `*_test.go`
- Coverage reports generated to `coverage.txt`
- CI runs tests automatically via GitHub Actions

## Common Patterns

### Error Handling
- Use error codes defined in `internal/ecode`
- Return errors with context using `pkg/errors.Wrap()`
- Handlers should return appropriate HTTP status codes

### Logging
- Use structured logging via `pkg/log` (wraps zap)
- Logger initialized globally, access via `log.GetLogger()`
- Include context and request IDs in logs

### Database Operations
- Use GORM for ORM operations
- Models defined in `internal/model`
- Repositories provide abstraction over data access
- Support for MySQL, PostgreSQL, MongoDB, ClickHouse

### Caching
- Redis integration via `pkg/redis`
- Local cache support via Ristretto
- Cache operations abstracted in `internal/cache`

### API Documentation
- Use Swagger annotations in handler functions
- Generate docs with `make docs`
- Main API info in `main.go` annotations

## Additional Notes

- Uses gitmoji for commit messages (see badge in README)
- Supports distributed tracing with OpenTelemetry/Jaeger
- Prometheus metrics exposed at `/metrics` endpoint
- Health check available at `/health` endpoint
- PProf profiling available when enabled (in development mode)
- Main branch is `master` (use this as base for PRs)
- Go version: 1.22+
