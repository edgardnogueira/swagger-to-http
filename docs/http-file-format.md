# HTTP File Format

The `.http` file format is used to define HTTP requests that can be executed and tested. This document explains the structure and syntax of these files as generated and understood by `swagger-to-http`.

## Overview

HTTP files are plain text files with the `.http` extension that contain one or more HTTP requests. Each request consists of:

- Request line (HTTP method and URL)
- Headers
- Empty line
- Body (optional)
- Separator line (three or more `#` characters)

## Basic Structure

Here's a simple example of an HTTP file with a single request:

```http
GET https://api.example.com/users
Accept: application/json
Authorization: Bearer token123

###
```

And an example with a request body:

```http
POST https://api.example.com/users
Content-Type: application/json
Accept: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "role": "admin"
}

###
```

## Multiple Requests

Multiple requests can be placed in a single file, separated by the `###` separator:

```http
GET https://api.example.com/users
Accept: application/json

###

POST https://api.example.com/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}

###

DELETE https://api.example.com/users/123
Accept: application/json

###
```

## Variables

HTTP files support variable substitution with the following syntax:

```http
GET https://api.example.com/users/{{userId}}
Authorization: Bearer {{authToken}}

###
```

Variables can be:
- Defined in environment files
- Provided at runtime
- Extracted from previous responses for sequential tests

## Comments

Comments start with `//` or `#` and can be placed anywhere in the file:

```http
// This is a comment
GET https://api.example.com/users

# Another comment
Accept: application/json

###
```

You can also use multi-line comments to provide descriptions:

```http
// Get a list of users
// This will return all users with pagination
GET https://api.example.com/users
Accept: application/json

###
```

## Request Parts

### Method and URL

The first line of each request specifies the HTTP method and URL:

```http
GET https://api.example.com/path?param=value
```

Supported HTTP methods include:
- GET
- POST
- PUT
- DELETE
- PATCH
- HEAD
- OPTIONS

### Headers

Headers follow the request line, with one header per line:

```http
Content-Type: application/json
Accept: application/json
Authorization: Bearer token123
X-API-Key: your-api-key
```

### Request Body

The request body comes after an empty line following the headers:

```http
POST https://api.example.com/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}

###
```

Different content types are supported:

#### JSON

```http
POST https://api.example.com/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "preferences": {
    "theme": "dark",
    "notifications": true
  },
  "tags": ["admin", "user"]
}

###
```

#### Form Data

```http
POST https://api.example.com/submit-form
Content-Type: application/x-www-form-urlencoded

name=John+Doe&email=john%40example.com&subscribe=true

###
```

#### Multipart Forms

```http
POST https://api.example.com/upload
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW

------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="text"

This is some text data
------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="file"; filename="example.txt"
Content-Type: text/plain

This is the content of the file
------WebKitFormBoundary7MA4YWxkTrZu0gW--

###
```

#### XML

```http
POST https://api.example.com/users
Content-Type: application/xml

<?xml version="1.0" encoding="UTF-8"?>
<user>
  <name>John Doe</name>
  <email>john@example.com</email>
  <roles>
    <role>admin</role>
    <role>user</role>
  </roles>
</user>

###
```

## Organization

`swagger-to-http` organizes HTTP files based on the Swagger/OpenAPI document structure:

1. Files are grouped by tags defined in the Swagger document
2. Each tag typically corresponds to a resource or controller
3. File names are derived from the tag names
4. Requests within files are grouped by endpoints

Example directory structure:

```
http-requests/
├── pets/
│   └── pets.http       # All pet-related endpoints
├── users/
│   └── users.http      # All user-related endpoints
└── default/
    └── default.http    # Endpoints without specific tags
```

## Compatibility

The HTTP file format is compatible with:

- JetBrains IDEs (IntelliJ IDEA, WebStorm, etc.)
- VS Code with the REST Client extension
- Other tools that support similar formats

## Next Steps

- Learn how to [run snapshot tests](snapshot-testing.md) using these HTTP files
- See [examples](examples/README.md) of HTTP files for different scenarios
