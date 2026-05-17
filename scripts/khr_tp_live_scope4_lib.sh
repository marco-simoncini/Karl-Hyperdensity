#!/usr/bin/env bash
# KHR-BE: shared guards for TP Live Scope-4 ResourceLease guarded apply.
set -euo pipefail

KHR_RUNTIME_CLUSTER_CONTEXT="${KHR_RUNTIME_CLUSTER_CONTEXT:-karl-metal-01@ovh}"
KHR_RUNTIME_SANDBOX_NS="${KHR_RUNTIME_SANDBOX_NS:-khr-runtime-sandbox}"
KHR_SCOPE4_APPLY_TIMEOUT_SEC="${KHR_SCOPE4_APPLY_TIMEOUT_SEC:-120}"

khr_scope4_root() {
  cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd
}

khr_scope4_apply_evidence_base() {
  echo "$(khr_scope4_root)/docs/evidence/khr-tp-live-scope4-guarded-apply"
}

khr_scope4_apply_run_id() {
  echo "${KHR_TP_LIVE_SCOPE4_APPLY_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
}

khr_scope4_apply_run_dir() {
  echo "$(khr_scope4_apply_evidence_base)/$(khr_scope4_apply_run_id)"
}

khr_scope4_log() {
  echo "[khr_tp_live_scope4] $*"
}

khr_scope4_assert_cluster_context() {
  local ctx current
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  current="$(kubectl config current-context 2>/dev/null || true)"
  if [[ "${current}" != "${ctx}" ]]; then
    echo "BLOCKED: current-context=${current:-<none>} required=${ctx}" >&2
    exit 2
  fi
  khr_scope4_log "cluster context OK: ${ctx}"
}

khr_scope4_assert_sandbox_namespace() {
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

khr_scope4_require_guarded_apply_confirmation() {
  if [[ "${KHR_TP_LIVE_SCOPE4_I_UNDERSTAND_GUARDED_APPLY:-}" != "true" ]]; then
    echo "BLOCKED: set KHR_TP_LIVE_SCOPE4_I_UNDERSTAND_GUARDED_APPLY=true" >&2
    exit 2
  fi
}

khr_scope4_assert_apply_flags_required() {
  if [[ "${KHR_TP_APPLY_RESOURCELEASE:-}" == "false" ]]; then
    echo "BLOCKED: guarded apply requires explicit apply-resourcelease=true" >&2
    exit 2
  fi
}

khr_scope4_production_snapshot() {
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

khr_scope4_sandbox_pod_restarts() {
  local ctx ns
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  ns="${KHR_RUNTIME_SANDBOX_NS}"
  kubectl --context "${ctx}" -n "${ns}" get pods -l app=khr-native-live-target \
    -o jsonpath='{range .items[*]}{.metadata.name}:{.status.containerStatuses[0].restartCount}{" "}{end}' 2>/dev/null || echo ""
}

khr_scope4_native_live_pod_uid() {
  local ctx ns
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  ns="${KHR_RUNTIME_SANDBOX_NS}"
  kubectl --context "${ctx}" -n "${ns}" get pods -l app=khr-native-live-target \
    -o jsonpath='{.items[0].metadata.uid}' 2>/dev/null || echo ""
}

khr_scope4_build_binaries() {
  local root
  root="$(khr_scope4_root)"
  mkdir -p "${root}/bin"
  (cd "${root}" && go build -o "${root}/bin/karl-host-runtime" ./cmd/karl-host-runtime)
  (cd "${root}" && go build -o "${root}/bin/khr-continuity-snapshot" ./cmd/khr-continuity-snapshot)
  (cd "${root}" && go build -o "${root}/bin/khr-continuity-proof" ./cmd/khr-continuity-proof)
  echo "${root}/bin/karl-host-runtime"
}
