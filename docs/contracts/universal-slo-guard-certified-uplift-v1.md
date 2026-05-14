# KARL Hyperdensity — Universal SLO Guard + Certified Performance Uplift v1

**Contract ID:** `hyperdensity_universal_slo_guard_certified_uplift_v1`  
**Milestone:** `hyperdensity_universal_slo_guard_certified_uplift_v1`  
**Release track:** `technical_preview`

## Product definition

> KARL Hyperdensity can universally protect enrolled workload SLOs and certify performance uplift only where CPU/RAM pressure is proven, while blocking regressions and keeping auto-apply disabled.

## Strict product distinctions

| Concept | Allowed | Forbidden |
|---------|---------|-----------|
| Universal SLO protection | Enrolled workloads with SLO profile + guard evaluation | Without SLO profile |
| Certified performance uplift | CPU/RAM-bound bottleneck proven; baseline + post-mutation; no regression | Universal performance improvement |
| Neutral / no-claim | Non-CPU/RAM bottleneck or insufficient evidence | Certified uplift |
| Regression blocked | Donor/receiver worsens beyond threshold | Certified uplift |

## Sprint 6 invariants

- `universalSloProtectionAllowed=true`
- `certifiedPerformanceUpliftAllowed=true` (scoped CPU/RAM-bound only)
- `universalPerformanceImprovementClaimed=false`
- `guaranteedSavingsClaimed=false`
- `autoApplyAllowed=false`
- `productionAutonomousApplyAllowed=false`

## Source-of-truth map

| Domain | Owner |
|--------|-------|
| SLO guard / performance proof contracts | Karl-Hyperdensity |
| Runtime performance measurements | FluidVirt |
| SLO/performance projection | Karl-Dashboard |
| Identity / signals | Karl-Inventory |

## Forbidden claims

- universal performance improvement
- guaranteed savings active
- production autonomous apply
- Windows total RAM hotplug / logical vCPU hotplug
- Dashboard as performance source of truth
- FluidVirt as performance claim authority
