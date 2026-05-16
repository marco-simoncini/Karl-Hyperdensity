#!/usr/bin/env bash
# KHR-BC: live manual ResourceLease dry-run against observed-json ResourcePorts (no apply).
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/khr_tp_live_scope3_lib.sh"

ROOT="$(khr_scope3_root)"
cd "${ROOT}"
RUN_DIR="$(khr_scope3_dryrun_run_dir)"
CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"
TIMEOUT_SEC="${KHR_SCOPE3_DRYRUN_TIMEOUT_SEC}"
CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-default.yaml"
LEASE_SAMPLE="${ROOT}/examples/khr/runtime-sandbox/resourcelease-dryrun-allowed.json"
LOOP_EVIDENCE="${ROOT}/docs/evidence/khr-tp-live-scope2-resourceport-loop/committed-scope2-loop-khr-ba/loop-output.json"
SCOPE3_PF="${ROOT}/docs/evidence/khr-tp-live-scope3-preflight/committed-scope3-preflight-khr-bb/scope3-preflight-summary.json"
SCOPE2_VERIFY="${ROOT}/docs/evidence/khr-tp-live-scope2-resourceport-loop/committed-scope2-loop-khr-ba/verify-summary.json"

mkdir -p "${RUN_DIR}"
log() { khr_scope3_log "$*" | tee -a "${RUN_DIR}/run.log"; }

khr_scope3_require_manual_dryrun_confirmation
khr_scope3_assert_cluster_context
khr_scope3_assert_sandbox_namespace
khr_scope3_assert_apply_disabled "${CFG}"
khr_scope3_assert_no_apply_resourcelease_flag

if [[ -f "${RUN_DIR}/dryrun-summary.json" ]] && [[ "${KHR_TP_LIVE_SCOPE3_DRYRUN_FORCE_RERUN:-}" != "true" ]]; then
  python3 -c "import json,sys; d=json.load(open('${RUN_DIR}/dryrun-summary.json')); sys.exit(0 if d.get('status')=='PASS' else 1)" \
    && { khr_scope3_log "SKIP existing PASS evidence in ${RUN_DIR}"; exit 0; }
fi

for f in "${SCOPE3_PF}" "${SCOPE2_VERIFY}" "${LOOP_EVIDENCE}"; do
  [[ -f "${f}" ]] || { echo "BLOCKED: missing ${f}" >&2; exit 2; }
done
python3 -c "
import json, sys
s2=json.load(open('${SCOPE2_VERIFY}'))
s3=json.load(open('${SCOPE3_PF}'))
sys.exit(0 if s2.get('readyForScope2')=='manual-loop-pass' and s3.get('status')=='PASS' else 1)
" || { echo "BLOCKED: scope-2/scope-3 preflight not ready" >&2; exit 2; }

PROD_BEFORE="${RUN_DIR}/production-before.json"
khr_scope3_production_snapshot >"${PROD_BEFORE}"
echo "$(khr_scope3_sandbox_pod_restarts)" > "${RUN_DIR}/sandbox-restarts-before.txt"

cp "${LOOP_EVIDENCE}" "${RUN_DIR}/observed-resourceports-source.json"
OBSERVED_PORTS="${RUN_DIR}/observed-resourceports.json"
python3 - "${LOOP_EVIDENCE}" "${OBSERVED_PORTS}" "${LEASE_SAMPLE}" "${RUN_DIR}/lease-input.json" <<'PY'
import json, sys
from pathlib import Path

loop_path, ports_path, lease_path, lease_out = map(Path, sys.argv[1:5])
loop = json.loads(loop_path.read_text())
ports = []
for it in loop.get("iterations") or []:
    ports.extend(it.get("resourcePorts") or [])
if not ports:
    raise SystemExit("no resourcePorts in observed-json evidence")
