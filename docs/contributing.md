# Contributing to swagger-to-http

Thank you for your interest in contributing to the `swagger-to-http` project! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Pull Request Process](#pull-request-process)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)

## Code of Conduct

Please be respectful and considerate of others when contributing to this project. We aim to foster an inclusive and welcoming community.

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git

### Setting Up the Development Environment

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/swagger-to-http.git
   cd swagger-to-http
   ```

3. Add the original repository as upstream:
   ```bash
   git remote add upstream https://github.com/edgardnogueira/swagger-to-http.git
   ```

4. Install dependencies:
   ```bash
   go mod download
   ```

### Building the Project

To build the project:

```bash
make build
```

This will create the `swagger-to-http` binary in the current directory.

### Running Tests

To run all tests:

```bash
make test
```

To run tests with coverage:

```bash
make test-coverage
```

## Project Structure

The project follows Clean Architecture principles and is organized as follows:

```
swagger-to-http/
├── cmd/                      # Command-line entry points
│   └── swagger-to-http/      # Main application
│       └── main.go
├── internal/                 # Private application code
│   ├── domain/               # Domain models and business rules
│   │   └── models/           # Core data structures
│   ├── application/          # Application services and use cases
│   │   ├── generator/        # HTTP file generation
│   │   ├── parser/           # Swagger parsing
│   │   └── snapshot/         # Snapshot testing
│   ├── cli/                  # Command-line interface logic
│   └── infrastructure/       # Infrastructure implementations
│       └── fs/               # File system operations
├── pkg/                      # Public packages for external use
├── docs/                     # Documentation
├── test/                     # Test files and fixtures
├── .github/                  # GitHub workflows and templates
└── ...
```

## Development Workflow

We follow a [Git Flow](https://nvie.com/posts/a-successful-git-branching-model/) branching model:

1. **Main Branch**: Production-ready code
2. **Develop Branch**: Latest development changes
3. **Feature Branches**: New features or improvements
4. **Release Branches**: Preparing for a new release
5. **Hotfix Branches**: Urgent fixes for production

### Creating a Feature Branch

```bash
git checkout develop
git pull upstream develop
git checkout -b feature/your-feature-name
```

### Committing Changes

Follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages:

```
<type>: <description>

[optional body]

[optional footer]
```

Types include:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Example:
```
feat: add XML formatter for snapshot comparison

Implement XML normalization and comparison for snapshot testing.
The formatter handles element order and whitespace differences.

Closes #123
```

## Coding Standards

### Go Coding Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines

### Clean Architecture Principles

- Keep domain models independent of frameworks
- Business logic should not depend on UI, database, or external agencies
- Organize code by layers: domain, application, infrastructure
- Use dependency injection and interfaces for loose coupling

### Package Organization

- Group related functionality in packages
- Keep package names short and meaningful
- Avoid package name collisions with standard library

## Pull Request Process

1. **Update your fork** with the latest changes from upstream:
   ```bash
   git checkout develop
   git pull upstream develop
   ```

2. **Create a feature branch** for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes** following the coding standards.

4. **Write tests** for your changes.

5. **Run the tests** to ensure they pass:
   ```bash
   make test
   ```

6. **Update documentation** if necessary.

7. **Commit your changes** with a descriptive commit message.

8. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

9. **Create a pull request** against the `develop` branch.

10. **Address review comments** if any.

11. Once approved, your PR will be merged.

## Testing Guidelines

### Unit Tests

- Write tests for all new functionality
- Use table-driven tests when appropriate
- Aim for high code coverage
- Keep tests independent and idempotent

### Test Organization

- Test files should be named `*_test.go`
- Tests should be in the same package as the code they test
- Use descriptive test names: `TestFunctionName_Scenario`

### Test Utilities

- Use `testify/assert` for assertions
- Use mocks for external dependencies
- Create test fixtures in the `test/` directory

## Documentation

Good documentation is crucial for the project:

- Update README.md for significant changes
- Add godoc comments to all exported functions, types, and methods
- Create or update documentation in the `docs/` directory
- Include examples where applicable

### Documentation Format

- Use Markdown for all documentation
- Include code examples where appropriate
- Keep documentation up to date with code changes

## Issue Tracking

If you're working on an issue, please follow these guidelines:

1. **Comment on the issue** you want to work on to avoid duplication
2. **Reference the issue number** in your commit messages and PR
3. **Update the issue status** as you make progress
4. **Close the issue** when it's resolved

## Questions and Help

If you need help or have questions:

- Check existing documentation
- Look for similar issues that may have been resolved
- Open a new issue with a clear description of your question

Thank you for contributing to `swagger-to-http`!
