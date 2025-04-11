# Snapshot Testing

This guide explains the snapshot testing functionality in `swagger-to-http` and how to use it effectively for API testing and validation.

## What is Snapshot Testing?

Snapshot testing is a technique for validating that the responses from your API don't change unexpectedly. It works as follows:

1. You execute an HTTP request defined in an `.http` file
2. The actual response is captured
3. This response is compared against a previously saved "snapshot"
4. If they match, the test passes; if they differ, the test fails

This approach helps you detect unintended changes in your API behavior, such as:
- Status code changes
- Response structure modifications
- Data type alterations
- Header changes
- Performance degradation

## Snapshot Test Workflow

The typical workflow for snapshot testing involves:

1. **Creation**: Generate initial snapshots for your API endpoints
2. **Verification**: Run tests to compare current responses with saved snapshots
3. **Updates**: Update snapshots when intentional API changes occur

## Getting Started

### Creating Initial Snapshots

To create initial snapshots for your HTTP files:

```bash
swagger-to-http snapshot update "api/*.http"
```

This command:
1. Executes all HTTP requests in the specified files
2. Saves the responses as snapshots in the `.snapshots` directory (or custom directory)
3. Creates a structured snapshot for each request

### Running Snapshot Tests

To test if your API responses match the saved snapshots:

```bash
swagger-to-http snapshot test "api/*.http"
```

This command:
1. Executes all HTTP requests in the specified files
2. Compares the responses with saved snapshots
3. Reports any differences or failures

### Viewing Snapshot Files

Snapshot files are stored in the `.snapshots` directory by default. Each snapshot is a JSON file with:

- Request metadata (method, path)
- Response status code
- Response headers
- Response body (formatted based on content type)

Example snapshot file:

```json
{
  "metadata": {
    "requestPath": "/api/users/1",
    "requestMethod": "GET",
    "contentType": "application/json",
    "statusCode": 200,
    "headers": {
      "Content-Type": ["application/json"],
      "Cache-Control": ["no-cache"]
    },
    "createdAt": "2025-04-10T15:30:00Z"
  },
  "content": "{\"id\":1,\"name\":\"John Doe\",\"email\":\"john@example.com\"}"
}
```

## Advanced Usage

### Update Modes

Different update modes are available:

- `none`: Never update snapshots (default)
- `all`: Always update all snapshots
- `failed`: Update only snapshots that failed comparison
- `missing`: Create snapshots only if they don't exist

For example, to update only snapshots that fail:

```bash
swagger-to-http snapshot test --update failed "api/*.http"
```

### Ignoring Headers

Some headers often change between requests (like timestamps or request IDs). You can ignore specific headers during comparison:

```bash
swagger-to-http snapshot test --ignore-headers "Date,Set-Cookie,X-Request-ID" "api/*.http"
```

You can also configure this in your configuration file:

```yaml
snapshots:
  ignore_headers:
    - Date
    - Set-Cookie
    - X-Request-ID
```

### Snapshot Directory

By default, snapshots are stored in the `.snapshots` directory. You can specify a custom directory:

```bash
swagger-to-http snapshot test --snapshot-dir "my-snapshots" "api/*.http"
```

### Failing on Missing Snapshots

In CI/CD environments, you may want to fail if snapshots are missing:

```bash
swagger-to-http snapshot test --fail-on-missing "api/*.http"
```

### Cleanup Unused Snapshots

To remove snapshots that are no longer associated with any HTTP requests:

```bash
swagger-to-http snapshot cleanup
```

Or automatically during testing:

```bash
swagger-to-http snapshot test --cleanup "api/*.http"
```

## Content Type Formatters

`swagger-to-http` includes content-type aware formatters for comparing different types of responses:

### JSON Formatter

- Normalizes JSON for consistent comparison
- Detects structural differences (missing/extra fields)
- Identifies type mismatches
- Shows differences in values
- Handles nested structures

### XML Formatter

- Compares XML structure
- Normalizes whitespace
- Detects element and attribute changes

### HTML Formatter

- Basic HTML structure comparison
- Normalizes formatting

### Plain Text Formatter

- Line-by-line text comparison
- Normalizes line endings

### Binary Formatter

- Detects binary data changes
- Shows hex representation for small files
- Reports size differences for large files

## Snapshot Comparison Process

When comparing snapshots, the tool follows this process:

1. **Status Comparison**: Checks if status codes match
2. **Header Comparison**: Compares headers (respecting ignored headers)
3. **Content Type Detection**: Determines the appropriate formatter based on content type
4. **Body Comparison**: Performs content-specific comparison
5. **Diff Generation**: Creates human-readable diffs for any differences

## Interpreting Test Results

Test results include:

- **Overall Status**: Pass/fail status for each request
- **Statistics**: Summary of tests run, passed, failed, updated, and created
- **Diff Details**: For failed tests, detailed information about what differed
- **Duration**: Time taken for each test and overall suite

Example output:

```
Running snapshot tests for api/users.http

Request 1: GET /api/users
  ✓ Snapshot matched

Request 2: GET /api/users/1
  ✓ Snapshot matched

Request 3: POST /api/users
  ✗ Snapshot comparison failed
    Status code: expected 201, got 400
    Body content differs (expected 27 bytes, got 42 bytes)
    Diff preview:
      {"id":123,"name":"New User"}
      ...
      {"error":"Invalid input","details":"Email is required"}

========================================
Snapshot Test Summary
========================================
Total tests:    3
Passed:         2
Failed:         1
Created:        0
Updated:        0
Duration:       0.35 seconds
========================================
```

## Best Practices

### Version Control

- Commit snapshot files to version control
- Review snapshot changes during code reviews
- Document intentional API changes that cause snapshot updates

### Test Organization

- Group related requests in the same HTTP file
- Use descriptive file names and comments
- Organize directories to match your API structure

### Maintenance

- Regularly clean up unused snapshots
- Update snapshots when API changes are intentional
- Consider using update modes wisely (e.g., `--update failed` to update only failing tests)

### CI/CD Integration

- Include snapshot testing in your CI/CD pipeline
- Use `--fail-on-missing` to ensure all endpoints have snapshots
- Generate test reports for easier debugging

## Troubleshooting

### Missing Snapshots

If tests fail due to missing snapshots:

```bash
# Create initial snapshots
swagger-to-http snapshot update "api/*.http"

# Or run tests with auto-creation
swagger-to-http snapshot test --update missing "api/*.http"
```

### Authentication Issues

If requests fail due to authentication:

```bash
# Include authentication in requests
swagger-to-http snapshot test "api/*.http" --auth --auth-token "your_token"

# Or configure in yaml file:
# snapshots:
#   auth_header: Authorization
#   auth_token: Bearer your_token
```

### Inconsistent Comparisons

If you're seeing inconsistent comparison results:

```bash
# Ignore variable headers
swagger-to-http snapshot test --ignore-headers "Date,Set-Cookie,X-Request-ID" "api/*.http"

# Use content-specific comparison
# For JSON, normalize with jq if needed
```

## Next Steps

- Review [example snapshots](examples/snapshots/) to understand the format
- Learn about [configuration options](configuration.md) for snapshot testing
- Explore [integrating with Git hooks](git-hooks.md) for automatic snapshots
