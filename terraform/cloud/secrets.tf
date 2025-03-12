# Secret Manager secrets
# Common configuration
resource "google_secret_manager_secret" "encryption_key" {
  secret_id = "encryption-key-${var.environment}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "jwt_key" {
  secret_id = "jwt-key-${var.environment}"
  replication {
    auto {}
  }
}

# Database configuration
resource "google_secret_manager_secret" "postgres_password" {
  secret_id = "postgres-password-${var.environment}"
  replication {
    auto {}
  }
}

# Redis configuration
resource "google_secret_manager_secret" "redis_password" {
  secret_id = "redis-password-${var.environment}"
  replication {
    auto {}
  }
}

# OAuth configuration
resource "google_secret_manager_secret" "google_oauth_client_id" {
  secret_id = "google-oauth-client-id-${var.environment}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "google_oauth_client_secret" {
  secret_id = "google-oauth-client-secret-${var.environment}"
  replication {
    auto {}
  }
}

# SMTP configuration
resource "google_secret_manager_secret" "smtp_password" {
  secret_id = "smtp-password-${var.environment}"
  replication {
    auto {}
  }
}

# Secret values
# Common configuration
resource "google_secret_manager_secret_version" "encryption_key" {
  secret      = google_secret_manager_secret.encryption_key.id
  secret_data = var.encryption_key
}

resource "google_secret_manager_secret_version" "jwt_key" {
  secret      = google_secret_manager_secret.jwt_key.id
  secret_data = var.jwt_key
}

# Database configuration
resource "google_secret_manager_secret_version" "postgres_password" {
  secret      = google_secret_manager_secret.postgres_password.id
  secret_data = google_sql_user.user.password
}

# Redis configuration
resource "google_secret_manager_secret_version" "redis_password" {
  secret      = google_secret_manager_secret.redis_password.id
  secret_data = google_redis_instance.default.auth_string
}

# OAuth configuration
resource "google_secret_manager_secret_version" "google_oauth_client_id" {
  secret      = google_secret_manager_secret.google_oauth_client_id.id
  secret_data = var.google_oauth_client_id
}

resource "google_secret_manager_secret_version" "google_oauth_client_secret" {
  secret      = google_secret_manager_secret.google_oauth_client_secret.id
  secret_data = var.google_oauth_client_secret
}

# SMTP configuration
resource "google_secret_manager_secret_version" "smtp_password" {
  secret      = google_secret_manager_secret.smtp_password.id
  secret_data = var.smtp_password
}

# IAM policy for Secret Manager
locals {
  secrets = {
    encryption_key              = google_secret_manager_secret.encryption_key.id
    jwt_key                    = google_secret_manager_secret.jwt_key.id
    postgres_password          = google_secret_manager_secret.postgres_password.id
    redis_password             = google_secret_manager_secret.redis_password.id
    google_oauth_client_id     = google_secret_manager_secret.google_oauth_client_id.id
    google_oauth_client_secret = google_secret_manager_secret.google_oauth_client_secret.id
    smtp_password              = google_secret_manager_secret.smtp_password.id
  }
}

resource "google_secret_manager_secret_iam_member" "secret_access" {
  for_each = local.secrets

  secret_id = each.value
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloud_run_sa.email}"
} 