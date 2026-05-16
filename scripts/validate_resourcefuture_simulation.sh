#!/usr/bin/env bash
# KHR-R: validate ResourceFuture simulation package and config.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"
go test ./pkg/khr/resourcefuture/... -count=1
grep -q 'resourceFutureSimulationEnabled: true' \
  examples/khr/runtime-sandbox/karl-host-runtime-config-resourcefuture-simulate.yaml
echo "[validate_resourcefuture_simulation] PASS"
