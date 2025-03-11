#!/bin/bash

# Exit on error
set -e

# Load environment variables from .env.terraform file
if [ -f "$(dirname "$0")/../.env.terraform" ]; then
  export $(cat "$(dirname "$0")/../.env.terraform" | grep -v '^#' | xargs)
else
  echo "Error: .env.terraform file not found. Please copy .env.terraform.example to .env.terraform and fill in the values."
  exit 1
fi

# Check if required environment variables are set
required_vars=(
  "GCP_PROJECT_ID"
  "GCP_REGION"
  "GCP_SQL_INSTANCE"
  "CLOUD_RUN_SERVICE_NAME"
)

for var in "${required_vars[@]}"; do
  if [ -z "${!var}" ]; then
    echo "Error: $var is not set in .env.terraform file"
    exit 1
  fi
done

# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Check if required tools are installed
if ! command_exists terraform; then
  echo "Error: terraform is not installed"
  exit 1
fi

if ! command_exists gcloud; then
  echo "Error: gcloud is not installed"
  exit 1
fi

# Check if user is authenticated with gcloud
if ! gcloud auth print-identity-token >/dev/null 2>&1; then
  echo "Error: Not authenticated with gcloud. Please run 'gcloud auth login'"
  exit 1
fi

# Set the working directory to the terraform directory
cd "$(dirname "$0")/../terraform" || exit 1

# Function to delete preview databases
delete_preview_databases() {
  echo "Deleting preview databases..."
  databases=$(gcloud sql databases list \
    --instance="$GCP_SQL_INSTANCE" \
    --filter="name ~ ^${CLOUD_RUN_SERVICE_NAME}_[0-9]+$" \
    --format="value(name)")

  for db in $databases; do
    echo "Deleting database: $db"
    gcloud sql databases delete "$db" \
      --instance="$GCP_SQL_INSTANCE" \
      --quiet || true
  done
}

# Function to delete preview Cloud Run revisions
delete_preview_revisions() {
  echo "Deleting preview Cloud Run revisions..."
  gcloud run services update-traffic "$CLOUD_RUN_SERVICE_NAME" \
    --region="$GCP_REGION" \
    --remove-tags="pr-.*" \
    --quiet || true
}

# Function to delete preview container images
delete_preview_images() {
  echo "Deleting preview container images..."
  images=$(gcloud container images list-tags "$GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/$CLOUD_RUN_SERVICE_NAME/$CLOUD_RUN_SERVICE_NAME" \
    --filter="tags ~ ^[0-9]+-.*$" \
    --format="value(digest)")

  for digest in $images; do
    echo "Deleting image digest: $digest"
    gcloud container images delete "$GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/$CLOUD_RUN_SERVICE_NAME/$CLOUD_RUN_SERVICE_NAME@$digest" \
      --quiet || true
  done
}

# Ask for confirmation
echo "WARNING: This will destroy all resources managed by Terraform and clean up preview environments."
echo "This action cannot be undone."
read -p "Are you sure you want to continue? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
  echo "Operation cancelled."
  exit 1
fi

# Clean up preview environments
echo "Cleaning up preview environments..."
delete_preview_databases
delete_preview_revisions
delete_preview_images

# Destroy Terraform-managed resources
echo "Destroying Terraform-managed resources..."
terraform destroy -auto-approve

echo "Reset completed successfully!"
echo "All resources have been destroyed." 