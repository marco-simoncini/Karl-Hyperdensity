# Controller Install Apply + Runtime Health Gate v1

**Milestone:** `hyperdensity_controller_install_apply_runtime_health_gate_v1`

## Purpose

Sprint 16 proved install eligibility through dry-run and cluster admission gates.  
Sprint 17 adds a fail-closed **controlled install apply + runtime health gate** so production install applied and runtime health can be claimed only with real cluster apply evidence and verified deployment health.

## Definitions

- **controlled controller install apply:** Kubernetes apply into a bounded namespace scope after Sprint 16 admission passes.
- **install apply request:** recorded apply intent with admission gate reference and apply mode.
- **install apply result:** recorded apply outcome with manifest counts and apply evidence.
- **applied manifest inventory:** inventory of resources actually applied.
- **runtime health gate:** decision layer combining post-apply health checks.
- **deployment availability:** Deployment observed generation and available/ready replica checks.
- **replica readiness:** desired vs ready/available replica alignment.
- **pod readiness:** pod selector readiness and restart stability.
- **probe verification:** liveness/readiness/startup probe health.
- **metrics endpoint verification:** metrics service reachability and forbidden auto gauge safety.
- **leader-election lease verification:** lease object observed with holder identity.
- **durable state read/write verification:** ConfigMap/CRD state store read/write/idempotency checks.
- **status condition verification:** controller status conditions observed.
- **Kubernetes event verification:** install and health audit events observed.
- **RBAC runtime verification:** runtime RBAC safety (no cluster-admin, no pods/exec).
- **rollback readiness verification:** rollback plan and dry-run readiness present.
- **degraded/fail-closed verification:** degraded mode and fail-closed safety remain active.
- **post-install audit event:** immutable install/health audit record.
- **production install health decision:** healthy vs blocked vs degraded vs not_applied.
- **source-of-truth map:** authority split across Hyperdensity, Dashboard, FluidVirt, Inventory.
- **claim boundaries:** explicit allowed and forbidden claims in Sprint 17.

## Core Product Rule

Runtime health gate pass is allowed only when all required conditions hold:

- Sprint 16 admission gate passed (`productionInstallEligible=true` at apply time)
- `installApplyMode=real_cluster_apply`
- `productionInstallApplied=true`
- `realApplyEvidencePresent=true`
- `deploymentAvailable=true`
- `podsReady=true`
- `probesHealthy=true`
- `metricsReachable=true`
- `leaderElectionLeaseObserved=true`
- `durableStateReadWriteVerified=true`
- `rollbackReady=true`
- `failClosedVerified=true`
- no blockers

And must also keep:

- `generalProductionAutoAllowed=false`
- `productionAutoWithPolicy=false`
- `dashboardExecutor=false`
- `fluidvirtPolicyAuthority=false`
- `fluidvirtInstallAuthority=false`
- `inventoryRuntimeExecutor=false`

## Source of Truth

| Concern | Authority |
|---|---|
| Install apply + runtime health model, validator, health decision | Karl-Hyperdensity |
| Runtime actuator compatibility evidence | FluidVirt |
| Read-only install apply/runtime health projection | Karl-Dashboard |
| Identity/signals | Karl-Inventory |

## Claim Boundaries

- Apply/health pass can prove controlled install applied and runtime healthy, not broad production mutation.
- `productionInstallApplied=true` requires `installApplyMode=real_cluster_apply` and real apply evidence.
- `runtimeHealthGatePassed=true` requires `productionInstallApplied=true`.
- No claim of general production auto, production_auto_with_policy, HA production proof, or large-fleet proof.
- Dashboard remains projection-only (non-executor).
- FluidVirt remains evidence provider only (no install/health authority).

## Allowed Sprint 17 claim

“KARL Hyperdensity can apply and verify its durable resource-market controller installation in a controlled Kubernetes scope, proving deployment health, probes, metrics, leader-election lease, durable state access and fail-closed safety while keeping general production auto disabled.”
