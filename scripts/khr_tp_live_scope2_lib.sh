#!/usr/bin/env bash
# KHR-BA: shared guards for TP Live Scope-2 ResourcePort manual loop.
set -euo pipefail

KHR_RUNTIME_CLUSTER_CONTEXT="${KHR_RUNTIME_CLUSTER_CONTEXT:-karl-metal-01@ovh}"
KHR_RUNTIME_SANDBOX_NS="${KHR_RUNTIME_SANDBOX_NS:-khr-runtime-sandbox}"
KHR_SCOPE2_LOOP_ITERATIONS="${KHR_SCOPE2_LOOP_ITERATIONS:-2}"
KHR_SCOPE2_LOOP_TIMEOUT_SEC="${KHR_SCOPE2_LOOP_TIMEOUT_SEC:-120}"
KHR_SCOPE2_LOOP_INTERVAL_MS="${KHR_SCOPE2_LOOP_INTERVAL_MS:-500}"

khr_scope2_root() {
  cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd
}

khr_scope2_evidence_base() {
  echo "$(khr_scope2_root)/docs/evidence/khr-tp-live-scope2-resourceport-loop"
}

khr_scope2_run_id() {
  echo "${KHR_TP_LIVE_SCOPE2_LOOP_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
}

khr_scope2_run_dir() {
  echo "$(khr_scope2_evidence_base)/$(khr_scope2_run_id)"
}

khr_scope2_log() {
  echo "[khr_tp_live_scope2] $*"
}

khr_scope2_assert_cluster_context() {
  local ctx current
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  current="$(kubectl config current-context 2>/dev/null || true)"
  if [[ "${current}" != "${ctx}" ]]; then
    echo "BLOCKED: current-context=${current:-<none>} required=${ctx}" >&2
    exit 2
  fi
  khr_scope2_log "cluster context OK: ${ctx}"
}

khr_scope2_assert_sandbox_namespace() {
  local ctx ns val
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  ns="${KHR_RUNTIME_SANDBOX_NS}"
  for blocked in karl karl-system kube-system default kube-public; do
    if [[ "${ns}" == "${blocked}" ]]; then
      echo "BLOCKED: forbidden namespace ${ns}" >&2
      exit 2
    fi
  done
  val="$(kubectl --context "${ctx}" get namespace "${ns}" \
    -o jsonpath='{.metadata.labels.khr\.karl\.io/sandbox}' 2>/dev/null || true)"
  if [[ "${val}" != "true" ]]; then
    echo "BLOCKED: namespace ${ns} missing label khr.karl.io/sandbox=true" >&2
    exit 2
  fi
}

khr_scope2_require_manual_loop_confirmation() {
  if [[ "${KHR_TP_LIVE_SCOPE2_I_UNDERSTAND_MANUAL_LOOP:-}" != "true" ]]; then
    echo "BLOCKED: set KHR_TP_LIVE_SCOPE2_I_UNDERSTAND_MANUAL_LOOP=true" >&2
    exit 2
  fi
}

khr_scope2_assert_config_safe() {
  local cfg="$1"
  python3 - "${cfg}" <<'PY'
import sys
from pathlib import Path
try:
    import yaml
except ImportError:
    yaml = None
path = Path(sys.argv[1])
text = path.read_text()
errors = []
if "sandboxApplyEnabled: true" in text or "sandboxApplyEnabled: true" in text.replace(" ", ""):
    errors.append("sandboxApplyEnabled must be false")
if "resourcePortLoopEnabled: true" not in text:
    errors.append("run config must enable resourcePortLoopEnabled for manual process only")
if yaml:
    doc = yaml.safe_load(text)
    spec = (doc or {}).get("spec") or {}
    if spec.get("sandboxApplyEnabled") is True:
        errors.append("sandboxApplyEnabled must be false (no ResourceLease apply path)")
if errors:
    for e in errors:
        print(f"BLOCKED: {e}", file=sys.stderr)
    sys.exit(2)
PY
}

khr_scope2_assert_no_resourcelease_path() {
  local cfg="$1"
  if grep -qE 'resourcelease|ResourceLease|sandboxApplyEnabled:\s*true' "${cfg}" 2>/dev/null; then
    if grep -q 'sandboxApplyEnabled: true' "${cfg}"; then
      echo "BLOCKED: ResourceLease apply path must not be enabled" >&2
      exit 2
    fi
  fi
}

khr_scope2_production_snapshot() {
  local ctx out
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  out="$(mktemp)"
  python3 - "${ctx}" "${out}" <<'PY'
import json, subprocess, sys
from datetime import datetime, timezone
ctx, out = sys.argv[1:3]
gens = {}
for ns in ("karl", "karl-system", "default", "kube-system"):
    r = subprocess.run(
        ["kubectl", "--context", ctx, "get", "deploy", "-n", ns,
         "-o", "jsonpath={.items[*].metadata.generation}"],
        capture_output=True, text=True,
    )
    gens[ns] = (r.stdout or "").strip() if r.returncode == 0 else "unavailable"
doc = {"clusterContext": ctx, "deployGenerations": gens,
       "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")}
open(out, "w").write(json.dumps(doc) + "\n")
PY
  cat "${out}"
  rm -f "${out}"
}

khr_scope2_assert_production_untouched() {
  local before="${1:-}" after="${2:-}"
  [[ -f "${before}" && -f "${after}" ]] || return 0
  python3 - "${before}" "${after}" <<'PY'
import json, sys
b, a = [json.load(open(p)) for p in sys.argv[1:3]]
if b.get("deployGenerations") != a.get("deployGenerations"):
    print("BLOCKED: production deploy generations changed", file=sys.stderr)
    sys.exit(2)
PY
  khr_scope2_log "production gateway namespaces unchanged"
}

khr_scope2_cluster_loop_disabled() {
  local ctx ns
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  ns="${KHR_RUNTIME_SANDBOX_NS}"
  local cm
  cm="$(kubectl --context "${ctx}" get configmap -n "${ns}" \
    -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' 2>/dev/null \
    | grep -E 'karl-host-runtime|scope1' | head -1 || true)"
  if [[ -z "${cm}" ]]; then
    echo "false"
    return 0
  fi
  if kubectl --context "${ctx}" -n "${ns}" get configmap "${cm}" \
    -o yaml 2>/dev/null | grep -q 'resourcePortLoopEnabled: true'; then
    echo "true"
  else
    echo "false"
  fi
}

khr_scope2_build_binary() {
  local root bin
  root="$(khr_scope2_root)"
  bin="${root}/bin/karl-host-runtime"
  mkdir -p "${root}/bin"
  (cd "${root}" && go build -o "${bin}" ./cmd/karl-host-runtime)
  echo "${bin}"
}
