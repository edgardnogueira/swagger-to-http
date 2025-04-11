# Frequently Asked Questions (FAQ)

This document answers common questions about using `swagger-to-http`.

## General Questions

### What is swagger-to-http?

`swagger-to-http` is a tool that converts Swagger/OpenAPI documentation into organized HTTP request files. It also provides snapshot testing capabilities to validate API responses against saved snapshots.

### What versions of Swagger/OpenAPI are supported?

The tool supports both Swagger/OpenAPI 2.0 and OpenAPI 3.0 specifications.

### Is this tool compatible with my IDE?

The generated `.http` files are compatible with:
- JetBrains IDEs (IntelliJ IDEA, WebStorm, PhpStorm, etc.)
- Visual Studio Code with the REST Client extension
- Other text editors that support HTTP request files

### What operating systems are supported?

`swagger-to-http` is written in Go and supports:
- Linux
- macOS
- Windows

## Installation Questions

### How do I install swagger-to-http?

You can install using Go:
```bash
go install github.com/edgardnogueira/swagger-to-http@latest
```

Or download a pre-built binary from the [Releases page](https://github.com/edgardnogueira/swagger-to-http/releases).

See the [Installation Guide](installation.md) for more details.

### Why can't I run the tool after installation?

Ensure the installation directory is in your system's PATH. For Go installations, this would typically be `$GOPATH/bin`.

### How do I update to the latest version?

If you installed via Go:
```bash
go install github.com/edgardnogueira/swagger-to-http@latest
```

Or download the latest release from the [Releases page](https://github.com/edgardnogueira/swagger-to-http/releases).

## Usage Questions

### How do I generate HTTP files?

```bash
swagger-to-http generate -f swagger.json -o http-requests
```

Or from a URL:
```bash
swagger-to-http generate -u https://petstore.swagger.io/v2/swagger.json
```

### How are the HTTP files organized?

Files are organized by tags defined in the Swagger document. Each tag typically represents a resource or controller in your API. Within each file, requests are grouped by endpoints.

### Can I customize the base URL for requests?

Yes, you can override the base URL defined in the Swagger document:
```bash
swagger-to-http generate -f swagger.json -b https://api.example.com/v1
```

### How do I include authentication headers?

Use the `--auth` flag with an optional token:
```bash
swagger-to-http generate -f swagger.json --auth --auth-token "Bearer your_token"
```

## Snapshot Testing Questions

### What is snapshot testing?

Snapshot testing involves executing HTTP requests, capturing the responses, and comparing them against previously saved "snapshots" to detect changes in API behavior.

### How do I create initial snapshots?

```bash
swagger-to-http snapshot update "api/*.http"
```

### How do I run snapshot tests?

```bash
swagger-to-http snapshot test "api/*.http"
```

### What if my API responses change intentionally?

You can update snapshots to match new responses:
```bash
swagger-to-http snapshot update "api/*.http"
```

Or update only failing snapshots:
```bash
swagger-to-http snapshot test --update failed "api/*.http"
```

### Some headers change between requests. How do I handle this?

You can ignore specific headers during comparison:
```bash
swagger-to-http snapshot test --ignore-headers "Date,Set-Cookie,X-Request-ID" "api/*.http"
```

### Where are snapshots stored?

By default, snapshots are stored in the `.snapshots` directory. You can specify a custom location:
```bash
swagger-to-http snapshot test --snapshot-dir "my-snapshots" "api/*.http"
```

### How do I clean up unused snapshots?

```bash
swagger-to-http snapshot cleanup
```

Or automatically during testing:
```bash
swagger-to-http snapshot test --cleanup "api/*.http"
```

## Configuration Questions

### How do I create a configuration file?

Create a `swagger-to-http.yaml` file in your project directory with your configuration options. See the [Configuration Guide](configuration.md) for details.

### Can I use environment variables for configuration?

Yes, you can use environment variables with the `STH_` prefix:
```bash
export STH_OUTPUT_DIRECTORY=api-requests
export STH_AUTH_TOKEN="Bearer mytoken123"
swagger-to-http generate -f swagger.json
```

### What's the order of precedence for configuration?

1. Command-line flags (highest precedence)
2. Environment variables
3. Configuration file
4. Default values

## Advanced Questions

### Can I use variables in HTTP requests?

Yes, you can use variables in the format `{{variableName}}`. These can be:
- Defined in environment files
- Provided at runtime
- Extracted from previous responses for sequential tests

### How do I handle different content types in snapshots?

The tool includes content-type aware formatters for:
- JSON (with normalization)
- XML
- HTML
- Plain text
- Binary data

These formatters automatically handle different content types during comparison.

### Can I integrate with CI/CD pipelines?

Yes, you can include snapshot testing in your CI/CD pipeline:
```bash
swagger-to-http snapshot test --fail-on-missing "api/*.http"
```

This will fail if any tests fail or if snapshots are missing.

### Is there a way to automate snapshot updates with Git hooks?

Integration with Git hooks is planned for a future release. This will allow automatic HTTP file generation when Swagger files change.

## Troubleshooting

### The tool reports "file not found" for my Swagger file

Ensure the path to your Swagger file is correct and accessible:
```bash
# Absolute path
swagger-to-http generate -f /path/to/swagger.json

# Relative to current directory
swagger-to-http generate -f ./swagger.json
```

### Snapshot tests are failing for non-deterministic responses

If your API returns different data each time (like timestamps or random IDs), you can:
1. Ignore specific headers: `--ignore-headers "Date,X-Request-ID"`
2. Use the JSON formatter which can handle structural differences

### How do I debug snapshot differences?

When tests fail, the tool shows a diff of what changed. For more details, you can:
1. Update snapshots: `swagger-to-http snapshot update "file.http"`
2. Compare the updated snapshot with the previous version in your version control system

### The tool hangs when running tests

This could indicate issues connecting to your API. Check that:
1. The API is accessible
2. Network connections are allowed
3. Authentication details are correct
4. Timeouts are appropriately configured

## Getting Help

### Where can I report bugs?

Please open an issue on the [GitHub repository](https://github.com/edgardnogueira/swagger-to-http/issues).

### How can I request a feature?

Feature requests can be submitted as issues on the [GitHub repository](https://github.com/edgardnogueira/swagger-to-http/issues).

### I want to contribute. How do I get started?

Great! Check out the [Contributing Guide](contributing.md) for instructions on setting up the development environment and our contribution workflow.
