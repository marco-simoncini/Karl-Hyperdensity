#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/khr_runtime_sandbox_lib.sh"

NS="${KHR_RUNTIME_SANDBOX_NS}"
MANIFEST_DIR="$(khr_runtime_sandbox_manifest_dir)"
BIN="$(khr_runtime_build_binary)"
CFG="${MANIFEST_DIR}/karl-host-runtime-config-default.yaml"
LEASE="${MANIFEST_DIR}/resourcelease-dry-run.json"
PORT="${MANIFEST_DIR}/resourceport-dry-run.json"
OUT="$(mktemp)"
PORT_OUT="$(mktemp)"

khr_runtime_assert_cluster_context
khr_runtime_assert_namespace_arg "${NS}"
khr_runtime_assert_sandbox_namespace_labels

"${BIN}" -mode=dry-run-lease -config="${CFG}" -lease-input="${LEASE}" -resource-port-input="${PORT}" > "${OUT}"
if ! jq -e '.allowed == true and .blocked != true' "${OUT}" >/dev/null 2>&1; then
  echo "BLOCKED: dry-run lease failed" >&2
  cat "${OUT}" >&2
  exit 1
fi

"${BIN}" -mode=emit-resourceport -config="${CFG}" \
  -shell-ref="${NS}/Shell/runtime-sandbox" \
  -cell-ref="${NS}/Cell/runtime-sandbox-linux" \
  -namespace="${NS}" -port-name=khr-runtime-sandbox-port > "${PORT_OUT}"

FR_OUT="$(mktemp)"
"${BIN}" -mode=register-host -config="${CFG}" > /dev/null
"${BIN}" -mode=report-capabilities -config="${CFG}" > /dev/null
"${BIN}" -mode=flight-recorder > "${FR_OUT}"

khr_runtime_artifact_json "dry-run-lease.json" "${OUT}"
khr_runtime_artifact_json "resourceport-candidate.json" "${PORT_OUT}"
khr_runtime_artifact_json "flight-recorder.json" "${FR_OUT}"
khr_runtime_artifact_text "dry-run.log" "dry-run PASS allowed=true blocked=false flight-recorder captured"

echo "[khr_runtime_sandbox_dry_run] PASS"
