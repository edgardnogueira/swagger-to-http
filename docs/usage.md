# Usage Guide

This guide explains how to use the `swagger-to-http` tool to convert Swagger/OpenAPI documents into HTTP files and perform snapshot testing.

## Table of Contents

- [Quick Start](#quick-start)
- [Command Line Interface](#command-line-interface)
- [Generate Command](#generate-command)
- [Snapshot Commands](#snapshot-commands)
- [Common Workflows](#common-workflows)

## Quick Start

Here are some common usage patterns to get you started:

### Generate HTTP Files from a Local Swagger File

```bash
swagger-to-http generate -f swagger.json -o http-requests
```

### Generate HTTP Files from a Remote Swagger URL

```bash
swagger-to-http generate -u https://petstore.swagger.io/v2/swagger.json -o http-requests
```

### Run Snapshot Tests

```bash
swagger-to-http snapshot test http-requests/*.http
```

### Update Snapshots

```bash
swagger-to-http snapshot update http-requests/*.http
```

## Command Line Interface

The `swagger-to-http` tool provides a command-line interface with several commands:

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

## Generate Command

The `generate` command converts Swagger/OpenAPI documents to HTTP files:

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

### Examples

#### Generate with Custom Base URL

```bash
swagger-to-http generate -f swagger.json -b https://api.example.com/v1
```

#### Generate with Authentication Header

```bash
swagger-to-http generate -f swagger.json --auth --auth-token "Bearer YOUR_TOKEN"
```

#### Generate from URL without Indentation

```bash
swagger-to-http generate -u https://petstore.swagger.io/v2/swagger.json -i=false
```

## Snapshot Commands

The `snapshot` command provides various subcommands for snapshot testing:

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

### Snapshot Test Command

```
Usage:
  swagger-to-http snapshot test [file-pattern]

Flags:
  --update string         Update mode: none, all, failed, missing (default "none")
  --ignore-headers string Comma-separated headers to ignore in comparison (default "Date,Set-Cookie")
  --snapshot-dir string   Directory for snapshot storage (default ".snapshots")
  --fail-on-missing       Fail when snapshot is missing
  --cleanup               Remove unused snapshots after testing
  -h, --help              help for test
```

### Snapshot Update Command

```
Usage:
  swagger-to-http snapshot update [file-pattern]

Flags:
  --snapshot-dir string   Directory for snapshot storage (default ".snapshots") 
  -h, --help              help for update
```

### Examples

#### Run Tests with Auto-Update for Failed Tests

```bash
swagger-to-http snapshot test --update failed "api/*.http"
```

#### Ignore Specific Headers

```bash
swagger-to-http snapshot test --ignore-headers "Date,Set-Cookie,X-Request-ID" "api/*.http"
```

#### Update All Snapshots

```bash
swagger-to-http snapshot update "api/*.http"
```

#### List and Cleanup Snapshots

```bash
# List all snapshots
swagger-to-http snapshot list

# Cleanup unused snapshots
swagger-to-http snapshot cleanup
```

## Common Workflows

### API Development Workflow

1. Start with your Swagger/OpenAPI document:
   ```bash
   swagger-to-http generate -f api-spec.yaml -o http-requests
   ```

2. Review the generated HTTP files to ensure they match expectations.

3. Create initial snapshots:
   ```bash
   swagger-to-http snapshot update http-requests/*.http
   ```

4. As you develop and make API changes, run tests to detect changes:
   ```bash
   swagger-to-http snapshot test http-requests/*.http
   ```

5. Update snapshots when API changes are intentional:
   ```bash
   swagger-to-http snapshot update http-requests/*.http
   ```

### CI/CD Integration

In your CI/CD pipeline:

1. Generate HTTP files:
   ```bash
   swagger-to-http generate -f api-spec.yaml -o http-requests
   ```

2. Run tests with strict mode:
   ```bash
   swagger-to-http snapshot test --fail-on-missing http-requests/*.http
   ```

3. (Optional) Clean up unused snapshots:
   ```bash
   swagger-to-http snapshot cleanup
   ```

## Next Steps

- Learn about [Configuration Options](configuration.md)
- Explore the [HTTP File Format](http-file-format.md)
- Read more about [Snapshot Testing](snapshot-testing.md)
