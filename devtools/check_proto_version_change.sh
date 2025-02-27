#!/bin/bash
set -e

# This script checks if the version in proto/ts/package.json has changed
# and outputs the result to GitHub Actions output file

# Check if the version line in package.json has changed
git diff HEAD^ HEAD -- proto/ts/package.json | grep '"version":' || echo "No version change detected"

# If the version has changed, set the output variable to true and capture the new version
if git diff HEAD^ HEAD -- proto/ts/package.json | grep '"version":'; then
  # Set output variable for use in the conditional for the publish job
  # This writes to the GitHub Actions output file which can be accessed by other steps/jobs
  echo "version_changed=true" >> $GITHUB_OUTPUT
  
  # Extract the new version number
  VERSION=$(grep '"version":' proto/ts/package.json | sed 's/.*"version": "\(.*\)",/\1/')
  
  # Set the version as an output variable
  echo "version=$VERSION" >> $GITHUB_OUTPUT
  
  # Log the new version for visibility in the workflow run
  echo "Version changed to: $VERSION"
else
  # If no version change, set output variable to false so publish job is skipped
  echo "version_changed=false" >> $GITHUB_OUTPUT
fi
