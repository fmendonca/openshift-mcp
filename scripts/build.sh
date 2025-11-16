#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

BINARY_NAME="openshift-mcp"
OUTPUT_DIR="$ROOT_DIR/build"

echo "[build] building ${BINARY_NAME}..."
mkdir -p "${OUTPUT_DIR}"

GOFLAGS="${GOFLAGS:-}"
CGO_ENABLED=0 GOFLAGS="${GOFLAGS}" go build -o "${OUTPUT_DIR}/${BINARY_NAME}" ./cmd/server

echo "[build] done -> ${OUTPUT_DIR}/${BINARY_NAME}"
