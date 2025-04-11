# Contributing to Swagger-to-HTTP

Thank you for considering contributing to Swagger-to-HTTP! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Set up the development environment
4. Create a feature branch for your changes
5. Make your changes
6. Submit a pull request

## Development Environment Setup

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, but recommended)

### Installation

1. Clone your fork:

```bash
git clone https://github.com/YOUR_USERNAME/swagger-to-http.git
cd swagger-to-http
```

2. Add the upstream repository:

```bash
git remote add upstream https://github.com/edgardnogueira/swagger-to-http.git
```

3. Install dependencies:

```bash
go mod download
```

4. Install development tools:

```bash
make lint-tools
```

## Development Workflow

We follow the Git Flow branching model:

- `main`: Production-ready code
- `develop`: Latest development changes
- `feature/*`: Feature branches (branched from `develop`)
- `release/*`: Release branches
- `hotfix/*`: Hotfix branches (branched from `main`)

### Creating a new feature

1. Ensure your fork is up-to-date:

```bash
git checkout develop
git pull upstream develop
```

2. Create a feature branch:

```bash
git checkout -b feature/your-feature-name
```

3. Make your changes, following the [coding standards](#coding-standards)

4. Commit your changes with semantic commit messages:

```bash
git commit -m "feat: add new feature"
```

5. Push your changes:

```bash
git push origin feature/your-feature-name
```

6. Create a pull request to the `develop` branch

## Coding Standards

We follow Clean Architecture principles and idiomatic Go coding practices:

### Clean Architecture

The project is organized into the following layers:

- **Domain** (`internal/domain`): Core business logic and entities
- **Application** (`internal/application`): Use cases and interfaces
- **Infrastructure** (`internal/infrastructure`): Implementation details
- **Interfaces** (`cmd/`, `internal/cli`): User interfaces and entry points

### Go Standards

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` or `goimports` to format your code
- Write meaningful godoc comments for exported functions and types
- Use meaningful variable and function names
- Keep functions small and focused

### Commit Messages

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code changes that neither fix bugs nor add features
- `perf`: Performance improvements
- `test`: Adding or modifying tests
- `chore`: Changes to the build process or auxiliary tools
- `ci`: Changes to CI configuration files and scripts

## Testing

- Write tests for all new features and bug fixes
- Ensure all tests pass before submitting a pull request
- Run tests:

```bash
make test
```

- Check test coverage:

```bash
make cover
```

## Submitting a Pull Request

1. Ensure your code passes all tests
2. Update documentation as needed
3. Create a pull request to the `develop` branch
4. Include a clear description of your changes
5. Link to any related issues

## Release Process

Releases are managed by maintainers following these steps:

1. Create a release branch from `develop`
2. Finalize and test the release
3. Merge to `main` and tag with a version
4. Update the `develop` branch with any changes

## Getting Help

If you need help, please:

- Check the documentation
- Open an issue on GitHub
- Ask questions in discussions

Thank you for contributing to Swagger-to-HTTP!
