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
- Execute HTTP requests and compare responses with snapshots (coming soon)
- Integration with Git hooks for automatic updates (coming soon)

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

## Quick Start

Generate HTTP files from a Swagger file:

```bash
swagger-to-http generate -f swagger.json -o http-requests
```

Or from a URL:

```bash
swagger-to-http generate -u https://petstore.swagger.io/v2/swagger.json -o http-requests
```

## Usage

```
Usage:
  swagger-to-http [command]

Available Commands:
  generate    Generate HTTP files from a Swagger/OpenAPI document
  help        Help about any command
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
  directory: snapshots
  update_on_difference: false
```

## Project Status

This project is currently in active development. The following features are planned:

- [x] Basic Swagger/OpenAPI parsing
- [x] HTTP file generation
- [x] CLI interface
- [ ] HTTP request execution
- [ ] Response snapshot comparison
- [ ] Git hooks integration
- [ ] Schema validation

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
