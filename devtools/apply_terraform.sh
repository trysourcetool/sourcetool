#!/bin/bash

# Exit on error
set -e

# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Function to check and refresh Google Cloud authentication
check_gcloud_auth() {
  echo "Checking Google Cloud authentication status..."
  
  # Try to get access token
  if ! gcloud auth print-access-token >/dev/null 2>&1; then
    echo "Google Cloud authentication required. Please run the following commands:"
    echo "1. gcloud auth login"
    echo "2. gcloud auth application-default login"
    exit 1
  fi
  
  # Check application default credentials
  if [ ! -f "${HOME}/.config/gcloud/application_default_credentials.json" ]; then
    echo "Google Cloud application default credentials not found."
    echo "Please run: gcloud auth application-default login"
    exit 1
  fi
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

# Check Google Cloud authentication
check_gcloud_auth

# Check if terraform.tfvars exists
if [ ! -f "$(dirname "$0")/../terraform/terraform.tfvars" ]; then
  echo "Error: terraform.tfvars file not found in the terraform directory."
  exit 1
fi

# Set the working directory to the terraform directory first
cd "$(dirname "$0")/../terraform" || exit 1

# Get project_id from terraform.tfvars
project_id=$(grep 'project_id' "terraform.tfvars" | cut -d'=' -f2 | tr -d ' "')
if [ -z "$project_id" ]; then
  echo "Error: Could not find project_id in terraform.tfvars"
  exit 1
fi

# Check if user has sufficient permissions for the project
if ! gcloud projects describe "$project_id" >/dev/null 2>&1; then
  echo "Error: You don't have sufficient permissions for project $project_id"
  echo "Please make sure you are authenticated with the correct account and have the necessary permissions."
  exit 1
fi

# Enable required GCP APIs
echo "Enabling required GCP APIs..."
required_apis=(
  "compute.googleapis.com"          # Compute Engine API
  "sqladmin.googleapis.com"        # Cloud SQL Admin API
  "run.googleapis.com"             # Cloud Run API
  "vpcaccess.googleapis.com"       # Serverless VPC Access API
  "servicenetworking.googleapis.com" # Service Networking API
  "dns.googleapis.com"             # Cloud DNS API
  "secretmanager.googleapis.com"   # Secret Manager API
  "certificatemanager.googleapis.com" # Certificate Manager API
)

for api in "${required_apis[@]}"; do
  if ! gcloud services list --project "$project_id" --filter="config.name:$api" --format="get(config.name)" | grep -q "^$api"; then
    echo "Enabling $api..."
    gcloud services enable "$api" --project "$project_id" --quiet
  else
    echo "$api is already enabled"
  fi
done

# Initialize Terraform
echo "Initializing Terraform..."
terraform init

# Plan the changes
echo "Planning Terraform changes..."
terraform plan

# Apply the changes
echo "Applying Terraform changes..."
terraform apply

echo "Infrastructure setup completed!"
echo "Important: After the setup is complete, check the nameservers output and update your domain registrar's nameserver settings accordingly."
echo "You can view the nameservers again by running: terraform output nameservers" 