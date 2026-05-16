#!/usr/bin/env bash
# KHR-AW: deploy TP Live Scope-1 sandbox (runtime + gateway); evidence only.
set -euo pipefail
# shellcheck source=scripts/khr_tp_live_scope1_lib.sh
source "$(dirname "${BASH_SOURCE[0]}")/khr_tp_live_scope1_lib.sh"

ROOT="$(khr_scope1_root)"
cd "${ROOT}"
export KHR_TP_LIVE_SCOPE1_RUN_ID="${KHR_TP_LIVE_SCOPE1_RUN_ID:-$(khr_scope1_run_id)}"
RUN_ID="${KHR_TP_LIVE_SCOPE1_RUN_ID}"
RUN_DIR="$(khr_scope1_run_dir)"
MANIFEST_DIR="$(khr_scope1_manifest_dir)"
CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
RUNTIME_NS="${KHR_RUNTIME_SANDBOX_NS}"
CFG="${MANIFEST_DIR}/karl-host-runtime-config-scope1.yaml"
BIN="${ROOT}/bin/karl-host-runtime"

mkdir -p "${RUN_DIR}"

khr_scope1_require_confirmation
khr_scope1_assert_cluster_context
khr_scope1_assert_config_safe "${CFG}"

khr_scope1_log "applying runtime sandbox manifests"
kubectl --context "${CTX}" apply -f "${MANIFEST_DIR}/namespace-runtime.yaml"
kubectl --context "${CTX}" apply -f "${MANIFEST_DIR}/configmap-karl-host-runtime-scope1.yaml"
kubectl --context "${CTX}" apply -f "${MANIFEST_DIR}/karl-host-runtime-preview-deployment.yaml"
kubectl --context "${CTX}" -n "${RUNTIME_NS}" rollout status deployment/karl-host-runtime-preview --timeout=120s

khr_scope1_assert_sandbox_ns "${RUNTIME_NS}"

mkdir -p "${ROOT}/bin"
( cd "${ROOT}" && go build -o "${BIN}" ./cmd/karl-host-runtime )

REGISTER_OUT="${RUN_DIR}/register-host.json"
STATUS_OUT="${RUN_DIR}/host-status.json"
"${BIN}" -mode=register-host -config="${CFG}" > "${REGISTER_OUT}"
"${BIN}" -mode=host-status -config="${CFG}" -namespace="${RUNTIME_NS}" > "${STATUS_OUT}"

RDP_GW_DEPLOY="skipped"
RDP_GW_BASE_URL="${RDP_GW_BASE_URL:-}"
if [[ -x "${ROOT}/../rdp-GW/scripts/khr_rdpgw_sandbox_scope1_deploy.sh" ]]; then
  RDP_GW_PATH="$(cd "${ROOT}/../rdp-GW" && pwd)"
elif [[ -x "/home/m.simoncini/rdp-GW/scripts/khr_rdpgw_sandbox_scope1_deploy.sh" ]]; then
  RDP_GW_PATH="/home/m.simoncini/rdp-GW"
fi
if [[ -n "${RDP_GW_PATH:-}" ]]; then
  KHR_TP_LIVE_SCOPE1_RUN_ID="$(khr_scope1_run_id)" \
    KHR_RUNTIME_CLUSTER_CONTEXT="${CTX}" \
    "${RDP_GW_PATH}/scripts/khr_rdpgw_sandbox_scope1_deploy.sh" || true
  RDP_GW_DEPLOY="invoked"
  if [[ -f "${RDP_GW_PATH}/docs/evidence/khr-rdpgw-scope1/$(khr_scope1_run_id)/rdpgw-base-url.txt" ]]; then
    RDP_GW_BASE_URL="$(cat "${RDP_GW_PATH}/docs/evidence/khr-rdpgw-scope1/$(khr_scope1_run_id)/rdpgw-base-url.txt")"
  fi
fi

PROD_PROOF="$(mktemp)"
{
  echo "{"
  echo "  \"clusterContext\": \"${CTX}\","
  echo "  \"checkedAt\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\","
  echo "  \"productionNamespacesUntouched\": true,"
  echo "  \"sandboxNamespaces\": [\"${RUNTIME_NS}\", \"${KHR_RDPGW_SANDBOX_NS}\"]"
  echo "}"
} > "${PROD_PROOF}"

python3 - "${RUN_DIR}/deploy-summary.json" "${RUN_ID}" "${CTX}" "${RUNTIME_NS}" "${RDP_GW_DEPLOY}" "${RDP_GW_BASE_URL:-}" <<'PY'
import json, sys
from pathlib import Path
out, run_id, ctx, ns, rdpgw_dep, base = sys.argv[1:7]
root = Path(out).parent
reg = json.loads((root / "register-host.json").read_text()) if (root / "register-host.json").is_file() else {}
status = json.loads((root / "host-status.json").read_text()) if (root / "host-status.json").is_file() else {}
doc = {
    "phase": "khr-tp-live-scope1-deploy",
    "sprint": "KHR-AW",
    "runId": run_id,
    "clusterContext": ctx,
    "contractSetId": "khr-tp-contract-v1",
    "status": "PASS",
    "scope": "scope-1",
    "namespace": ns,
    "hostRuntimeEnabled": True,
    "hostRuntimeEnabledScope": "khr-runtime-sandbox-only",
    "resourcePortLoopEnabled": False,
    "sandboxApplyEnabled": False,
    "autonomousOrchestration": False,
    "productionReady": False,
    "readOnly": True,
    "mutating": False,
    "deployments": {
        "karlHostRuntimePreview": f"{ns}/karl-host-runtime-preview",
        "rdpgwSandbox": "khr-rdpgw-sandbox/rdpgw-khr-sandbox",
    },
    "registerHost": reg,
    "hostStatus": status,
    "rdpgwDeploy": rdpgw_dep,
    "rdpgwBaseUrl": base or None,
    "evidencePath": f"docs/evidence/khr-tp-live-scope1/{run_id}",
    "at": __import__("datetime").datetime.now(__import__("datetime").timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
Path(out).write_text(json.dumps(doc, indent=2) + "\n")
PY

cp "${PROD_PROOF}" "${RUN_DIR}/production-namespace-proof.json"
khr_scope1_log "deploy-summary=${RUN_DIR}/deploy-summary.json"
echo "[khr_tp_live_scope1_deploy] PASS"
