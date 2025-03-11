# VPC Network
resource "google_compute_network" "vpc" {
  name                    = "sourcetool-vpc-${var.environment}"
  auto_create_subnetworks = false
}

# Subnet
resource "google_compute_subnetwork" "subnet" {
  name          = "sourcetool-subnet-${var.environment}"
  ip_cidr_range = "10.0.0.0/24"
  network       = google_compute_network.vpc.id
  region        = var.region
}

# Cloud SQL private IP
resource "google_compute_global_address" "private_ip_address" {
  name          = "sourcetool-private-ip-${var.environment}"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.vpc.id
}

# VPC peering
resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.vpc.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}

# Serverless VPC Access connector
resource "google_vpc_access_connector" "connector" {
  name          = "vpc-connector-${var.environment}"
  ip_cidr_range = "10.8.0.0/28"
  network       = google_compute_network.vpc.self_link
  region        = var.region
  machine_type  = var.vpc_connector_machine_type
  min_instances = var.vpc_connector_min_instances
  max_instances = var.vpc_connector_max_instances

  depends_on = [
    google_project_service.required_apis["vpcaccess.googleapis.com"]
  ]
} 