#!/usr/bin/env bash
# KHR-AK: validate canonical KHR contract manifest and TP anchors.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

FAIL=0
MANIFEST="docs/contracts/khr/khr-contract-manifest.yaml"
SELF="scripts/khr_contract_manifest_check.sh"

say() { echo "[khr_contract_manifest_check] $*"; }
fail() { say "FAIL: $*"; FAIL=1; }

require_file() {
  [[ -f "${ROOT}/$1" ]] || fail "missing: $1"
}

say "checking contract manifest..."
require_file "${MANIFEST}"

rg -q 'contractSetId:\s*khr-tp-contract-v1' "${ROOT}/${MANIFEST}" || fail "contractSetId khr-tp-contract-v1"
rg -q 'contractVersion:\s*"1.0.0"' "${ROOT}/${MANIFEST}" || fail "contractVersion 1.0.0"
rg -q 'productionReady:\s*false' "${ROOT}/${MANIFEST}" || fail "productionReady must be false"
rg -q 'autonomousOrchestration:\s*false' "${ROOT}/${MANIFEST}" || fail "autonomousOrchestration must be false"

for crd in hosts.runtime.karl.io resourceports.runtime.karl.io resourceleases.hyperdensity.karl.io; do
  rg -q "${crd}" "${ROOT}/${MANIFEST}" || fail "expected CRD ${crd} in manifest"
done

for doc in \
  docs/khr/RESOURCELEASE_TP_FREEZE_CANDIDATE.md \
  docs/khr/RESOURCEPORT_TP_FREEZE_CANDIDATE.md \
  docs/khr/NATIVE_LIVE_TP_FREEZE_CANDIDATE.md; do
  require_file "${doc}"
done

for crd in \
  api/crds/runtime.karl.io/host.yaml \
  api/crds/runtime.karl.io/resourceport.yaml \
  api/crds/hyperdensity.karl.io/resourcelease.yaml; do
  require_file "${crd}"
done

require_file "docs/evidence/khr-native-live-lane/certification-summary.json"

# Wording guard on manifest only
while IFS= read -r line; do
  [[ "${line}" == *"NOT production"* ]] && continue
  [[ "${line}" == *"productionReady: false"* ]] && continue
  [[ "${line}" == *"autonomousOrchestration: false"* ]] && continue
  if echo "${line}" | rg -qi 'generally available|production[- ]ready[^:]*true|autonomous orchestration.{0,20}enabled'; then
    fail "forbidden wording in manifest: ${line}"
  fi
done < "${ROOT}/${MANIFEST}"

# Optional cross-repo inventory schema
INV_SCHEMA="../Karl-Inventory/docs/contracts/khr/runtime-posture.schema.json"
if [[ -f "${INV_SCHEMA}" ]]; then
  say "OK: inventory runtime-posture schema present"
else
  say "INFO: optional inventory schema not found at ${INV_SCHEMA}"
fi

if [[ "${FAIL}" -ne 0 ]]; then
  exit 1
fi
say "PASS"
