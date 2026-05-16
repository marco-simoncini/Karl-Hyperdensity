# Hyperdensity Parent Fabric — workload helpers copy deferred (Sprint 48)

## Decision

**`copy-deferred`** for `hyperdensity_parent_fabric_workload_helpers.go`.

No `helpers.go`, golden, or copy-contract in Karl-Hyperdensity for Sprint 48.

## Blockers

| Blocker | Detail |
|---------|--------|
| **Cross-package coupling** | 35+ functions call `server` symbols defined outside this file (execution actions, nested map helpers, pilot types, deployment state builder). |
| **KubeVirt / K8s paths** | Six functions build REST paths containing `kubevirt.io` or Kubernetes API prefixes as string literals. |
| **Runtime observation** | Functions read untyped `map[string]interface{}` workload/pod shapes and produce `HyperdensityPilotObservedState` — runtime-bound, not DTO-only. |
| **Apply semantics** | Execution mode / ready-reason / mechanism helpers are part of the apply orchestration surface. |

## What would be required before copy

1. ~~Extract stdlib-only nested-map + quantity primitives~~ — **Sprint 49 done** (`parentfabric/primitives`); not wired to Dashboard yet.
2. **Sprint 50:** adapter **boundary** + interface **proposal** + re-audit **criteria** documented; Dashboard classification fixture (46 functions).
3. **Sprint 51:** Dashboard **test-only** adapter stubs (`*_test.go`) — **no production wiring**.
4. **Sprint 52:** **three pure-candidates** copied to `parentfabric/workload` with golden + manifest — **full file still `copy-deferred`**.
5. Re-audit for remaining 43 functions only after production adapter sprint.
6. **Sprint 64:** observation surface re-audit — broad `ObservationWiredV1` stays **false**; granular subflags only.
7. **Sprint 65:** apply-observation proposal only — `apply.go` legacy; `ApplyObservationWiredV1` placeholder **false**.
8. **Sprint 66:** apply-observation shadow matrix — candidate helpers test-only; `ApplyObservationCandidateRuntimeUsedV1` **false**.
9. **Sprint 67:** apply-observation staged wrappers — `apply_observation_wiring_v1.go`; `apply.go` does not call wrappers yet.

## Placeholder

`pkg/hyperdensity/parentfabric/workload/` — **three-function** pure-core copy-contract (Sprint 52); remaining helpers deferred.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_HELPERS_AUDIT.md`
