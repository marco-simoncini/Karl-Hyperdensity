# KARL Hyperdensity

[![CI](https://github.com/marco-simoncini/Karl-Hyperdensity/actions/workflows/ci.yml/badge.svg)](https://github.com/marco-simoncini/Karl-Hyperdensity/actions/workflows/ci.yml)

KARL Hyperdensity is the resource-control plane ("Grande Padre") for Linux shells in KARL.

> Phase 0 note: runtime implementation remains in `marco-simoncini/Karl-Dashboard` (branch `The-Father`). This repository is the canonical product/core home for contracts, schemas, blocker taxonomy, evidence model, overcommit model, and extraction planning.

## What Hyperdensity Is

Hyperdensity governs live, evidence-driven CPU and RAM control under strict non-disruptive guarantees.

Core shell kinds:
- VM Linux shell
- Container Linux shell

## Current Product Truth

Current implementation truth (owned by `Karl-Dashboard` runtime):
- `hyperdensity_linux_shell_compliance_v1`: operational
- `hyperdensity_resource_equilibrium_v1`: operational in `recommendation_only`
- `hyperdensity_fleet_equilibrium_onboarding_v1`: deployed and validated with mixed compliant/remediable/blocked/excluded fleet classes
- VM Linux reference: live CPU and RAM up/down, same runtime, no reboot/recreate/destructive migration path
- Container Linux reference: live CPU and RAM up/down through pod resize and cgroup evidence, no restart/rollout path
- Autonomous mode remains off
- Windows Hyperdensity lane is frozen/removed from active support claims

## Architecture Pillars

1. Shell compliance first (prove live-safe behavior before support claims)
2. Runtime overlay model (Declared vs Runtime vs Observed separated)
3. Evidence-driven guardrails (blockers, risk states, rollback proof)
4. Recommendation-first policy (`recommendation_only`)
5. Controlled execution (`operator_controlled`)
6. Safe overcommit via equilibrium and bounded risk

## Live Linux Shell Guarantees

Compliance requires proof of:
- CPU up live
- CPU down live
- RAM up live
- RAM down live
- no reboot
- no VMI recreate
- no rollout
- no destructive migration
- same-runtime continuity
- rollback capability

## Resource Equilibrium And Safe Overcommit

Equilibrium compares:
- declared envelope/baseline
- runtime target and runtime actual
- observed usage

It derives:
- floor, baseline, burst-step, ceiling
- reclaimable resources
- burst headroom
- donor/receiver candidates
- safe overcommit budget
- risk state and blockers

## Repository Ownership Map

- `Karl-Dashboard`: current runtime/API/UI owner (active implementation)
- `Karl-Hyperdensity`: future product/core contracts, policy logic, taxonomy, executors library
- `Karl-OS-ISO`: platform ownership (KubeVirt/Longhorn/Kube-OVN, runtime authority patch lane)
- `Karl-Migration-Factory`: import/remediation/source-attestation ownership
- Windows Hyperdensity track: separate/frozen for now (no active support claim)

## Extraction Roadmap

- Phase 0: repository creation + contracts/schemas/docs/examples (this milestone)
- Phase 1: extract pure contracts and blocker taxonomy packages
- Phase 2: extract pure compliance/equilibrium logic
- Phase 3: extract executor libraries
- Phase 4: `Karl-Dashboard` imports Hyperdensity module
- Phase 5: optional standalone Hyperdensity controller/service

Detailed plan: `docs/migration/dashboard-to-hyperdensity-extraction-plan.md`.

## Anti-goals (Current)

- no runtime move in Phase 0
- no autonomous mode
- no broad fleet support claim
- no Windows support claim
- no fake compliance
- no restart-bound support

## Validation

Run baseline validation locally before PRs:

- `./scripts/validate.sh`
- `go test ./...`

Validation includes:
- Go tests
- JSON parsing for `schemas/*.json` and `examples/*.json`
- schema metadata checks (`$schema`, `$id` or `title`, `type`)
- required docs structure checks
