# Swagger to HTTP

A tool to convert Swagger/OpenAPI documentation into organized HTTP request files with snapshot testing capabilities.

[![Go Report Card](https://goreportcard.com/badge/github.com/edgardnogueira/swagger-to-http)](https://goreportcard.com/report/github.com/edgardnogueira/swagger-to-http)
[![GoDoc](https://godoc.org/github.com/edgardnogueira/swagger-to-http?status.svg)](https://godoc.org/github.com/edgardnogueira/swagger-to-http)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- Convert Swagger/OpenAPI 2.0 and 3.0 documents to .http files
- Organize requests by tags and endpoints
- Support for path, query, and body parameters
- Automatically generate example request bodies based on schemas
- Execute HTTP requests and compare responses with snapshots
- Snapshot testing with content-type aware comparison
- Integration with Git hooks for automatic updates (coming soon)

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Configuration](#configuration)
- [Snapshot Testing](#snapshot-testing)
- [Examples](#examples)
- [Project Status](#project-status)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

## Installation

### Using Go

```bash
go install github.com/edgardnogueira/swagger-to-http@latest
```

### Binary Releases

Download the pre-built binaries from the [Releases](https://github.com/edgardnogueira/swagger-to-http/releases) page.

### Using Docker

```bash
docker pull edgardnogueira/swagger-to-http
docker run -v $(pwd):/app edgardnogueira/swagger-to-http generate -f /app/swagger.json -o /app/http-requests
```

For more detailed installation instructions, see the [Installation Guide](docs/installation.md).

## Quick Start

Generate HTTP files from a Swagger file:

```bash
swagger-to-http generate -f swagger.json -o http-requests
```

Or from a URL:

```bash
swagger-to-http generate -u https://petstore.swagger.io/v2/swagger.json -o http-requests
```

Run snapshot tests on your API:

```bash
swagger-to-http snapshot test http-requests/*.http
```

## Usage

```
Usage:
  swagger-to-http [command]

Available Commands:
  generate    Generate HTTP files from a Swagger/OpenAPI document
  help        Help about any command
  snapshot    Snapshot testing commands
  version     Print the version information

Flags:
  -h, --help      help for swagger-to-http
  -v, --version   version for swagger-to-http
```

### Generate Command

```
Usage:
  swagger-to-http generate [flags]

Flags:
  -f, --file string         Swagger/OpenAPI file to process (required if url not provided)
  -u, --url string          URL to Swagger/OpenAPI document (required if file not provided)
  -o, --output string       Output directory for HTTP files (default "http-requests")
  -b, --base-url string     Base URL for requests (overrides the one in the Swagger doc)
  -t, --default-tag string  Default tag for operations without tags (default "default")
  -i, --indent-json         Indent JSON in request bodies (default true)
      --auth                Include authentication header in requests
      --auth-header string  Authentication header name (default "Authorization")
      --auth-token string   Authentication token value
  -h, --help                help for generate
```

### Snapshot Commands

```
Usage:
  swagger-to-http snapshot [command]

Available Commands:
  test        Run snapshot tests
  update      Update snapshots
  list        List snapshots
  cleanup     Cleanup snapshots

Flags:
  -h, --help   help for snapshot
```

For more detailed usage information, see the [Usage Guide](docs/usage.md).

## Configuration

swagger-to-http uses the following configuration file lookup paths:
- `./swagger-to-http.yaml`
- `$HOME/.swagger-to-http/swagger-to-http.yaml`
- `/etc/swagger-to-http/swagger-to-http.yaml`

You can also use environment variables with the prefix `STH_` (e.g., `STH_OUTPUT_DIRECTORY`).

Example configuration file:

```yaml
output:
  directory: http-requests
generator:
  indent_json: true
  include_auth: false
  auth_header: Authorization
  default_tag: default
snapshots:
  directory: .snapshots
  update_mode: none
  ignore_headers:
    - Date
    - Set-Cookie
```

For more configuration options, see the [Configuration Guide](docs/configuration.md).

## Snapshot Testing

Snapshot testing allows you to:

1. Execute HTTP requests in .http files
2. Save responses as snapshots
3. Compare future responses against stored snapshots
4. Detect changes in API behavior

### Creating Snapshots

```bash
# Create or update all snapshots
swagger-to-http snapshot update "api/*.http"

# Create specific snapshots
swagger-to-http snapshot update "api/users.http"
```

### Running Tests

```bash
# Test all HTTP files
swagger-to-http snapshot test "api/*.http"

# Test specific files
swagger-to-http snapshot test "api/users.http"

# Update failed snapshots automatically
swagger-to-http snapshot test --update failed "api/*.http"
```

For more information on snapshot testing, see the [Snapshot Testing Guide](docs/snapshot-testing.md).

## Examples

We provide various examples to help you get started:

- [Basic Usage](docs/examples/basic-usage/README.md)
- [Authentication](docs/examples/auth/)
- [Complex Parameters](docs/examples/parameters/)
- [Snapshot Testing](docs/examples/snapshot-testing/)

## Project Status

This project is in active development. The following features are implemented or planned:

- [x] Basic Swagger/OpenAPI parsing
- [x] HTTP file generation
- [x] CLI interface
- [x] HTTP request execution
- [x] Response snapshot comparison
- [ ] Git hooks integration
- [ ] Schema validation

## Documentation

Comprehensive documentation is available in the [docs](docs/) directory:

- [Installation Guide](docs/installation.md)
- [Usage Guide](docs/usage.md)
- [Configuration Guide](docs/configuration.md)
- [HTTP File Format](docs/http-file-format.md)
- [Snapshot Testing](docs/snapshot-testing.md)
- [Examples](docs/examples/)
- [API Reference](docs/api-reference.md)
- [Contributing Guide](docs/contributing.md)
- [FAQ](docs/faq.md)

## Contributing

Contributions are welcome! Please read our [Contributing Guide](docs/contributing.md) for details on how to submit pull requests, the development process, and coding standards.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
