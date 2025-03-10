#!/bin/sh
set -e

# Create debug output directory if it doesn't exist
mkdir -p /debug_output

# Debug: Show the contents of the static directories
echo "Contents of /app/static:"
ls -la /app/static
echo "Contents of /app/static (recursive):"
find /app/static -type f | sort

echo "Contents of /app/static-full:"
ls -la /app/static-full
echo "Contents of /app/static-full (recursive):"
find /app/static-full -type f | sort

# Debug: Check for index.html file
echo "Checking for index.html file:"
if [ -f "/app/static/index.html" ]; then
  echo "index.html found in /app/static"
  cat /app/static/index.html | head -10
else
  echo "index.html NOT found in /app/static"
fi

if [ -f "/app/static-full/client/index.html" ]; then
  echo "index.html found in /app/static-full/client"
  cat /app/static-full/client/index.html | head -10
else
  echo "index.html NOT found in /app/static-full/client"
fi

# Debug: Check if localization files exist
echo "Checking localization files:"
ls -la /app/static/locales/en/common.json || echo "English localization file not found"
ls -la /app/static/locales/ja/common.json || echo "Japanese localization file not found"

# Save debug info to file
echo "Contents of /app/static:" > /debug_output/static_files.txt
ls -la /app/static >> /debug_output/static_files.txt
echo "Contents of /app/static (recursive):" >> /debug_output/static_files.txt
find /app/static -type f | sort >> /debug_output/static_files.txt

echo "Contents of /app/static-full:" >> /debug_output/static_files.txt
ls -la /app/static-full >> /debug_output/static_files.txt
echo "Contents of /app/static-full (recursive):" >> /debug_output/static_files.txt
find /app/static-full -type f | sort >> /debug_output/static_files.txt

# Check for React Router assets
echo "Checking for React Router assets:" >> /debug_output/static_files.txt
find /app/static -name "*react-router*" >> /debug_output/static_files.txt
find /app/static-full -name "*react-router*" >> /debug_output/static_files.txt

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until PGPASSWORD=$POSTGRES_PASSWORD psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB -c '\q'; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done
echo "PostgreSQL is up - executing migrations"

# Run migrations
cd /app/migrations
migrate -path=. -database "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable" up

# Start the server
echo "Starting the server..."
echo "Environment variables:"
echo "STATIC_FILES_DIR=$STATIC_FILES_DIR"

# Save environment variables to file
echo "Environment variables:" > /debug_output/environment.txt
env >> /debug_output/environment.txt

exec /app/server
