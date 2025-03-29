#!/bin/bash

# Get the absolute path to the project root directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Change to scripts directory to ensure proper operation
cd "$SCRIPT_DIR"

# Main deployment script that takes environment as a parameter

# Check if environment parameter is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 [dev|prod]"
    echo "Example: $0 dev"
    exit 1
fi

ENV=$1
PORT=${2:-8080}
REGION=${3:-us-east1}

# Source the base script
source "$SCRIPT_DIR/deploy-base.sh"

# Call the function with the specified environment
deploy_to_environment "$ENV" "$PORT" "$REGION"
