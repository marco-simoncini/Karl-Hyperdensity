#!/usr/bin/env bash
# KHR-BC: shared guards for TP Live Scope-3 ResourceLease manual dry-run.
set -euo pipefail

KHR_RUNTIME_CLUSTER_CONTEXT="${KHR_RUNTIME_CLUSTER_CONTEXT:-karl-metal-01@ovh}"
KHR_RUNTIME_SANDBOX_NS="${KHR_RUNTIME_SANDBOX_NS:-khr-runtime-sandbox}"
KHR_SCOPE3_DRYRUN_TIMEOUT_SEC="${KHR_SCOPE3_DRYRUN_TIMEOUT_SEC:-90}"

khr_scope3_root() {
  cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd
}

khr_scope3_dryrun_evidence_base() {
  echo "$(khr_scope3_root)/docs/evidence/khr-tp-live-scope3-dryrun"
}

khr_scope3_dryrun_run_id() {
  echo "${KHR_TP_LIVE_SCOPE3_DRYRUN_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
}

khr_scope3_dryrun_run_dir() {
  echo "$(khr_scope3_dryrun_evidence_base)/$(khr_scope3_dryrun_run_id)"
}

khr_scope3_log() {
  echo "[khr_tp_live_scope3] $*"
}

khr_scope3_assert_cluster_context() {
  local ctx current
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  current="$(kubectl config current-context 2>/dev/null || true)"
  if [[ "${current}" != "${ctx}" ]]; then
    echo "BLOCKED: current-context=${current:-<none>} required=${ctx}" >&2
    exit 2
  fi
  khr_scope3_log "cluster context OK: ${ctx}"
}

khr_scope3_assert_sandbox_namespace() {
  local ctx ns val
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  ns="${KHR_RUNTIME_SANDBOX_NS}"
  for blocked in karl karl-system kube-system default; do
    [[ "${ns}" == "${blocked}" ]] && { echo "BLOCKED: forbidden namespace ${ns}" >&2; exit 2; }
  done
  val="$(kubectl --context "${ctx}" get namespace "${ns}" \
    -o jsonpath='{.metadata.labels.khr\.karl\.io/sandbox}' 2>/dev/null || true)"
  [[ "${val}" == "true" ]] || { echo "BLOCKED: missing sandbox label on ${ns}" >&2; exit 2; }
}

khr_scope3_require_manual_dryrun_confirmation() {
  if [[ "${KHR_TP_LIVE_SCOPE3_I_UNDERSTAND_MANUAL_DRYRUN:-}" != "true" ]]; then
    echo "BLOCKED: set KHR_TP_LIVE_SCOPE3_I_UNDERSTAND_MANUAL_DRYRUN=true" >&2
    exit 2
  fi
}

khr_scope3_assert_apply_disabled() {
  local cfg="$1"
  if grep -qE 'sandboxApplyEnabled:\s*true' "${cfg}" 2>/dev/null; then
    echo "BLOCKED: sandboxApplyEnabled must be false" >&2
    exit 2
  fi
  if grep -qE 'resourcePortLoopEnabled:\s*true' "${cfg}" 2>/dev/null; then
    echo "BLOCKED: resourcePortLoopEnabled must be false for dry-run-only run" >&2
    exit 2
  fi
}

khr_scope3_assert_no_apply_resourcelease_flag() {
  if [[ "${KHR_TP_APPLY_RESOURCELEASE:-}" == "true" ]]; then
    echo "BLOCKED: KHR_TP_APPLY_RESOURCELEASE must not be true for dry-run" >&2
    exit 2
  fi
}

khr_scope3_production_snapshot() {
  local ctx out
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  out="$(mktemp)"
  python3 - "${ctx}" "${out}" <<'PY'
import json, subprocess, sys
from datetime import datetime, timezone
ctx, out = sys.argv[1:3]
gens, reps = {}, {}
for ns in ("karl", "karl-system", "default", "kube-system"):
    r = subprocess.run(
        ["kubectl", "--context", ctx, "get", "deploy", "-n", ns,
         "-o", "jsonpath={.items[*].metadata.generation}"],
        capture_output=True, text=True,
    )
    gens[ns] = (r.stdout or "").strip() if r.returncode == 0 else "unavailable"
for ns in ("khr-runtime-sandbox",):
    r = subprocess.run(
        ["kubectl", "--context", ctx, "get", "deploy", "-n", ns,
         "-o", "jsonpath={.items[*].metadata.generation}"],
        capture_output=True, text=True,
    )
    reps[ns] = (r.stdout or "").strip() if r.returncode == 0 else "unavailable"
open(out, "w").write(json.dumps({
    "clusterContext": ctx,
    "productionDeployGenerations": gens,
    "sandboxDeployGenerations": reps,
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}, indent=2) + "\n")
PY
  cat "${out}"
  rm -f "${out}"
}

khr_scope3_sandbox_pod_restarts() {
  local ctx ns
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  ns="${KHR_RUNTIME_SANDBOX_NS}"
  kubectl --context "${ctx}" -n "${ns}" get pods \
    -o jsonpath='{range .items[*]}{.metadata.name}:{.status.containerStatuses[0].restartCount}{" "}{end}' 2>/dev/null || echo ""
}

khr_scope3_build_binary() {
  local root bin
  root="$(khr_scope3_root)"
  bin="${root}/bin/karl-host-runtime"
  mkdir -p "${root}/bin"
  (cd "${root}" && go build -o "${bin}" ./cmd/karl-host-runtime)
  echo "${bin}"
}
