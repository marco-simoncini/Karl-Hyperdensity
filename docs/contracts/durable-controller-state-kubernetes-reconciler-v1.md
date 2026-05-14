# Durable Controller State + Kubernetes Reconciler v1

**Milestone:** `hyperdensity_durable_controller_state_kubernetes_reconciler_v1`

## Purpose

Sprint 14 moves the Sprint 13 live reconciliation loop from in-memory reference state to **durable Kubernetes-backed state** with a minimal reconciler model, status conditions, events, leader-election readiness, RBAC manifests, metrics, and recovery semantics.

## Implemented vs reference-defined

| Component | Status |
|---|---|
| ConfigMap-backed state store | **implemented** (fake Kubernetes client tests) |
| CRD-backed state | **reference_defined** (manifest + schema only) |
| Kubernetes manifests | **reference_only** (`not_production_install=true`) |
| Leader election | **configured** (not HA production proven) |

## Core product rule

Sprint 13 proved live reconciliation with `in_memory_reference`. Sprint 14 makes controller state **durable** and Kubernetes-reconciler ready while keeping general production auto disabled.

## Permitted store types

- `kubernetes_configmap` (implemented)
- `kubernetes_crd` (reference_defined)

## Sprint 14 allowed flags

- `durableStateStoreEnabled=true`
- `kubernetesReconcilerEnabled=true`
- `configMapBackedStateEnabled=true`
- `crdBackedStateDefined=true`
- `fakeClientTestsEnabled=true`
- `controllerStatusConditionsEnabled=true`
- `kubernetesEventsEnabled=true`
- `leaderElectionReady=true`
- `rbacManifestsDefined=true`
- `metricsExportDefined=true`
- `recoverySemanticsDefined=true`

## Sprint 14 forbidden flags

- `generalProductionAutoAllowed=false`
- `productionAutoWithPolicy=false`
- `projectedCompressionCountedAsRealized=false`
- `estimatedIdleCountedAsMoved=false`
- `dashboardExecutor=false`
- `fluidvirtPolicyAuthority=false`

## Source of truth

| Concern | Authority |
|---|---|
| Durable state / reconciler | Karl-Hyperdensity |
| Runtime actuator evidence | FluidVirt |
| Operator projection | Karl-Dashboard (read-only) |

## Claim boundaries

- Stale `resourceVersion` writes are rejected (optimistic lock)
- Idempotency records persist across reconciles
- Degraded mode is fail-closed with execution selection disabled
- Manifests are reference only unless explicitly marked production-ready
- Leader election configured does not prove HA production
