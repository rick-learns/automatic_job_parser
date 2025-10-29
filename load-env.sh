#!/bin/bash
# Helper script to load environment variables from .env file
# Usage: source load-env.sh or . load-env.sh

set -euo pipefail

if [ -f .env ]; then
    set -a
    . ./.env
    set +a
    echo "✅ Environment variables loaded from .env"
else
    echo "❌ Error: .env file not found"
    echo "Please create .env from .env.example:"
    echo "  cp .env.example .env"
    exit 1
fi

