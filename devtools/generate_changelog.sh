#!/bin/bash
set -e

# This script generates a changelog from git history
# Usage: ./generate_changelog.sh [previous_tag] [path]
# If previous_tag is not provided, it will be automatically detected
# If path is provided, only commits affecting that path will be included

# Get the previous tag if not provided
if [ -z "$1" ]; then
  PREV_TAG=$(git describe --tags --abbrev=0 HEAD^ 2>/dev/null || echo "")
else
  PREV_TAG=$1
fi

# Get the path filter if provided
PATH_FILTER=$2

# Get the current tag and date
CURRENT_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "Initial Release")
RELEASE_DATE=$(date +"%Y-%m-%d")

# Generate changelog content
{
  if [ -n "$PREV_TAG" ]; then
    echo "## ${CURRENT_TAG} - ${RELEASE_DATE}"
    if [ -n "$PATH_FILTER" ]; then
      git log --pretty=format:"* %s" $PREV_TAG..HEAD -- "$PATH_FILTER"
    else
      git log --pretty=format:"* %s" $PREV_TAG..HEAD
    fi
  else
    echo "## ${CURRENT_TAG} - ${RELEASE_DATE}"
    if [ -n "$PATH_FILTER" ]; then
      git log --pretty=format:"* %s" -- "$PATH_FILTER"
    else
      git log --pretty=format:"* %s"
    fi
  fi
  echo
} 