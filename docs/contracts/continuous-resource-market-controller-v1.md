# KARL Hyperdensity — Continuous Resource Market Controller v1

**Contract ID:** `hyperdensity_continuous_resource_market_controller_v1`  
**Milestone:** `hyperdensity_continuous_resource_market_controller_v1`

## Product definition

> KARL Hyperdensity can continuously maintain a bounded runtime resource market by indexing idle donors, pressure receivers, risk, rollback and SLO readiness, generating prioritized action slate and resource futures without full N×N pairing, while keeping general production auto disabled.

## Core rule

No full N×N pairing. Use shard scope, incremental indices, top-K donors/receivers, priority queue, and policy gates.

## Sprint 12 invariants

- `generalProductionAutoAllowed=false`
- `productionAutoWithPolicy=false`
- `noFullNxNPairing=true`
- Projected compression is labeled projected, not realized
