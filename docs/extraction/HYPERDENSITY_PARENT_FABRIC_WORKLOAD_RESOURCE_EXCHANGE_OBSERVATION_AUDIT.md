# Hyperdensity Parent Fabric — resource exchange observation audit (Sprint 78)

## Summary

**Sprint 78** is **audit/proposal-only** for `resource_exchange_*` observation call sites. **Sprint 79** adds shadow matrix (candidate parity); still **no** production wiring. **No** wrappers, candidates, flags, or runtime changes. Apply-observation track (Sprint 65–77) is **complete and separate**. Broad observation remains **disabled**.

---

## 1. Scope

| Item | Sprint 78 |
|------|-----------|
| Inventory observation call sites in `resource_exchange_*` | **Yes** (Dashboard audit script + golden) |
| Risk classification + future flag proposal | **Yes** (doc-only) |
| Wrapper / candidate creation | **No** |
| `resource_exchange_*` runtime changes | **No** |
| Broad observation flip | **No** |

---

## 2. Non-goals

- `ObservationWiredV1 = true` or `ProductionWiredV1 = true`.
- Introducing `hyperdensityWorkloadAdapterResourceExchangeObservationWiredV1` in Go (doc-only proposal).
- Reusing `ApplyObservationWiredV1` for resource_exchange.
- Wiring rollback, VM runtime, admission_guard, or `apply.go`.
- Hyperdensity Go adapter code or Dashboard `pkg/hyperdensity/parentfabric` import.
- Changing `workload_helpers.go` verdict (stays **`copy-deferred`**).
- API response, JSON ordering, or Sprint 56–77 wiring changes.

---

## 3. Current state after Sprint 77

| Flag / track | Value |
|--------------|-------|
| `ApplyObservationWiredV1` | **`true`** (apply track complete) |
| `ApplyObservationCandidateRuntimeUsedV1` | **`true`** |
| `ObservationWiredV1` | **`false`** |
| `ProductionWiredV1` | **`false`** |
| `ResourceExchangeObservationWiredV1` | **does not exist** |
| resource_exchange observation | **legacy** |
| rollback / VM runtime observation | **legacy** (out of scope) |
| `workload_helpers.go` verdict | **`copy-deferred`** |

resource_exchange must **not** reuse apply observation flags or apply wrappers/candidates.

---

## 4. Resource exchange files to audit

Production files matching `hyperdensity_parent_fabric_resource_exchange_*.go` (exclude `*_test.go`):

| File | Listed-helper call sites (Sprint 78 scan) |
|------|-------------------------------------------|
| `hyperdensity_parent_fabric_resource_exchange_stage_apply.go` | 4 × `hyperdensityObservedPodCPURequest` |
| `hyperdensity_parent_fabric_resource_exchange_stage_apply_chain.go` | 4 × `hyperdensityObservedPodCPURequest` |
| `hyperdensity_parent_fabric_resource_exchange_stage_apply_history.go` | 0 |
| `hyperdensity_parent_fabric_resource_exchange_v1.go` | 0 |

**Total listed-helper call sites:** **8**.

**Additional local observation-like helpers** (outside Sprint 78 listed inventory): e.g. `hyperdensityObservedPodContainerRestartCount`, `hyperdensityObservedPodContainerReady` in stage_apply files — classify in future shadow matrix; not counted in Sprint 78 golden.

---

## 5. Expected observation call-site classes

| Class | Helpers | Sprint 78 status |
|-------|---------|------------------|
| Pod CPU request (legacy) | `hyperdensityObservedPodCPURequest` | **8** call sites |
| Pod UID / limits / memory / container ID | listed spec helpers | **0** |
| FromObservation variants | listed spec helpers | **0** |
| Pilot observed state | `hyperdensityPilotObservedStateFrom*` | **0** |
| Apply wrappers / candidates | `hyperdensityWorkloadApplyObserved*` | **0** (forbidden) |

---

## 6. Risk classification

| Risk | Level | Notes |
|------|-------|-------|
| Stage-apply CPU observation during exchange admission | **High** | 8 direct legacy calls; wrong flip affects exchange plans |
| Confusion with apply observation track | **High** | Must use **dedicated** flag; never `ApplyObservationWiredV1` |
| Broad observation accidental flip | **Critical** | `ObservationWiredV1` must stay false |
| Bundling rollback / VM in same sprint | **Critical** | Safety surfaces; separate tracks |
| Local container-ready/restart helpers | **Medium** | Not in listed inventory; need shadow matrix |
| parentfabric import in Dashboard runtime | **High** | Forbidden; remains absent |

---

## 7. Proposed future flag

```go
// Doc-only proposal — NOT introduced in Sprint 78.
hyperdensityWorkloadAdapterResourceExchangeObservationWiredV1 = false
```

| Rule | Detail |
|------|--------|
| Introduce in Sprint 78? | **No** (documentation only) |
| Separate from apply? | **Yes** — must not alias `ApplyObservationWiredV1` |
| Default when introduced | **`false`** until dedicated activation sprint |

---

## 8. Required future shadow matrix

Before any wiring sprint (Sprint 79+):

1. Wrapper + candidate helpers scoped **only** to resource_exchange surfaces.
2. Legacy ≡ wrapper ≡ candidate matrix per helper class (including local container helpers).
3. Golden fixture: resource_exchange-only; no apply/rollback/VM drift.
4. Audit script extension: fail on wiring without flag approval.
5. Explicit exclusion of apply wrappers/candidates in `resource_exchange_*`.

---

## 9. Rollback / no-change posture

| Surface | Sprint 78 |
|---------|-----------|
| resource_exchange runtime | **unchanged** |
| apply observation | **unchanged** (flags stay true) |
| rollback observed-state | **untouched** |
| VM runtime observed-state | **untouched** |
| broad observation | **false** |

Rollback of apply track is documented in Sprint 77 boundary; Sprint 78 does not alter it.

---

## 10. Risks

- Treating Sprint 78 audit as permission to wire resource_exchange in Sprint 79 without shadow matrix.
- Sharing apply observation wrappers with stage-apply exchange paths.
- Omitting local container helpers from future matrix.
- Inferring `ObservationWiredV1=true` from apply track completion or resource_exchange inventory.

---

## 11. Recommended next sprint

**Sprint 79:** `resource_exchange_observation_shadow_matrix` — **complete** (see `SHADOW_MATRIX.md`). **Sprint 80:** `resource_exchange_observation_staged_wrappers` — proposal only; **no** wired flip.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_CRITERIA.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_MIGRATION_BOUNDARY.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_OBSERVATION_NEXT_SURFACE_DECISION.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_AUDIT_M73.md`
