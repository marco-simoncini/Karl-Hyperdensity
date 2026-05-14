# KARL Hyperdensity — Dashboard Enterprise Cleanup + GA Release Gate v1

**Contract ID:** `hyperdensity_dashboard_enterprise_cleanup_ga_release_gate_v1`  
**Milestone:** `hyperdensity_dashboard_enterprise_cleanup_ga_release_gate_v1`  
**Release track:** `enterprise`

## Product definition

> KARL Hyperdensity exposes an enterprise-grade cockpit with GA/Preview/Lab separation and a strict release gate that blocks unproven claims, reference-only evidence, synthetic proof, unsafe auto modes and Dashboard runtime controls.

## Surface classification

| Track | Executive visibility | Production proof |
|-------|---------------------|------------------|
| GA | Allowed with evidence | Yes if gated |
| Preview | Labeled Preview | No |
| production_canary | Labeled Production Canary | Canary scope only |
| lab/debug/archived | Hidden by default | No |
| reference_only | Hidden | Never |
| synthetic_shadow | Lab/archive only | Never |

## Sprint 10 invariants

- `generalProductionAutoAllowed=false`
- `productionAutoWithPolicy=false`
- `guaranteedSavingsClaimed=false`
- `universalPerformanceImprovementClaimed=false`
- `referencePayloadCountedAsProduction=false`
- `syntheticProofCountedAsProduction=false`
- `dashboardRuntimeControlsExposed=false`
- `dashboardExecutor=false`
- `releaseDecision` is `canary_only` or `ga_blocked` unless all GA gates pass

## Forbidden claims

- general production auto / production_auto_with_policy
- guaranteed savings active / universal performance improvement
- reference ConfigMap as production proof / synthetic as production proof
- Dashboard executor / raw runtime controls
- FluidVirt policy authority / Inventory runtime executor
