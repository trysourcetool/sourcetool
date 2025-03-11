# External IP address
resource "google_compute_global_address" "default" {
  name = "lb-ip-${var.environment}"
}

# HTTPS certificate
resource "google_compute_managed_ssl_certificate" "default" {
  name = "cert-${var.environment}"
  managed {
    domains = local.domains
  }
}

# URL map
resource "google_compute_url_map" "default" {
  name = "url-map-${var.environment}"
  default_service = google_compute_backend_service.default.id
}

# HTTPS proxy
resource "google_compute_target_https_proxy" "default" {
  name             = "https-proxy-${var.environment}"
  url_map         = google_compute_url_map.default.id
  ssl_certificates = [google_compute_managed_ssl_certificate.default.id]
}

# Forwarding rule
resource "google_compute_global_forwarding_rule" "https" {
  name       = "https-lb-${var.environment}"
  target     = google_compute_target_https_proxy.default.id
  port_range = "443"
  ip_address = google_compute_global_address.default.address
}

# Backend service
resource "google_compute_backend_service" "default" {
  name = "backend-${var.environment}"

  protocol    = "HTTP"
  port_name   = "http"
  timeout_sec = 30

  backend {
    group = google_compute_region_network_endpoint_group.default.id
  }
}

# Network Endpoint Group
resource "google_compute_region_network_endpoint_group" "default" {
  name                  = "neg-${var.environment}"
  network_endpoint_type = "SERVERLESS"
  region               = var.region
  cloud_run {
    service = google_cloud_run_service.default.name
  }
} 