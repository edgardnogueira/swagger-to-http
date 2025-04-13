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
- Variable substitution in URLs, headers, and request bodies
- Environment variable support for configurable requests
- Integration with Git hooks for automatic updates (coming soon)
- **NEW**: Schema validation against OpenAPI specifications
- **NEW**: Sequential tests with dependencies and variable extraction
- **NEW**: Test assertions with customizable validation rules
- **NEW**: Continuous testing in watch mode

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Configuration](#configuration)
- [HTTP Executor](#http-executor)
- [Snapshot Testing](#snapshot-testing)
- [Advanced Testing Features](#advanced-testing-features)
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

Run tests with schema validation:

```bash
swagger-to-http test validate --swagger-file swagger.json http-requests/*.http
```

Execute sequence tests with variable extraction:

```bash
swagger-to-http test sequence sequence-tests/*.json
```

## Usage

```
Usage:
  swagger-to-http [command]

Available Commands:
  generate    Generate HTTP files from a Swagger/OpenAPI document
  help        Help about any command
  snapshot    Snapshot testing commands
  test        Run HTTP tests
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

### Test Commands

```
Usage:
  swagger-to-http test [command]

Available Commands:
  list        List available HTTP tests
  validate    Validate responses against OpenAPI schema
  sequence    Run test sequences with dependency support

Flags:
  --update string         Update mode: none, all, failed, missing (default "none")
  --ignore-headers string Comma-separated headers to ignore in comparison (default "Date,Set-Cookie")
  --snapshot-dir string   Directory for snapshot storage (default ".snapshots")
  --fail-on-missing       Fail when snapshot is missing
  --cleanup               Remove unused snapshots after testing
  --timeout duration      HTTP request timeout (default 30s)
  --parallel              Run tests in parallel
  --max-concurrent int    Maximum number of concurrent tests (default 5)
  --stop-on-failure       Stop testing after first failure
  --tags strings          Filter tests by tags
  --methods strings       Filter tests by HTTP methods
  --paths strings         Filter tests by request paths
  --names strings         Filter tests by test names
  --report-format string  Report format: console, json, html, junit (default "console")
  --report-output string  Path to write report file
  --detailed              Include detailed information in report
  --watch                 Run in continuous (watch) mode
  --watch-interval int    Interval between watch checks in milliseconds (default 1000)
  -h, --help              help for test
```

### Schema Validation Command

```
Usage:
  swagger-to-http test validate [file-patterns]

Flags:
  --swagger-file string    Path to Swagger/OpenAPI file (required)
  --ignore-props string    Comma-separated properties to ignore in validation
  --ignore-add-props       Ignore additional properties not in schema
  --ignore-formats         Ignore format validation (e.g., date, email)
  --ignore-patterns        Ignore pattern validation
  --req-props-only         Validate only required properties
  --ignore-nullable        Ignore nullable field validation
```

### Test Sequence Command

```
Usage:
  swagger-to-http test sequence [file-patterns]

Flags:
  --variables-path string  Path to load/save variables
  --save-vars              Save extracted variables to file
  --var-format string      Variable format (default: ${varname})
  --fail-fast              Stop sequence on first failure
  --validate-schema        Validate responses against schema
  --swagger-file string    Path to Swagger/OpenAPI file
```

## Configuration

swagger-to-http uses the following configuration file lookup paths:
- `./swagger-to-http.yaml`
- `$HOME/.swagger-to-http/swagger-to-http.yaml`
- `/etc/swagger-to-http/swagger-to-http.yaml`

You can also use environment variables with the prefix `STH_`  (e.g., `STH_OUTPUT_DIRECTORY`).

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
advanced_testing:
  validate_schema: false
  schema_validation:
    ignore_additional_properties: false
    ignore_formats: false
    ignore_patterns: false
    required_properties_only: false
  test_sequences:
    variable_format: "${%s}"
    save_variables: false
    fail_fast: false
```

## HTTP Executor

The HTTP Executor is responsible for executing HTTP requests defined in `.http` files. It supports:

- Variable substitution in URLs, headers, and request bodies using `{{VARIABLE_NAME}}` syntax
- Environment variables support with system environment variables (prefixed with `HTTP_`)
- All standard HTTP methods (GET, POST, PUT, DELETE, PATCH, etc.)
- Custom headers and authentication
- Timeout configuration
- Content-type aware processing

Example of an HTTP file with variables:

```http
# Get user details
@name GetUser
@tag users

GET {{BASE_URL}}/api/users/{{USER_ID}}
Accept: application/json
Authorization: Bearer {{TOKEN}}

```

Key features of the HTTP Executor include:

1. **Request Execution**: Execute individual or multiple HTTP requests from `.http` files
2. **Variable Handling**: Support for environment variables and request-specific variables with priority order:
   - Request-specific variables (highest priority)
   - Variables passed to the executor
   - System environment variables (lowest priority)
3. **Content-Type Management**: Proper handling of different content types for requests/responses
4. **Error Handling**: Comprehensive error reporting for network issues, timeouts, and HTTP errors
5. **Integration with Snapshot System**: Seamless connection with snapshot testing functionality

The HTTP Executor can be used programmatically:

```go
// Create an HTTP executor with a 30-second timeout
executor := http.NewExecutor(30*time.Second, nil)

// Execute a request
response, err := executor.Execute(ctx, request, nil)
if err != nil {
    // Handle error
}

// Use response
fmt.Printf("Status: %d\n", response.StatusCode)
fmt.Printf("Body: %s\n", string(response.Body))
```

For more details on the HTTP Executor, see the [HTTP Executor Documentation](docs/http-executor.md).

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

### Update Modes

- `none`: Do not update any snapshots (default)
- `all`: Update all snapshots regardless of test result
- `failed`: Update only snapshots that failed comparison
- `missing`: Create snapshots only if they don't exist

### Snapshot Formatters

The snapshot system includes content-type aware formatters for:

- JSON - Normalizes and prettifies JSON for reliable comparison
- XML - Handles XML structure comparison
- HTML - Compares HTML responses
- Plain text - Basic text comparison
- Binary data - Compares binary data with diff visualization

### Ignoring Headers

Some HTTP headers will change between requests (like timestamps). You can ignore specific headers:

```bash
swagger-to-http snapshot test --ignore-headers "Date,Set-Cookie,X-Request-ID" "api/*.http"
```

### Managing Snapshots

```bash
# List all snapshots
swagger-to-http snapshot list

# List snapshots in a specific directory
swagger-to-http snapshot list api/users

# Clean up unused snapshots
swagger-to-http snapshot cleanup

```

For more configuration options, see the [Configuration Guide](docs/configuration.md).

For more information on snapshot testing, see the [Snapshot Testing Guide](docs/snapshot-testing.md).

## Advanced Testing Features

This new section includes features added in the latest update that enhance the testing capabilities of the tool.

### Schema Validation

Validate HTTP responses against OpenAPI/Swagger schema definitions:

```bash
swagger-to-http test validate --swagger-file swagger.json http-requests/*.http
```

Features:
- Type validation (string, number, boolean, object, array)
- Required property checking
- Pattern matching
- Format validation (date, email, etc.)
- Configurable validation rules with options to ignore specific validations

### Test Sequences

Run tests in a specific order with dependencies:

```bash
swagger-to-http test sequence sequence-tests/*.json
```

Test sequence files are JSON files that define:
- A series of HTTP requests in a specific order
- Variable extractions to pass data between requests
- Wait times between requests
- Conditional execution based on previous results
- Schema validation and assertions

Example sequence file:
```json
{
  "name": "User Registration Flow",
  "description": "Test the complete user registration and login flow",
  "steps": [
    {
      "name": "Register User",
      "request": {
        "method": "POST",
        "url": "https://api.example.com/users",
        "headers": {
          "Content-Type": ["application/json"]
        },
        "body": "{\"email\":\"test@example.com\",\"password\":\"password123\"}"
      },
      "variables": [
        {
          "name": "userId",
          "source": "body",
          "path": "id"
        }
      ],
      "expectedStatus": 201
    },
    {
      "name": "Login User",
      "request": {
        "method": "POST",
        "url": "https://api.example.com/login",
        "headers": {
          "Content-Type": ["application/json"]
        },
        "body": "{\"email\":\"test@example.com\",\"password\":\"password123\"}"
      },
      "variables": [
        {
          "name": "token",
          "source": "body",
          "path": "token"
        }
      ],
      "expectedStatus": 200
    },
    {
      "name": "Get User Profile",
      "request": {
        "method": "GET",
        "url": "https://api.example.com/users/${userId}",
        "headers": {
          "Authorization": ["Bearer ${token}"]
        }
      },
      "expectedStatus": 200,
      "schemaValidate": true,
      "assertions": [
        {
          "type": "equals",
          "source": "body",
          "path": "email",
          "value": "test@example.com"
        }
      ]
    }
  ]
}
```

### Variable Extraction

Extract variables from HTTP responses for use in subsequent requests:

- JSON path extraction (`$.user.id`, `user.profile.name`)
- Header extraction
- Regular expression extraction
- Automatic variable substitution in URLs, headers, and bodies

### Test Assertions

Define assertions on HTTP responses:

```json
"assertions": [
  {
    "type": "equals",
    "source": "body",
    "path": "status",
    "value": "active"
  },
  {
    "type": "contains",
    "source": "body",
    "path": "items",
    "value": "product1"
  },
  {
    "type": "matches",
    "source": "header",
    "path": "Content-Type",
    "value": "application/json.*"
  }
]
```

Supported assertion types:
- equals - Exact value matching
- contains - String contains
- matches - Regular expression matching
- exists - Value exists
- notExists - Value doesn't exist
- in - Value is in a list
- lt/lessthan - Numeric less than
- gt/greaterthan - Numeric greater than
- null - Value is null

### Continuous Testing in Watch Mode

Run tests continuously as files change:

```bash
swagger-to-http test --watch http-requests/*.http
```

Options:
- --watch-interval - Milliseconds between file checks
- --watch-paths - Specific paths to watch for changes

## Examples

We provide various examples to help you get started:

- [Basic Usage](docs/examples/basic-usage/README.md)
- [Authentication](docs/examples/auth/)
- [Complex Parameters](docs/examples/parameters/)
- [Snapshot Testing](docs/examples/snapshot-testing/)
- [Schema Validation](docs/examples/schema-validation/)
- [Test Sequences](docs/examples/test-sequences/)

## Project Status

This project is in active development. The following features are implemented or planned:


- [x] Basic Swagger/OpenAPI parsing
- [x] HTTP file generation
- [x] CLI interface
- [x] HTTP request execution
- [x] Response snapshot comparison
- [x] Schema validation
- [x] Test sequences and variable extraction
- [ ] Git hooks integration

## Documentation

Comprehensive documentation is available in the [docs](docs/) directory:

- [Installation Guide](docs/installation.md)
- [Usage Guide](docs/usage.md)
- [Configuration Guide](docs/configuration.md)
- [HTTP File Format](docs/http-file-format.md)
- [HTTP Executor](docs/http-executor.md)
- [Snapshot Testing](docs/snapshot-testing.md)
- [Advanced Testing Features](docs/advanced-testing.md)
- [Examples](docs/examples/)
- [API Reference](docs/api-reference.md)
- [Contributing Guide](docs/contributing.md)
- [FAQ](docs/faq.md)

## Contributing

Contributions are welcome! Please read our [Contributing Guide](docs/contributing.md) for details on how to submit pull requests, the development process, and coding standards.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
