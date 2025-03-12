# VPC Network
resource "google_compute_network" "vpc" {
  name                    = "vpc-${var.environment}"
  auto_create_subnetworks = false
}

# Subnet
resource "google_compute_subnetwork" "subnet" {
  name          = "subnet-${var.environment}"
  ip_cidr_range = "10.0.0.0/24"
  network       = google_compute_network.vpc.id
  region        = var.region
}

# Cloud SQL private IP
resource "google_compute_global_address" "private_ip_address" {
  name          = "private-ip-${var.environment}"
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

# Cloud NAT IP address
resource "google_compute_address" "nat_ip" {
  name         = "nat-ip-${var.environment}"
  region       = var.region
  address_type = "EXTERNAL"
}

# Cloud Router
resource "google_compute_router" "router" {
  name    = "router-${var.environment}"
  region  = var.region
  network = google_compute_network.vpc.id
}

# Cloud NAT
resource "google_compute_router_nat" "nat" {
  name                               = "nat-${var.environment}"
  router                             = google_compute_router.router.name
  region                             = var.region
  nat_ip_allocate_option             = "MANUAL_ONLY"
  nat_ips                            = [google_compute_address.nat_ip.self_link]
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"

  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
} 