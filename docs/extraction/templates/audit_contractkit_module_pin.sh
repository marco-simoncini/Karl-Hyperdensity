#!/usr/bin/env bash
# =============================================================================
# TEMPLATE: audit_contractkit_module_pin.sh (Sprint 42)
# -----------------------------------------------------------------------------
# Copy into a consumer repo and wire from CI — NOT installed automatically.
#
# - Set EXPECTED_CONTRACTKIT_VERSION only when Hyperdensity publishes a new
#   nested-module tag / bumps ContractKitModuleVersion (explicit sprint).
# - Fails on Go pseudo-versions for CONTRACTKIT_MODULE.
# - Fails on superseded semvers listed in FORBIDDEN_VERSIONS.
# - Run with module root (directory containing go.mod) as current working dir.
# =============================================================================
set -euo pipefail

CONTRACTKIT_MODULE="github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit"
EXPECTED_CONTRACTKIT_VERSION="v0.1.9-khr-m1-m19"

FORBIDDEN_VERSIONS=(
  "v0.1.5-khr-m1-m16"
  "v0.1.7-khr-m1-m18"
)

if [[ ! -f go.mod ]]; then
  echo "audit_contractkit_module_pin (template): ERROR: go.mod missing in $(pwd)" >&2
  exit 1
fi

for bad in "${FORBIDDEN_VERSIONS[@]}"; do
  if grep -qF "${CONTRACTKIT_MODULE} ${bad}" go.mod; then
    echo "audit_contractkit_module_pin (template): ERROR: forbidden pin ${bad}" >&2
    exit 1
  fi
done

# Escape dots for CONTRACTKIT_MODULE in grep -E
req_lines="$(grep -E "^[[:blank:]]*github\\.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit[[:blank:]]+" go.mod || true)"
req_count="$(printf '%s\n' "${req_lines}" | sed '/^$/d' | wc -l | tr -d ' ')"
if [[ "${req_count}" != "1" ]]; then
  echo "audit_contractkit_module_pin (template): ERROR: expected exactly one require line, found ${req_count}" >&2
  printf '%s\n' "${req_lines}" >&2
  exit 1
fi

version="$(printf '%s' "${req_lines}" | awk '{print $2}' | sed 's|//.*$||' | tr -d '[:space:]')"
if [[ "${version}" != "${EXPECTED_CONTRACTKIT_VERSION}" ]]; then
  echo "audit_contractkit_module_pin (template): ERROR: expected ${EXPECTED_CONTRACTKIT_VERSION}, got ${version}" >&2
  exit 1
fi

if printf '%s' "${version}" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+-0\.[0-9]{14}-[0-9a-fA-F]+$'; then
  echo "audit_contractkit_module_pin (template): ERROR: pseudo-version not allowed: ${version}" >&2
  exit 1
fi
if printf '%s' "${version}" | grep -qE '^v0\.0\.0-[0-9]{8,14}-'; then
  echo "audit_contractkit_module_pin (template): ERROR: pseudo-version not allowed: ${version}" >&2
  exit 1
fi

echo "Contractkit module pin PASS (template)"
