# Advanced Testing Features

This document covers the advanced testing features available in swagger-to-http, including:

- Schema validation against OpenAPI specifications
- Sequential tests with dependencies
- Variable extraction for use in subsequent tests
- Test assertions with customizable validation rules
- Continuous testing in watch mode

## Schema Validation

Schema validation allows you to validate HTTP responses against OpenAPI/Swagger schema definitions.

### Usage

```bash
swagger-to-http test validate --swagger-file swagger.json http-requests/*.http
```

### Options

| Flag | Description |
|------|-------------|
| `--swagger-file` | Path to the Swagger/OpenAPI file (required) |
| `--ignore-props` | Comma-separated properties to ignore in validation |
| `--ignore-add-props` | Ignore additional properties not defined in schema |
| `--ignore-formats` | Ignore format validation (e.g., date, email) |
| `--ignore-patterns` | Ignore pattern validation |
| `--req-props-only` | Validate only required properties |
| `--ignore-nullable` | Ignore nullable field validation |

### Example

```bash
# Validate all HTTP responses with more lenient validation
swagger-to-http test validate --swagger-file swagger.json --ignore-formats --ignore-patterns http-requests/*.http
```

## Test Sequences

Test sequences allow you to run tests in a specific order with dependencies between them, enabling you to test multi-step workflows.

### Usage

```bash
swagger-to-http test sequence sequence-tests/*.json
```

### Options

| Flag | Description |
|------|-------------|
| `--variables-path` | Path to load/save variables |
| `--save-vars` | Save extracted variables to file |
| `--var-format` | Variable format (default: ${varname}) |
| `--fail-fast` | Stop sequence on first failure |
| `--validate-schema` | Validate responses against schema |
| `--swagger-file` | Path to Swagger/OpenAPI file (when using schema validation) |

### Test Sequence File Format

Test sequence files are JSON files that define:
- A series of HTTP requests in a specific order
- Variable extractions to pass data between requests
- Wait times between requests
- Conditional execution based on previous results
- Schema validation and assertions

Example:

