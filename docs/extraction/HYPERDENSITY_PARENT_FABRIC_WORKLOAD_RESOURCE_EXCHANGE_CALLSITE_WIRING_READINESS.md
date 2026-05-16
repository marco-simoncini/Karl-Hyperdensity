# Hyperdensity Parent Fabric — resource exchange call-site wiring readiness (Sprint 82)

## Summary

**Sprint 82** certifies readiness for a future **Sprint 83** that may wire production `resource_exchange_*` call sites to staged wrappers. **Readiness-only** — no production changes.

---

## 1. Scope

| Item | Sprint 82 |
|------|-----------|
| Readiness golden + tests | **Yes** |
| Full-helper staged wrappers complete | **Yes** |
| Production call-site wiring | **No** |
| Flag activation | **No** |

---

## 2. Readiness decision

| Decision | Value |
|----------|-------|
| `readyForResourceExchangeCallsiteWiring` | **true** |
| `futureWiringDecision` | **`full-helper`** |
| `cpuOnlyWiringAllowed` | **`false`** |

---

## 3. Call-site inventory

| Helper | Production call sites |
|--------|----------------------|
| `hyperdensityObservedPodCPURequest` | **8** |
| `hyperdensityObservedPodContainerReady` | **12** |
| `hyperdensityObservedPodContainerRestartCount` | **12** |
| **Total** | **32** |

Files: `resource_exchange_stage_apply.go`, `resource_exchange_stage_apply_chain.go` (and history/v1 with zero listed-helper sites).

---

## 4. Required future replacements

Sprint 83 must replace:

| Legacy | Wrapper |
|--------|---------|
| `hyperdensityObservedPodCPURequest` | `hyperdensityWorkloadResourceExchangeObservedPodCPURequestV1` |
| `hyperdensityObservedPodContainerReady` | `hyperdensityWorkloadResourceExchangeObservedPodContainerReadyV1` |
| `hyperdensityObservedPodContainerRestartCount` | `hyperdensityWorkloadResourceExchangeObservedPodContainerRestartCountV1` |

**Partial CPU-only replacement is forbidden.**

---

## 5. Required future flags

| Flag | Sprint 83 first wiring |
|------|------------------------|
| `ResourceExchangeObservationWiredV1` | **`false`** (like apply Sprint 70) |
| `ResourceExchangeObservationCandidateRuntimeUsedV1` | **`false`** |
| `ObservationWiredV1` | **`false`** |

Candidate-runtime flip and wired flip are **later** dedicated sprints.

---

## 6. Required tests/audits

- Parity tests green after wiring.
- `audit_workload_resource_exchange_observation.sh` updated counts.
- No direct candidate calls in production.
- wrapper == legacy with flags false.

---

## 7. Rollback

Restore legacy calls in `resource_exchange_*`; set flags false; re-run parity.

---

## 8. Risks

- Incomplete wiring (CPU only).
- Premature `WiredV1=true` before candidate-runtime staging.
- Broad observation inferred from resource_exchange progress.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_FULL_HELPER_STAGED_WRAPPERS.md`


---

## Sprint 83 (call-site wiring)

Sprint 83 wires all 32 production `resource_exchange_*` observation call sites to full-helper staged wrappers (8 CPU + 12 ready + 12 restart). `ResourceExchangeObservationWiredV1` and `ResourceExchangeObservationCandidateRuntimeUsedV1` remain **false**; effective runtime path is **legacy**. Direct candidate calls in production remain forbidden. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING.md` (Hyperdensity) and `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING_M83.md` (Dashboard).


---

## Sprint 84 (candidate-runtime staging)

Sprint 84 sets `ResourceExchangeObservationCandidateRuntimeUsedV1=true` while `ResourceExchangeObservationWiredV1=false`. AND gate keeps effective runtime on legacy; candidate branch inactive. Production call-sites remain wrappers from Sprint 83. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CANDIDATE_RUNTIME_STAGING.md`.


---

## Sprint 85 (activation readiness)

Sprint 85 is readiness-only for `ResourceExchangeObservationWiredV1=true`. No flag changes. Sprint 86 may execute activation flip if approved. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_ACTIVATION_READINESS.md`.
