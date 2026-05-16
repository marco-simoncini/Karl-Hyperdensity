# Hyperdensity Parent Fabric — apply observation branch logic (Sprint 73)

## Summary

**Sprint 73** documented branch logic. **Sprint 74** flipped `CandidateRuntimeUsedV1=true`; AND gate keeps legacy path while `ApplyObservationWiredV1=false`. in apply observation staged wrappers (`apply_observation_wiring_v1.go`). **No code changes** in Sprint 73.

---

## Expected wrapper pattern (all four helpers)

Each of the four apply observation wrappers follows:

```go
if hyperdensityWorkloadAdapterApplyObservationWiredV1 &&
   hyperdensityWorkloadAdapterApplyObservationCandidateRuntimeUsedV1 {
    return candidate(...)
}
return legacy(...)
```

Concrete functions:

| Wrapper | Candidate | Legacy |
|---------|-----------|--------|
| `hyperdensityWorkloadApplyObservedPodMemoryRequestV1` | `hyperdensityWorkloadApplyCandidateObservedPodMemoryRequestV1` | `hyperdensityObservedPodMemoryRequest` |
| `hyperdensityWorkloadApplyObservedPodMemoryLimitV1` | `hyperdensityWorkloadApplyCandidateObservedPodMemoryLimitV1` | `hyperdensityObservedPodMemoryLimit` |
| `hyperdensityWorkloadApplyObservedPodCPURequestV1` | `hyperdensityWorkloadApplyCandidateObservedPodCPURequestV1` | `hyperdensityObservedPodCPURequest` |
| `hyperdensityWorkloadApplyObservedPodCPULimitV1` | `hyperdensityWorkloadApplyCandidateObservedPodCPULimitV1` | `hyperdensityObservedPodCPULimit` |

---

## Truth table (runtime path)

| `ApplyObservationWiredV1` | `ApplyObservationCandidateRuntimeUsedV1` | Runtime path |
|---------------------------|------------------------------------------|--------------|
| `false` | `false` | **legacy** ← current |
| `false` | `true` | **legacy** |
| `true` | `false` | **legacy** |
| `true` | `true` | **candidate** |

**Implication:** A future flip of **only** `ApplyObservationCandidateRuntimeUsedV1=true` does **not** change production behavior until `ApplyObservationWiredV1=true` in a separate approved sprint.

---

## Call graph (current)

```text
apply.go
  → hyperdensityWorkloadApplyObservedPod*V1 (wrapper)
      → [both flags true] → hyperdensityWorkloadApplyCandidateObservedPod*V1
      → [otherwise]       → hyperdensityObservedPod* (legacy)
```

`apply.go` must **not** call candidate helpers directly.

---

## Contrast with live observation (Sprint 63)

Live observation branch swap uses `LiveObservationWiredV1` and `CandidateRuntimeUsedV1` on **live** wrappers in `live.go` — a **different** surface from apply observation.

Apply observation uses **apply-specific** flags only; do not reuse live flags in apply wrappers.

---

## Future flip constraints

| Change | Allowed in candidate-runtime flip sprint? |
|--------|-------------------------------------------|
| Set `ApplyObservationCandidateRuntimeUsedV1=true` | **Yes** (dedicated sprint) |
| Set `ApplyObservationWiredV1=true` | **Separate sprint** (recommended) |
| Change AND to OR in wrapper branch | **No** without explicit architecture sprint |
| Change `apply.go` call sites | **No** unless rollback/wiring sprint |

---

## Verification (Sprint 73)

- Dashboard test reads `apply_observation_wiring_v1.go` and asserts AND pattern × 4.
- Shadow tests continue to prove candidate ≡ legacy for all fixtures.
- Post-wiring hardening golden: `branchSwapAllowed: false`.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_CANDIDATE_RUNTIME_READINESS.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_STAGED_WRAPPERS.md`
