# KARL Hyperdensity — Realized Savings Ledger v1

**Contract ID:** `hyperdensity_realized_savings_ledger_v1`  
**Milestone:** `hyperdensity_realized_savings_ledger_v1`  
**Release track:** `technical_preview`

## Product definition

> KARL Hyperdensity can account realized value from operator-approved resource movements using lease duration, resource amount, unit price, SLO preservation, rollback impact and confidence, while keeping guaranteed savings disabled.

## Strict value distinctions

| Class | Rules |
|-------|-------|
| Estimated opportunity | Projection only; never realized; never guaranteed |
| Realized value | Mutation observed + post-verify passed or evidence-gated with explicit non-guaranteed classification |
| Eligible for future guarantee | Classification only in Sprint 5; `guaranteeEligibleForFuture=true` allowed; `guaranteedSavingsAllowed=false` |
| Excluded / non-guaranteed | Synthetic, estimated-only, missing price/duration, SLO failed, rollback unavailable, post-verify failed, Windows evidence-gated |

## Value formula (Sprint 5)

```
grossValue = normalizedResourceAmount × durationHours × unitPrice
netRealizedValue = grossValue - controlPlaneCost - rollbackCost - riskAdjustment
```

Conservative rules: no round-up; negative net value → `realized_non_guaranteed` with zero future guarantee eligibility; missing unit price → `excluded_missing_unit_price`.

## Sprint 5 invariants

- `guaranteedSavingsAllowed=false`
- `guaranteedSavingsClaimed=false`
- `estimatedValueCountedAsGuaranteed=false`
- `syntheticValueCountedAsProduction=false`

## Source-of-truth map

| Domain | Owner |
|--------|-------|
| Ledger contracts / classification | Karl-Hyperdensity |
| Movement measurements | FluidVirt |
| Ledger projection | Karl-Dashboard |
| Identity / signals | Karl-Inventory |

## Forbidden claims

- guaranteed savings active / guaranteed savings claimed
- universal performance improvement
- production autonomous apply
- Windows total RAM hotplug / logical vCPU hotplug
- estimated value counted as guaranteed
- synthetic/shadow value counted as production
- Dashboard as accounting source of truth
- FluidVirt as accounting authority
