# Hyperdensity Parent Fabric — resource exchange local helper shadow matrix (Sprint 81)

## Summary

**Sprint 81** extends shadow coverage to ready/restart candidates. **Sprint 82** adds ready/restart **staged wrappers** and full-helper wiring readiness — see `FULL_HELPER_STAGED_WRAPPERS.md`. classified in Sprint 79: `hyperdensityObservedPodContainerReady` and `hyperdensityObservedPodContainerRestartCount`. **Candidate parity only** — no wrappers, no production wiring.

---

## 1. Scope

| Item | Sprint 81 |
|------|-----------|
| Ready candidate helper | **Yes** (test file) |
| Restart candidate helper | **Yes** (test file) |
| 8-case ready matrix | **Yes** |
| 8-case restart matrix | **Yes** |
| Ready/restart staged wrappers | **No** |
| Production call-site changes | **No** |

---

## 2. Non-goals

- `ResourceExchangeObservationWiredV1 = true` or `CandidateRuntimeUsedV1 = true`.
- Wrappers for ready/restart (Sprint 80 CPU wrapper unchanged, not wired).
- Replacing legacy calls in `resource_exchange_*`.
- Hyperdensity copy of local `ContainerReady` or `workload_helpers` restart.
- KARL-native storage/network primitives — not in scope; no contradictory decisions.

---

## 3. Preconditions from Sprint 79–80

| Sprint | State |
|--------|-------|
| **79** | Ready/restart classified (12+12 call sites); CPU shadow matrix |
| **80** | CPU staged wrapper; production still legacy |

---

## 4. Local helper semantics

| Helper | Definition | Call sites |
|--------|------------|------------|
| `hyperdensityObservedPodContainerReady` | `status.containerStatuses[].ready` | **12** (local to `resource_exchange_stage_apply.go`) |
| `hyperdensityObservedPodContainerRestartCount` | `status.containerStatuses[].restartCount` | **12** (delegates via `workload_helpers.go`) |

---

## 5. Candidate helper design

- `hyperdensityWorkloadResourceExchangeCandidateObservedPodContainerReadyV1` → delegates to legacy `bool`.
- `hyperdensityWorkloadResourceExchangeCandidateObservedPodContainerRestartCountV1` → delegates to legacy `int64`.

---

## 6. Matrix cases

**ContainerReady (8):** ready true/false, status missing, container missing, empty name, multi-container, malformed statuses, nil pod.

**ContainerRestartCount (8):** count 0, positive, float64 numeric, status missing, container missing, empty name, malformed statuses, nil pod.

---

## 7. Expected invariants

| Metric | Value |
|--------|-------|
| CPU legacy call sites | **8** |
| Ready legacy call sites | **12** |
| Restart legacy call sites | **12** |
| Production candidate/wrapper use | **0** |
| `ObservationWiredV1` | **false** |

---

## 8. No-touch surfaces

- CPU wrapper file (Sprint 80) — no ready/restart wrappers added.
- apply track, rollback, VM runtime, admission_guard.
- `workload_helpers.go` verdict: **`copy-deferred`**.

---

## 9. Future wrapper/wiring gates

1. Local helper **staged wrappers** (dedicated sprint), or
2. **Call-site wiring readiness** with explicit CPU-only vs full-helper decision.

---

## 10. Risks

- Wiring CPU only while leaving ready/restart legacy → exchange gate regressions.
- Copying local `ContainerReady` to Hyperdensity prematurely.
- Accidental candidate use in production `resource_exchange_*`.

---

## 11. Recommended next sprint

**Option A:** `resource_exchange_local_helper_staged_wrappers` (proposal/test only).  
**Option B:** `resource_exchange_callsite_wiring_readiness` with CPU-only vs full-helper decision record.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_LOCAL_HELPER_CRITERIA.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_LOCAL_HELPER_CLASSIFICATION.md`
