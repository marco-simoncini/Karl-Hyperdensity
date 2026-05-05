# Technical Preview Release Notes v1

Release: **KARL Hyperdensity Technical Preview**  
Program: **Grande Padre resource control plane**

## Release summary

KARL Hyperdensity Technical Preview delivers a governed live resource market for CPU/RAM across Hyperdensity-ready shells, with:
- policy-governed decision boundaries
- operator-controlled execution paths
- evidence-backed support claims
- one product surface for live runtime control (`KARL Live Resource Authority`)

This is not a Kubernetes dashboard and not a simple autoscaler.

## What is included

- Linux container CPU/RAM live scaling within supported TP paths.
- Linux VM CPU/RAM evidence-backed live scaling in object-specific lanes.
- VM RAM lane via runtime overlay (`virtio-mem/QMP/QOM requested-size`) where proven.
- No generic KubeVirt template memory mutation claim.
- Resource Exchange donor/receiver planning and transfer recommendations.
- Transfer dry-run and staged/chained apply+rollback proof in evidence namespace.
- Shell Factory.
- Shell Claim Generator.
- Shell Claim Template/Profile Pack.
- Admission Guard audit/classification.
- Mutate Preview.
- Enforce Simulation.
- Mutate Preview Apply Dry-Run.
- Policy Pack.
- Policy Consistency Checker.
- Release Support Matrix.
- Evidence Bundle / Demo Scenario Pack.
- KARL Live Resource Authority.

## Approved claims

- **Linux container claim**: In TP, bounded Linux container CPU/RAM live adjustment with verification and rollback in supported paths.
- **Linux VM claim**: In TP, Linux VM CPU/RAM support is evidence-backed/object-specific; VM RAM uses runtime overlay through virtio-mem/QMP/QOM requested-size where proven.
- **Live Resource Authority claim**: KARL exposes one governed/auditable live CPU/RAM control surface over multiple runtime drivers.
- **Resource Exchange claim**: KARL can model donor/receiver liquidity-demand, validate, stage, apply (operator-controlled), verify, and rollback with audit history.
- **Admission remediation claim**: Classification, mutate preview, enforce simulation, and apply dry-run remediation chain is available without enabling enforcement.
- **Official creation path claim**: No raw resource creation; only Hyperdensity-ready shell creation through template/profile + claim flow.
- **Governance/support boundary claim**: Policy Pack, Consistency, Support Matrix, Evidence Bundle, and Authority define bounded TP support.

## Explicit non-claims

- Not GA.
- No autonomous production movement.
- No enforcement enabled.
- No production workload mutation.
- No Windows support.
- No universal workload support claim.
- No generic VM RAM live resize claim via KubeVirt template mutation.
- No universal no-disruption guarantee for every workload.
- Dry-run success is not production readiness.
- `warming_up` is not ready.

## Known limitations

- VM lane remains evidence-backed/object-specific.
- Historical proof references require refresh before Beta/GA.
- Control Room UI redesign is deferred.
- Packaging/install/upgrade is not GA-ready.
- Telemetry freshness/confidence remains required for stronger claims.
- Migration/snapshot/restore interaction proof is still expanding.
- Runtime overlay reconciliation after restart/recreate requires explicit handling and revalidation.

## Upgrade path / next release direction

- Use the Operator Runbook and Readiness Gate as operational baseline.
- Keep TP wording and safety boundaries explicit.
- Plan Control Room UI redesign in separate milestone.
- Plan Reactivity Optimization in later milestone.
- Plan Action Slate / Resource Futures in later milestone.
