#!/bin/bash
set -e

# This script syncs the sdk/go directory to a mirror repository
# Usage: ./sync_sdk_go_mirror.sh <github_pat> <mirror_repo_url>

# Get the GitHub PAT from the first argument
GH_PAT=$1
MIRROR_REPO_URL=$2

if [ -z "$GH_PAT" ] || [ -z "$MIRROR_REPO_URL" ]; then
  echo "Error: GitHub PAT and mirror repository URL are required"
  echo "Usage: ./sync_sdk_go_mirror.sh <github_pat> <mirror_repo_url>"
  exit 1
fi

# Create a temporary directory
TEMP_DIR=$(mktemp -d)
echo "Created temporary directory: $TEMP_DIR"

# Clone the mirror repository
echo "Cloning mirror repository..."
git clone "https://${GH_PAT}@${MIRROR_REPO_URL}" "$TEMP_DIR"

# Remove all files except .git directory
echo "Cleaning mirror repository..."
find "$TEMP_DIR" -mindepth 1 -maxdepth 1 -not -path "$TEMP_DIR/.git" -exec rm -rf {} \;

# Copy all files from the sdk/go directory
echo "Copying sdk/go files to mirror repository..."
cp -R sdk/go/* "$TEMP_DIR/"

# Replace README.md with MIRROR_README.md
if [ -f "$TEMP_DIR/MIRROR_README.md" ]; then
  echo "Replacing README.md with MIRROR_README.md..."
  mv "$TEMP_DIR/MIRROR_README.md" "$TEMP_DIR/README.md"
fi

# Change to the mirror repo directory
cd "$TEMP_DIR"

# Check if there are any changes
if [[ -z $(git status --porcelain) ]]; then
  echo "No changes detected. Exiting."
  exit 0
fi

# Get the latest commit message that affected sdk/go
LATEST_COMMIT_MSG=$(git log -1 --pretty=%B -- sdk/go | head -n 1)

# Add and commit the changes
echo "Committing changes..."
git add .
git commit -m "$LATEST_COMMIT_MSG"

# Push the changes
echo "Pushing changes to mirror repository..."
git push origin main

echo "Sync completed successfully!"
