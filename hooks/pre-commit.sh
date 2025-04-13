#!/bin/bash

# Pre-commit hook for swagger-to-http
# Validates Swagger/OpenAPI files and updates HTTP files if needed

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

# Find staged Swagger files
get_staged_swagger_files() {
    git diff --cached --name-only --diff-filter=ACMR | grep -E "(swagger|openapi)\\.(json|yaml|yml)$"
}

STAGED_SWAGGER_FILES=$(get_staged_swagger_files)

# If no Swagger files are staged, exit early
if [ -z "${STAGED_SWAGGER_FILES}" ]; then
    echo -e "${ECHO_PREFIX} No Swagger/OpenAPI files staged, skipping validation"
    exit 0
fi

echo -e "${ECHO_PREFIX} Found staged Swagger/OpenAPI files:"
echo "${STAGED_SWAGGER_FILES}"

# Validate Swagger files if enabled
if [ "${VALIDATE_SWAGGER}" = "true" ]; then
    echo -e "${ECHO_PREFIX} Validating Swagger/OpenAPI files..."
    
    # Validate each file
    EXIT_CODE=0
    for FILE in ${STAGED_SWAGGER_FILES}; do
        echo -e "${ECHO_PREFIX} Validating ${FILE}..."
        
        # Check if file exists (it might have been deleted)
        if [ ! -f "${FILE}" ]; then
            echo -e "${WARNING_PREFIX} File ${FILE} no longer exists, skipping"
            continue
        fi
        
        # Use swagger-to-http to validate the file
        if ! "${SWAGGER_TO_HTTP_BIN}" validate "${FILE}" 2>/dev/null; then
            echo -e "${ERROR_PREFIX} ${FILE} is not a valid Swagger/OpenAPI file"
            EXIT_CODE=1
        else
            echo -e "${ECHO_PREFIX} ${FILE} is valid"
        fi
    done
    
    # If any file is invalid, abort the commit
    if [ ${EXIT_CODE} -ne 0 ]; then
        echo -e "${ERROR_PREFIX} Aborting commit due to invalid Swagger/OpenAPI files"
        exit ${EXIT_CODE}
    fi
fi

# Generate HTTP files
echo -e "${ECHO_PREFIX} Generating HTTP files..."

# Generate HTTP files for each Swagger file
for FILE in ${STAGED_SWAGGER_FILES}; do
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
    
    # Stage the generated files
    git add "${OUTPUT_DIR}"/*.http
    
    echo -e "${ECHO_PREFIX} HTTP files for ${FILE} generated and staged"
done

echo -e "${ECHO_PREFIX} Pre-commit hook completed successfully"
exit 0
