#!/bin/sh

set -e

# --- Configuration ---
IMAGE_NAME="edsa:1.0-builder"
BINARY_PATH_IN_IMAGE="/app/main"
OUTPUT_BINARY="edsa_api"
OUTPUT_DIR="./bin"
LOG_FILE="${OUTPUT_DIR}/build_$(date +'%Y-%m-%d_%H-%M-%S').log"

# --- Build Process ---
mkdir -p $OUTPUT_DIR

{
    echo "🔨 Starting build process..."

    echo "   (1/3) Building builder image..."
    docker build -t $IMAGE_NAME .

    echo "   (2/3) Creating temporary container to copy binary..."
    CONTAINER_ID=$(docker create $IMAGE_NAME)

    echo "   (3/3) Copying binary and cleaning up container..."
    docker cp "${CONTAINER_ID}:${BINARY_PATH_IN_IMAGE}" "${OUTPUT_DIR}/${OUTPUT_BINARY}"
    docker rm -f $CONTAINER_ID > /dev/null

    echo "✅ Done! Binary '${OUTPUT_BINARY}' has been created."
} 2>&1 | tee "$LOG_FILE"