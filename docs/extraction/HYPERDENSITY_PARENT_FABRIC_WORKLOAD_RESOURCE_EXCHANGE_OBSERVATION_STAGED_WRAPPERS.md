# Hyperdensity Parent Fabric — resource exchange observation staged wrappers (Sprint 80)

## Summary

**Sprint 80** introduces CPU staged wrapper. **Sprint 82** completes **full-helper** staged wrappers (CPU + ready + restart). Production remains legacy. Wrappers exist in a dedicated wiring file and are validated by tests only. **Production `resource_exchange_*` call sites remain legacy** (`hyperdensityObservedPodCPURequest`). **`ResourceExchangeObservationWiredV1` and `CandidateRuntimeUsedV1` remain false.**

---

## 1. Scope

| Item | Sprint 80 |
|------|-----------|
| Staged wrapper file | **Yes** (`hyperdensityWorkloadResourceExchangeObservedPodCPURequestV1`) |
| 8-case wrapper ≡ candidate ≡ legacy matrix | **Yes** |
| Production call-site wiring | **No** |
| Local ready/restart wrappers | **No** (deferred) |
| Broad observation | **No** |

---

## 2. Non-goals

- `ResourceExchangeObservationWiredV1 = true` or `CandidateRuntimeUsedV1 = true`.
- Replacing legacy calls in `resource_exchange_stage_apply*.go`.
- Apply wrappers/candidates in resource_exchange.
- Rollback, VM runtime, admission_guard.
- Hyperdensity copy or Dashboard `parentfabric` import.
- KARL-native storage/network primitives (EphemeralDisk, KARLNetwork, etc.) — **not in scope**; no decisions contrary to roadmap.

---

## 3. Preconditions from Sprint 78–79

| Sprint | Deliverable |
|--------|-------------|
| **78** | Audit — **8** `hyperdensityObservedPodCPURequest` call sites |
| **79** | Shadow matrix — candidate parity, not wired |

---

## 4. Staged wrapper design

```text
hyperdensityWorkloadResourceExchangeObservedPodCPURequestV1(pod, containerName)
  IF WiredV1 AND CandidateRuntimeUsedV1 → candidate helper
  ELSE → hyperdensityObservedPodCPURequest (legacy)
```

With both flags **false**, runtime path is **legacy** (production-equivalent when wired in future).

---

## 5. Flag behavior

| Flag | Value |
|------|-------|
| `ResourceExchangeObservationCandidateV1` | **true** |
| `ResourceExchangeObservationCandidateRuntimeUsedV1` | **false** |
| `ResourceExchangeObservationShadowReadyV1` | **true** |
| `ResourceExchangeObservationWiredV1` | **false** |
| `ObservationWiredV1` | **false** |
| `ApplyObservationWiredV1` | **true** (unchanged) |

---

## 6. Matrix coverage

Same **8** cases as Sprint 79 shadow matrix: normal CPU, missing container, empty name, multi main/sidecar, missing resources, malformed resources, nil pod.

**Invariant:** wrapper == candidate == legacy for all cases (flags false → wrapper uses legacy branch, still equivalent to candidate).

---

## 7. Local helper deferral

| Helper | Call sites | Sprint 80 |
|--------|------------|-------------|
| `hyperdensityObservedPodContainerReady` | 12 | **not wrapped** |
| `hyperdensityObservedPodContainerRestartCount` | 12 | **not wrapped** |

Future: local helper shadow matrix or explicit decision before call-site wiring.

---

## 8. No-touch surfaces

- All `resource_exchange_*` production observation call sites.
- apply track, pilot/live path wiring, rollback, VM runtime.
- `workload_helpers.go` verdict: **`copy-deferred`**.

---

## 9. Future wiring gates

Before production uses wrappers:

1. Call-site wiring readiness sprint.
2. Optional local helper matrix.
3. `CandidateRuntimeUsedV1` staging flip (dedicated sprint).
4. `WiredV1` activation (dedicated sprint).
5. Parity + audit green.

---

## 10. Rollback / no-change posture

Removing wiring file has **no** production effect while call sites stay legacy.

---

## 11. Risks

- Accidental wrapper use in `resource_exchange_*` before readiness.
- Flipping `WiredV1` without `CandidateRuntimeUsedV1` staging.
- CPU-only wrappers missing ready/restart regression coverage.

---

## 12. Recommended next sprint

1. **resource_exchange local helper shadow matrix**, or  
2. **resource_exchange call-site wiring readiness** (proposal only).

**Not recommended:** `ResourceExchangeObservationWiredV1=true` without readiness + local helper policy.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_STAGED_WRAPPER_CRITERIA.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_SHADOW_MATRIX.md`
