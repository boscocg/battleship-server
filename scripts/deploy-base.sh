#!/bin/bash

# Get the absolute path to the project root directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Function to deploy to a specific environment
deploy_to_environment() {
    local ENV=$1
    local PORT=${2:-8080}
    local REGION=${3:-us-east1}

    echo "Deploying to environment: $ENV"
    
    # Check if we're running from the scripts directory
    if [[ $(basename "$(pwd)") != "scripts" ]]; then
        echo "Changing to scripts directory..."
        cd "$SCRIPT_DIR"
    fi

    # Create secure temp directory
    TEMP_DIR=$(mktemp -d)
    echo "Created temporary directory: $TEMP_DIR"

    # Copy key file to temp directory
    cp "$PROJECT_ROOT/user-key.json" "$TEMP_DIR/key.json"
    KEY_FILE="$TEMP_DIR/key.json"

    # Verify key file exists
    if [ ! -f "$KEY_FILE" ]; then
        echo "Error: Key file $KEY_FILE not found!"
        exit 1
    fi

    # Extract service account email from key file
    SERVICE_ACCOUNT=$(grep -o '"client_email": "[^"]*' "$KEY_FILE" | cut -d '"' -f4)
    if [ -z "$SERVICE_ACCOUNT" ]; then
        echo "Error: Could not extract client_email from key file!"
        exit 1
    fi
    echo "Service Account: $SERVICE_ACCOUNT"

    # Extract information from key file
    echo "Extracting project information from service account key..."
    PROJECT_ID=$(grep -o '"project_id": "[^"]*' "$KEY_FILE" | cut -d '"' -f4)
    if [ -z "$PROJECT_ID" ]; then
        echo "Error: Could not extract project_id from key file!"
        exit 1
    fi
    echo "Project ID: $PROJECT_ID"

    echo "Using cloudbuild.yaml to handle environment variables..."

    # Navigate to project root for build context
    cd "$PROJECT_ROOT"

    # Debug: List files in the build context
    echo "Files in build context:"
    ls -la

    # Authenticate with service account using the key in the temp directory
    echo "Authenticating with service account..."
    gcloud auth activate-service-account --key-file="$KEY_FILE"

    # Set project
    gcloud config set project "$PROJECT_ID"

    # Run Cloud Build
    echo "Starting build process..."
    gcloud builds submit \
        --project="${PROJECT_ID}" \
        --impersonate-service-account="${SERVICE_ACCOUNT}" \
        --substitutions=_ENV_FILE=.env,_ENV=${ENV},_PORT=${PORT} \
        --config=cloudbuild.yaml

    # Check build result
    if [ $? -eq 0 ]; then
        echo "Build completed successfully!"
        
        # Deploy to Cloud Run
        echo "Deploying to Cloud Run..."
        gcloud run deploy "battledak-server-${ENV}" \
            --image "gcr.io/${PROJECT_ID}/battledak-server-${ENV}" \
            --platform managed \
            --allow-unauthenticated \
            --region $REGION
        
        echo "Deployment completed!"
    else
        echo "Error during build process!"
        # Clean up before exiting
        rm -rf "$TEMP_DIR"
        exit 1
    fi

    # Clean up temporary directory
    echo "Cleaning up temporary files..."
    rm -rf "$TEMP_DIR"
    echo "Cleanup complete"
}
