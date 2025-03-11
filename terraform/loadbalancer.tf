# External IP address
resource "google_compute_global_address" "default" {
  name = "lb-ip-${var.environment}"
}

# DNS Authorization
resource "google_certificate_manager_dns_authorization" "default" {
  name        = "dns-auth-${var.environment}"
  description = "DNS Authorization for ${var.domain_name}"
  domain      = var.domain_name
}

# Certificate
resource "google_certificate_manager_certificate" "default" {
  name        = "cert-${var.environment}"
  description = "Certificate for ${var.domain_name}"
  scope       = "DEFAULT"

  managed {
    domains = [
      "*.${var.domain_name}"
    ]
    dns_authorizations = [google_certificate_manager_dns_authorization.default.id]
  }
}

# Certificate Map
resource "google_certificate_manager_certificate_map" "default" {
  name        = "cert-map-${var.environment}"
  description = "Certificate map for ${var.domain_name}"
}

# Certificate Map Entry
resource "google_certificate_manager_certificate_map_entry" "default" {
  name        = "cert-map-entry-${var.environment}"
  description = "Certificate map entry for ${var.domain_name}"
  map         = google_certificate_manager_certificate_map.default.name
  certificates = [google_certificate_manager_certificate.default.id]
  hostname    = "*.${var.domain_name}"
}

# URL map
resource "google_compute_url_map" "default" {
  name            = "url-map-${var.environment}"
  default_service = google_compute_backend_service.default.id
}

# HTTPS proxy
resource "google_compute_target_https_proxy" "default" {
  name             = "https-proxy-${var.environment}"
  url_map          = google_compute_url_map.default.id
  certificate_map  = "//certificatemanager.googleapis.com/${google_certificate_manager_certificate_map.default.id}"
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