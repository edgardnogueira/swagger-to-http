# Examples

This directory contains practical examples for using the `swagger-to-http` tool in various scenarios.

## Example Categories

- [Basic Usage](#basic-usage)
- [Authentication](#authentication)
- [Complex Parameters](#complex-parameters)
- [File Uploads](#file-uploads)
- [Snapshot Testing](#snapshot-testing)

## Basic Usage

### Example 1: Generating HTTP Files from Swagger

See [basic-usage/README.md](basic-usage/README.md) for details on:
- Converting a simple Swagger/OpenAPI document to HTTP files
- Understanding the generated file structure
- Running basic requests

### Example 2: Using Environment Variables

See [environment-variables/README.md](environment-variables/README.md) for:
- Using variables in HTTP requests
- Defining environment files
- Switching between environments (dev, test, prod)

## Authentication

### Example 3: Basic Authentication

See [auth/basic-auth.md](auth/basic-auth.md) for:
- Adding Basic Authentication to requests
- Handling authentication headers

### Example 4: OAuth/Bearer Token

See [auth/oauth.md](auth/oauth.md) for:
- Setting up Bearer token authentication
- Refreshing tokens
- Using token variables

### Example 5: API Keys

See [auth/api-keys.md](auth/api-keys.md) for:
- Using API keys in headers or query parameters
- Keeping API keys secure

## Complex Parameters

### Example 6: Path Parameters

See [parameters/path-params.md](parameters/path-params.md) for:
- Using path parameters in requests
- Handling required vs. optional parameters

### Example 7: Query Parameters

See [parameters/query-params.md](parameters/query-params.md) for:
- Adding query parameters to requests
- Handling arrays and complex objects in query strings

### Example 8: Request Bodies

See [parameters/request-bodies.md](parameters/request-bodies.md) for:
- Structuring request bodies
- Working with different content types (JSON, XML, form data)
- Handling nested objects and arrays

## File Uploads

### Example 9: Single File Upload

See [uploads/single-file.md](uploads/single-file.md) for:
- Uploading single files via multipart/form-data
- Setting file metadata

### Example 10: Multiple File Upload

See [uploads/multiple-files.md](uploads/multiple-files.md) for:
- Uploading multiple files in one request
- Handling file arrays

## Snapshot Testing

### Example 11: Basic Snapshot Testing

See [snapshot-testing/basic.md](snapshot-testing/basic.md) for:
- Creating initial snapshots
- Running snapshot tests
- Understanding test results

### Example 12: Selective Updates

See [snapshot-testing/selective-updates.md](snapshot-testing/selective-updates.md) for:
- Updating only failed snapshots
- Handling missing snapshots
- Using different update modes

### Example 13: CI/CD Integration

See [snapshot-testing/ci-cd.md](snapshot-testing/ci-cd.md) for:
- Integrating snapshot tests in CI/CD pipelines
- Handling failures properly
- Generating reports
