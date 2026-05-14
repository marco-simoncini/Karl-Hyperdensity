# Controller Runtime Reconciliation Smoke Tick v1

**Milestone:** `hyperdensity_controller_runtime_reconciliation_smoke_tick_v1`

## Purpose

Sprint 17 proved the controller is installed and runtime healthy.  
Sprint 18 adds a fail-closed **runtime reconciliation smoke tick** so the installed controller can prove a safe first market tick without broad production movement.

## Definitions

- **runtime reconciliation smoke tick:** controlled post-install reconciliation tick in bounded scope.
- **installed controller tick:** tick executed by the installed Hyperdensity controller pod.
- **leader-held tick:** tick observed only while leader election lease is held.
- **durable state read/write:** ConfigMap/CRD state access with optimistic lock and idempotency.
- **market snapshot collection:** fleet shell observation for donors/receivers/idle/pressure.
- **bounded top-K pair evaluation:** no full N×N pairing; evaluated pairs ≤ topKDonors × topKReceivers.
- **action slate generation:** bounded candidate actions without execution.
- **resource future generation:** bounded futures without execution.
- **execution safety boundary:** no production movement, no forbidden auto, no executor authority.
- **smoke tick health decision:** passed vs blocked vs failed vs not_run.
- **source-of-truth map:** authority split across Hyperdensity, Dashboard, FluidVirt, Inventory.
- **claim boundaries:** explicit allowed and forbidden claims in Sprint 18.

## Core Product Rule

Smoke tick pass is allowed only when all required conditions hold:

- `installedControllerSmokeTickEnabled=true`
- `smokeTickRequested=true`, `smokeTickExecuted=true`
- `leaderHeldTickObserved=true`
- `durableStateReadVerified=true`, `durableStateWriteVerified=true`
- `marketSnapshotCollected=true`, `indicesRefreshed=true`
- `boundedPairingVerified=true`, `noFullNxNPairing=true`
- `actionSlateGenerated=true`, `resourceFuturesGenerated=true`
- `tickMetricsEmitted=true`, `tickEventsEmitted=true`
- `dashboardProjectionUpdated=true`
- `productionMovementExecuted=false`, `broadProductionMutationExecuted=false`
- `generalProductionAutoAllowed=false`, `productionAutoWithPolicy=false`

## Allowed Sprint 18 claim

“KARL Hyperdensity can execute a first controlled runtime reconciliation smoke tick from the installed controller, proving leader-held state access, market snapshot generation, bounded action/future creation, metrics emission and Dashboard projection while keeping general production auto disabled.”
