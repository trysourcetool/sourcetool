# Domain Name
output "domain_name" {
  description = "The domain name for the application"
  value       = var.domain_name
}

# Load Balancer IP
output "load_balancer_ip" {
  description = "IP address of the load balancer"
  value       = google_compute_global_address.default.address
}

# ACME Challenge
output "acme_challenge" {
  description = "ACME challenge record for SSL certificate validation"
  value       = google_certificate_manager_dns_authorization.default.dns_record[0].data
}

# DNS Record Name for ACME Challenge
output "acme_challenge_record_name" {
  description = "Record name for ACME challenge"
  value       = google_certificate_manager_dns_authorization.default.dns_record[0].name
} 