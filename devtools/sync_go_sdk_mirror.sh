#!/bin/bash
set -e

# This script syncs the sdk/go directory to a mirror repository
# Usage: ./sync_go_sdk_mirror.sh <github_pat> <mirror_repo_url>

# Get the GitHub PAT from the first argument
GH_PAT=$1
MIRROR_REPO_URL=$2

# Check if the GitHub PAT is provided
if [ -z "$GH_PAT" ]; then
  echo "Error: GitHub PAT is not provided or is empty"
  echo "Please ensure the GH_PAT secret is properly configured in your repository settings"
  exit 1
fi

# Check if the mirror repository URL is provided
if [ -z "$MIRROR_REPO_URL" ]; then
  echo "Error: Mirror repository URL is not provided"
  echo "Usage: ./sync_go_sdk_mirror.sh <github_pat> <mirror_repo_url>"
  exit 1
fi

# Get the latest commit message that affected sdk/go from the source repository
# Do this before changing directories
LATEST_COMMIT_MSG=$(git log -1 --pretty=%B -- sdk/go | head -n 1)
COMMIT_AUTHOR_NAME=$(git log -1 --pretty="%an" -- sdk/go)
COMMIT_AUTHOR_EMAIL=$(git log -1 --pretty="%ae" -- sdk/go)
echo "Latest commit message: $LATEST_COMMIT_MSG"
echo "Commit author: $COMMIT_AUTHOR_NAME <$COMMIT_AUTHOR_EMAIL>"

# If commit message is empty, use a default message
if [ -z "$LATEST_COMMIT_MSG" ]; then
  LATEST_COMMIT_MSG="Update Go SDK files"
  echo "Using default commit message: $LATEST_COMMIT_MSG"
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

# Configure git with the original committer's information
git config user.name "$COMMIT_AUTHOR_NAME"
git config user.email "$COMMIT_AUTHOR_EMAIL"

# Check if there are any changes
if [[ -z $(git status --porcelain) ]]; then
  echo "No changes detected. Exiting."
  exit 0
fi

# Add and commit the changes
echo "Committing changes..."
git add .
git commit -m "$LATEST_COMMIT_MSG" || git commit -m "Update Go SDK files"

# Push the changes
echo "Pushing changes to mirror repository..."
git push origin main

echo "Sync completed successfully!"
