#!/bin/bash

# Exit on error
set -e

# Show usage
show_usage() {
  echo "Usage: $0 [OPTIONS]"
  echo ""
  echo "Options:"
  echo "  -h, --help     Show this help message"
  echo "  -p, --plan     Only show plan without applying"
  echo ""
  echo "Example:"
  echo "  $0            # Apply cloud infrastructure"
  echo "  $0 -p        # Show plan for cloud infrastructure"
}

# Parse arguments
PLAN_ONLY=false

while [[ $# -gt 0 ]]; do
  case $1 in
    -h|--help)
      show_usage
      exit 0
      ;;
    -p|--plan)
      PLAN_ONLY=true
      shift
      ;;
    *)
      echo "Error: Unknown argument $1"
      show_usage
      exit 1
      ;;
  esac
done

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

# Set the working directory
TERRAFORM_DIR="$(dirname "$0")/../terraform/cloud"
if [ ! -d "$TERRAFORM_DIR" ]; then
  echo "Error: Directory $TERRAFORM_DIR does not exist"
  exit 1
fi

# Check if terraform.tfvars exists
if [ ! -f "$TERRAFORM_DIR/terraform.tfvars" ]; then
  echo "Error: terraform.tfvars file not found in $TERRAFORM_DIR"
  exit 1
fi

# Set the working directory
cd "$TERRAFORM_DIR" || exit 1

# Get project_id and environment from terraform.tfvars
project_id=$(grep 'project_id' "terraform.tfvars" | cut -d'=' -f2 | tr -d ' "')
environment=$(grep 'environment' "terraform.tfvars" | cut -d'=' -f2 | tr -d ' "')

if [ -z "$project_id" ]; then
  echo "Error: Could not find project_id in terraform.tfvars"
  exit 1
fi

if [ -z "$environment" ]; then
  echo "Error: Could not find environment in terraform.tfvars"
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
  "compute.googleapis.com"            # Compute Engine API
  "run.googleapis.com"                # Cloud Run API
  "vpcaccess.googleapis.com"          # Serverless VPC Access API
  "servicenetworking.googleapis.com"  # Service Networking API
  "sqladmin.googleapis.com"           # Cloud SQL Admin API
  "secretmanager.googleapis.com"      # Secret Manager API
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

# Apply the changes if not in plan-only mode
if [ "$PLAN_ONLY" = false ]; then
  echo "Applying Terraform changes..."
  terraform apply
  
  echo "Infrastructure setup completed!"
  echo ""
  echo "Important: Please complete the following steps to finish the setup:"
  echo ""
  echo "1. DNS Configuration:"
  echo "   a. Add an A record for your domain pointing to the Load Balancer IP:"
  echo "      Domain: $(terraform output -raw domain_name 2>/dev/null || echo '<your-domain>')"
  echo "      Type: A"
  echo "      Value: $(terraform output -raw load_balancer_ip 2>/dev/null || echo '<run: terraform output load_balancer_ip>')"
  echo ""
  echo "2. SSL Certificate Setup:"
  echo "   Add the following TXT record for ACME challenge:"
  echo "      Name: $(terraform output -raw acme_challenge_record_name 2>/dev/null || echo '<run: terraform output acme_challenge_record_name>')"
  echo "      Type: CNAME"
  echo "      Value: $(terraform output -raw acme_challenge 2>/dev/null || echo '<run: terraform output acme_challenge>')"
  echo ""
  echo "Note: DNS propagation can take time (up to 48 hours). You can check the certificate status with:"
  echo "gcloud certificate-manager certificates describe cert-${environment} --project=$project_id"
fi 