#!/bin/bash
set -e

# This script deletes a git tag both locally and remotely
# Usage: ./delete_tag.sh <tag_name>

if [ -z "$1" ]; then
  echo "Error: Tag name is required"
  echo "Usage: ./delete_tag.sh <tag_name>"
  exit 1
fi

TAG=$1

# Delete local tag
echo "Deleting local tag $TAG..."
git tag -d $TAG || echo "Local tag $TAG not found"

# Delete remote tag from origin
echo "Deleting remote tag $TAG from origin..."
git push origin :refs/tags/$TAG || echo "Remote tag $TAG not found or already deleted"

PREFIXES=("sdk/go" "mcp/docs-mcp-server")

# Delete additional tags
for PREFIX in "${PREFIXES[@]}"; do
  ADDITIONAL_TAG="$PREFIX/$TAG"
  echo "Deleting additional tag $ADDITIONAL_TAG..."
  git tag -d "$ADDITIONAL_TAG" || echo "Local tag $ADDITIONAL_TAG not found"
  git push origin ":refs/tags/$ADDITIONAL_TAG" || echo "Remote tag $ADDITIONAL_TAG not found or already deleted"
done

# If SDK repo exists, delete tag there too
if [ -d "../sourcetool-go" ]; then
  echo "SDK repository found, deleting tag from sourcetool-go..."
  (
    cd ../sourcetool-go
    git tag -d $TAG || echo "Local SDK tag $TAG not found"
    git push origin :refs/tags/$TAG || echo "Remote SDK tag $TAG not found or already deleted"
  )
else
  echo "SDK repository not found at ../sourcetool-go"
  echo "To delete SDK tag, run these commands in the SDK repository:"
  echo "  git tag -d $TAG"
  echo "  git push origin :refs/tags/$TAG"
fi 