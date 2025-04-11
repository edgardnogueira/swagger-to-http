# Basic Usage Example

This example demonstrates the basic usage of `swagger-to-http` to convert a Swagger/OpenAPI specification to HTTP files and run simple requests.

## Example Swagger Document

Let's start with a simple Swagger document (`petstore.json`):

```json
{
  "swagger": "2.0",
  "info": {
    "title": "Petstore API",
    "description": "A sample API that uses a petstore as an example to demonstrate features",
    "version": "1.0.0"
  },
  "host": "petstore.example.com",
  "basePath": "/api/v1",
  "schemes": ["https"],
  "consumes": ["application/json"],
  "produces": ["application/json"],
  "paths": {
    "/pets": {
      "get": {
        "summary": "List all pets",
        "operationId": "listPets",
        "tags": ["pets"],
        "parameters": [
          {
            "name": "limit",
            "in": "query",
            "description": "How many items to return at one time (max 100)",
            "required": false,
            "type": "integer",
            "format": "int32"
          }
        ],
        "responses": {
          "200": {
            "description": "A paged array of pets",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/Pet"
              }
            }
          }
        }
      },
      "post": {
        "summary": "Create a pet",
        "operationId": "createPets",
        "tags": ["pets"],
        "parameters": [
          {
            "name": "pet",
            "in": "body",
            "description": "Pet to add to the store",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Pet"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Created"
          }
        }
      }
    },
    "/pets/{petId}": {
      "get": {
        "summary": "Info for a specific pet",
        "operationId": "showPetById",
        "tags": ["pets"],
        "parameters": [
          {
            "name": "petId",
            "in": "path",
            "required": true,
            "description": "The id of the pet to retrieve",
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "Expected response to a valid request",
            "schema": {
              "$ref": "#/definitions/Pet"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Pet": {
      "type": "object",
      "required": ["id", "name"],
      "properties": {
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "name": {
          "type": "string"
        },
        "tag": {
          "type": "string"
        }
      }
    }
  }
}
```

## Step 1: Generate HTTP Files

```bash
swagger-to-http generate -f petstore.json -o http-requests
```

This will create the following structure:

```
http-requests/
└── pets/
    └── pets.http
```

The `pets.http` file will contain:

```http
### List all pets
GET https://petstore.example.com/api/v1/pets
Accept: application/json

###

### Create a pet
POST https://petstore.example.com/api/v1/pets
Content-Type: application/json
Accept: application/json

{
  "id": 0,
  "name": "string",
  "tag": "string"
}

###

### Info for a specific pet
GET https://petstore.example.com/api/v1/pets/{{petId}}
Accept: application/json

###
```

## Step 2: Modify the Generated Requests

You might want to customize the generated requests. For example, to add a specific petId and create a real pet:

```http
### List all pets with limit
GET https://petstore.example.com/api/v1/pets?limit=10
Accept: application/json

###

### Create a pet
POST https://petstore.example.com/api/v1/pets
Content-Type: application/json
Accept: application/json

{
  "id": 1,
  "name": "Fluffy",
  "tag": "cat"
}

###

### Get a specific pet
GET https://petstore.example.com/api/v1/pets/1
Accept: application/json

###
```

## Step 3: Use a Different Base URL

If you want to use a different server than what's specified in the Swagger document:

```bash
swagger-to-http generate -f petstore.json -o http-requests -b https://dev-api.example.com
```

This will override the base URL in the generated files:

```http
### List all pets
GET https://dev-api.example.com/api/v1/pets
Accept: application/json

###
```

## Step 4: Run Snapshot Tests

Once you've generated the HTTP files, you can create snapshots:

```bash
swagger-to-http snapshot update http-requests/pets/pets.http
```

This will execute the requests and save the responses as snapshots.

To test if your API returns the same responses:

```bash
swagger-to-http snapshot test http-requests/pets/pets.http
```

## Complete Example Workflow

```bash
# Generate HTTP files from Swagger
swagger-to-http generate -f petstore.json -o http-requests

# Create initial snapshots
swagger-to-http snapshot update http-requests/pets/pets.http

# Make changes to your API

# Test to see if responses have changed
swagger-to-http snapshot test http-requests/pets/pets.http

# If intended changes, update snapshots
swagger-to-http snapshot update http-requests/pets/pets.http
```

## Inspecting the Generated Files

The generated HTTP files follow a standardized format:

1. Request line with method and URL
2. Headers (including Accept, Content-Type, etc.)
3. Empty line
4. Body (for POST, PUT, PATCH)
5. Separator line (###)

This format is compatible with tools like:
- JetBrains IDEs (IntelliJ, WebStorm, etc.)
- VS Code with the REST Client extension

## Next Steps

- Learn about [authentication](../auth/basic-auth.md)
- Try [complex parameters](../parameters/request-bodies.md)
- Explore [snapshot testing](../snapshot-testing/basic.md)
