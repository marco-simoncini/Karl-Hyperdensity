#!/usr/bin/env bash
# KHR-AS/AT: verify rdp-GW access graph continuity evidence bundle (read-only).
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

pick_summary() {
  local base="$1"
  local evidence_dir="${base}/docs/evidence/khr-accessgraph-continuity"
  [[ -d "${evidence_dir}" ]] || return 1
  python3 - "${evidence_dir}" <<'PY'
import json, os, sys
root = sys.argv[1]
candidates = []
for name in sorted(os.listdir(root)):
    path = os.path.join(root, name)
    summary = os.path.join(path, "summary.json")
    if not os.path.isdir(path) or not os.path.isfile(summary):
        continue
    try:
        with open(summary) as f:
            s = json.load(f)
    except (OSError, json.JSONDecodeError):
        continue
    if s.get("status") != "PASS":
        continue
    src = s.get("source", "")
    trust = s.get("trustLevel", s.get("trust_level", ""))
    # prefer live-readonly over fixture
    rank = 0
    if src == "live-readonly" or trust == "live-readonly":
        rank = 2
    elif src == "fixture-readonly" or trust == "fixture":
        rank = 1
    candidates.append((rank, name, summary, src, trust))
if not candidates:
    sys.exit(1)
candidates.sort(key=lambda x: (x[0], x[1]), reverse=True)
print(candidates[0][2])
PY
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

if [[ -n "${RDP_GW}" ]]; then
  SUMMARY="$(pick_summary "${RDP_GW}" 2>/dev/null || true)"
fi
if [[ -z "${SUMMARY}" ]]; then
  RELAY="${ROOT}/docs/evidence/khr-accessgraph-continuity-relay"
  SUMMARY="$(pick_summary "${RELAY}" 2>/dev/null || true)"
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
source = s.get("source")
trust = s.get("trustLevel", s.get("trust_level"))
allowed_source = ("live-readonly", "fixture-readonly")
allowed_trust = ("live-readonly", "fixture")
if source not in allowed_source:
    errors.append(f"source={source!r} want one of {allowed_source}")
if trust not in allowed_trust:
    errors.append(f"trustLevel={trust!r} want one of {allowed_trust}")
if source == "live-readonly" and trust != "live-readonly":
    errors.append(f"live source requires trustLevel=live-readonly (got {trust!r})")
if source == "fixture-readonly" and trust != "fixture":
    errors.append(f"fixture source requires trustLevel=fixture (got {trust!r})")
if errors:
    for e in errors:
        print("FAIL:", e, file=sys.stderr)
    sys.exit(1)
level = "live-readonly" if source == "live-readonly" else "fixture-readonly"
print(
    "bundle check OK",
    f"evidenceLevel={level}",
    f"trustLevel={trust}",
    f"runId={s.get('runId')}",
)
if source == "fixture-readonly":
    print(
        "NOTE: fixture-readonly evidence — live-readonly preferred when rdp-GW sandbox is reachable",
        file=sys.stderr,
    )
PY

log "PASS"
exit "${FAIL}"
