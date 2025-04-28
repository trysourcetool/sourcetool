#!/bin/bash
set -e

# This script creates and pushes a git tag
# Usage: ./create_tag.sh <tag_name>

if [ -z "$1" ]; then
  echo "Error: Tag name is required"
  echo "Usage: ./create_tag.sh <tag_name>"
  exit 1
fi

TAG=$1

# Validate tag format
if ! echo "$TAG" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$"; then
  echo "Error: Invalid version format: $TAG"
  echo "Version must follow semver format: vX.Y.Z"
  exit 1
fi

# Check if there are any uncommitted changes
if ! git diff-index --quiet HEAD --; then
  echo "Error: There are uncommitted changes in the repository"
  echo "Please commit or stash them before creating a tag"
  exit 1
fi

# Pull latest changes
git pull origin main

PREFIXES=("sdk/go" "mcp/docs-mcp-server")

# Create and push additional tags
for PREFIX in "${PREFIXES[@]}"; do
  ADDITIONAL_TAG="$PREFIX/$TAG"
  echo "Creating tag $ADDITIONAL_TAG..."
  git tag $ADDITIONAL_TAG
  git push origin $ADDITIONAL_TAG
done

# Create and push main tag
echo "Creating tag $TAG..."
git tag $TAG
git push origin $TAG

echo ""
echo "✨ Tag $TAG has been created and pushed successfully!"
echo "GitHub Actions will handle the release process automatically." 