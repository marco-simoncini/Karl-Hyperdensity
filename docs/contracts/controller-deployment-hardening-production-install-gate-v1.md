# Controller Deployment Hardening + Production Install Gate v1

**Milestone:** `hyperdensity_controller_deployment_hardening_production_install_gate_v1`

## Purpose

Sprint 14 delivered durable controller state and Kubernetes reconciler behavior with reference-only manifests. Sprint 15 introduces a **controlled production-install gate** with hardened manifests, namespace-scoped RBAC, probes, metrics, leader election wiring, upgrade/rollback safety, and fail-closed install decisions.

## Definitions

- **Controller deployment hardening:** deployment boundary checks for security context, probes, resources, service account, and forbidden auto flags.
- **Production install gate:** fail-closed contract that allows installation only when all mandatory checks pass and blockers are empty.
- **Install candidate:** manifest set and surface payload proposed for controlled production installation.
- **Install readiness:** aggregate of manifest lint, RBAC hardening, probes, metrics, leader election, durable state, upgrade and rollback checks.
- **Manifest hardening:** hardened, non-reference manifests under `deploy/kubernetes/controller/install-gate`.
- **RBAC hardening:** namespace-scoped Role/RoleBinding without cluster-admin, pods/exec, nodes write, wildcard verbs/resources, raw runtime controls.
- **Service account boundary:** dedicated service account with constrained role bindings and no elevated authority.
- **Deployment spec boundary:** non-privileged pod security context, probes configured, requests/limits configured, and forbidden auto flags disabled.
- **Probe policy:** liveness/readiness/startup probes required and validated.
- **Metrics endpoint:** explicit metrics service/port/path contract with forbidden-auto gauges fixed at zero.
- **Leader election wiring:** leases API and wiring enabled; HA production proof remains false unless evidence exists.
- **Durable state wiring:** Kubernetes ConfigMap/CRD state wiring with optimistic lock and migration/restore guardrails.
- **ConfigMap state mode:** `kubernetes_configmap` supported install mode.
- **CRD state mode:** `kubernetes_crd` supported where CRD is defined and wiring checks pass.
- **Upgrade safety:** bounded rollout strategy with backup/migration prerequisites and fail-closed readiness behavior.
- **Rollback safety:** rollback plan + previous image/config retention + fail-closed on rollback failure.
- **Install audit event:** immutable install-gate event log entries.
- **Install blocker:** structured blocker object that prevents install while unresolved.
- **Install warning:** non-blocking condition requiring operator attention.
- **Production install decision:** `install_allowed | install_blocked | reference_only | preview_only`.
- **Fail-closed install behavior:** install defaults to blocked when evidence is missing/unsafe.
- **Source-of-truth map:** authority map across Hyperdensity, FluidVirt, Dashboard, Inventory.
- **Claim boundaries:** explicit claims allowed/forbidden for Sprint 15.

## Core Product Rule

Controller installation is allowed only when:

- `productionInstallGateEnabled=true`
- `hardenedManifestsDefined=true`
- namespace-scoped RBAC passes;
- `clusterAdmin=false`;
- `podsExecAllowed=false`;
- `nodesWriteAllowed=false`;
- no raw runtime controls;
- no direct libvirt/cgroup access;
- probes exist;
- metrics endpoint exists;
- leader election is wired;
- durable state is wired;
- status/events permissions exist;
- rollback plan exists;
- upgrade strategy exists;
- install audit is enabled;
- blockers list is empty.

## Mandatory Sprint 15 invariants

- `installFailClosed=true`
- `generalProductionAutoAllowed=false`
- `productionAutoWithPolicy=false`
- `universalGuaranteedSavingsAllowed=false`
- `universalGuaranteedSavingsClaimed=false`
- `estimatedIdleCountedAsMoved=false`
- `projectedCompressionCountedAsRealized=false`
- `syntheticFleetCountedAsProduction=false`
- `referenceFleetCountedAsProduction=false`
- `dashboardExecutor=false`
- `fluidvirtPolicyAuthority=false`
- `inventoryRuntimeExecutor=false`

## Forbidden claims

- general production auto enabled;
- production auto with policy enabled;
- HA production proven (without dedicated HA evidence);
- large-fleet production proven;
- universal guaranteed savings;
- projected compression counted as realized;
- estimated idle counted as moved;
- Dashboard executor authority;
- FluidVirt policy/reconciler authority;
- Inventory runtime executor authority.

## Source of Truth

| Concern | Authority |
|---|---|
| Install gate model, validators, hardened manifests | Karl-Hyperdensity |
| Runtime actuator compatibility evidence | FluidVirt |
| Install-gate projection UI/API | Karl-Dashboard (read-only) |
| Identity and signals | Karl-Inventory |

## Allowed Sprint 15 claim

“KARL Hyperdensity can package and gate its durable resource-market controller for controlled Kubernetes installation with hardened RBAC, probes, metrics, leader-election wiring, upgrade/rollback safety and production-install blockers, while keeping general production auto disabled.”
