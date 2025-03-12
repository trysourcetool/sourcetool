#!/bin/bash

# Source common functions and configurations
source "$(dirname "$0")/common.sh"

# Check required environment variables
check_required_vars VITE_API_BASE_URL VITE_DOMAIN

log "Starting build process..."

# Build and push migration image
log "Building migration image..."
docker build -t "gcr.io/$PROJECT_ID/$SERVICE_NAME-migrate:$COMMIT_SHA" \
    -f docker/cloud/Dockerfile.migrate .

log "Pushing migration image..."
docker push "gcr.io/$PROJECT_ID/$SERVICE_NAME-migrate:$COMMIT_SHA"

# Build and push main application image
log "Building main application image..."
docker build -t "gcr.io/$PROJECT_ID/$SERVICE_NAME:$COMMIT_SHA" \
    --build-arg "VITE_API_BASE_URL=$VITE_API_BASE_URL" \
    --build-arg "VITE_DOMAIN=$VITE_DOMAIN" .

log "Pushing main application image..."
docker push "gcr.io/$PROJECT_ID/$SERVICE_NAME:$COMMIT_SHA"

log "Build process completed successfully!" 