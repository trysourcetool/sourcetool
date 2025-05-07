#!/bin/bash

set -e

# Check for required commands
REQUIRED_TOOLS=("docker" "docker compose" "head" "base64" "sed" "cat")
MISSING_TOOLS=()

for tool in "${REQUIRED_TOOLS[@]}"; do
  # For tools with spaces (like "docker compose"), check the first word
  if [[ "$tool" == *" "* ]]; then
    IFS=' ' read -r cmd subcmd <<< "$tool"
    if ! command -v "$cmd" >/dev/null 2>&1 || ! docker compose version >/dev/null 2>&1; then
      MISSING_TOOLS+=("$tool")
    fi
  else
    if ! command -v "$tool" >/dev/null 2>&1; then
      MISSING_TOOLS+=("$tool")
    fi
  fi
done

if [ ${#MISSING_TOOLS[@]} -ne 0 ]; then
  echo "The following required tools are missing:"
  for tool in "${MISSING_TOOLS[@]}"; do
    echo "  - $tool"
  done
  echo "Please install them before proceeding."
  exit 1
fi

echo "All required tools are installed."

# Check if .env already exists
if [ -f .env ]; then
  read -p ".env file already exists. Overwrite? (y/N): " yn
  case "$yn" in
    [yY]*) ;;
    *) echo "Aborted."; exit 1;;
  esac
fi

# Copy .env.example to .env
cp .env.example .env

# Generate ENCRYPTION_KEY and JWT_KEY
ENCRYPTION_KEY=$(head -c 32 /dev/urandom | base64 | tr -d '\n')
JWT_KEY=$(cat /dev/urandom | base64 | head -c 256)

# Inject keys into .env
sed -i '' "s|ENCRYPTION_KEY=<encryption-key>|ENCRYPTION_KEY=${ENCRYPTION_KEY}|" .env
sed -i '' "s|JWT_KEY=<jwt-key>|JWT_KEY=${JWT_KEY}|" .env

# Ask if Google OAuth will be used
read -p "Do you want to use Google OAuth? (y/N): " use_google_oauth
if [[ "$use_google_oauth" =~ ^[yY]$ ]]; then
  read -p "Enter your Google OAuth CLIENT_ID: " google_client_id
  read -p "Enter your Google OAuth CLIENT_SECRET: " google_client_secret
  sed -i '' "s|GOOGLE_OAUTH_CLIENT_ID=<google-oauth-client-id>|GOOGLE_OAUTH_CLIENT_ID=${google_client_id}|" .env
  sed -i '' "s|GOOGLE_OAUTH_CLIENT_SECRET=<google-oauth-client-secret>|GOOGLE_OAUTH_CLIENT_SECRET=${google_client_secret}|" .env
  echo "Google OAuth CLIENT_ID and CLIENT_SECRET have been set in .env."
else
  echo "Google OAuth values in .env are left as placeholders."
fi

# Ask if SMTP will be configured now
read -p "Do you want to configure SMTP settings now? (y/N): " use_smtp
if [[ "$use_smtp" =~ ^[yY]$ ]]; then
  read -p "Enter your SMTP HOST: " smtp_host
  read -p "Enter your SMTP PORT: " smtp_port
  read -p "Enter your SMTP USERNAME: " smtp_username
  read -p "Enter your SMTP PASSWORD: " smtp_password
  read -p "Enter your SMTP FROM EMAIL: " smtp_from_email
  read -p "Use TLS for SMTP? (true/false): " smtp_use_tls
  sed -i '' "s|SMTP_HOST=<smtp-host>|SMTP_HOST=${smtp_host}|" .env
  sed -i '' "s|SMTP_PORT=<smtp-port>|SMTP_PORT=${smtp_port}|" .env
  sed -i '' "s|SMTP_USERNAME=<smtp-username>|SMTP_USERNAME=${smtp_username}|" .env
  sed -i '' "s|SMTP_PASSWORD=<smtp-password>|SMTP_PASSWORD=${smtp_password}|" .env
  sed -i '' "s|SMTP_FROM_EMAIL=<smtp-from-email>|SMTP_FROM_EMAIL=${smtp_from_email}|" .env
  sed -i '' "s|SMTP_USE_TLS=true|SMTP_USE_TLS=${smtp_use_tls}|" .env
  echo "SMTP settings have been set in .env."
else
  echo "SMTP values in .env are left as placeholders."
fi

# Print instructions
echo ".env file has been created and keys have been generated."
echo "You can edit Google OAuth and SMTP values in .env later if needed."
echo ""
echo "You can now start the local environment with:"
echo "  make start"