# KARL Hyperdensity — Guarded Auto Policy Engine v1

**Contract ID:** `hyperdensity_guarded_auto_policy_engine_v1`  
**Milestone:** `hyperdensity_guarded_auto_policy_engine_v1`  
**Release track:** `technical_preview`

## Product definition

> KARL Hyperdensity can classify operator-ready actions as guarded-auto candidates under explicit policy, blast-radius, SLO, rollback, ledger and kill-switch constraints, while keeping auto-apply execution disabled.

## Policy modes

| Mode | Sprint 7 |
|------|----------|
| `recommendation_only` | allowed |
| `operator_controlled` | allowed |
| `guarded_auto_candidate` | allowed (classification only) |
| `guarded_auto_sandbox_ready` | allowed (classification only) |
| `guarded_auto_nonprod_ready` | allowed (classification only) |
| `production_canary_blocked` | allowed |
| `production_auto_blocked` | allowed |
| `production_auto_with_policy` | **forbidden** |

## Candidate gate requirements

A `guarded_auto_candidate` requires all gates to pass:

- shell passport exists
- action slate entry exists
- dry-run valid
- rollback ready
- SLO guard passed
- no-regression certified
- donor health preserved
- receiver health preserved or neutral/no-claim accepted
- realized ledger record or explicit ledger not-required reason
- risk score within policy
- blast radius budget available
- kill switch clear
- circuit breaker closed
- rate limit available
- cooldown expired
- resource family allowed
- no blocked/protected shell unless manual-only override
- no Windows evidence-gated auto unless remediation-only or manual-only

## Sprint 7 invariants

- `autoApplyExecutionEnabled=false`
- `productionAutonomousApplyAllowed=false`
- `productionScope=false`
- `productionMutationAllowed=false` (implicit)
- no auto-apply execution in this sprint

## Source-of-truth map

| Domain | Owner |
|--------|-------|
| Guarded auto policy / eligibility | Karl-Hyperdensity |
| Runtime constraint evidence | FluidVirt |
| Policy projection | Karl-Dashboard (read-only) |
| Identity / signals | Karl-Inventory |

## Forbidden claims

- auto-apply executed
- production autonomous apply / autonomous production mutation
- guaranteed savings active
- universal performance improvement
- Windows total RAM hotplug / logical vCPU hotplug
- Dashboard as policy source of truth
- FluidVirt as policy authority
- Inventory/Warden runtime apply
