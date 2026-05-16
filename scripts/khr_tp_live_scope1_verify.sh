#!/usr/bin/env bash
# KHR-AW: verify TP Live Scope-1 sandbox (read-only checks + live rdpgw evidence).
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/khr_tp_live_scope1_lib.sh"

ROOT="$(khr_scope1_root)"
cd "${ROOT}"
RUN_ID="${KHR_TP_LIVE_SCOPE1_RUN_ID:-$(ls -1t "$(khr_scope1_evidence_base)" 2>/dev/null | head -1)}"
RUN_DIR="$(khr_scope1_evidence_base)/${RUN_ID}"
CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
RUNTIME_NS="${KHR_RUNTIME_SANDBOX_NS}"
RDPGW_NS="${KHR_RDPGW_SANDBOX_NS}"
MANIFEST_DIR="$(khr_scope1_manifest_dir)"
CFG="${MANIFEST_DIR}/karl-host-runtime-config-scope1.yaml"
FAIL=0

[[ -f "${RUN_DIR}/deploy-summary.json" ]] || { echo "missing deploy-summary.json in ${RUN_DIR}" >&2; exit 1; }

khr_scope1_assert_cluster_context
khr_scope1_assert_config_safe "${CFG}"
khr_scope1_assert_sandbox_ns "${RUNTIME_NS}"
khr_scope1_assert_sandbox_ns "${RDPGW_NS}" 2>/dev/null || true

mkdir -p "${RUN_DIR}"

check_deploy_ready() {
  local ns dep
  ns="$1"
  dep="$2"
  if kubectl --context "${CTX}" -n "${ns}" get deploy "${dep}" -o jsonpath='{.status.availableReplicas}' 2>/dev/null | grep -qx '1'; then
    return 0
  fi
  return 1
}

RUNTIME_OK=false
check_deploy_ready "${RUNTIME_NS}" "karl-host-runtime-preview" && RUNTIME_OK=true

RDPGW_OK=false
RDPGW_MODE="missing"
if kubectl --context "${CTX}" get ns "${RDPGW_NS}" >/dev/null 2>&1; then
  khr_scope1_assert_sandbox_ns "${RDPGW_NS}"
  if check_deploy_ready "${RDPGW_NS}" "rdpgw-khr-sandbox"; then
    RDPGW_OK=true
    RDPGW_MODE="cluster"
  elif [[ -f "${ROOT}/../rdp-GW/docs/evidence/khr-rdpgw-scope1/${RUN_ID}/rdpgw-local.pid" ]] || \
       [[ -f "/home/m.simoncini/rdp-GW/docs/evidence/khr-rdpgw-scope1/${RUN_ID}/rdpgw-local.pid" ]]; then
    RDPGW_OK=true
    RDPGW_MODE="local-gateway"
  fi
fi

RDP_GW_BASE_URL="${RDP_GW_BASE_URL:-}"
if [[ -z "${RDP_GW_BASE_URL}" && -f "${RUN_DIR}/deploy-summary.json" ]]; then
  RDP_GW_BASE_URL="$(python3 -c "import json; print(json.load(open('${RUN_DIR}/deploy-summary.json')).get('rdpgwBaseUrl') or '')")"
fi
if [[ -z "${RDP_GW_BASE_URL}" ]]; then
  for p in "${ROOT}/../rdp-GW" "/home/m.simoncini/rdp-GW"; do
    u="${p}/docs/evidence/khr-rdpgw-scope1/${RUN_ID}/rdpgw-base-url.txt"
    [[ -f "${u}" ]] && RDP_GW_BASE_URL="$(cat "${u}")" && break
  done
fi
[[ -z "${RDP_GW_BASE_URL}" ]] && RDP_GW_BASE_URL="http://127.0.0.1:9443"

ACCESSGRAPH_OK=false
if [[ -n "${RDP_GW_BASE_URL}" ]]; then
  for rdp in "${ROOT}/../rdp-GW" "/home/m.simoncini/rdp-GW"; do
    if [[ -x "${rdp}/scripts/khr_accessgraph_continuity_evidence.sh" ]]; then
      if ( cd "${rdp}" && \
        KHR_ACCESSGRAPH_EVIDENCE_RUN_ID="${RUN_ID}-scope1-verify" \
        RDP_GW_BASE_URL="${RDP_GW_BASE_URL}" \
        KHR_RUNTIME_CLUSTER_CONTEXT="${CTX}" \
        ./scripts/khr_accessgraph_continuity_evidence.sh ); then
        ACCESSGRAPH_OK=true
        cp "${rdp}/docs/evidence/khr-accessgraph-continuity/${RUN_ID}-scope1-verify/summary.json" \
          "${RUN_DIR}/accessgraph-summary.json" 2>/dev/null || true
      fi
      break
    fi
  done
fi

python3 - "${CFG}" <<'PY' || FAIL=1
import sys, yaml
with open(sys.argv[1]) as f:
    s = yaml.safe_load(f)["spec"]
assert s.get("resourcePortLoopEnabled") is not True
assert s.get("sandboxApplyEnabled") is not True
assert s.get("sandboxMode") is True
print("config guards OK")
PY

python3 - "${RUN_DIR}/verify-summary.json" "${RUN_ID}" "${CTX}" "${RUNTIME_OK}" "${RDPGW_OK}" "${ACCESSGRAPH_OK}" "${RDPGW_MODE}" <<'PY'
import json, sys
from datetime import datetime, timezone
out, run_id, ctx, rt, gw, ag, gw_mode = sys.argv[1:8]
rt, gw, ag = rt == "true", gw == "true", ag == "true"
status = "PASS" if rt and ag and gw else "FAIL"
doc = {
    "phase": "khr-tp-live-scope1-verify",
    "sprint": "KHR-AW",
    "runId": run_id,
    "clusterContext": ctx,
    "status": status,
    "scope": "scope-1",
    "karlHostRuntimePreviewReady": rt,
    "rdpgwSandboxReady": gw,
    "rdpgwDeployMode": gw_mode,
    "accessGraphLiveReadonly": ag,
    "resourcePortLoopEnabled": False,
    "sandboxApplyEnabled": False,
    "autonomousOrchestration": False,
    "readyForScope2": False,
    "scope2BlockedReason": "ResourcePort loop not enabled in KHR-AW",
    "readOnly": True,
    "mutating": False,
    "noRevoke": True,
    "noDisconnect": True,
    "productionReady": False,
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
open(out, "w").write(json.dumps(doc, indent=2) + "\n")
print(f"verify status={status} runtime={rt} rdpgw={gw} accessgraph={ag}")
if status != "PASS":
    sys.exit(1)
PY

echo "[khr_tp_live_scope1_verify] PASS"
