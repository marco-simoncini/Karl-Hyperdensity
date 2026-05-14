# KARL Hyperdensity — Idle Time Compression + Fleet Value Coverage v1

**Contract ID:** `hyperdensity_idle_time_compression_fleet_value_coverage_v1`  
**Milestone:** `hyperdensity_idle_time_compression_fleet_value_coverage_v1`

## Product definition

> KARL Hyperdensity can measure fleet idle value, eligible idle value, moved idle value, unmoved idle value, idle compression rate and guarantee coverage across enrolled runtime shells, separating real production evidence from reference, synthetic and estimated opportunity.

## Key metrics

| Metric | Formula |
|--------|---------|
| Idle compression rate | movedIdleValue / eligibleIdleValue |
| Guarantee coverage % | guaranteedEligibleSavingsTotal / realizedMovedIdleValue |
| Fleet liquidity rate | eligibleIdleValue / totalIdleValue |
| Movement density | successfulMovementCount / enrolledShellCount |

## Sprint 11B invariants

- `universalGuaranteedSavingsAllowed=false`
- `estimatedIdleCountedAsMoved=false`
- `syntheticFleetCountedAsProduction=false`
- `referenceFleetCountedAsProduction=false`
- Production canary ≠ general production proof
