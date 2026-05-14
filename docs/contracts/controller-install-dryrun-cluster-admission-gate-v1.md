# Controller Install Dry-Run + Cluster Admission Gate v1

**Milestone:** `hyperdensity_controller_install_dryrun_cluster_admission_gate_v1`

## Purpose

Sprint 15 hardened controller install manifests and production install gating.  
Sprint 16 adds a fail-closed **dry-run + cluster admission gate** so production install eligibility is granted only after modeled client/server dry-run and mandatory preflight checks pass.

## Definitions

- **production install dry-run:** modeled Kubernetes install simulation before real apply.
- **cluster admission gate:** decision layer combining dry-run and preflight checks.
- **namespace admission:** namespace existence and policy compatibility preflight.
- **RBAC admission:** namespace-scoped RBAC safety decision (no cluster-admin, no pods/exec, no nodes write).
- **manifest admission:** deployment/manifests policy and labels/env/probe safety validation.
- **CRD admission:** CRD presence/scope/schema/status-subresource preflight.
- **ConfigMap state admission:** durable state mode safety and lock/backup/restore readiness.
- **leader election admission:** leader-election wiring and lease RBAC preflight.
- **probe admission:** liveness/readiness/startup health preflight.
- **metrics admission:** metrics endpoint/service and forbidden auto gauge safety preflight.
- **rollback admission:** rollback-plan and rollback dry-run readiness preflight.
- **server-side dry-run result:** modeled API-server validation result.
- **client-side dry-run result:** modeled client rendering/validation result.
- **admission policy decision:** `admitted_for_install | blocked | dryrun_only | reference_only`.
- **dry-run audit event:** immutable dry-run/admission event record.
- **admission blocker:** structured fail-closed blocker preventing install eligibility.
- **admission warning:** non-blocking condition requiring operator attention.
- **production install eligibility:** install eligibility after successful dry-run + admission; does not imply apply.
- **fail-closed behavior:** missing evidence or failed preflight means blocked.
- **source-of-truth map:** authority split across Hyperdensity, Dashboard, FluidVirt, Inventory.
- **claim boundaries:** explicit allowed and forbidden claims in Sprint 16.

## Core Product Rule

Production install eligibility is allowed only when all required conditions hold:

- `productionInstallDryRunEnabled=true`
- `clientSideDryRunModeled=true`
- `serverSideDryRunModeled=true`
- `clusterAdmissionGateEnabled=true`
- `productionInstallDryRunPassed=true`
- `clusterAdmissionGatePassed=true`
- `namespacePreflightPassed=true`
- `rbacPreflightPassed=true`
- `manifestAdmissionPassed=true`
- `crdPreflightPassed=true`
- `durableStatePreflightPassed=true`
- `leaderElectionPreflightPassed=true`
- `probesPreflightPassed=true`
- `metricsPreflightPassed=true`
- `rollbackPreflightPassed=true`
- no blockers

And must also keep:

- `productionInstallApplied=false` unless real apply evidence exists
- `generalProductionAutoAllowed=false`
- `productionAutoWithPolicy=false`
- `dashboardExecutor=false`
- `fluidvirtPolicyAuthority=false`
- `fluidvirtAdmissionAuthority=false`
- `inventoryRuntimeExecutor=false`

## Source of Truth

| Concern | Authority |
|---|---|
| Install dry-run + admission model, validator, eligibility policy | Karl-Hyperdensity |
| Runtime actuator compatibility evidence | FluidVirt |
| Read-only dry-run/admission projection | Karl-Dashboard |
| Identity/signals | Karl-Inventory |

## Claim Boundaries

- Dry-run/admission pass can prove **eligibility**, not actual install apply.
- `productionInstallApplied=true` requires real apply evidence.
- No claim of general production auto, production_auto_with_policy, HA production proof, or large-fleet proof.
- Dashboard remains projection-only (non-executor).
- FluidVirt remains evidence provider only (no admission authority).

## Allowed Sprint 16 claim

“KARL Hyperdensity can validate its durable resource-market controller installation through Kubernetes dry-run, namespace/RBAC/CRD/state/probe/metrics/leader-election/rollback admission gates before allowing production install, while keeping general production auto disabled.”
