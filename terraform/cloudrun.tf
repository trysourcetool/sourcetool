# Cloud Run service
resource "google_cloud_run_service" "default" {
  name     = var.cloud_run_service_name
  location = var.region

  template {
    spec {
      service_account_name = google_service_account.cloud_run_sa.email
      containers {
        image = var.container_image
        
        resources {
          limits = {
            cpu    = var.cloud_run_cpu
            memory = var.cloud_run_memory
          }
        }

        # Common configuration
        env {
          name  = "ENV"
          value = var.environment
        }

        env {
          name  = "DOMAIN"
          value = var.domain_name
        }

        env {
          name = "ENCRYPTION_KEY"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.encryption_key.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "JWT_KEY"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.jwt_key.secret_id
              key  = "latest"
            }
          }
        }

        # Database configuration
        env {
          name = "POSTGRES_HOST"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.postgres_host.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "POSTGRES_DB"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.postgres_db.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "POSTGRES_USER"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.postgres_user.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "POSTGRES_PASSWORD"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.postgres_password.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "POSTGRES_PORT"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.postgres_port.secret_id
              key  = "latest"
            }
          }
        }

        # Redis configuration
        env {
          name = "REDIS_HOST"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.redis_host.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "REDIS_PASSWORD"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.redis_password.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "REDIS_PORT"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.redis_port.secret_id
              key  = "latest"
            }
          }
        }

        # OAuth configuration
        env {
          name = "GOOGLE_OAUTH_CLIENT_ID"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.google_oauth_client_id.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "GOOGLE_OAUTH_CLIENT_SECRET"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.google_oauth_client_secret.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "GOOGLE_OAUTH_CALLBACK_URL"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.google_oauth_callback_url.secret_id
              key  = "latest"
            }
          }
        }

        # SMTP configuration
        env {
          name = "SMTP_HOST"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.smtp_host.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "SMTP_PORT"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.smtp_port.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "SMTP_USERNAME"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.smtp_username.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "SMTP_PASSWORD"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.smtp_password.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "SMTP_FROM_EMAIL"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.smtp_from_email.secret_id
              key  = "latest"
            }
          }
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/minScale"        = var.cloud_run_min_instances
        "autoscaling.knative.dev/maxScale"        = var.cloud_run_max_instances
        "run.googleapis.com/vpc-access-connector"  = google_vpc_access_connector.connector.name
        "run.googleapis.com/vpc-access-egress"    = "private-ranges-only"
        "run.googleapis.com/cloudsql-instances"   = google_sql_database_instance.instance.connection_name
        "run.googleapis.com/cpu-throttling"       = "false"
        "run.googleapis.com/execution-environment" = "gen2"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [
    google_project_service.required_apis["run.googleapis.com"],
    google_secret_manager_secret_version.encryption_key,
    google_secret_manager_secret_version.jwt_key,
    google_secret_manager_secret_version.postgres_host,
    google_secret_manager_secret_version.postgres_db,
    google_secret_manager_secret_version.postgres_user,
    google_secret_manager_secret_version.postgres_password,
    google_secret_manager_secret_version.postgres_port,
    google_secret_manager_secret_version.redis_host,
    google_secret_manager_secret_version.redis_password,
    google_secret_manager_secret_version.redis_port,
    google_secret_manager_secret_version.google_oauth_client_id,
    google_secret_manager_secret_version.google_oauth_client_secret,
    google_secret_manager_secret_version.google_oauth_callback_url,
    google_secret_manager_secret_version.smtp_host,
    google_secret_manager_secret_version.smtp_port,
    google_secret_manager_secret_version.smtp_username,
    google_secret_manager_secret_version.smtp_password,
    google_secret_manager_secret_version.smtp_from_email,
  ]
}

# IAM policy for Cloud Run service
resource "google_cloud_run_service_iam_member" "default" {
  location = google_cloud_run_service.default.location
  project  = google_cloud_run_service.default.project
  service  = google_cloud_run_service.default.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.cloud_run_sa.email}"
} 