#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "[fmt] gofmt ./cmd ./internal ./pkg"
gofmt -w ./cmd ./internal ./pkg || true

if command -v goimports >/dev/null 2>&1; then
  echo "[fmt] goimports ./cmd ./internal ./pkg"
  goimports -w ./cmd ./internal ./pkg || true
fi
