#!/usr/bin/env bash
# KHR-BE: live manual ResourceLease guarded apply (cpu.max only; sandbox native-live).
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/khr_tp_live_scope4_lib.sh"

ROOT="$(khr_scope4_root)"
cd "${ROOT}"
RUN_DIR="$(khr_scope4_apply_run_dir)"
CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"
TIMEOUT_SEC="${KHR_SCOPE4_APPLY_TIMEOUT_SEC}"
CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-guarded-apply.yaml"
LEASE_SAMPLE="${ROOT}/examples/khr/runtime-sandbox/resourcelease-native-live-cpu.json"
WORKLOAD="${ROOT}/examples/khr/runtime-sandbox/native-live-workload.yaml"
LOOP_EVIDENCE="${ROOT}/docs/evidence/khr-tp-live-scope2-resourceport-loop/committed-scope2-loop-khr-ba/loop-output.json"
SCOPE3_VERIFY="${ROOT}/docs/evidence/khr-tp-live-scope3-dryrun/committed-scope3-dryrun-khr-bc/verify-summary.json"
SCOPE4_PF="${ROOT}/docs/evidence/khr-tp-live-scope4-preflight/committed-scope4-preflight-khr-bd/scope4-preflight-summary.json"
BASELINE_ID="scope4-$(basename "${RUN_DIR}")"

mkdir -p "${RUN_DIR}/sandbox-work"
log() { khr_scope4_log "$*" | tee -a "${RUN_DIR}/run.log"; }

khr_scope4_require_guarded_apply_confirmation
khr_scope4_assert_cluster_context
khr_scope4_assert_sandbox_namespace
export KHR_TP_APPLY_RESOURCELEASE=true
khr_scope4_assert_apply_flags_required

if [[ -f "${RUN_DIR}/apply-summary.json" ]] && [[ "${KHR_TP_LIVE_SCOPE4_APPLY_FORCE_RERUN:-}" != "true" ]]; then
  python3 -c "import json,sys; d=json.load(open('${RUN_DIR}/apply-summary.json')); sys.exit(0 if d.get('status')=='PASS' else 1)" \
    && { khr_scope4_log "SKIP existing PASS apply evidence in ${RUN_DIR}"; exit 0; }
fi

for f in "${SCOPE4_PF}" "${SCOPE3_VERIFY}" "${LOOP_EVIDENCE}"; do
  [[ -f "${f}" ]] || { echo "BLOCKED: missing ${f}" >&2; exit 2; }
done
python3 -c "
import json, sys
s3=json.load(open('${SCOPE3_VERIFY}'))
s4=json.load(open('${SCOPE4_PF}'))
sys.exit(0 if s3.get('readyForScope3')=='manual-dryrun-pass' and s4.get('status')=='PASS' else 1)
" || { echo "BLOCKED: scope-3/scope-4 preflight not ready" >&2; exit 2; }

# Native-live target must exist in sandbox
if ! kubectl --context "${CTX}" -n "${NS}" get deploy khr-native-live-target >/dev/null 2>&1; then
  log "deploying native-live workload"
  kubectl --context "${CTX}" apply -f "${WORKLOAD}" >>"${RUN_DIR}/workload-apply.log" 2>&1
  kubectl --context "${CTX}" -n "${NS}" rollout status deploy/khr-native-live-target --timeout=120s \
    >>"${RUN_DIR}/workload-apply.log" 2>&1 || true
fi
kubectl --context "${CTX}" -n "${NS}" get pods -l app=khr-native-live-target -o wide \
  >"${RUN_DIR}/native-live-pods.txt" 2>&1 || { echo "BLOCKED: native-live target missing" >&2; exit 2; }

PROD_BEFORE="${RUN_DIR}/production-before.json"
khr_scope4_production_snapshot >"${PROD_BEFORE}"
echo "$(khr_scope4_sandbox_pod_restarts)" > "${RUN_DIR}/sandbox-restarts-before.txt"
echo "$(khr_scope4_native_live_pod_uid)" > "${RUN_DIR}/native-live-pod-uid-before.txt"
echo "$(kubectl --context "${CTX}" -n "${NS}" get deploy khr-native-live-target -o jsonpath='{.metadata.generation}' 2>/dev/null || echo 0)" \
  > "${RUN_DIR}/deploy-generation-before.txt"

cp "${LOOP_EVIDENCE}" "${RUN_DIR}/observed-resourceports-source.json"
python3 - "${LOOP_EVIDENCE}" "${LEASE_SAMPLE}" "${RUN_DIR}/lease-input.json" <<'PY'
import json, sys
from pathlib import Path

loop_path, lease_path, lease_out = map(Path, sys.argv[1:4])
loop = json.loads(loop_path.read_text())
ports = []
for it in loop.get("iterations") or []:
    ports.extend(it.get("resourcePorts") or [])
if not ports:
    raise SystemExit("no resourcePorts in observed-json evidence")
port = ports[0]
labels = port.get("metadata", {}).get("labels") or {}
if labels.get("khr.karl.io/native-live") != "true" and labels.get("khr.karl.io/native_live") != "true":
    # observed-json uses khr.karl.io/native-live on pod; port may inherit via loop emission
    pass
pname = port.get("metadata", {}).get("name", "")
if "native-live" not in pname:
    raise SystemExit(f"port {pname} is not native-live target")
