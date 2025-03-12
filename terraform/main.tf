terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

provider "google-beta" {
  project = var.project_id
  region  = var.region
}

# Enable required APIs
resource "google_project_service" "required_apis" {
  for_each = toset([
    "compute.googleapis.com",           # Compute Engine API
    "run.googleapis.com",               # Cloud Run API
    "vpcaccess.googleapis.com",         # Serverless VPC Access API
    "servicenetworking.googleapis.com", # Service Networking API
    "sqladmin.googleapis.com",          # Cloud SQL Admin API
    "secretmanager.googleapis.com",     # Secret Manager API
    "certificatemanager.googleapis.com" # Certificate Manager API
  ])

  project = var.project_id
  service = each.value

  disable_on_destroy = false
} 