#!/bin/bash

# Source common functions and configurations
source "$(dirname "$0")/common.sh"

# Check if DEPLOY_ENV is set
if [ -z "${DEPLOY_ENV:-}" ]; then
    echo "Error: DEPLOY_ENV environment variable is not set"
    echo "Expected values: staging, prod"
    exit 1
fi

# Validate DEPLOY_ENV value
if [[ "$DEPLOY_ENV" != "staging" && "$DEPLOY_ENV" != "prod" ]]; then
    echo "Error: Invalid DEPLOY_ENV value: $DEPLOY_ENV"
    echo "Expected values: staging, prod"
    exit 1
fi

log "Starting deployment to $DEPLOY_ENV environment..."

# Update and execute migration job
log "Updating and executing migration job..."
gcloud run jobs update "$SERVICE_NAME-$DEPLOY_ENV-migrate" \
    --image "gcr.io/$PROJECT_ID/$SERVICE_NAME-migrate:$COMMIT_SHA" \
    --region "$REGION" \
    --execute-now \
    --wait

# Deploy to Cloud Run
log "Deploying to Cloud Run..."
gcloud run deploy "$SERVICE_NAME-$DEPLOY_ENV" \
    --image "gcr.io/$PROJECT_ID/$SERVICE_NAME:$COMMIT_SHA" \
    --region "$REGION" \
    --platform managed

log "Deployment to $DEPLOY_ENV completed successfully!"