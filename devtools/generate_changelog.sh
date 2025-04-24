#!/bin/bash
set -e

# This script generates a changelog from git history
# Usage: ./generate_changelog.sh [previous_tag]
# If previous_tag is not provided, it will be automatically detected

# Get the previous tag if not provided
if [ -z "$1" ]; then
  PREV_TAG=$(git describe --tags --abbrev=0 HEAD^ 2>/dev/null || echo "")
else
  PREV_TAG=$1
fi

# Generate changelog content
{
  if [ -n "$PREV_TAG" ]; then
    echo "## Changes since $PREV_TAG"
    git log --pretty=format:"* %s" $PREV_TAG..HEAD
  else
    echo "## Initial Release"
    git log --pretty=format:"* %s"
  fi
  echo
} 