# Git Hooks Example

This example demonstrates how to use the Git hooks integration with swagger-to-http to automatically generate and update HTTP files when Swagger/OpenAPI files change.

## Prerequisites

- Git repository initialized (`git init`)
- swagger-to-http installed

## Setup

1. Install the Git hooks:

```bash
# Using the CLI
swagger-to-http hooks install

# OR using the install script directly
chmod +x hooks/install.sh
./hooks/install.sh
```

2. Check that the hooks are installed properly:

```bash
swagger-to-http hooks status
```

## How It Works

### Pre-commit Hook

When you commit changes to Swagger/OpenAPI files, the pre-commit hook will:

1. Detect which Swagger/OpenAPI files are being committed
2. Validate those files for correctness
3. Generate or update corresponding HTTP files
4. Automatically stage the generated HTTP files

Example workflow:

```bash
# Edit a Swagger file
vim api/swagger.json

# Stage the changes
git add api/swagger.json

# Commit - this will trigger the pre-commit hook
git commit -m "Updated API endpoints"

# The HTTP files will be automatically generated and included in the commit
```

### Post-merge Hook

When you pull changes or merge branches, the post-merge hook will:

1. Detect which Swagger/OpenAPI files changed in the merge
2. Generate or update corresponding HTTP files
3. Notify you that these files need to be committed

Example workflow:

```bash
# Pull changes from remote
git pull origin main

# If Swagger files changed, HTTP files will be updated
# You'll need to commit these changes manually
git status
git add http/
git commit -m "Updated HTTP files after pull"
```

## Configuration

The Git hooks can be configured by editing the `.swagger-to-http/hooks.config` file:

```bash
# Enable or disable hooks
HOOKS_ENABLED=true

# Swagger/OpenAPI file patterns to watch
SWAGGER_FILE_PATTERNS="**/swagger.json **/swagger.yaml **/openapi.json **/openapi.yaml"

# Output directory for HTTP files
HTTP_OUTPUT_DIR="http"

# Whether to validate Swagger files
VALIDATE_SWAGGER=true

# Whether to use selective updates for better performance
SELECTIVE_UPDATES=true
```

## Temporarily Disabling Hooks

To temporarily disable the hooks for a specific command:

```bash
# For a single command
SWAGGER_TO_HTTP_DISABLE_HOOKS=true git commit -m "Skip hooks for this commit"

# For the current session
export SWAGGER_TO_HTTP_DISABLE_HOOKS=true
```

Or use the CLI:

```bash
# Disable hooks
swagger-to-http hooks disable

# Enable hooks again
swagger-to-http hooks enable
```

## Troubleshooting

### Hooks not running

1. Check if hooks are installed and enabled:

```bash
swagger-to-http hooks status
```

2. Ensure executable permissions:

```bash
chmod +x .git/hooks/pre-commit .git/hooks/post-merge
```

3. Check that the swagger-to-http binary is in your PATH:

```bash
which swagger-to-http
```

### Error messages

If you see error messages about Swagger validation, the hook is working correctly but has detected issues with your Swagger/OpenAPI files that need to be fixed before committing.
