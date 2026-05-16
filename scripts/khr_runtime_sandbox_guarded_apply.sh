#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/khr_runtime_sandbox_lib.sh"

NS="${KHR_RUNTIME_SANDBOX_NS}"
MANIFEST_DIR="$(khr_runtime_sandbox_manifest_dir)"
BIN="$(khr_runtime_build_binary)"
CFG_DEFAULT="${MANIFEST_DIR}/karl-host-runtime-config-default.yaml"
CFG_APPLY="${MANIFEST_DIR}/karl-host-runtime-config-apply.yaml"
LEASE="${MANIFEST_DIR}/resourcelease-dry-run.json"
PORT="${MANIFEST_DIR}/resourceport-dry-run.json"
SANDBOX_DIR="$(mktemp -d)"
BLOCKED_OUT="$(mktemp)"
APPLY_OUT="$(mktemp)"

khr_runtime_assert_cluster_context
khr_runtime_assert_namespace_arg "${NS}"
khr_runtime_assert_sandbox_namespace_labels
khr_runtime_block_production_namespace "${NS}"

khr_runtime_log "proving apply blocked with sandboxApplyEnabled=false"
"${BIN}" -mode=apply-lease -config="${CFG_DEFAULT}" -lease-input="${LEASE}" -resource-port-input="${PORT}" \
  -namespace="${NS}" -sandbox-dir="${SANDBOX_DIR}" > "${BLOCKED_OUT}"
if ! jq -e '.blocked == true and .applied == false' "${BLOCKED_OUT}" >/dev/null 2>&1; then
  echo "BLOCKED: expected default apply to be blocked" >&2
  cat "${BLOCKED_OUT}" >&2
  exit 1
fi
khr_runtime_artifact_json "guarded-apply-blocked-default.json" "${BLOCKED_OUT}"

khr_runtime_log "guarded apply with explicit sandboxApplyEnabled=true"
"${BIN}" -mode=apply-lease -config="${CFG_APPLY}" -lease-input="${LEASE}" -resource-port-input="${PORT}" \
  -namespace="${NS}" -sandbox-dir="${SANDBOX_DIR}" > "${APPLY_OUT}"
if ! jq -e '.applied == true and .blocked == false' "${APPLY_OUT}" >/dev/null 2>&1; then
  echo "BLOCKED: guarded sandbox apply failed" >&2
  cat "${APPLY_OUT}" >&2
  exit 1
fi
if [[ ! -f "${SANDBOX_DIR}/apply-marker.txt" ]]; then
  echo "BLOCKED: sandbox marker missing" >&2
  exit 1
fi

khr_runtime_artifact_json "guarded-apply-sandbox.json" "${APPLY_OUT}"
khr_runtime_artifact_text "apply-marker.txt" "$(cat "${SANDBOX_DIR}/apply-marker.txt")"
khr_runtime_artifact_text "guarded-apply.log" "default blocked; explicit apply wrote marker only"

echo "[khr_runtime_sandbox_guarded_apply] PASS"
