#!/usr/bin/env bash
set -euo pipefail

services=("notification-service") # "next-service"

echo "Starting local build and test for all services..."

for service in "${services[@]}"; do
  echo "----------------------------------------------"
  echo "Processing ${service}..."
  
  pushd "${service}" > /dev/null
  
  echo "Building ${service}..."
  go build -o "${service}" ./cmd/main.go
  
  echo "Running tests for ${service}..."
  go test ./... -v
  
  popd > /dev/null
  
  echo "${service} built and tests passed successfully."
done

echo "----------------------------------------------"
echo "All services built and tested successfully."
