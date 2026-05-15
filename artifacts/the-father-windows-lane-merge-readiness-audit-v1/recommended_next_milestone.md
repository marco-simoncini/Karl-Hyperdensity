# Recommended Next Milestone

## Recommendation
`hyperdensity_windows_fluidvirt_core_contracts_port_v1`

## Why
Audit outcome confirms selective integration is feasible if started with minimal core contracts/models and strict safety boundary. The lane is not direct-merge ready but is suitable for gated preview onboarding.

## Entry Criteria
- this audit artifacts set approved
- no direct merge decision accepted
- claim boundary accepted by product/safety stakeholders
- integration branch created from `main`

## Exit Criteria (for next milestone)
- PR 1 delivered (core contracts/models + minimal fixtures)
- no runtime mutation path enabled
- no Dashboard/UI porting included
- safety invariants explicitly asserted in tests/docs

## Fallback Milestone
If governance rejects immediate selective onboarding:
- `hyperdensity_windows_lane_cleanup_before_port_v1`
- focus on reducing artifact load, deduplicating contracts, and hardening claim boundary before any port.
