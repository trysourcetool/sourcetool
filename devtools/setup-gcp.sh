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
  # GCP Project Configuration
  "GCP_PROJECT_ID"
  "GCP_REGION"
  "GCP_ENVIRONMENT"
  
  # Domain Configuration
  "DOMAIN_NAME"
  
  # Database Configuration
  "GCP_SQL_INSTANCE"
  "DB_NAME"
  "DB_USER"
  "DB_PASSWORD"
  "DB_VERSION"
  "DB_TIER"
  "DB_BACKUP_ENABLED"
  "DB_MAX_CONNECTIONS"
  
  # Cloud Run Configuration
  "CLOUD_RUN_SERVICE_NAME"
  "CONTAINER_IMAGE"
  "CLOUD_RUN_CPU"
  "CLOUD_RUN_MEMORY"
  "CLOUD_RUN_MIN_INSTANCES"
  "CLOUD_RUN_MAX_INSTANCES"
  "CLOUD_RUN_CONCURRENCY"
  "CLOUD_RUN_TIMEOUT"
  
  # Network Configuration
  "VPC_CONNECTOR_MACHINE_TYPE"
  "VPC_CONNECTOR_MIN_INSTANCES"
  "VPC_CONNECTOR_MAX_INSTANCES"
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

# Check if application default credentials are set
if [ ! -f "${HOME}/.config/gcloud/application_default_credentials.json" ]; then
  echo "Error: Google Cloud application default credentials not found."
  echo "Please run 'gcloud auth application-default login' to set up authentication."
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
)

for api in "${required_apis[@]}"; do
  echo "Enabling $api..."
  gcloud services enable "$api" --project "$GCP_PROJECT_ID"
done

# Set the working directory to the terraform directory
cd "$(dirname "$0")/../terraform" || exit 1

# Initialize Terraform
echo "Initializing Terraform..."
terraform init

# Create a terraform.tfvars file
cat > terraform.tfvars << EOF
# GCP Project Configuration
project_id  = "$GCP_PROJECT_ID"
region      = "$GCP_REGION"
environment = "$GCP_ENVIRONMENT"

# Domain Configuration
domain_name = "$DOMAIN_NAME"

# Database Configuration
db_instance_name    = "$GCP_SQL_INSTANCE"
db_name            = "$DB_NAME"
db_user            = "$DB_USER"
db_password        = "$DB_PASSWORD"
db_version         = "$DB_VERSION"
db_tier            = "$DB_TIER"
db_backup_enabled  = $DB_BACKUP_ENABLED
db_max_connections = $DB_MAX_CONNECTIONS

# Cloud Run Configuration
cloud_run_service_name   = "$CLOUD_RUN_SERVICE_NAME"
container_image         = "$CONTAINER_IMAGE"
cloud_run_cpu           = "$CLOUD_RUN_CPU"
cloud_run_memory        = "$CLOUD_RUN_MEMORY"
cloud_run_min_instances = $CLOUD_RUN_MIN_INSTANCES
cloud_run_max_instances = $CLOUD_RUN_MAX_INSTANCES
cloud_run_concurrency   = $CLOUD_RUN_CONCURRENCY
cloud_run_timeout       = $CLOUD_RUN_TIMEOUT

# Network Configuration
vpc_connector_machine_type   = "$VPC_CONNECTOR_MACHINE_TYPE"
vpc_connector_min_instances = $VPC_CONNECTOR_MIN_INSTANCES
vpc_connector_max_instances = $VPC_CONNECTOR_MAX_INSTANCES
EOF

# Plan the changes
echo "Planning Terraform changes..."
terraform plan -var-file="terraform.tfvars"

# Apply the changes
echo "Applying Terraform changes..."
terraform apply -var-file="terraform.tfvars"

# Get the Cloud Run service URL
echo "Getting Cloud Run service URL..."
SERVICE_URL=$(gcloud run services describe "$CLOUD_RUN_SERVICE_NAME" \
  --region "$GCP_REGION" \
  --format='value(status.url)')

echo "Setup completed successfully!"
echo "Cloud Run service URL: $SERVICE_URL"

# Print next steps
echo "
Next steps:
1. Set up GitHub repository secrets:
   - GCP_PROJECT_ID: $GCP_PROJECT_ID
   - GCP_REGION: $GCP_REGION
   - GCP_SA_KEY: (Create a service account key with necessary permissions)
   - GCP_SQL_INSTANCE: $GCP_SQL_INSTANCE
   - DB_HOST: (Get from Cloud SQL instance)
   - DB_USER: $DB_USER
   - DB_PASSWORD: $DB_PASSWORD

2. Enable required APIs in GCP:
   - Cloud Run API
   - Cloud SQL Admin API
   - Compute Engine API
   - Service Networking API
   - Artifact Registry API
"

echo "Infrastructure setup completed!"
echo "Important: After the setup is complete, check the nameservers output and update your domain registrar's nameserver settings accordingly."
echo "You can view the nameservers again by running: terraform output nameservers" 