#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
go test ./pkg/khr/resourceport/... -run 'Loop' -count=1
if [[ -f "${ROOT}/docs/evidence/khr-resourceport-loop/summary.json" ]]; then
  jq -e '.status == "PASS"' "${ROOT}/docs/evidence/khr-resourceport-loop/summary.json" >/dev/null
fi
echo "[validate_resourceport_loop] PASS"
