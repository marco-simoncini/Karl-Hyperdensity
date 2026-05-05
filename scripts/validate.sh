#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

required_files=(
  "README.md"
  "docs/architecture/hyperdensity-overview.md"
  "docs/architecture/runtime-overlay-model.md"
  "docs/architecture/control-room-ui-concept.md"
  "docs/contracts/linux-shell-compliance-v1.md"
  "docs/contracts/resource-equilibrium-v1.md"
  "docs/contracts/fleet-equilibrium-onboarding-v1.md"
  "docs/contracts/shell-factory-v1.md"
  "docs/contracts/shell-claim-v1.md"
  "docs/contracts/policy-pack-v1.md"
  "docs/contracts/policy-pack-consistency-checker-v1.md"
  "docs/contracts/admission-guard-enforce-simulation-v1.md"
  "docs/overcommit/resource-equilibrium-and-safe-overcommit.md"
  "docs/migration/dashboard-to-hyperdensity-extraction-plan.md"
)

for required in "${required_files[@]}"; do
  if [[ ! -f "$required" ]]; then
    echo "[validate] ERROR: missing required file: $required" >&2
    exit 1
  fi
done

go test ./...
python3 scripts/validate_json.py

schema_count="$(ls -1 schemas/*.json | wc -l | tr -d ' ')"
example_count="$(ls -1 examples/*.json | wc -l | tr -d ' ')"
doc_count="${#required_files[@]}"

echo "[validate] SUCCESS: go tests + JSON validation passed (schemas=${schema_count}, examples=${example_count}, required_docs=${doc_count})"
