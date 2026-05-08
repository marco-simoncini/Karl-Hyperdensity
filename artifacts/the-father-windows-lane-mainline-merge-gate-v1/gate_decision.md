# gate_decision

- decision: `merge_gate_passed_with_conditions`
- rationale: Hyperdensity Windows FluidVirt integration branch and Inventory fluidShell witness branch are operational for gated Technical Preview, with strict safety constraints and explicit defer of Dashboard stale branch and OS-ISO packaging.

## Conditions

- `directMergeAllowed=false`
- `selectiveMergeRequired=true`
- `dashboardMergeAllowed=false`
- `osIsoMergeAllowed=false`
- no GA/production-ready/execution-ready claims
- no runtime mutation, controlled apply, executor runtime, or actuator runtime enablement
