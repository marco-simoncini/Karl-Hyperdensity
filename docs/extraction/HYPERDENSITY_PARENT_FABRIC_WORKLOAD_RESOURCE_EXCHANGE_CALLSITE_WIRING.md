# Hyperdensity Parent Fabric — resource exchange call-site wiring (Sprint 83)

## Summary

**Sprint 83** wires all **32** production `resource_exchange_*` observation call sites to **full-helper staged wrappers** (8 CPU + 12 ready + 12 restart). Wiring follows the Sprint 82 decision: **FULL-HELPER**, not CPU-only. Flags remain **false** — effective runtime path stays **legacy**; behavior expected **unchanged**.

---

## 1. Scope

| Item | Sprint 83 |
|------|-----------|
| Production call-site repoint (32 sites) | **Yes** |
| `ResourceExchangeObservationWiredV1 = true` | **No** |
| `ResourceExchangeObservationCandidateRuntimeUsedV1 = true` | **No** |
| Broad `ObservationWiredV1` / `ProductionWiredV1` | **No** |
| Apply observation track | **Unchanged** (complete) |

---

## 2. Non-goals

- Flag activation or candidate-runtime flip.
- Broad observation.
- Direct candidate calls in production.
- CPU-only partial wiring.
- Changes to ready/restart **legacy helper definitions**.
- `apply.go`, apply wrappers, admission_guard, rollback, VM runtime.
- Dashboard import of Hyperdensity `parentfabric`.
- Storage/network primitive implementation.
- Moving `workload_helpers.go` (verdict stays **`copy-deferred`**).

---

## 3. Preconditions from Sprint 78–82

| Sprint | Deliverable |
|--------|-------------|
| 78 | CPU observation audit (8 sites) |
| 79 | CPU candidate + shadow matrix |
| 80 | CPU staged wrapper |
| 81 | Ready/restart candidates + local shadow matrix |
| 82 | Full-helper staged wrappers + wiring readiness |

---

## 4. Full-helper call-site replacement

| Legacy | Wrapper | Count |
|--------|---------|-------|
| `hyperdensityObservedPodCPURequest` | `hyperdensityWorkloadResourceExchangeObservedPodCPURequestV1` | **8** |
| `hyperdensityObservedPodContainerReady` | `hyperdensityWorkloadResourceExchangeObservedPodContainerReadyV1` | **12** |
| `hyperdensityObservedPodContainerRestartCount` | `hyperdensityWorkloadResourceExchangeObservedPodContainerRestartCountV1` | **12** |

Files: `resource_exchange_stage_apply.go`, `resource_exchange_stage_apply_chain.go`.

**CPU-only wiring is forbidden.** Local `func hyperdensityObservedPodContainerReady` definition is **not** replaced.

---

## 5. Flag state

| Flag | Sprint 83 |
|------|-----------|
| `ResourceExchangeObservationWiredV1` | **`false`** |
| `ResourceExchangeObservationCandidateRuntimeUsedV1` | **`false`** |
| `ObservationWiredV1` | **`false`** |
| `ProductionWiredV1` | **`false`** |
| `ApplyObservationWiredV1` | **`true`** (unchanged) |
| `ApplyObservationCandidateRuntimeUsedV1` | **`true`** (unchanged) |

---

## 6. Runtime behavior statement

With both resource-exchange flags false, wrappers use the **AND gate** and take the **legacy** branch. Runtime behavior is **expected unchanged** vs pre-Sprint-83 direct legacy calls.

---

## 7. Production source invariants

- Legacy production call-site counts: **0 / 0 / 0** (CPU / ready / restart).
- Wrapper production call-site counts: **8 / 12 / 12**.
- Direct candidate production calls: **0**.
- No apply observation wrappers/candidates in `resource_exchange_*`.
- No `pkg/hyperdensity/parentfabric` import in Dashboard runtime.

---

## 8. Rollback

Restore legacy helper calls in `resource_exchange_*` production files; keep flags false; re-run parity and shadow matrices.

---

## 9. Risks

Accidental flag flip without re-validation. Confusion that call-site wiring implies `ResourceExchangeObservationWiredV1=true`. Partial CPU-only rollback.

---

## 10. Recommended next sprint

Dedicated sprint to stage **candidate-runtime flip** (`ResourceExchangeObservationCandidateRuntimeUsedV1`) with shadow re-validation — **not** combined with `WiredV1` activation unless criteria explicitly allow.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING_READINESS.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_FULL_HELPER_STAGED_WRAPPERS.md`
