# Enable Redis API
resource "google_project_service" "redis" {
  service = "redis.googleapis.com"

  disable_on_destroy = false
}

# Redis instance
resource "google_redis_instance" "default" {
  name           = "redis-${var.environment}"
  tier           = "STANDARD_HA"
  memory_size_gb = 1

  region                  = var.region
  location_id            = "${var.region}-a"
  alternative_location_id = "${var.region}-b"

  authorized_network = google_compute_network.vpc.id
  connect_mode       = "PRIVATE_SERVICE_ACCESS"
  auth_enabled       = true

  redis_version     = "REDIS_6_X"
  display_name      = "Redis for ${var.environment}"

  maintenance_policy {
    weekly_maintenance_window {
      day = "SUNDAY"
      start_time {
        hours   = 2
        minutes = 0
        seconds = 0
        nanos   = 0
      }
    }
  }

  depends_on = [
    google_project_service.redis,
    google_service_networking_connection.private_vpc_connection
  ]
} 