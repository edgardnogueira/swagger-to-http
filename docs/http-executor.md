# HTTP Executor

The HTTP Executor is responsible for executing HTTP requests from `.http` files and returning the responses. It's a core component of the snapshot testing system, enabling the tool to make real HTTP requests and capture responses for snapshot testing.

## Features

- Execute individual HTTP requests from `.http` files
- Support for all common HTTP methods (GET, POST, PUT, DELETE, PATCH, etc.)
- Variable substitution in URLs, headers, and body content
- Environment variable loading
- Header management
- Proper content-type handling
- Timeout configuration
- Comprehensive error handling

## Usage

### Basic Usage

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

### With Variables

```go
// Create environment variables
env := map[string]string{
    "BASE_URL": "https://api.example.com",
    "API_KEY": "my-api-key",
}

// Create executor with environment
executor := http.NewExecutor(30*time.Second, env)

// Execute with request-specific variables
requestVars := map[string]string{
    "USER_ID": "123",
}

// Variables will be substituted in URL, headers, and body
response, err := executor.Execute(ctx, request, requestVars)
```

### Executing Files

You can execute all requests in an HTTP file:

```go
// Parse the HTTP file
parser := http.NewParser()
httpFile, err := parser.ParseFile("api/users.http")
if err != nil {
    // Handle error
}

// Execute all requests in the file
responses, err := executor.ExecuteFile(ctx, httpFile, nil)
if err != nil {
    // Handle error
}

// Process responses
for i, resp := range responses {
    fmt.Printf("Response %d: %d %s\n", i+1, resp.StatusCode, resp.Status)
}
```

## Variable Substitution

The HTTP executor supports variable substitution in the following formats:

- `{{VARIABLE_NAME}}` - Will be replaced with the value of VARIABLE_NAME from environment or request variables

Variables can be used in:

- URLs: `GET {{BASE_URL}}/api/users/{{USER_ID}}`
- Headers: `Authorization: Bearer {{TOKEN}}`
- Request bodies: `{"name": "{{NAME}}", "email": "{{EMAIL}}"}`

## Environment Variables

The HTTP executor can load environment variables from:

1. System environment variables with the `HTTP_` prefix (e.g., `HTTP_BASE_URL` becomes `BASE_URL`)
2. Variables passed directly to the executor
3. Request-specific variables passed to the `Execute` method

Priority order (if the same variable exists in multiple sources):
1. Request-specific variables (highest priority)
2. Variables passed to the executor
3. System environment variables (lowest priority)

## HTTP File Format

The HTTP executor works with `.http` files in the following format:

```http
# This is a comment
@name GetUser
@tag users

GET https://api.example.com/users/123
Accept: application/json
Authorization: Bearer {{TOKEN}}

###

# This is another request
@name CreateUser
@tag users

POST https://api.example.com/users
Content-Type: application/json
Accept: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}
```

## Implementation Details

The HTTP executor is implemented in the following files:

- `internal/infrastructure/http/executor.go` - Main implementation
- `internal/infrastructure/http/parser.go` - HTTP file parser
- `internal/cli/snapshot_command.go` - CLI integration

The executor follows clean architecture principles:

- Domain models are defined in `internal/domain/models`
- The executor interface is defined in `internal/application/interfaces.go`
- The implementation is in the infrastructure layer

## Integration with Snapshot Testing

The HTTP executor is integrated with the snapshot testing system to:

1. Parse `.http` files into request objects
2. Execute the requests to get responses
3. Save responses as snapshots or compare with existing snapshots

This workflow enables automated testing of APIs against known good responses.

## Error Handling

The HTTP executor provides detailed error messages for various failure scenarios:

- Network errors
- Timeout errors
- Invalid request format
- HTTP status errors
- Body parsing errors

All error messages include context to help diagnose the issue.

## Testing

The HTTP executor is thoroughly tested with unit tests covering:

- Request execution
- Variable substitution
- Error handling
- HTTP file parser
- Multi-request file execution

Tests can be found in:
- `internal/infrastructure/http/executor_test.go`
- `internal/infrastructure/http/parser_test.go`

## Future Improvements

Planned enhancements for the HTTP executor:

1. Session management
   - Cookie persistence between requests
   - OAuth flows

2. Advanced variable handling
   - Extract variables from responses
   - Use variables in subsequent requests
   - Support for complex variable expressions

3. Performance optimizations
   - Connection pooling
   - Parallel request execution
   - Caching responses

4. Improved error reporting
   - Detailed response analysis
   - Context-aware error messages
   - Suggestions for common issues

5. Extended protocol support
   - WebSockets
   - gRPC
   - GraphQL

## CLI Commands

The HTTP executor can be used via the following CLI commands:

```bash
# Run tests for all .http files
swagger-to-http snapshot test

# Run tests for specific file or pattern
swagger-to-http snapshot test api/users.http

# Update snapshots
swagger-to-http snapshot update api/products.http

# List existing snapshots
swagger-to-http snapshot list
```

See the [CLI Documentation](cli.md) for more details on available commands and options.
