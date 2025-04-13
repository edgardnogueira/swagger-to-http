#!/bin/bash

# Install Git hooks for swagger-to-http

set -e

ECHO_PREFIX="\033[1;36m[swagger-to-http hooks]\033[0m"

echo -e "${ECHO_PREFIX} Installing Git hooks..."

# Create .git/hooks directory if it doesn't exist
mkdir -p .git/hooks

# Install pre-commit hook
cp hooks/pre-commit.sh .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
echo -e "${ECHO_PREFIX} Installed pre-commit hook"

# Install post-merge hook
cp hooks/post-merge.sh .git/hooks/post-merge
chmod +x .git/hooks/post-merge
echo -e "${ECHO_PREFIX} Installed post-merge hook"

# Add husky configuration if it exists in the project
if command -v npx &> /dev/null && [ -f package.json ]; then
    echo -e "${ECHO_PREFIX} Detected Node.js project, configuring Husky..."
    
    # Install husky if not already installed
    if ! grep -q "\"husky\"" package.json; then
        npm install --save-dev husky
        npx husky install
    fi
    
    # Add husky hooks
    npx husky add .husky/pre-commit "hooks/pre-commit.sh"
    npx husky add .husky/post-merge "hooks/post-merge.sh"
    
    echo -e "${ECHO_PREFIX} Husky configuration complete"
fi

# Create config directory for hook settings
mkdir -p .swagger-to-http

# Create default configuration if it doesn't exist
if [ ! -f .swagger-to-http/hooks.config ]; then
    cat > .swagger-to-http/hooks.config << EOF
# swagger-to-http Git hooks configuration

# Set to false to disable hooks temporarily
HOOKS_ENABLED=true

# Swagger/OpenAPI file patterns (space separated)
SWAGGER_FILE_PATTERNS="**/swagger.json **/swagger.yaml **/openapi.json **/openapi.yaml"

# Output directory for HTTP files
HTTP_OUTPUT_DIR="http"

# Whether to validate Swagger/OpenAPI files before generating HTTP files
VALIDATE_SWAGGER=true

# Whether to regenerate all HTTP files on changes or only affected ones
SELECTIVE_UPDATES=true
EOF
    echo -e "${ECHO_PREFIX} Created default configuration in .swagger-to-http/hooks.config"
fi

echo -e "${ECHO_PREFIX} Git hooks installation complete!"
echo -e "${ECHO_PREFIX} You can customize hook behavior by editing .swagger-to-http/hooks.config"
