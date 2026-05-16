# Hyperdensity Parent Fabric — workload observation re-audit (Sprint 64)

## Summary

**Sprint 64** re-audits remaining legacy observation surfaces after Sprint 63 live branch swap. **Sprint 65** adds apply-observation **proposal-only**. **Sprint 66** shadow matrix. **Sprint 67** apply staged wrappers. **Sprint 68** wrapper hardening. **Sprint 69** wiring readiness. **Sprint 70** wired `apply.go` call sites to wrappers. **Sprint 71** post-wiring hardening. **Sprint 72** flip criteria. **Sprint 73–74** candidate readiness/flip. **Sprint 75** activation. **Sprint 76** post-activation hardening. **Sprint 77** migration boundary — apply track complete. **Sprint 78** audit — **8** CPU call sites. **Sprint 79** CPU shadow. **Sprint 80** CPU staged wrappers. **Sprint 81** ready/restart shadow — **no** production wiring. **`hyperdensityWorkloadAdapterObservationWiredV1`** remains **`false`** by deliberate policy.

---

## 1. Scope

| Item | Sprint 64 |
|------|-----------|
| Remaining observation inventory | **Yes** (Dashboard audit script) |
| Broad observation policy decision | **Yes** |
| New runtime wiring | **No** |
| Hyperdensity Go adapter code | **No** |

---

## 2. Current wired surfaces

| Surface | Flag | Status |
|---------|------|--------|
| Path helpers | `PathWiredV1` | **`true`** (Sprint 56) |
| Pilot observed-state | `PilotObservationWiredV1` | **`true`** (Sprint 57) |
| Live observed-state | `LiveObservationWiredV1` + candidate runtime | **`true`** (Sprint 61–63) |
| Broad observation | `ObservationWiredV1` | **`false`** |

---

## 3. Remaining legacy surfaces

| Category | Examples | Policy |
|----------|----------|--------|
| **legacy-apply** | `hyperdensity_parent_fabric_apply.go` | **Forbidden** until dedicated sprint |
| **legacy-resource-exchange** | `resource_exchange_stage_apply*.go` | **Forbidden** |
| **legacy-rollback** | `rollback.go` | **Forbidden** |
| **legacy-vm-runtime** | `vm_linux_*` | **Forbidden** |
| **legacy-admission** | `admission_guard_*` | **Forbidden** |
| **other-review** | e.g. `usage.go` | Review per sprint |
| **already-wired-pilot** | `pilot.go` (pod enrichment legacy) | In-scope pilot; not broad |
| **already-wired-live** | `live.go` (wrappers only) | In-scope live |

---

## 4. Risk classification

| Risk | Level | Mitigation |
|------|-------|------------|
| Accidental `ObservationWiredV1=true` | High | Policy test + audit |
| apply/resource-exchange wiring | High | Explicit forbidden categories |
| False sense of completion after live | Medium | Re-audit documents 32+ remaining call sites |

---

## 5. Broad observation policy

See **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_BROAD_OBSERVATION_DECISION.md`**. Recommendation: keep **`ObservationWiredV1=false`**; use granular subflags; broad flip only when **all** surfaces migrated and shadowed.

---

## 6. Rollback posture

Sprint 64 is docs/audit only — no runtime rollback required.

---

## 7. Test/audit requirements

| Artifact | Owner |
|----------|-------|
| `audit_workload_observation_remaining.sh` | Dashboard |
| `hyperdensity_parent_fabric_workload_observation_remaining_audit_test.go` | Dashboard |
| `hyperdensity_parent_fabric_workload_broad_observation_policy_test.go` | Dashboard |

---

## 8. Risks

Premature broad observation flip would bypass granular governance established in Sprints 56–63.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_BROAD_OBSERVATION_DECISION.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_REAUDIT_CRITERIA.md`


---

## Sprint 83 (call-site wiring)

Sprint 83 wires all 32 production `resource_exchange_*` observation call sites to full-helper staged wrappers (8 CPU + 12 ready + 12 restart). `ResourceExchangeObservationWiredV1` and `ResourceExchangeObservationCandidateRuntimeUsedV1` remain **false**; effective runtime path is **legacy**. Direct candidate calls in production remain forbidden. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING.md` (Hyperdensity) and `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING_M83.md` (Dashboard).


---

## Sprint 84 (candidate-runtime staging)

Sprint 84 sets `ResourceExchangeObservationCandidateRuntimeUsedV1=true` while `ResourceExchangeObservationWiredV1=false`. AND gate keeps effective runtime on legacy; candidate branch inactive. Production call-sites remain wrappers from Sprint 83. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CANDIDATE_RUNTIME_STAGING.md`.


---

## Sprint 85 (activation readiness)

Sprint 85 is readiness-only for `ResourceExchangeObservationWiredV1=true`. No flag changes. Sprint 86 may execute activation flip if approved. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_ACTIVATION_READINESS.md`.
