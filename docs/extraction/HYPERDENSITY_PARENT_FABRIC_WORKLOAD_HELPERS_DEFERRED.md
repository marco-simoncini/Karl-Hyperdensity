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
10. **Sprint 68:** apply wrapper hardening — 8×4 shadow matrix; still no `apply.go` wiring.
11. **Sprint 69:** apply wiring readiness — `readyForApplyGoCallSiteWiring`.
12. **Sprint 70:** apply call-site wiring — 4 wrappers in `apply.go`; flags **false**.
13. **Sprint 71:** apply post-wiring hardening — 8×4; flags **false**.
14. **Sprint 72:** apply flip criteria — docs-only; flags **false**; no runtime changes.
15. **Sprint 73:** candidate-runtime readiness — branch logic verified; no flip.
16. **Sprint 74:** candidate-runtime staging flip — `CandidateRuntimeUsedV1=true`; Wired **false**; runtime legacy-equivalent.
17. **Sprint 75:** apply observation activation — `ApplyObservationWiredV1=true`; candidate branch active; candidate ≡ legacy.
18. **Sprint 76:** post-activation hardening — 8×4; flags unchanged; broad observation **false**.
19. **Sprint 77:** migration boundary — apply track complete; broad observation **false**.
20. **Sprint 78:** resource_exchange observation audit — **full file verdict unchanged: `copy-deferred`**; no helpers copy.
21. **Sprint 79:** resource_exchange shadow matrix — candidate in test file only; **full file verdict unchanged: `copy-deferred`**.
22. **Sprint 80:** resource_exchange staged wrappers — wiring file test-only; **full file verdict unchanged: `copy-deferred`**.
23. **Sprint 81:** local helper shadow matrix — ready/restart candidates test-only; **full file verdict unchanged: `copy-deferred`**.
24. **Sprint 82:** full-helper staged wrappers + wiring readiness — **full file verdict unchanged: `copy-deferred`**.

## Placeholder

`pkg/hyperdensity/parentfabric/workload/` — **three-function** pure-core copy-contract (Sprint 52); remaining helpers deferred.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_HELPERS_AUDIT.md`


---

## Sprint 83 (call-site wiring)

Sprint 83 wires all 32 production `resource_exchange_*` observation call sites to full-helper staged wrappers (8 CPU + 12 ready + 12 restart). `ResourceExchangeObservationWiredV1` and `ResourceExchangeObservationCandidateRuntimeUsedV1` remain **false**; effective runtime path is **legacy**. Direct candidate calls in production remain forbidden. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING.md` (Hyperdensity) and `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING_M83.md` (Dashboard).


---

## Sprint 84 (candidate-runtime staging)

Sprint 84 sets `ResourceExchangeObservationCandidateRuntimeUsedV1=true` while `ResourceExchangeObservationWiredV1=false`. AND gate keeps effective runtime on legacy; candidate branch inactive. Production call-sites remain wrappers from Sprint 83. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CANDIDATE_RUNTIME_STAGING.md`.


---

## Sprint 85 (activation readiness)

Sprint 85 is readiness-only for `ResourceExchangeObservationWiredV1=true`. No flag changes. Sprint 86 may execute activation flip if approved. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_ACTIVATION_READINESS.md`.


---

## Sprint 86 (resource_exchange activation)

Sprint 86 sets ResourceExchangeObservationWiredV1=true. Candidate branch active in resource_exchange wrappers only. ObservationWiredV1/ProductionWiredV1 remain false. See ACTIVATION.md and POST_ACTIVATION_HARDENING.md.


---

## Sprint 87 (resource_exchange boundary closure)

Sprint 87 closes resource_exchange observation Sprint 78–86 as boundary complete. No flag/runtime changes. Broad observation remains false. Next phase: KHR architecture memory and storage/network semantics. See MIGRATION_BOUNDARY.md, REMAINING_SURFACE_DECISION.md, KHR_ROADMAP_TRANSITION_NOTE.md.
