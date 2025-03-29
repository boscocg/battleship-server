#!/bin/bash

# Get script directory path
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source the base script
source "$SCRIPT_DIR/deploy-base.sh"

# Call the function with production environment
deploy_to_environment "prod" "8080" "us-east1"
