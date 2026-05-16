#!/usr/bin/env bash
# KHR-AS: verify rdp-GW access graph continuity evidence bundle (read-only).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

RDP_GW=""
for candidate in \
  "${KHR_RDP_GW_PATH:-}" \
  "${ROOT}/../rdp-GW" \
  "${ROOT}/../../rdp-GW" \
  "/home/m.simoncini/rdp-GW"; do
  [[ -n "${candidate}" && -d "${candidate}" ]] || continue
  RDP_GW="$(cd "${candidate}" && pwd)"
  break
done
FAIL=0

log() { echo "[khr_access_graph_continuity_bundle_check] $*"; }

find_summary() {
  local base="$1"
  if [[ ! -d "${base}/docs/evidence/khr-accessgraph-continuity" ]]; then
    return 1
  fi
  local latest
  latest="$(find "${base}/docs/evidence/khr-accessgraph-continuity" -mindepth 1 -maxdepth 1 -type d 2>/dev/null | sort | tail -1)"
  [[ -n "${latest}" && -f "${latest}/summary.json" ]] || return 1
  echo "${latest}/summary.json"
}

SUMMARY=""
if [[ -z "${RDP_GW}" ]]; then
  log "WARN: rdp-GW repo not found; set KHR_RDP_GW_PATH"
else
if [[ -x "${RDP_GW}/scripts/khr_accessgraph_continuity_evidence_test.sh" ]]; then
  log "running rdp-GW fixture evidence script..."
  if ! (cd "${RDP_GW}" && ./scripts/khr_accessgraph_continuity_evidence_test.sh); then
    log "WARN: rdp-GW evidence script failed; checking existing summary"
    FAIL=1
  fi
fi
fi

SUMMARY=""
if [[ -n "${RDP_GW}" ]]; then
  SUMMARY="$(find_summary "${RDP_GW}" 2>/dev/null || true)"
fi
if [[ -z "${SUMMARY}" ]]; then
  RELAY="${ROOT}/docs/evidence/khr-accessgraph-continuity-relay"
  SUMMARY="$(find_summary "${RELAY}" 2>/dev/null || true)"
fi

if [[ -z "${SUMMARY}" ]]; then
  log "FAIL: no accessgraph continuity summary.json found"
  exit 1
fi

log "summary=${SUMMARY}"

python3 - "${SUMMARY}" <<'PY'
import json, sys
path = sys.argv[1]
with open(path) as f:
    s = json.load(f)
required = {
    "status": "PASS",
    "readOnly": True,
    "mutating": False,
    "noDisconnect": True,
    "noRevoke": True,
    "continuityObserved": True,
    "noSessionMutation": True,
    "productionReady": False,
}
errors = []
for k, want in required.items():
    if s.get(k) != want:
        errors.append(f"{k}={s.get(k)!r} want {want}")
if s.get("source") not in ("live", "fixture-readonly"):
    errors.append(f"source={s.get('source')!r}")
if errors:
    for e in errors:
        print("FAIL:", e, file=sys.stderr)
    sys.exit(1)
print("bundle check OK", s.get("source"), s.get("runId"))
PY

log "PASS"
exit "${FAIL}"
