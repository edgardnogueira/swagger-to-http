# swagger-to-http Development Documentation

This document provides key information about the project structure, architecture, and ongoing development to facilitate continuing work across multiple chat sessions.

## Project Overview

swagger-to-http is a tool that converts Swagger/OpenAPI specifications into .http files for easy API testing. It also provides snapshot testing capabilities to compare API responses.

## Architecture

The project follows Clean Architecture principles with the following layers:

- **Domain**: Core business models and interfaces (internal/domain/models)
- **Application**: Use cases and business logic (internal/application)
- **Infrastructure**: Implementation details for external systems (internal/infrastructure)
- **CLI**: Command-line interface (internal/cli)

## Key Components

1. **Swagger Parser**: Parses Swagger/OpenAPI documents into domain models
2. **HTTP Generator**: Generates HTTP files from Swagger documents
3. **HTTP Executor**: Executes HTTP requests and returns responses (issue #10)
4. **Snapshot Manager**: Manages response snapshots for testing
5. **CLI Commands**: User interface for interacting with the tool

## Current Development Status

### Completed Features
- Swagger/OpenAPI parser
- HTTP file generation
- Basic CLI commands
- Snapshot system

### In Progress
- HTTP Executor (issue #10)
  - Implementation of the HTTP execution engine
  - HTTP file parser
  - CLI integration

### Upcoming Features
- Test interface and reporting (issue #12)
- Advanced testing features (issue #13)
- Git hooks integration (issue #8)

## Development Workflow

1. Use Git Flow with feature branches
2. Branch naming pattern: `feature/issue-X-description`
3. Run tests before committing: `go test ./...`
4. Format code: `go fmt ./...`
5. Commit with conventional commit messages: `feat:`, `fix:`, `docs:`, etc.
6. Create pull requests for merging features

## Key Interfaces

The core interfaces are defined in `internal/application/interfaces.go`:

- `SwaggerParser`: Parsing Swagger/OpenAPI documents
- `HTTPGenerator`: Generating HTTP requests from Swagger
- `FileWriter`: Writing HTTP files to the filesystem
- `HTTPExecutor`: Executing HTTP requests
- `SnapshotManager`: Managing response snapshots
- `ConfigProvider`: Retrieving configuration

## Project Structure

```
swagger-to-http/
├── cmd/                      # Entry points
│   └── swagger-to-http/      # Main application
├── internal/                 # Private application code
│   ├── domain/               # Domain models
│   │   └── models/           # Core data structures
│   ├── application/          # Application use cases
│   │   ├── parser/           # Swagger parsing
│   │   ├── generator/        # HTTP generation
│   │   └── snapshot/         # Snapshot management
│   ├── infrastructure/       # External adapters
│   │   ├── http/             # HTTP execution (new)
│   │   ├── fs/               # Filesystem operations
│   │   └── config/           # Configuration
│   └── cli/                  # Command-line interface
└── docs/                     # Documentation
```
