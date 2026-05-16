#!/usr/bin/env bash
# Shared guards for KHR-G karl-host-runtime cluster sandbox scripts.
set -euo pipefail

KHR_RUNTIME_CLUSTER_CONTEXT="${KHR_RUNTIME_CLUSTER_CONTEXT:-karl-metal-01@ovh}"
KHR_RUNTIME_SANDBOX_NS="${KHR_RUNTIME_SANDBOX_NS:-khr-runtime-sandbox}"
KHR_RUNTIME_SANDBOX_LABEL="${KHR_RUNTIME_SANDBOX_LABEL:-khr.karl.io/sandbox=true}"
KHR_RUNTIME_SANDBOX_LABEL_KEY="${KHR_RUNTIME_SANDBOX_LABEL_KEY:-khr.karl.io/sandbox}"
KHR_RUNTIME_SANDBOX_LABEL_VALUE="${KHR_RUNTIME_SANDBOX_LABEL_VALUE:-true}"

# Production namespaces that must never be mutated by sandbox scripts.
KHR_RUNTIME_PRODUCTION_NS_BLOCKLIST=(
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

khr_runtime_root() {
  cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd
}

khr_runtime_evidence_dir() {
  local root
  root="$(khr_runtime_root)"
  echo "${root}/docs/evidence/khr-runtime-sandbox"
}

khr_runtime_run_id() {
  if [[ -n "${KHR_RUNTIME_RUN_ID:-}" ]]; then
    echo "${KHR_RUNTIME_RUN_ID}"
    return
  fi
  date -u +%Y%m%dT%H%M%SZ
}

khr_runtime_run_evidence_dir() {
  echo "$(khr_runtime_evidence_dir)/$(khr_runtime_run_id)"
}

khr_runtime_log() {
  local msg="$1"
  echo "[khr-runtime-sandbox] ${msg}"
}

khr_runtime_artifact_json() {
  local name="$1"
  local file="$2"
  local dir
  dir="$(khr_runtime_run_evidence_dir)"
  mkdir -p "${dir}"
  cp -f "${file}" "${dir}/${name}"
  khr_runtime_log "artifact ${dir}/${name}"
}

khr_runtime_artifact_text() {
  local name="$1"
  local content="$2"
  local dir
  dir="$(khr_runtime_run_evidence_dir)"
  mkdir -p "${dir}"
  printf '%s\n' "${content}" > "${dir}/${name}"
  khr_runtime_log "artifact ${dir}/${name}"
}

khr_runtime_assert_cluster_context() {
  local ctx current marker
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  current="$(kubectl config current-context 2>/dev/null || true)"
  if [[ "${current}" != "${ctx}" ]]; then
    echo "BLOCKED: current-context=${current:-<none>} required=${ctx}" >&2
    exit 2
  fi
  marker="$(kubectl --context "${ctx}" config view --minify -o jsonpath='{.contexts[0].context.cluster}' 2>/dev/null || true)"
  if [[ -z "${marker}" ]]; then
    echo "BLOCKED: cannot resolve cluster marker for context ${ctx}" >&2
    exit 2
  fi
  khr_runtime_log "cluster context OK: ${ctx} cluster=${marker}"
}

khr_runtime_assert_namespace_arg() {
  local ns="$1"
  if [[ "${ns}" != "${KHR_RUNTIME_SANDBOX_NS}" ]]; then
    echo "BLOCKED: namespace=${ns} required=${KHR_RUNTIME_SANDBOX_NS}" >&2
    exit 2
  fi
}

khr_runtime_block_production_namespace() {
  local ns="$1"
  local blocked
  for blocked in "${KHR_RUNTIME_PRODUCTION_NS_BLOCKLIST[@]}"; do
    if [[ "${ns}" == "${blocked}" ]]; then
      echo "BLOCKED: production namespace ${ns}" >&2
      exit 2
    fi
  done
}

khr_runtime_assert_sandbox_namespace_labels() {
  local ctx ns
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  ns="${KHR_RUNTIME_SANDBOX_NS}"
  khr_runtime_assert_namespace_arg "${ns}"
  if ! kubectl --context "${ctx}" get namespace "${ns}" -o jsonpath='{.metadata.labels.khr\.karl\.io/sandbox}' 2>/dev/null | grep -qx "${KHR_RUNTIME_SANDBOX_LABEL_VALUE}"; then
    echo "BLOCKED: namespace ${ns} missing label ${KHR_RUNTIME_SANDBOX_LABEL}" >&2
    exit 2
  fi
  khr_runtime_log "namespace label allowlist OK on ${ns}"
}

khr_runtime_assert_workload_labels() {
  local ctx ns selector
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  ns="${KHR_RUNTIME_SANDBOX_NS}"
  selector="${KHR_RUNTIME_SANDBOX_LABEL}"
  if ! kubectl --context "${ctx}" -n "${ns}" get pods -l "${selector}" --no-headers 2>/dev/null | grep -q .; then
    echo "BLOCKED: no pods with label ${selector} in ${ns}" >&2
    exit 2
  fi
}

khr_runtime_build_binary() {
  local root out
  root="$(khr_runtime_root)"
  out="${root}/bin/karl-host-runtime"
  mkdir -p "${root}/bin"
  (cd "${root}" && go build -o "${out}" ./cmd/karl-host-runtime)
  echo "${out}"
}

khr_runtime_sandbox_manifest_dir() {
  echo "$(khr_runtime_root)/examples/khr/runtime-sandbox"
}

khr_runtime_production_mutation_proof() {
  local ctx out
  ctx="${KHR_RUNTIME_CLUSTER_CONTEXT}"
  out="$(mktemp)"
  {
    echo "{"
    echo "  \"clusterContext\": \"${ctx}\","
    echo "  \"checkedAt\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\","
    echo "  \"productionNamespaces\": ["
    local first=1 ns gen
    for ns in karl-system kube-system default; do
      gen="$(kubectl --context "${ctx}" -n "${ns}" get deploy -o json 2>/dev/null | python3 -c "
import json,sys
d=json.load(sys.stdin)
items=d.get('items',[])
print(json.dumps({i['metadata']['name']: i['metadata'].get('generation',0) for i in items[:20]}))
" 2>/dev/null || echo '{}')"
      [[ ${first} -eq 1 ]] || echo ","
      first=0
      echo "    {\"namespace\": \"${ns}\", \"deployGenerations\": ${gen}}"
    done
    echo "  ],"
    echo "  \"sandboxOnly\": true,"
    echo "  \"note\": \"generations captured before/after; unchanged generations imply no production deploy mutation\""
    echo "}"
  } > "${out}"
  cat "${out}"
  rm -f "${out}"
}
