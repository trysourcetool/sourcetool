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

# Extract version number without 'v'
VERSION_NO_V=${TAG#v}

# Validate that the version is present in the relevant files
INDEX_TS="mcp/docs-mcp-server/src/index.ts"
COMMON_JSON="frontend/app/locales/en/common.json"
RUNTIME_GO="sdk/go/runtime.go"

if ! grep -q "version: '${VERSION_NO_V}'" "$INDEX_TS"; then
  echo "Error: $INDEX_TS does not contain version: '${VERSION_NO_V}'"
  exit 1
fi

if ! grep -q "\"components_layout_app_version\": \"${TAG}\"" "$COMMON_JSON"; then
  echo "Error: $COMMON_JSON does not contain components_layout_app_version: ${TAG}"
  exit 1
fi

if ! grep -q "SdkVersion: \"${VERSION_NO_V}\"" "$RUNTIME_GO"; then
  echo "Error: $RUNTIME_GO does not contain SdkVersion: \"${VERSION_NO_V}\""
  exit 1
fi

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