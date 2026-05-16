#!/usr/bin/env bash
# KHR-Q: validate lane discovery package tests and example config.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

go test ./pkg/khr/lanediscovery/... -count=1

grep -q 'laneDiscoveryEnabled: true' examples/khr/runtime-sandbox/karl-host-runtime-config-lane-discovery.yaml
echo "[validate_lane_discovery] PASS"
