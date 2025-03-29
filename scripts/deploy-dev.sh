#!/bin/bash

# Get script directory path
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source the base script
source "$SCRIPT_DIR/deploy-base.sh"

# Call the function with development environment
deploy_to_environment "dev" "8080" "us-east1"
