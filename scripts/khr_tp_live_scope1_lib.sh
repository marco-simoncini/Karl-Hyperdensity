#!/usr/bin/env bash
# KHR-AW: shared guards for TP Live Scope-1 (sandbox only).
set -euo pipefail

KHR_RUNTIME_CLUSTER_CONTEXT="${KHR_RUNTIME_CLUSTER_CONTEXT:-karl-metal-01@ovh}"
KHR_RUNTIME_SANDBOX_NS="${KHR_RUNTIME_SANDBOX_NS:-khr-runtime-sandbox}"
KHR_RDPGW_SANDBOX_NS="${KHR_RDPGW_SANDBOX_NS:-khr-rdpgw-sandbox}"
KHR_RUNTIME_SANDBOX_LABEL="${KHR_RUNTIME_SANDBOX_LABEL:-khr.karl.io/sandbox=true}"
KHR_TP_LIVE_SCOPE1_MANIFEST_DIR="${KHR_TP_LIVE_SCOPE1_MANIFEST_DIR:-examples/khr/tp-live-scope1}"

KHR_RUNTIME_PRODUCTION_NS_BLOCKLIST=(
  karl
  karl-system
  kube-system
  kube-public
  kube-node-lease
  default
  ingress
  longhorn-system
  kubevirt
  cdi
)

khr_scope1_root() {
  cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd
}

khr_scope1_evidence_base() {
  echo "$(khr_scope1_root)/docs/evidence/khr-tp-live-scope1"
}

khr_scope1_run_id() {
  echo "${KHR_TP_LIVE_SCOPE1_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
}

khr_scope1_run_dir() {
  echo "$(khr_scope1_evidence_base)/$(khr_scope1_run_id)"
}

khr_scope1_log() {
  echo "[khr_tp_live_scope1] $*"
}

khr_scope1_assert_cluster_context() {
  local ctx current
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  current="$(kubectl config current-context 2>/dev/null || true)"
  if [[ "${current}" != "${ctx}" ]]; then
    echo "BLOCKED: current-context=${current:-<none>} required=${ctx}" >&2
    exit 2
  fi
  khr_scope1_log "cluster context OK: ${ctx}"
}

khr_scope1_block_namespace() {
  local ns="$1"
  local blocked
  for blocked in "${KHR_RUNTIME_PRODUCTION_NS_BLOCKLIST[@]}"; do
    if [[ "${ns}" == "${blocked}" ]]; then
      echo "BLOCKED: production/forbidden namespace ${ns}" >&2
      exit 2
    fi
  done
}

khr_scope1_assert_sandbox_ns() {
  local ns="$1"
  khr_scope1_block_namespace "${ns}"
  if [[ "${ns}" != "${KHR_RUNTIME_SANDBOX_NS}" && "${ns}" != "${KHR_RDPGW_SANDBOX_NS}" ]]; then
    echo "BLOCKED: namespace=${ns} not in sandbox allowlist" >&2
    exit 2
  fi
  local val
  val="$(kubectl --context "${KHR_RUNTIME_CLUSTER_CONTEXT}" get namespace "${ns}" \
    -o jsonpath='{.metadata.labels.khr\.karl\.io/sandbox}' 2>/dev/null || true)"
  if [[ "${val}" != "true" ]]; then
    echo "BLOCKED: namespace ${ns} missing label ${KHR_RUNTIME_SANDBOX_LABEL}" >&2
    exit 2
  fi
}

khr_scope1_assert_config_safe() {
  local cfg="$1"
  python3 - "${cfg}" <<'PY'
import sys, yaml
path = sys.argv[1]
with open(path) as f:
    doc = yaml.safe_load(f)
spec = doc.get("spec") or {}
errors = []
if spec.get("resourcePortLoopEnabled") is True:
    errors.append("resourcePortLoopEnabled must be false for scope-1")
if spec.get("sandboxApplyEnabled") is True:
    errors.append("sandboxApplyEnabled must be false for scope-1")
if spec.get("sandboxMode") is not True:
    errors.append("sandboxMode must be true")
if errors:
    for e in errors:
        print(f"BLOCKED: {e}", file=sys.stderr)
    sys.exit(2)
PY
}

khr_scope1_require_confirmation() {
  if [[ "${KHR_TP_LIVE_SCOPE1_I_UNDERSTAND_SANDBOX:-}" != "true" ]]; then
    echo "BLOCKED: set KHR_TP_LIVE_SCOPE1_I_UNDERSTAND_SANDBOX=true" >&2
    exit 2
  fi
}

khr_scope1_manifest_dir() {
  echo "$(khr_scope1_root)/${KHR_TP_LIVE_SCOPE1_MANIFEST_DIR}"
}
