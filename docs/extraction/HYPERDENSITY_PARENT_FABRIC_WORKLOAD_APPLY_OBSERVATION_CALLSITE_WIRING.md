# Hyperdensity Parent Fabric — apply observation call-site wiring (Sprint 70)

## Summary

**Sprint 70** is the **first real wiring** on `apply.go`: four legacy observation call sites repointed to **staged wrappers**. Flags remain **false** — wrappers delegate to legacy helpers (legacy-equivalent behavior).

---

## 1. Scope

| Item | Sprint 70 |
|------|-----------|
| `apply.go` call-site repoint (4 sites) | **Yes** |
| `ApplyObservationWiredV1 = true` | **No** |
| `ApplyObservationCandidateRuntimeUsedV1 = true` | **No** |
| Broad `ObservationWiredV1` | **No** |

---

## 2. Non-goals

- Flag flips (`ApplyObservationWiredV1`, candidate runtime).
- resource_exchange, rollback, VM, admission_guard.
- Other `apply.go` call sites.
- Hyperdensity Go code or Dashboard `parentfabric` import.

---

## 3. Preconditions from Sprint 65–69

Proposal, shadow matrix, staged wrappers, hardening 8×4, wiring readiness certified.

---

## 4. Call-site replacements

| Legacy | Wrapper |
|--------|---------|
| `hyperdensityObservedPodMemoryRequest` | `hyperdensityWorkloadApplyObservedPodMemoryRequestV1` |
| `hyperdensityObservedPodMemoryLimit` | `hyperdensityWorkloadApplyObservedPodMemoryLimitV1` |
| `hyperdensityObservedPodCPURequest` | `hyperdensityWorkloadApplyObservedPodCPURequestV1` |
| `hyperdensityObservedPodCPULimit` | `hyperdensityWorkloadApplyObservedPodCPULimitV1` |

---

## 5. Flag behavior

With `ApplyObservationWiredV1=false` and `ApplyObservationCandidateRuntimeUsedV1=false`, wrapper `if` branches take the **legacy** path.

---

## 6. Behavior equivalence statement

Runtime behavior is **expected unchanged** vs direct legacy calls because wrappers invoke the same legacy helpers when flags are false.

---

## 7. Rollback

Restore four legacy helper calls in `apply.go` (lines ~271–279). No flag changes required if flags were not flipped.

---

## 8. Test/audit coverage

- `hyperdensity_parent_fabric_workload_apply_observation_callsite_wiring_test.go`
- `audit_workload_apply_observation.sh` Sprint 70 guards
- Updated branch-swap guard (apply.go allowed for wrappers)

---

## 9. Risks

Future accidental flag flip without shadow re-validation. Confusion that call-site wiring implies `ApplyObservationWiredV1=true`.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_WIRING_READINESS.md`
