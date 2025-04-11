#!/bin/sh
#
# Install Git hooks script
# This script installs the Git hooks into the .git/hooks directory of the repository

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Check if we're in a git repository
if [ ! -d ".git" ]; then
  echo "${RED}Error: This script must be run from the root of a Git repository${NC}"
  exit 1
fi

# Create hooks directory if it doesn't exist
mkdir -p .git/hooks

# Get the scripts directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
HOOKS_DIR="$SCRIPT_DIR/hooks"

# Check if hooks directory exists
if [ ! -d "$HOOKS_DIR" ]; then
  echo "${RED}Error: Hooks directory not found at $HOOKS_DIR${NC}"
  exit 1
fi

# Install each hook
echo "${YELLOW}Installing Git hooks...${NC}"

# List of hooks to install
HOOKS="pre-commit post-checkout post-merge"

for HOOK in $HOOKS; do
  HOOK_PATH="$HOOKS_DIR/$HOOK"
  TARGET_PATH=".git/hooks/$HOOK"
  
  # Check if hook script exists
  if [ ! -f "$HOOK_PATH" ]; then
    echo "${YELLOW}Warning: $HOOK script not found, skipping${NC}"
    continue
  fi
  
  # Copy hook script to .git/hooks directory
  cp "$HOOK_PATH" "$TARGET_PATH"
  
  # Make hook executable
  chmod +x "$TARGET_PATH"
  
  echo "${GREEN}Installed $HOOK hook${NC}"
done

echo "${GREEN}Git hooks installation complete${NC}"
echo "${YELLOW}You can disable hooks for a specific command with:${NC}"
echo "  SKIP_GIT_HOOKS=1 git commit"

exit 0
