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
resource "google_secret_manager_secret" "postgres_host" {
  secret_id = "postgres-host-${var.environment}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "postgres_db" {
  secret_id = "postgres-db-${var.environment}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "postgres_user" {
  secret_id = "postgres-user-${var.environment}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "postgres_password" {
  secret_id = "postgres-password-${var.environment}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "postgres_port" {
  secret_id = "postgres-port-${var.environment}"
  replication {
    auto {}
  }
}

# Redis configuration
resource "google_secret_manager_secret" "redis_host" {
  secret_id = "redis-host-${var.environment}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "redis_password" {
  secret_id = "redis-password-${var.environment}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "redis_port" {
  secret_id = "redis-port-${var.environment}"
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

resource "google_secret_manager_secret" "google_oauth_callback_url" {
  secret_id = "google-oauth-callback-url-${var.environment}"
  replication {
    auto {}
  }
}

# SMTP configuration
resource "google_secret_manager_secret" "smtp_host" {
  secret_id = "smtp-host-${var.environment}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "smtp_port" {
  secret_id = "smtp-port-${var.environment}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "smtp_username" {
  secret_id = "smtp-username-${var.environment}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "smtp_password" {
  secret_id = "smtp-password-${var.environment}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "smtp_from_email" {
  secret_id = "smtp-from-email-${var.environment}"
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
resource "google_secret_manager_secret_version" "postgres_host" {
  secret      = google_secret_manager_secret.postgres_host.id
  secret_data = google_sql_database_instance.instance.connection_name
}

resource "google_secret_manager_secret_version" "postgres_db" {
  secret      = google_secret_manager_secret.postgres_db.id
  secret_data = var.db_name
}

resource "google_secret_manager_secret_version" "postgres_user" {
  secret      = google_secret_manager_secret.postgres_user.id
  secret_data = var.db_user
}

resource "google_secret_manager_secret_version" "postgres_password" {
  secret      = google_secret_manager_secret.postgres_password.id
  secret_data = var.db_password
}

resource "google_secret_manager_secret_version" "postgres_port" {
  secret      = google_secret_manager_secret.postgres_port.id
  secret_data = "5432"
}

# Redis configuration
resource "google_secret_manager_secret_version" "redis_host" {
  secret      = google_secret_manager_secret.redis_host.id
  secret_data = var.redis_host
}

resource "google_secret_manager_secret_version" "redis_password" {
  secret      = google_secret_manager_secret.redis_password.id
  secret_data = var.redis_password
}

resource "google_secret_manager_secret_version" "redis_port" {
  secret      = google_secret_manager_secret.redis_port.id
  secret_data = var.redis_port
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

resource "google_secret_manager_secret_version" "google_oauth_callback_url" {
  secret      = google_secret_manager_secret.google_oauth_callback_url.id
  secret_data = var.google_oauth_callback_url
}

# SMTP configuration
resource "google_secret_manager_secret_version" "smtp_host" {
  secret      = google_secret_manager_secret.smtp_host.id
  secret_data = var.smtp_host
}

resource "google_secret_manager_secret_version" "smtp_port" {
  secret      = google_secret_manager_secret.smtp_port.id
  secret_data = var.smtp_port
}

resource "google_secret_manager_secret_version" "smtp_username" {
  secret      = google_secret_manager_secret.smtp_username.id
  secret_data = var.smtp_username
}

resource "google_secret_manager_secret_version" "smtp_password" {
  secret      = google_secret_manager_secret.smtp_password.id
  secret_data = var.smtp_password
}

resource "google_secret_manager_secret_version" "smtp_from_email" {
  secret      = google_secret_manager_secret.smtp_from_email.id
  secret_data = var.smtp_from_email
}

# IAM policy for Secret Manager
resource "google_secret_manager_secret_iam_member" "secret_access" {
  for_each = toset([
    google_secret_manager_secret.encryption_key.id,
    google_secret_manager_secret.jwt_key.id,
    google_secret_manager_secret.postgres_host.id,
    google_secret_manager_secret.postgres_db.id,
    google_secret_manager_secret.postgres_user.id,
    google_secret_manager_secret.postgres_password.id,
    google_secret_manager_secret.postgres_port.id,
    google_secret_manager_secret.redis_host.id,
    google_secret_manager_secret.redis_password.id,
    google_secret_manager_secret.redis_port.id,
    google_secret_manager_secret.google_oauth_client_id.id,
    google_secret_manager_secret.google_oauth_client_secret.id,
    google_secret_manager_secret.google_oauth_callback_url.id,
    google_secret_manager_secret.smtp_host.id,
    google_secret_manager_secret.smtp_port.id,
    google_secret_manager_secret.smtp_username.id,
    google_secret_manager_secret.smtp_password.id,
    google_secret_manager_secret.smtp_from_email.id,
  ])

  secret_id = each.key
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloud_run_sa.email}"
} 