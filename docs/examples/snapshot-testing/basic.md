# Basic Snapshot Testing

This example demonstrates how to use the snapshot testing feature of `swagger-to-http` to validate API responses.

## What is Snapshot Testing?

Snapshot testing is a technique where you:
1. Execute HTTP requests
2. Save the responses as "snapshots"
3. Later compare new responses against these snapshots to detect changes

This is particularly useful for:
- Regression testing
- API versioning
- Ensuring backward compatibility
- Detecting unintended changes in your API

## Example Setup

For this example, we'll use the following HTTP file (`users.http`):

```http
### Get all users
GET https://api.example.com/users
Accept: application/json

###

### Get user by ID
GET https://api.example.com/users/1
Accept: application/json

###

### Create a new user
POST https://api.example.com/users
Content-Type: application/json
Accept: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "role": "user"
}

###
```

## Step 1: Create Initial Snapshots

First, let's create the initial snapshots:

```bash
swagger-to-http snapshot update users.http
```

This command:
1. Executes all requests in the file
2. Saves the responses as snapshots in the `.snapshots` directory
3. Creates a structured snapshot file for each request

The console output will look something like:

```
Updating snapshots for users.http

Request 1: GET /users
  ⟳ Snapshot created

Request 2: GET /users/1
  ⟳ Snapshot created

Request 3: POST /users
  ⟳ Snapshot created

========================================
Snapshot Update Summary
========================================
Total requests: 3
Created:        3
Updated:        0
Duration:       0.85 seconds
========================================
```

## Step 2: Examine the Snapshots

Let's look at the generated snapshot files in the `.snapshots` directory:

```
.snapshots/
└── users_get.snap.json
└── users_1_get.snap.json
└── users_post.snap.json
```

A snapshot file (`users_get.snap.json`) might look like:

```json
{
  "metadata": {
    "requestPath": "/users",
    "requestMethod": "GET",
    "contentType": "application/json",
    "statusCode": 200,
    "headers": {
      "Content-Type": ["application/json"],
      "Cache-Control": ["no-cache"]
    },
    "createdAt": "2025-04-11T10:15:23Z"
  },
  "content": "[{\"id\":1,\"name\":\"John Doe\",\"email\":\"john@example.com\"},{\"id\":2,\"name\":\"Jane Smith\",\"email\":\"jane@example.com\"}]"
}
```

## Step 3: Run Snapshot Tests

Now that we have our baseline snapshots, we can run tests to compare current responses:

```bash
swagger-to-http snapshot test users.http
```

If nothing has changed in the API, you'll see:

```
Running snapshot tests for users.http

Request 1: GET /users
  ✓ Snapshot matched

Request 2: GET /users/1
  ✓ Snapshot matched

Request 3: POST /users
  ✓ Snapshot matched

========================================
Snapshot Test Summary
========================================
Total tests:    3
Passed:         3
Failed:         0
Created:        0
Updated:        0
Duration:       0.78 seconds
========================================
```

## Step 4: Handling API Changes

Let's say the API now returns a different response for the user with ID 1. When you run the tests:

```bash
swagger-to-http snapshot test users.http
```

You'll see a failure:

```
Running snapshot tests for users.http

Request 1: GET /users
  ✓ Snapshot matched

Request 2: GET /users/1
  ✗ Snapshot comparison failed
    Body content differs (expected 58 bytes, got 65 bytes)
    Diff preview:
      {"id":1,"name":"John Doe","email":"john@example.com"}
      {"id":1,"name":"John Doe Updated","email":"john@example.com"}

Request 3: POST /users
  ✓ Snapshot matched

========================================
Snapshot Test Summary
========================================
Total tests:    3
Passed:         2
Failed:         1
Created:        0
Updated:        0
Duration:       0.80 seconds
========================================
1 of 3 tests failed
```

## Step 5: Updating Failed Snapshots

If the change is expected, you can update only the failed snapshots:

```bash
swagger-to-http snapshot test --update failed users.http
```

Or update all snapshots:

```bash
swagger-to-http snapshot update users.http
```

## Step 6: Ignoring Certain Headers

Some headers change between requests (like timestamps or request IDs). You can ignore them:

```bash
swagger-to-http snapshot test --ignore-headers "Date,X-Request-ID" users.http
```

## Working with JSON Responses

The snapshot system automatically handles JSON normalization, so you don't need to worry about whitespace or property order differences.

Original JSON:
```json
{
  "name": "John",
  "id": 1
}
```

Equivalent for comparison:
```json
{"id":1,"name":"John"}
```

## Detecting Structural Changes

The snapshot system can detect structural changes in JSON:

```
Diff preview:
  Missing fields:
    - user.address
  Extra fields:
    + user.phone
  Different types:
    user.age: expected "number", got "string"
  Different values:
    user.name: expected "John", got "Johnny"
```

## Practical Use Cases

### CI/CD Pipeline

```bash
# In your CI/CD pipeline
swagger-to-http snapshot test --fail-on-missing api/*.http
```

### Development Workflow

```bash
# After making API changes
swagger-to-http snapshot test api/*.http

# If changes are intentional
swagger-to-http snapshot update api/*.http
```

### API Versioning

```bash
# Test against v1 snapshots
swagger-to-http snapshot test --snapshot-dir .snapshots-v1 api/*.http

# Test against v2 snapshots
swagger-to-http snapshot test --snapshot-dir .snapshots-v2 api/*.http
```

## Working with Different Content Types

The snapshot system handles different content types appropriately:

### JSON

```http
GET https://api.example.com/users/1
Accept: application/json
```

The JSON formatter normalizes and provides structural comparison.

### XML

```http
GET https://api.example.com/users/1
Accept: application/xml
```

The XML formatter handles XML structure comparison.

### HTML

```http
GET https://api.example.com/home
Accept: text/html
```

The HTML formatter provides basic HTML comparison.

### Plain Text

```http
GET https://api.example.com/status
Accept: text/plain
```

The text formatter does line-by-line comparison.

### Binary Data

```http
GET https://api.example.com/image.png
```

The binary formatter reports size differences and hex representation for small files.

## Advanced Tips

### Use Comments to Document Test Cases

```http
// Test case: Get user with valid ID
// Expect: Status 200 with user details
GET https://api.example.com/users/1
Accept: application/json

###

// Test case: Get user with invalid ID
// Expect: Status 404 with error message
GET https://api.example.com/users/999
Accept: application/json

###
```

### Group Related Tests

Organize your HTTP files by resource or functionality:

- `users.http` - All user-related endpoints
- `auth.http` - Authentication endpoints
- `products.http` - Product management endpoints

### Use Cleanup After Testing

To remove unused snapshots after testing:

```bash
swagger-to-http snapshot test --cleanup api/*.http
```

Or separately:

```bash
swagger-to-http snapshot cleanup
```

## Next Steps

- Try [selective updates](selective-updates.md) for more control over snapshots
- Learn about [CI/CD integration](ci-cd.md) for automated testing
- Explore the [advanced snapshot comparison options](../advanced/custom-formatters.md)
