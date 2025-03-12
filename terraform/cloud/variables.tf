# GCP Project Configuration
variable "project_id" {
  description = "The GCP project ID"
  type        = string
}

variable "region" {
  description = "The GCP region"
  type        = string
}

variable "environment" {
  description = "Environment name (e.g., prod, staging)"
  type        = string
}

# Database Configuration
variable "db_instance_name" {
  description = "Name of the Cloud SQL instance"
  type        = string
}

variable "db_version" {
  description = "Cloud SQL database version"
  type        = string
}

variable "db_tier" {
  description = "Cloud SQL instance tier"
  type        = string
}

variable "db_name" {
  description = "Name of the database to create"
  type        = string
}

variable "db_user" {
  description = "Database user name"
  type        = string
}

variable "db_password" {
  description = "Database user password"
  type        = string
  sensitive   = true
}

variable "db_backup_enabled" {
  description = "Whether to enable database backups"
  type        = bool
  default     = true
}

variable "db_max_connections" {
  description = "Maximum number of database connections"
  type        = number
  default     = 100
}

# Cloud Run Configuration
variable "cloud_run_service_name" {
  description = "Name of the Cloud Run service"
  type        = string
}

variable "container_image" {
  description = "Container image to deploy"
  type        = string
}

variable "job_container_image" {
  description = "Container image to deploy for the job"
  type        = string
}

variable "cloud_run_cpu" {
  description = "CPU allocation for Cloud Run service"
  type        = string
  default     = "1"
}

variable "cloud_run_memory" {
  description = "Memory allocation for Cloud Run service"
  type        = string
  default     = "512Mi"
}

variable "cloud_run_min_instances" {
  description = "Minimum number of Cloud Run instances"
  type        = number
  default     = 0
}

variable "cloud_run_max_instances" {
  description = "Maximum number of Cloud Run instances"
  type        = number
  default     = 10
}

variable "cloud_run_concurrency" {
  description = "Maximum number of concurrent requests per instance"
  type        = number
  default     = 80
}

variable "cloud_run_timeout" {
  description = "Maximum request timeout in seconds"
  type        = number
  default     = 300
}

# Network Configuration
variable "vpc_connector_machine_type" {
  description = "The machine type to use for the VPC connector"
  type        = string
  default     = "e2-micro"
}

variable "vpc_connector_min_instances" {
  description = "Minimum number of instances for the VPC connector"
  type        = number
  default     = 2
}

variable "vpc_connector_max_instances" {
  description = "Maximum number of instances for the VPC connector"
  type        = number
  default     = 3
}

# Domain Configuration
variable "domain_name" {
  description = "The base domain name (e.g., stg.trysourcetool.com)"
  type        = string
}

# Common Configuration
variable "encryption_key" {
  description = "Encryption key for sensitive data"
  type        = string
  sensitive   = true
}

variable "jwt_key" {
  description = "Key for JWT signing"
  type        = string
  sensitive   = true
}

# SMTP Configuration
variable "smtp_host" {
  description = "SMTP host"
  type        = string
}

variable "smtp_port" {
  description = "SMTP port"
  type        = string
}

variable "smtp_username" {
  description = "SMTP username"
  type        = string
}

variable "smtp_password" {
  description = "SMTP password"
  type        = string
  sensitive   = true
}

variable "smtp_from_email" {
  description = "SMTP from email address"
  type        = string
}

# OAuth Configuration
variable "google_oauth_client_id" {
  description = "Google OAuth client ID"
  type        = string
}

variable "google_oauth_client_secret" {
  description = "Google OAuth client secret"
  type        = string
  sensitive   = true
}

variable "google_oauth_callback_url" {
  description = "Google OAuth callback URL"
  type        = string
}
