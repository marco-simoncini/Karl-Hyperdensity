# Hyperdensity Parent Fabric — resource exchange observation shadow matrix (Sprint 79)

## Summary

**Sprint 79** is **CPU shadow-matrix / test-only**. **Sprint 80** CPU staged wrappers. **Sprint 81** local ready/restart shadow — see `LOCAL_HELPER_SHADOW_MATRIX.md`. Dashboard introduces a **candidate** helper equivalent to `hyperdensityObservedPodCPURequest` with an **8-case** parity matrix. **No** production wiring in `resource_exchange_*`. **`ResourceExchangeObservationWiredV1` remains false or absent.**

---

## 1. Scope

| Item | Sprint 79 |
|------|-----------|
| Candidate helper (test-only file) | **Yes** |
| Shadow matrix golden + tests | **Yes** |
| Local helper classification doc/test | **Yes** |
| `resource_exchange_*` runtime call-site changes | **No** |
| Staged wrappers | **No** |
| Broad observation flip | **No** |

---

## 2. Non-goals

- `ResourceExchangeObservationWiredV1 = true`.
- `ObservationWiredV1 = true` or `ProductionWiredV1 = true`.
- Replacing `hyperdensityObservedPodCPURequest` in resource_exchange production files.
- Apply wrappers/candidates in resource_exchange.
- Rollback, VM runtime, admission_guard wiring.
- Hyperdensity Go copy or Dashboard `parentfabric` import.
- Changing `workload_helpers.go` verdict (**`copy-deferred`**).

---

## 3. Sprint 78 input audit

| Metric | Value |
|--------|-------|
| Listed-helper call sites | **8** |
| Helper | `hyperdensityObservedPodCPURequest` |
| `stage_apply.go` | 4 |
| `stage_apply_chain.go` | 4 |
| Apply wrappers in resource_exchange | **0** |

---

## 4. Shadow helper classes

| Class | Sprint 79 |
|-------|-----------|
| Pod CPU request (legacy + candidate) | **1** helper in matrix |
| Local container ready/restart | **classified only** — not in matrix |

Candidate: `hyperdensityWorkloadResourceExchangeCandidateObservedPodCPURequestV1` → delegates to legacy (equivalence).

---

## 5. Matrix cases

| # | Case |
|---|------|
| 1 | Normal container with CPU request (`250m`) |
| 2 | Missing container name |
| 3 | Empty container name |
| 4 | Multi-container — main |
| 5 | Multi-container — sidecar (`100m`) |
| 6 | Missing resources block |
| 7 | Malformed resources (wrong type) |
| 8 | Nil pod map |

**Invariant:** candidate == legacy (`reflect.DeepEqual`) for all cases.

---

## 6. Expected invariants

| Flag | Value |
|------|-------|
| `ResourceExchangeObservationCandidateV1` | **true** |
| `ResourceExchangeObservationCandidateRuntimeUsedV1` | **false** |
| `ResourceExchangeObservationShadowReadyV1` | **true** |
| `ResourceExchangeObservationWiredV1` | **false** |
| `ApplyObservationWiredV1` | **true** (unchanged) |
| `ObservationWiredV1` | **false** |
| Production resource_exchange CPU call sites | **8** (unchanged) |

---

## 7. No-touch surfaces

- `apply.go`, apply wrappers, apply candidate semantics.
- `resource_exchange_*` production observation call sites.
- rollback, VM runtime, admission_guard.
- Sprint 56–77 path/pilot/live/apply wiring.

---

## 8. Future wrapper/wiring gates

Before `ResourceExchangeObservationWiredV1=true`:

1. Staged wrappers sprint (Sprint 80+ proposal).
2. Wrapper ≡ candidate ≡ legacy matrix extended if local helpers wired.
3. Dedicated activation sprint.
4. Audit script + parity green.

---

## 9. Rollback / no-change posture

| Surface | Sprint 79 |
|---------|-----------|
| resource_exchange runtime | **legacy** |
| Candidate file | removable without runtime effect |
| Apply track | **unchanged** |

---

## 10. Risks

- Accidental candidate call from `resource_exchange_*` production files.
- Treating shadow PASS as permission to flip wired flag without wrappers sprint.
- Ignoring local helpers in future matrix.
- Reusing apply observation flags for resource_exchange.

---

## 11. Recommended next sprint

**Sprint 80:** staged wrappers — **complete** (see `STAGED_WRAPPERS.md`).

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_AUDIT.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_LOCAL_HELPER_CLASSIFICATION.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_SHADOW_MATRIX_M75.md`


---

## Sprint 83 (call-site wiring)

Sprint 83 wires all 32 production `resource_exchange_*` observation call sites to full-helper staged wrappers (8 CPU + 12 ready + 12 restart). `ResourceExchangeObservationWiredV1` and `ResourceExchangeObservationCandidateRuntimeUsedV1` remain **false**; effective runtime path is **legacy**. Direct candidate calls in production remain forbidden. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING.md` (Hyperdensity) and `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING_M83.md` (Dashboard).


---

## Sprint 84 (candidate-runtime staging)

Sprint 84 sets `ResourceExchangeObservationCandidateRuntimeUsedV1=true` while `ResourceExchangeObservationWiredV1=false`. AND gate keeps effective runtime on legacy; candidate branch inactive. Production call-sites remain wrappers from Sprint 83. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CANDIDATE_RUNTIME_STAGING.md`.
