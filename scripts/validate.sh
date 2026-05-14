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
  "docs/contracts/shell-claim-template-profile-pack-v1.md"
  "docs/contracts/release-support-matrix-v1.md"
  "docs/contracts/evidence-bundle-demo-scenario-pack-v1.md"
  "docs/contracts/live-resource-authority-v1.md"
  "docs/contracts/action-slate-v1.md"
  "docs/contracts/guarded-auto-sandbox-v1.md"
  "docs/contracts/auto-rollback-controller-v1.md"
  "docs/contracts/blast-radius-policy-v1.md"
  "docs/contracts/production-kernel-boundary-v1.md"
  "docs/contracts/shell-passport-factory-v1.md"
  "docs/contracts/resource-lease-action-slate-readiness-v1.md"
  "docs/contracts/operator-controlled-apply-gate-v1.md"
  "docs/contracts/realized-savings-ledger-v1.md"
  "docs/contracts/universal-slo-guard-certified-uplift-v1.md"
  "docs/contracts/guarded-auto-policy-engine-v1.md"
  "docs/contracts/guarded-auto-apply-sandbox-nonprod-v1.md"
  "docs/contracts/production-canary-auto-apply-v1.md"
  "docs/contracts/dashboard-enterprise-cleanup-ga-release-gate-v1.md"
  "docs/contracts/guaranteed-eligible-savings-activation-v1.md"
  "docs/contracts/idle-time-compression-fleet-value-coverage-v1.md"
  "docs/contracts/continuous-resource-market-controller-v1.md"
  "docs/contracts/policy-pack-v1.md"
  "docs/contracts/policy-pack-consistency-checker-v1.md"
  "docs/contracts/admission-guard-enforce-simulation-v1.md"
  "docs/contracts/mutate-preview-apply-dry-run-v1.md"
  "docs/runbooks/operator-runbook-v1.md"
  "docs/releases/technical-preview-release-notes-v1.md"
  "docs/releases/technical-preview-readiness-gate-v1.md"
  "docs/releases/technical-preview-release-candidate-gate-v1.md"
  "docs/demos/technical-preview-demo-guide-v1.md"
  "docs/releases/technical-preview-documentation-pack-v1.md"
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
