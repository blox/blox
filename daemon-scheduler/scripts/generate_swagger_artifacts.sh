#!/bin/bash
set -eu

# Normalize to working directory being build root (up one level from ./scripts)
ROOT=$( cd "$( dirname "${BASH_SOURCE[0]}"  )/.." && pwd  )
GENERATED_DIR="${ROOT}/generated/v1"
cd "${GENERATED_DIR}"
rm -rf ./models ./client
swagger generate client model -f swagger.json -A blox_daemon_scheduler
cd "${ROOT}"
