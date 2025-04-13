#!/bin/bash

# Post-merge hook for swagger-to-http
# Updates HTTP files when Swagger/OpenAPI files have changed after a merge/pull

set -e

ECHO_PREFIX="\033[1;36m[swagger-to-http]\033[0m"
ERROR_PREFIX="\033[1;31m[swagger-to-http]\033[0m"
WARNING_PREFIX="\033[1;33m[swagger-to-http]\033[0m"

# Check if the environment variable to disable hooks is set
if [ "${SWAGGER_TO_HTTP_DISABLE_HOOKS}" = "true" ]; then
    echo -e "${WARNING_PREFIX} Git hooks are disabled via SWAGGER_TO_HTTP_DISABLE_HOOKS"
    exit 0
fi

# Load configuration
CONFIG_FILE=".swagger-to-http/hooks.config"
if [ -f "${CONFIG_FILE}" ]; then
    source "${CONFIG_FILE}"
else
    echo -e "${WARNING_PREFIX} Configuration file not found at ${CONFIG_FILE}"
    echo -e "${WARNING_PREFIX} Using default settings"
    HOOKS_ENABLED=true
    SWAGGER_FILE_PATTERNS="**/swagger.json **/swagger.yaml **/openapi.json **/openapi.yaml"
    HTTP_OUTPUT_DIR="http"
    VALIDATE_SWAGGER=true
    SELECTIVE_UPDATES=true
fi

# Check if hooks are enabled in config
if [ "${HOOKS_ENABLED}" != "true" ]; then
    echo -e "${WARNING_PREFIX} Git hooks are disabled in ${CONFIG_FILE}"
    exit 0
fi

# Get the swagger-to-http binary path
SWAGGER_TO_HTTP_BIN="swagger-to-http"
if ! command -v "${SWAGGER_TO_HTTP_BIN}" &> /dev/null; then
    # Try to find it in the project
    if [ -f "./bin/swagger-to-http" ]; then
        SWAGGER_TO_HTTP_BIN="./bin/swagger-to-http"
    elif [ -f "./swagger-to-http" ]; then
        SWAGGER_TO_HTTP_BIN="./swagger-to-http"
    else
        echo -e "${ERROR_PREFIX} swagger-to-http binary not found"
        echo -e "${ERROR_PREFIX} Please make sure it's installed and in your PATH"
        exit 1
    fi
fi

# Check if this was a pull/merge that updated files
get_updated_swagger_files() {
    git diff-tree -r --name-only --no-commit-id ORIG_HEAD HEAD | grep -E "(swagger|openapi)\\.(json|yaml|yml)$"
}

UPDATED_SWAGGER_FILES=$(get_updated_swagger_files)

# If no Swagger files were updated, exit early
if [ -z "${UPDATED_SWAGGER_FILES}" ]; then
    echo -e "${ECHO_PREFIX} No Swagger/OpenAPI files changed, skipping update"
    exit 0
fi

echo -e "${ECHO_PREFIX} Swagger/OpenAPI files changed in pull/merge:"
echo "${UPDATED_SWAGGER_FILES}"

# Generate HTTP files for the updated Swagger files
echo -e "${ECHO_PREFIX} Updating HTTP files..."

for FILE in ${UPDATED_SWAGGER_FILES}; do
    echo -e "${ECHO_PREFIX} Generating HTTP files for ${FILE}..."
    
    # Check if file exists (it might have been deleted)
    if [ ! -f "${FILE}" ]; then
        echo -e "${WARNING_PREFIX} File ${FILE} no longer exists, skipping"
        continue
    fi
    
    # Determine output directory
    OUTPUT_DIR="${HTTP_OUTPUT_DIR}/$(dirname "${FILE}")"
    mkdir -p "${OUTPUT_DIR}"
    
    # Generate HTTP files
    if [ "${SELECTIVE_UPDATES}" = "true" ]; then
        # Use selective update mode
        "${SWAGGER_TO_HTTP_BIN}" generate "${FILE}" --output "${OUTPUT_DIR}" --selective
    else
        # Regenerate all files
        "${SWAGGER_TO_HTTP_BIN}" generate "${FILE}" --output "${OUTPUT_DIR}"
    fi
    
    echo -e "${ECHO_PREFIX} HTTP files for ${FILE} updated"
done

# Notify the user that they may need to commit the changes
echo -e "${ECHO_PREFIX} HTTP files have been updated based on pulled changes"
echo -e "${ECHO_PREFIX} Please review the changes and commit them if needed"

exit 0
