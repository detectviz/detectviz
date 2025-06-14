#!/bin/bash

# DetectViz Server Startup Script
# zh: DetectViz 伺服器啟動腳本

set -e

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "Starting DetectViz Server..."
echo "Project root: $PROJECT_ROOT"

# Change to project root
cd "$PROJECT_ROOT"

# Set environment variables
export DETECTVIZ_COMPOSITION="compositions/minimal-platform/composition.yaml"

# Build and run server
echo "Building server..."
go build -o bin/detectviz-server ./apps/server

echo "Starting server..."
./bin/detectviz-server