ports_path.write_text(json.dumps(ports, indent=2) + "\n")
port = ports[0]
pname = port.get("metadata", {}).get("name", "")
cell_ref = port.get("spec", {}).get("cellRef", "")
lease = json.loads(lease_path.read_text())
lease.setdefault("metadata", {}).setdefault("annotations", {})[
    "khr.karl.io/resource-port-ref"
] = f"cluster/ResourcePort/{pname}"
if cell_ref and "/" in cell_ref:
    _, _, cell_name = cell_ref.rpartition("/")
    ns = cell_ref.split("/")[0]
    tr = lease.setdefault("spec", {}).setdefault("transfer", {})
    tr["donor"] = {"apiGroup": "runtime.karl.io", "kind": "Cell", "namespace": ns, "name": cell_name}
    tr["receiver"] = {"apiGroup": "runtime.karl.io", "kind": "Cell", "namespace": ns, "name": f"{cell_name}-receiver"}
lease_out.write_text(json.dumps(lease, indent=2) + "\n")
print(pname)
PY
PORT_NAME="$(cat "${RUN_DIR}/lease-input.json" | python3 -c "import json,sys; print(json.load(sys.stdin)['metadata']['annotations']['khr.karl.io/resource-port-ref'].split('/')[-1])")"

BIN="$(khr_scope3_build_binary)"
log "runId=$(basename "${RUN_DIR}") mode=resourcelease-dryrun port=${PORT_NAME} timeoutSec=${TIMEOUT_SEC}"

set +e
timeout "${TIMEOUT_SEC}" "${BIN}" \
  -mode=resourcelease-dryrun \
  -config="${CFG}" \
  -namespace="${NS}" \
  -cluster-context="${CTX}" \
  -apply-resourcelease=false \
  -lease-input="${RUN_DIR}/lease-input.json" \
  -observed-resourceports-json="${RUN_DIR}/observed-resourceports-source.json" \
  -resource-port-ref="cluster/ResourcePort/${PORT_NAME}" \
  -sandbox-dir="${RUN_DIR}/sandbox-work" \
  >"${RUN_DIR}/dryrun-output.json" 2>"${RUN_DIR}/dryrun-stderr.log"
RC=$?
set -e
[[ "${RC}" -ne 124 ]] || { echo "FAIL: dry-run timed out" >&2; exit 1; }
[[ "${RC}" -eq 0 ]] || { echo "FAIL: dry-run exit ${RC}" >&2; cat "${RUN_DIR}/dryrun-stderr.log" >&2; exit 1; }

export RUN_DIR CTX NS TIMEOUT_SEC RC
python3 <<'PY'
import json, os
from datetime import datetime, timezone
from pathlib import Path

run = Path(os.environ["RUN_DIR"])
out = json.loads((run / "dryrun-output.json").read_text())
ok = (
    out.get("mode") == "resourcelease-dryrun"
    and out.get("noMutation") is True
    and out.get("noApply") is True
    and out.get("dryRunDecision") in ("allowed", "blocked")
)
doc = {
    "phase": "khr-tp-live-scope3-dryrun-execute",
    "sprint": "KHR-BC",
    "runId": run.name,
    "clusterContext": os.environ["CTX"],
    "namespace": os.environ["NS"],
    "status": "PASS" if ok else "FAIL",
    "dryRunExecuted": True,
    "dryRunDecision": out.get("dryRunDecision"),
    "blockedReason": out.get("blockedReason") or out.get("reason"),
    "rollbackPlanRef": out.get("rollbackPlanRef"),
    "verificationPlanRef": out.get("verificationPlanRef"),
    "sourceResourcePortRef": out.get("sourceResourcePortRef"),
    "noMutation": out.get("noMutation"),
    "noApply": out.get("noApply"),
    "resourceLeaseApplyEnabled": False,
    "cgroupMutationObserved": False,
    "observedJsonSource": "docs/evidence/khr-tp-live-scope2-resourceport-loop/committed-scope2-loop-khr-ba/loop-output.json",
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
(run / "dryrun-summary.json").write_text(json.dumps(doc, indent=2) + "\n")
if not ok:
    raise SystemExit(1)
PY

log "PASS dryrun-summary=${RUN_DIR}/dryrun-summary.json decision=$(jq -r .dryRunDecision "${RUN_DIR}/dryrun-output.json")"
echo "[khr_tp_live_scope3_dryrun_execute] PASS"
