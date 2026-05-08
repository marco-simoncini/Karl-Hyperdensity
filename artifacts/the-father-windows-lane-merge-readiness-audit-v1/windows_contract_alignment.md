# Windows Contract Alignment

## Alignment Objective
Map Windows FluidVirt contracts to Linux Hyperdensity governance vocabulary while preserving strict support boundary.

## Lifecycle Alignment

| Capability | Linux Hyperdensity Baseline | Windows FluidVirt Mapping | Alignment Status |
|---|---|---|---|
| prepare | shell eligibility + lease preparation | `EvaluateWindowsHyperdensityTarget` + `PrepareWindowsFluidResourceLease` | aligned (preview) |
| dry-run | governance/read-only prevalidation | `BuildWindowsFluidActionSlate` + `EvaluateWindowsComplianceReplay` + actuator dry-run step | aligned (read-only) |
| apply | gated, explicit operator intent | `WindowsFluidControlledApplyPlan` (policy gated; default blocked) | partial (planning-first) |
| verify | guest/runtime post-check | guest ACK + workload verify plan + continuity evidence | aligned (preview) |
| rollback | mandatory safety path | rollback target in lease + rollback blockers | aligned |
| return-to-floor | mandatory restoration path | return-to-floor target + explicit blocker model | aligned |
| audit | immutable evidence chain | replay hash + bundle hash chain + attestation envelope (future-signable) | aligned |
| savings semantics | avoid overclaim | maintain estimated vs realized separation | required guardrail |
| grant/authorization | explicit governance controls | allowlist/TTL/kill-switch/manual approval gates | aligned |
| witness | guest truth/evidence plane | `fluidShell` in `Karl-Inventory` | aligned, separate repo |
| readiness | promote only with full evidence | compliance phase model + controlled apply blockers | aligned, not GA |
| support boundary | no unsafe claims | no hotplug, no pool scaling, no autonomous apply | aligned |

## Naming Convergence Recommendation
Use Parent Fabric oriented naming for integration artifacts:
- `windowsFluidVirtProductModel`
- `windowsFluidVirtActionSlate`
- `windowsFluidVirtDriverAdapterContract`
- `windowsFluidVirtControlledApplyPlan`

This keeps Linux GA path stable while introducing Windows as gated preview lane.