```json
{
  "name": "User Registration Flow",
  "description": "Test the complete user registration and login flow",
  "tags": ["user", "auth"],
  "variables": {
    "baseUrl": "https://api.example.com"
  },
  "steps": [
    {
      "name": "Register User",
      "description": "Create a new user account",
      "request": {
        "method": "POST",
        "url": "${baseUrl}/users",
        "headers": {
          "Content-Type": ["application/json"]
        },
        "body": "{\"email\":\"test@example.com\",\"password\":\"password123\"}"
      },
      "variables": [
        {
          "name": "userId",
          "source": "body",
          "path": "id",
          "required": true
        }
      ],
      "expectedStatus": 201,
      "assertions": [
        {
          "type": "exists",
          "source": "body",
          "path": "id"
        }
      ]
    },
    {
      "name": "Login User",
      "waitBefore": "1s",
      "request": {
        "method": "POST",
        "url": "${baseUrl}/login",
        "headers": {
          "Content-Type": ["application/json"]
        },
        "body": "{\"email\":\"test@example.com\",\"password\":\"password123\"}"
      },
      "variables": [
        {
          "name": "token",
          "source": "body",
          "path": "token",
          "required": true
        }
      ],
      "expectedStatus": 200,
      "stopOnFail": true
    },
    {
      "name": "Get User Profile",
      "request": {
        "method": "GET",
        "url": "${baseUrl}/users/${userId}",
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

### Step Options

Each step in a test sequence supports the following options:

| Option | Description |
|--------|-------------|
| `name` | Name of the test step (required) |
| `description` | Description of the test step |
| `request` | HTTP request object (required) |
| `expectedStatus` | Expected HTTP status code |
| `variables` | Array of variable extraction definitions |
| `waitBefore` | Duration to wait before executing the step |
| `waitAfter` | Duration to wait after executing the step |
| `skip` | Boolean flag to skip this step |
| `skipCondition` | Condition for skipping this step |
| `stopOnFail` | Boolean flag to stop sequence on failure |
| `schemaValidate` | Boolean flag to validate response against schema |
| `assertions` | Array of test assertions |

## Variable Extraction

Variable extraction allows you to extract values from responses and use them in subsequent requests.

### Variable Extraction Definition

```json
{
  "name": "token",         // Name of the variable (required)
  "source": "body",        // Source: body, header, status (required)
  "path": "data.token",    // JSON path for body source, header name for header source
  "regexp": "Bearer (.*)", // Regular expression with capture group (optional)
  "default": "",           // Default value if extraction fails
  "required": true         // Whether extraction is required
}
```

### Sources

- `body`: Extract from response body using JSON path
- `header`: Extract from response header
- `status`: Extract the HTTP status code

### Examples

Extract a token from JSON response:
```json
{
  "name": "token",
  "source": "body",
  "path": "data.auth.token"
}
```

Extract a value with regex:
```json
{
  "name": "csrfToken",
  "source": "body",
  "regexp": "name=\"csrf\" value=\"([^\"]+)\""
}
```

Extract a header:
```json
{
  "name": "requestId",
  "source": "header",
  "path": "X-Request-ID"
}
```

## Test Assertions

Assertions allow you to validate specific aspects of HTTP responses.

### Assertion Definition

```json
{
  "type": "equals",       // Assertion type (required)
  "source": "body",       // Source: body, header, status (required)
  "path": "user.active",  // Path within source (for body and header)
  "value": "true",        // Value to check against
  "values": ["a", "b"],   // Array of values (for 'in' assertion)
  "not": false,           // Invert the assertion
  "ignoreCase": true      // Case-insensitive comparison
}
```

### Assertion Types

| Type | Description |
|------|-------------|
| `equals` | Value exactly matches expected value |
| `contains` | Value contains expected substring |
| `matches` | Value matches regular expression |
| `exists` | Value exists (is not null or undefined) |
| `notExists` | Value does not exist |
| `in` | Value is in the provided list of values |
| `lessthan` / `lt` | Numeric value is less than expected |
| `greaterthan` / `gt` | Numeric value is greater than expected |
| `null` | Value is null |

### Examples

Check if status is "active":
```json
{
  "type": "equals",
  "source": "body",
  "path": "status",
  "value": "active"
}
```

Check if array contains an item:
```json
{
  "type": "contains",
  "source": "body",
  "path": "permissions",
  "value": "admin"
}
```

Check if header matches a pattern:
```json
{
  "type": "matches",
  "source": "header",
  "path": "Content-Type",
  "value": "application/json.*",
  "ignoreCase": true
}
```

Check if value is in a list:
```json
{
  "type": "in",
  "source": "body",
  "path": "status",
  "values": ["active", "pending", "approved"]
}
```

## Continuous Testing in Watch Mode

Watch mode allows you to run tests continuously as files change.

### Usage

```bash
swagger-to-http test --watch http-requests/*.http
```

### Options

| Flag | Description |
|------|-------------|
| `--watch` | Enable watch mode |
| `--watch-interval` | Milliseconds between file checks (default: 1000) |
| `--watch-paths` | Specific paths to watch for changes |

### Example

```bash
# Run tests continuously with a 2-second interval
swagger-to-http test --watch --watch-interval 2000 http-requests/*.http
```

## Extending Your Tests

The advanced testing features can be combined to create sophisticated test scenarios:

1. Create a test sequence that performs a complete workflow
2. Use variable extraction to pass data between steps
3. Add assertions to validate specific response aspects
4. Enable schema validation to ensure responses conform to your API schema
5. Run tests in watch mode during development to get immediate feedback

This enables true end-to-end testing of your API with capabilities like:

- Testing authentication flows
- Testing dependent operations (create → read → update → delete)
- Testing complex business processes
- Validating response schemas against your API definition
- Verifying specific response properties

## Integration with CI/CD

These advanced testing features work well in continuous integration environments:

### Example GitHub Actions Workflow

```yaml
name: API Tests

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.18
    
    - name: Install swagger-to-http
      run: go install github.com/edgardnogueira/swagger-to-http@latest
    
    - name: Run API Schema Validation
      run: swagger-to-http test validate --swagger-file api/swagger.json tests/*.http
      
    - name: Run API Sequence Tests
      run: swagger-to-http test sequence tests/sequences/*.json
```

## Best Practices

When using the advanced testing features, consider the following best practices:

1. **Organize Test Sequences**: Group related steps into sequence files by feature or workflow
2. **Use Clear Naming**: Give sequences and steps descriptive names
3. **Add Timeouts**: Set appropriate timeouts for HTTP requests
4. **Test Progressive Complexity**: Start with basic tests before building complex sequences
5. **Validate Essential Data**: Focus assertions on critical data points
6. **Manage Variables**: Save variables between test runs for stateful testing
7. **Schema Validation**: Use schema validation for structure, assertions for specific values
8. **Handle Failures**: Use `stopOnFail` for critical steps, but allow non-critical steps to continue

## Troubleshooting

### Common Issues and Solutions

**Issue**: Variable extraction fails from response body
- **Solution**: Check content type, verify JSON path, try using regex extraction

**Issue**: Schema validation fails
- **Solution**: Use `--ignore-formats` or `--ignore-add-props` for more lenient validation

**Issue**: Sequence tests hang or timeout
- **Solution**: Check server availability, increase timeout values

**Issue**: Watch mode doesn't detect file changes
- **Solution**: Increase `--watch-interval`, check file permissions

**Issue**: Tests work locally but fail in CI
- **Solution**: Check environment variables, ensure correct schema path

## Upgrading from Basic Testing

If you've been using the basic snapshot testing features, you can gradually upgrade to the advanced features:

1. Start with schema validation for existing tests
2. Convert critical test workflows to sequence tests
3. Add variable extraction for data that flows between requests
4. Add assertions for detailed validation
5. Finally, implement watch mode for development efficiency
