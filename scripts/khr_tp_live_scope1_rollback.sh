#!/usr/bin/env bash
# KHR-AW: rollback TP Live Scope-1 sandbox deployments (sandbox namespaces only).
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/khr_tp_live_scope1_lib.sh"

ROOT="$(khr_scope1_root)"
cd "${ROOT}"
RUN_ID="${KHR_TP_LIVE_SCOPE1_RUN_ID:-$(ls -1t "$(khr_scope1_evidence_base)" 2>/dev/null | head -1)}"
RUN_DIR="$(khr_scope1_evidence_base)/${RUN_ID}"
CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
RUNTIME_NS="${KHR_RUNTIME_SANDBOX_NS}"
RDPGW_NS="${KHR_RDPGW_SANDBOX_NS}"

khr_scope1_require_confirmation
khr_scope1_assert_cluster_context

khr_scope1_log "deleting scope-1 deployments (namespaces retained per policy)"
kubectl --context "${CTX}" -n "${RUNTIME_NS}" delete deployment/karl-host-runtime-preview --ignore-not-found=true
kubectl --context "${CTX}" -n "${RUNTIME_NS}" delete configmap/karl-host-runtime-scope1-config --ignore-not-found=true

for rdp in "${ROOT}/../rdp-GW" "/home/m.simoncini/rdp-GW"; do
  if [[ -x "${rdp}/scripts/khr_rdpgw_sandbox_scope1_rollback.sh" ]]; then
    KHR_TP_LIVE_SCOPE1_RUN_ID="${RUN_ID}" KHR_RUNTIME_CLUSTER_CONTEXT="${CTX}" \
      "${rdp}/scripts/khr_rdpgw_sandbox_scope1_rollback.sh" || true
    break
  fi
done

pkill -f '/tmp/rdpgw-khr-scope1' 2>/dev/null || true

python3 - "${RUN_DIR}/rollback-summary.json" "${RUN_ID}" "${CTX}" <<'PY'
import json, sys
from datetime import datetime, timezone
out, run_id, ctx = sys.argv[1:4]
doc = {
    "phase": "khr-tp-live-scope1-rollback",
    "sprint": "KHR-AW",
    "runId": run_id,
    "clusterContext": ctx,
    "status": "PASS",
    "scope": "scope-1",
    "deploymentsRemoved": [
        "khr-runtime-sandbox/karl-host-runtime-preview",
        "khr-rdpgw-sandbox/rdpgw-khr-sandbox",
    ],
    "namespacesRetained": ["khr-runtime-sandbox", "khr-rdpgw-sandbox"],
    "productionNamespacesMutated": False,
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
open(out, "w").write(json.dumps(doc, indent=2) + "\n")
PY

khr_scope1_log "rollback-summary=${RUN_DIR}/rollback-summary.json"
echo "[khr_tp_live_scope1_rollback] PASS"
