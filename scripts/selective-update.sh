#!/bin/sh
#
# Selective update script
# This script selectively updates HTTP files based on changes in Swagger files

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Check for arguments
if [ "$#" -lt 1 ]; then
  echo "${YELLOW}Usage:${NC} $0 <swagger-file1> [<swagger-file2> ...]"
  echo "       $0 --all                      # Process all swagger files"
  echo "       $0 --modified                 # Process modified swagger files only"
  exit 1
fi

# Function to process a single swagger file
process_swagger_file() {
  SWAGGER_FILE="$1"
  
  # Check if file exists
  if [ ! -f "$SWAGGER_FILE" ]; then
    echo "${RED}Error: File not found: $SWAGGER_FILE${NC}"
    return 1
  fi
  
  # Check if file is a swagger file
  if ! grep -q "swagger\|openapi" "$SWAGGER_FILE"; then
    echo "${YELLOW}Warning: File does not appear to be a Swagger/OpenAPI file: $SWAGGER_FILE${NC}"
    echo "${YELLOW}Processing anyway...${NC}"
  fi
  
  # Default output directory is http-requests in the same directory as the Swagger file
  OUTPUT_DIR="$(dirname "$SWAGGER_FILE")/http-requests"
  
  # Check for custom output directory in config
  if [ -f ".swagger-to-http.yml" ]; then
    # Read output directory from config, skipping for simplicity in this example
    :
  fi
  
  # Check if swagger-to-http is installed
  if ! command -v swagger-to-http >/dev/null 2>&1; then
    # Check if it exists in the local bin directory
    if [ -x "./bin/swagger-to-http" ]; then
      SWAGGER_TO_HTTP="./bin/swagger-to-http"
    else
      echo "${RED}Error: swagger-to-http command not found${NC}"
      echo "${YELLOW}Run 'make build' to build the tool locally${NC}"
      return 1
    fi
  else
    SWAGGER_TO_HTTP="swagger-to-http"
  fi
  
  echo "${YELLOW}Processing $SWAGGER_FILE...${NC}"
  echo "${YELLOW}Generating HTTP files to $OUTPUT_DIR...${NC}"
  
  # Create output directory if it doesn't exist
  mkdir -p "$OUTPUT_DIR"
  
  # Generate HTTP files
  if ! $SWAGGER_TO_HTTP generate -f "$SWAGGER_FILE" -o "$OUTPUT_DIR"; then
    echo "${RED}Error: Failed to generate HTTP files from $SWAGGER_FILE${NC}"
    return 1
  fi
  
  echo "${GREEN}Successfully generated HTTP files for $SWAGGER_FILE${NC}"
  return 0
}

# Process based on arguments
case "$1" in
  --all)
    echo "${YELLOW}Processing all Swagger files...${NC}"
    # Find all swagger files in the repository
    SWAGGER_FILES=$(find . -type f \( -name "*.yaml" -o -name "*.yml" -o -name "*.json" \) | grep -i swagger)
    
    if [ -z "$SWAGGER_FILES" ]; then
      echo "${YELLOW}No Swagger files found in the repository${NC}"
      exit 0
    fi
    
    SUCCESS=0
    FAILURE=0
    
    for FILE in $SWAGGER_FILES; do
      if process_swagger_file "$FILE"; then
        SUCCESS=$((SUCCESS + 1))
      else
        FAILURE=$((FAILURE + 1))
      fi
    done
    
    echo "${GREEN}Successfully processed $SUCCESS Swagger files${NC}"
    if [ "$FAILURE" -gt 0 ]; then
      echo "${RED}Failed to process $FAILURE Swagger files${NC}"
      exit 1
    fi
    ;;
    
  --modified)
    echo "${YELLOW}Processing modified Swagger files...${NC}"
    # Get modified swagger files
    SWAGGER_FILES=$(git diff --name-only --diff-filter=ACMR | grep -E '\.(yaml|yml|json)$' | grep -i swagger)
    
    if [ -z "$SWAGGER_FILES" ]; then
      echo "${YELLOW}No modified Swagger files found${NC}"
      exit 0
    fi
    
    SUCCESS=0
    FAILURE=0
    
    for FILE in $SWAGGER_FILES; do
      if process_swagger_file "$FILE"; then
        SUCCESS=$((SUCCESS + 1))
      else
        FAILURE=$((FAILURE + 1))
      fi
    done
    
    echo "${GREEN}Successfully processed $SUCCESS modified Swagger files${NC}"
    if [ "$FAILURE" -gt 0 ]; then
      echo "${RED}Failed to process $FAILURE Swagger files${NC}"
      exit 1
    fi
    ;;
    
  *)
    # Process specific files
    SUCCESS=0
    FAILURE=0
    
    for FILE in "$@"; do
      if process_swagger_file "$FILE"; then
        SUCCESS=$((SUCCESS + 1))
      else
        FAILURE=$((FAILURE + 1))
      fi
    done
    
    echo "${GREEN}Successfully processed $SUCCESS Swagger files${NC}"
    if [ "$FAILURE" -gt 0 ]; then
      echo "${RED}Failed to process $FAILURE Swagger files${NC}"
      exit 1
    fi
    ;;
esac

exit 0
