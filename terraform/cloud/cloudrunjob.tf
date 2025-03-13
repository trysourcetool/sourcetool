# Cloud Run Job for database migrations
resource "google_cloud_run_v2_job" "migrate" {
  name     = var.cloud_run_job_migrate_name
  location = var.region

  template {
    template {
      containers {
        image = var.job_container_image

        resources {
          limits = {
            cpu    = "1000m"
            memory = "512Mi"
          }
        }

        # Common configuration
        env {
          name  = "DOMAIN"
          value = var.domain_name
        }

        env {
          name = "ENCRYPTION_KEY"
          value_source {
            secret_key_ref {
              secret = google_secret_manager_secret.encryption_key.secret_id
              version = "latest"
            }
          }
        }

        env {
          name = "JWT_KEY"
          value_source {
            secret_key_ref {
              secret = google_secret_manager_secret.jwt_key.secret_id
              version = "latest"
            }
          }
        }

        # Database configuration
        env {
          name  = "POSTGRES_HOST"
          value = var.postgres_host
        }

        env {
          name  = "POSTGRES_DB"
          value = var.postgres_db
        }

        env {
          name  = "POSTGRES_USER"
          value = var.postgres_user
        }

        env {
          name = "POSTGRES_PASSWORD"
          value_source {
            secret_key_ref {
              secret = google_secret_manager_secret.postgres_password.secret_id
              version = "latest"
            }
          }
        }

        env {
          name  = "POSTGRES_PORT"
          value = var.postgres_port
        }

        # Redis configuration
        env {
          name  = "REDIS_HOST"
          value = var.redis_host
        }

        env {
          name = "REDIS_PASSWORD"
          value_source {
            secret_key_ref {
              secret = google_secret_manager_secret.redis_password.secret_id
              version = "latest"
            }
          }
        }

        env {
          name  = "REDIS_PORT"
          value = var.redis_port
        }

        # OAuth configuration
        env {
          name = "GOOGLE_OAUTH_CLIENT_ID"
          value_source {
            secret_key_ref {
              secret = google_secret_manager_secret.google_oauth_client_id.secret_id
              version = "latest"
            }
          }
        }

        env {
          name = "GOOGLE_OAUTH_CLIENT_SECRET"
          value_source {
            secret_key_ref {
              secret = google_secret_manager_secret.google_oauth_client_secret.secret_id
              version = "latest"
            }
          }
        }

        env {
          name  = "GOOGLE_OAUTH_CALLBACK_URL"
          value = var.google_oauth_callback_url
        }

        # SMTP configuration
        env {
          name  = "SMTP_HOST"
          value = var.smtp_host
        }

        env {
          name  = "SMTP_PORT"
          value = var.smtp_port
        }

        env {
          name  = "SMTP_USERNAME"
          value = var.smtp_username
        }

        env {
          name = "SMTP_PASSWORD"
          value_source {
            secret_key_ref {
              secret = google_secret_manager_secret.smtp_password.secret_id
              version = "latest"
            }
          }
        }

        env {
          name  = "SMTP_FROM_EMAIL"
          value = var.smtp_from_email
        }
      }

      service_account = google_service_account.cloud_run_sa.email
      
      vpc_access {
        connector = google_vpc_access_connector.connector.id
        egress    = "PRIVATE_RANGES_ONLY"
      }
    }
  }

  depends_on = [
    google_project_service.required_apis["run.googleapis.com"],
    google_secret_manager_secret_version.encryption_key,
    google_secret_manager_secret_version.jwt_key,
    google_secret_manager_secret_version.postgres_password,
    google_secret_manager_secret_version.redis_password,
    google_secret_manager_secret_version.google_oauth_client_id,
    google_secret_manager_secret_version.google_oauth_client_secret,
    google_secret_manager_secret_version.smtp_password,
  ]
} 