#!/bin/bash

set -euo pipefail

# Required environment variables check
check_required_vars() {
    local vars=("$@")
    for var in "${vars[@]}"; do
        if [ -z "${!var:-}" ]; then
            echo "Error: Required environment variable $var is not set"
            exit 1
        fi
    done
}

# Log function
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1"
} 