#!/usr/bin/env bash
# Static dependency guard for pkg/hyperdensity/parentfabric (Sprint 45–46).
# Recursively checks all subpackages (e.g. executiontypes/).
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PF_DIR="${ROOT_DIR}/pkg/hyperdensity/parentfabric"

if [[ ! -d "${PF_DIR}" ]]; then
  echo "[validate_parentfabric_pure_deps] ERROR: missing ${PF_DIR}" >&2
  exit 1
fi

mapfile -t gofiles < <(find "${PF_DIR}" -type f -name '*.go' | sort)
if [[ "${#gofiles[@]}" -eq 0 ]]; then
  echo "[validate_parentfabric_pure_deps] ERROR: no .go files under parentfabric" >&2
  exit 1
fi

failures=0
check_pattern() {
  local pattern="$1"
  local label="$2"
  local hits
  hits="$(grep -RIn --include='*.go' -F "${pattern}" "${PF_DIR}" 2>/dev/null || true)"
  if [[ -n "${hits}" ]]; then
    echo "[validate_parentfabric_pure_deps] FORBIDDEN ${label}:" >&2
    echo "${hits}" >&2
    failures=$((failures + 1))
  fi
}

check_pattern 'k8s.io/' 'k8s.io import'
check_pattern 'kubevirt.io/' 'kubevirt.io import'
check_pattern 'github.com/gorilla/' 'gorilla import'
check_pattern '"net/http"' 'net/http import'
check_pattern 'github.com/openshift/console' 'openshift console import'
check_pattern 'Karl-Dashboard' 'Karl-Dashboard string/import'
check_pattern 'client-go' 'client-go string/import'
check_pattern 'context.Context' 'context.Context reference'
check_pattern 'http.Request' 'http.Request reference'
check_pattern 'rest.Config' 'rest.Config reference'
check_pattern 'dynamic.Interface' 'dynamic.Interface reference'
check_pattern 'clientset' 'clientset reference'

if [[ ! -d "${PF_DIR}/executiontypes" ]]; then
  echo "[validate_parentfabric_pure_deps] ERROR: missing ${PF_DIR}/executiontypes (Sprint 46)" >&2
  exit 1
fi

if [[ ! -d "${PF_DIR}/workload" ]]; then
  echo "[validate_parentfabric_pure_deps] ERROR: missing ${PF_DIR}/workload (Sprint 48)" >&2
  exit 1
fi

if [[ "${failures}" -ne 0 ]]; then
  echo "[validate_parentfabric_pure_deps] FAIL: ${failures} forbidden pattern group(s)" >&2
  exit 1
fi

echo "[validate_parentfabric_pure_deps] PASS: checked ${#gofiles[@]} file(s) under pkg/hyperdensity/parentfabric (includes executiontypes, workload)"
