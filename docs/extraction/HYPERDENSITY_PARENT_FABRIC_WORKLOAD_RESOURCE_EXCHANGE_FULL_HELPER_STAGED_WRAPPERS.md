# Hyperdensity Parent Fabric — resource exchange full-helper staged wrappers (Sprint 82)

## Summary

**Sprint 82** completes staged wrappers for **all three** resource_exchange observation helpers: CPU (Sprint 80), ready, and restart (Sprint 81 candidates → Sprint 82 wrappers). Production `resource_exchange_*` remains **fully legacy**. Future wiring must be **FULL-HELPER** (8 + 12 + 12 call sites), not CPU-only.

---

## 1. Scope

| Item | Sprint 82 |
|------|-----------|
| CPU staged wrapper | **Yes** (from Sprint 80) |
| Ready staged wrapper | **Yes** (new) |
| Restart staged wrapper | **Yes** (new) |
| 8+8+8 triple parity matrix | **Yes** |
| Call-site wiring readiness doc | **Yes** |
| Production wiring | **No** |

---

## 2. Non-goals

- `ResourceExchangeObservationWiredV1 = true` or `CandidateRuntimeUsedV1 = true`.
- Sprint 83 call-site wiring (readiness only).
- CPU-only partial wiring (explicitly rejected).
- Hyperdensity copy; `parentfabric` import; broad observation.

---

## 3. Preconditions from Sprint 78–81

| Sprint | Deliverable |
|--------|-------------|
| **78** | 8 CPU legacy call sites |
| **79** | CPU candidate shadow |
| **80** | CPU staged wrapper |
| **81** | Ready/restart candidate shadow (12+12) |

---

## 4. Full-helper staged wrapper design

All three wrappers share AND gate:

```text
IF WiredV1 AND CandidateRuntimeUsedV1 → candidate
ELSE → legacy
```

| Wrapper | Legacy |
|---------|--------|
| `hyperdensityWorkloadResourceExchangeObservedPodCPURequestV1` | `hyperdensityObservedPodCPURequest` |
| `hyperdensityWorkloadResourceExchangeObservedPodContainerReadyV1` | `hyperdensityObservedPodContainerReady` (local) |
| `hyperdensityWorkloadResourceExchangeObservedPodContainerRestartCountV1` | `hyperdensityObservedPodContainerRestartCount` |

---

## 5. Why CPU-only wiring is rejected

| Risk | Detail |
|------|--------|
| Mixed semantics | 8 CPU wrappers + 12+12 legacy ready/restart → inconsistent exchange gates |
| Regression | Stage-apply checks ready/restart for donor/receiver health |
| Audit drift | Partial wiring breaks full-helper parity assumptions |

**Decision:** Future Sprint 83+ must replace **all 32** observation call sites (8+12+12), not CPU alone.

---

## 6. Flag behavior

Unchanged from Sprint 80–81: both flags **false** → legacy path at runtime.

---

## 7. Matrix coverage

- CPU: 8 cases (Sprint 79/80)
- Ready: 8 cases (Sprint 81)
- Restart: 8 cases (Sprint 81)

**Invariant:** wrapper == candidate == legacy per case.

---

## 8. No-touch surfaces

- Production `resource_exchange_*` call sites.
- apply track, rollback, VM, admission_guard.
- `workload_helpers.go` verdict: **`copy-deferred`**.

---

## 9. Future call-site wiring gates

See **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING_READINESS.md`**.

Sprint 83 may wire wrappers with **`ResourceExchangeObservationWiredV1` still false** (apply Sprint 70 pattern).

---

## 10. Rollback / no-change posture

Removing wiring file changes has no production effect while call sites stay legacy.

---

## 11. Risks

- Sprint 83 wires CPU only despite policy.
- Local `ContainerReady` vs shared `ContainerRestartCount` split complicates Hyperdensity extraction later.

---

## 12. Recommended next sprint

**Sprint 83:** `resource_exchange_callsite_wiring` — replace 8+12+12 legacy calls with wrappers; flags remain false until dedicated flip sprints.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING_READINESS.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_LOCAL_HELPER_SHADOW_MATRIX.md`
