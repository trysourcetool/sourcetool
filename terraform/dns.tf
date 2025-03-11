# DNS zone
resource "google_dns_managed_zone" "default" {
  name        = "dns-zone-${var.environment}"
  dns_name    = "${var.domain_name}."  # Trailing dot is required for DNS names
  description = "DNS zone for ${var.domain_name}"
}

# A record for wildcard subdomains
resource "google_dns_record_set" "wildcard" {
  name         = "*.${var.domain_name}."  # Trailing dot is required for DNS names
  managed_zone = google_dns_managed_zone.default.name
  type         = "A"
  ttl          = 300

  rrdatas = [google_compute_global_address.default.address]
}

# Output the nameservers
output "nameservers" {
  description = "The nameservers for the DNS zone"
  value       = google_dns_managed_zone.default.name_servers
} 