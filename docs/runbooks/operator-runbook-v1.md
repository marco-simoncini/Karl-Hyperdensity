# Operator Runbook v1

Runbook ID: `hyperdensity_operator_runbook_v1`  
Release track: `technical_preview`

## Purpose

KARL Hyperdensity (Grande Padre) is a governed resource exchange/control plane for live CPU/RAM management across Hyperdensity-ready Linux shells.

In Technical Preview, operators can safely:
- inspect governance/safety/support surfaces
- generate and dry-run shell creation paths
- review remediation and apply-dry-run outputs
- review Resource Exchange plans and staged transfer evidence
- validate KARL Live Resource Authority support boundaries

In Technical Preview, operators must not:
- perform raw resource creation
- expose or execute raw runtime controls
- enable enforcement or autonomous apply
- mutate production workloads
- treat Windows as supported

## Core operating principles

- No raw resource creation.
- Only Hyperdensity-ready shell creation.
- Operator-controlled only.
- No autonomous apply.
- No enforcement.
- No production workload mutation.
- Evidence namespace only for proof paths.
- Windows is out-of-scope/frozen.

## Key surfaces

- **Policy Pack**: canonical safety/governance policy baseline.
- **Policy Pack Consistency Checker**: confirms policy invariants are consistent and not drifting.
- **Release Support Matrix**: official support boundary and approved claims.
- **Evidence Bundle / Demo Scenario Pack**: curated TP claim/evidence narrative and scenario walk.
- **KARL Live Resource Authority**: unified product-level live CPU/RAM control contract surface.
- **Shell Claim Template/Profile Pack**: approved shell profile templates.
- **Shell Claim Generator**: manifest generation and profile-aligned claim projection.
- **Shell Claim Dry-Run Create**: server-side create preflight/dry-run guard.
- **Shell Claim Evidence Create / History**: evidence-scoped create path and audit history.
- **Shell Factory**: official Hyperdensity-ready creation source surface.
- **Admission Guard**: classification and audit safety posture.
- **Mutate Preview**: non-mutating patch preview suggestions.
- **Admission Guard Enforce Simulation**: simulated enforcement impact without enabling enforcement.
- **Mutate Preview Apply Dry-Run**: server-side apply dry-run path with no-mutation verification.
- **Resource Exchange**: donor/receiver market plan and transfer recommendations.
- **Transfer Dry-Run**: transfer viability validation before staged apply.
- **Staged Transfer Apply**: operator-controlled staged transfer proof path.
- **Stage Apply History**: historical transfer/apply/rollback audit trace.

## Operational states

- `ready`: all required gates pass for current scoped action.
- `factory_warming_up`: shell factory still materializing readiness signals.
- `warming_up`: data/surface initialization in progress.
- `partial`: some sources present, but not enough for safe operational confidence.
- `blocked`: mandatory gate violation exists.
- `degraded`: surfaced with missing/invalid source dependencies.
- `evidence_only`: claim is bounded to evidence-backed/object-specific scope.
- `preview_only`: not GA-grade; TP controls and boundaries apply.
- `simulation_only`: non-enforcing simulation path.
- `dry_run_only`: no apply lane; validation-only path.
- `recommendation_only`: advisory planning without autonomous execution.
- `operator_controlled_only`: explicit operator action required.
- `out_of_scope`: not supported in current release boundary.

Readiness rule:
- `warming_up` is not ready.
- `partial` is not ready.
- `blocked` is not ready.

## Standard operator flows

### A) Validate release boundary

1. Inspect `hyperdensityReleaseSupportMatrix`.
2. Confirm support claim mode is evidence-bounded for TP.
3. Confirm `hyperdensityPolicyPackConsistency.consistent=true`.
4. Confirm profile pack alignment (`hyperdensity_shell_claim_templates_profile_pack_v1`).
5. Confirm:
   - `enforcementMode=disabled`
   - `autonomousApplyAllowed=false`
   - `productionMutationAllowed=false`

### B) Create Hyperdensity-ready shell manifest

1. Start from Template/Profile Pack.
2. Use Shell Claim Generator for profile-aligned manifest generation.
3. Run Dry-Run Create before any create path.
4. If using evidence-scoped create path, keep operator-controlled mode and cleanup plan.
5. Capture evidence create/history outputs.

### C) Investigate raw/non-conforming object

1. Review Admission Guard classification (`wouldReject` / rationale).
2. Review Mutate Preview patch hints.
3. Review Enforce Simulation impact (still simulation only).
4. Run Apply Dry-Run and review rejections/remediation.
5. Use remediation hints; do not apply via preview/dry-run lane.

### D) Validate Resource Exchange plan

1. Confirm donor liquidity and receiver demand are present.
2. Review transfer recommendation quality and scope.
3. Run transfer dry-run.
4. Check staged apply readiness and rollback source requirements.
5. Verify history/audit references for staged/chained apply and rollback.

### E) Validate live CPU/RAM claim

1. Confirm Linux container CPU/RAM lane within supported TP path.
2. Confirm Linux VM CPU lane is evidence-backed/object-specific.
3. Confirm Linux VM RAM wording is runtime overlay via virtio-mem/QMP/QOM requested-size.
4. Verify no rollout/restart/recreate expectations for proven paths.
5. Confirm warning-clean and verification checks are present.

## Troubleshooting

- **Policy Pack missing/inconsistent**: restore Policy Pack source, re-run consistency checker, block apply lanes.
- **Support Matrix degraded/blocked**: restore upstream surfaces and boundary links.
- **Profile Pack mismatch**: align profile pack ID and supported profiles with Shell Claim Generator.
- **Shell Claim dry-run rejected**: inspect rejection reason; fix profile/namespace/resource constraints.
- **Evidence create cleanup failed**: run explicit cleanup and verify no residual evidence objects.
- **Admission Guard discovery empty**: validate discovery inputs and namespace/object filters.
- **Mutate Preview patch unavailable**: use recreate-via-shell-claim path where indicated.
- **Apply Dry-Run rejected**: inspect safety gates and remediation hints; do not force.
- **Resource Exchange no liquidity**: wait for usage confidence or broaden eligible donor set within policy.
- **Sustained idle evidence missing**: gather additional telemetry history window.
- **Safe-band exceeded**: reduce target/delta to bounded lane.
- **Runtime value drift**: revalidate runtime lease/overlay and reconciliation state.
- **Warning events present**: clear warnings before claiming success.
- **No-mutation proof failed**: block progression and inspect dry-run hash/event evidence.
- **VM RAM overlay capability missing**: verify virtio-mem, guest/kernel capability, and QMP/QOM path availability.
- **QMP/QOM path unavailable**: treat as unsupported for that object; remain evidence_only.
- **Guest memory signal unavailable**: retain bounded claim; require alternate proof signals.
- **Auth/OIDC live GET failure**: refresh OIDC flow/cookie and verify ingress route.
- **Test gate failure**: stop release claim; fix tests/validation before proceeding.

## Safety checklist (pre-action)

- Evidence namespace only.
- Dry-run completed.
- Rollback source captured.
- No production namespace.
- No autonomous mode.
- No enforcement mode.
- No warning events.
- No rollout/restart/recreate expected in selected path.
- Cleanup plan present.
- Artifact capture planned.

## Artifact capture checklist

- `parent_fabric_raw.json`
- relevant surface extract(s)
- safety gates extract
- preservation extract
- `live_auth_validation.log`
- tests/typecheck/bridge logs
- build/deploy logs (if applicable)
- cleanup logs (if evidence objects are created)