cell_ref = port.get("spec", {}).get("cellRef", "")
lease = json.loads(lease_path.read_text())
lease.setdefault("metadata", {}).setdefault("labels", {})["khr.karl.io/sandbox"] = "true"
lease["metadata"]["labels"]["khr.karl.io/native-live"] = "true"
lease.setdefault("metadata", {}).setdefault("annotations", {})[
    "khr.karl.io/resource-port-ref"
] = f"cluster/ResourcePort/{pname}"
if cell_ref and "/" in cell_ref:
    _, _, cell_name = cell_ref.rpartition("/")
    ns = cell_ref.split("/")[0]
    tr = lease.setdefault("spec", {}).setdefault("transfer", {})
    tr["donor"] = {"apiGroup": "runtime.karl.io", "kind": "Cell", "namespace": ns, "name": cell_name}
    tr["receiver"] = {"apiGroup": "runtime.karl.io", "kind": "Cell", "namespace": ns, "name": f"{cell_name}-receiver"}
gov = lease.setdefault("spec", {}).setdefault("governance", {})
gov["dryRunOnly"] = False
gov.setdefault("rollbackPlanRef", {"name": "khr-sandbox-rollback", "namespace": "khr-runtime-sandbox"})
lease_out.write_text(json.dumps(lease, indent=2) + "\n")
print(pname)
PY
PORT_NAME="$(python3 -c "import json; print(json.load(open('${RUN_DIR}/lease-input.json'))['metadata']['annotations']['khr.karl.io/resource-port-ref'].split('/')[-1])")"

BIN="$(khr_scope4_build_binaries)"
SNAP="${ROOT}/bin/khr-continuity-snapshot"
"${SNAP}" -cluster-context="${CTX}" -namespace="${NS}" -out="${RUN_DIR}/continuity-before.json"

log "runId=$(basename "${RUN_DIR}") mode=resourcelease-guarded-apply port=${PORT_NAME} baseline=${BASELINE_ID} timeoutSec=${TIMEOUT_SEC}"

set +e
timeout "${TIMEOUT_SEC}" "${BIN}" \
  -mode=resourcelease-guarded-apply \
  -config="${CFG}" \
  -namespace="${NS}" \
  -cluster-context="${CTX}" \
  -apply-resourcelease=true \
  -i-understand-this-is-sandbox \
  -lease-input="${RUN_DIR}/lease-input.json" \
  -observed-resourceports-json="${RUN_DIR}/observed-resourceports-source.json" \
  -resource-port-ref="cluster/ResourcePort/${PORT_NAME}" \
  -sandbox-dir="${RUN_DIR}/sandbox-work" \
  -baseline-id="${BASELINE_ID}" \
  >"${RUN_DIR}/apply-output.json" 2>"${RUN_DIR}/apply-stderr.log"
RC=$?
set -e
[[ "${RC}" -ne 124 ]] || { echo "FAIL: guarded apply timed out" >&2; exit 1; }
[[ "${RC}" -eq 0 ]] || { echo "FAIL: guarded apply exit ${RC}" >&2; cat "${RUN_DIR}/apply-stderr.log" >&2; exit 1; }

export RUN_DIR CTX NS BASELINE_ID RC
python3 <<'PY'
import json, os
from datetime import datetime, timezone
from pathlib import Path

run = Path(os.environ["RUN_DIR"])
out = json.loads((run / "apply-output.json").read_text())
ver = out.get("verification") or {}
ok = (
    out.get("mode") == "resourcelease-guarded-apply"
    and out.get("applied") is True
    and out.get("applyState") == "applied"
    and ver.get("state") == "pass"
    and out.get("dryRun", {}).get("dryRunDecision") == "allowed"
    and out.get("dryRun", {}).get("resource") == "cpu"
    and ver.get("noRestart") is True
    and ver.get("noRollout") is True
    and ver.get("noRecreate") is True
)
doc = {
    "phase": "khr-tp-live-scope4-guarded-apply-execute",
    "sprint": "KHR-BE",
    "runId": run.name,
    "clusterContext": os.environ["CTX"],
    "namespace": os.environ["NS"],
    "status": "PASS" if ok else "FAIL",
    "guardedApplyExecuted": bool(ok),
    "cgroupMutationObserved": bool(ok),
    "mutationScope": "cpu.max",
    "applyScope": "sandbox-only",
    "lane": "native-live",
    "baselineId": os.environ["BASELINE_ID"],
    "cgroupPath": out.get("cgroupPath"),
    "observedCpuMax": ver.get("observedCpuMax"),
    "expectedCpuMax": ver.get("expectedCpuMax"),
    "dryRunDecision": out.get("dryRun", {}).get("dryRunDecision"),
    "rollbackPlanRef": out.get("dryRun", {}).get("rollbackPlanRef"),
    "verificationPlanRef": out.get("dryRun", {}).get("verificationPlanRef"),
    "sourceResourcePortRef": out.get("dryRun", {}).get("sourceResourcePortRef"),
    "noAutonomousOrchestration": True,
    "noPersistentRuntimeLoop": True,
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
(run / "apply-summary.json").write_text(json.dumps(doc, indent=2) + "\n")
if not ok:
    raise SystemExit(1)
PY

log "PASS apply-summary=${RUN_DIR}/apply-summary.json cgroup=$(jq -r .cgroupPath "${RUN_DIR}/apply-output.json")"
echo "[khr_tp_live_scope4_guarded_apply_execute] PASS"
