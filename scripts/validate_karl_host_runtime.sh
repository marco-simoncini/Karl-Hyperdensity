#!/usr/bin/env bash
# Sprint KHR-F: build karl-host-runtime and run package tests.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

go build -o /tmp/karl-host-runtime ./cmd/karl-host-runtime
go test ./pkg/khr/host/... ./pkg/khr/resourceport/... ./pkg/khr/resourcelease/... ./pkg/khr/flightrecorder/... -count=1

echo "[validate_karl_host_runtime] PASS"
