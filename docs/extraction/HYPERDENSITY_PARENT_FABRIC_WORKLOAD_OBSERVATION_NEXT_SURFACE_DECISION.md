# Hyperdensity Parent Fabric — observation next-surface decision (Sprint 77)

## Summary

**Sprint 77** classifies **remaining** legacy observation surfaces after apply-observation track completion. **`ObservationWiredV1` remains `false`.** Next sprint must **not** be broad observation.

**Sprint 78–81:** audit, CPU shadow, CPU staged wrappers, local helper shadow complete. Production resource_exchange not wired.

---

## 1. resource_exchange

| Field | Value |
|-------|-------|
| **Status** | legacy |
| **Risk** | **high** |
| **Files** | `resource_exchange_stage_apply*.go` |
| **Next allowed action** | audit / proposal only |
| **Forbidden until** | dedicated shadow matrix + criteria sprint |

Stage-apply observation is **not** the same as the four `apply.go` pod helpers. Do not reuse `ApplyObservationWiredV1` for exchange surfaces.

---

## 2. rollback

| Field | Value |
|-------|-------|
| **Status** | legacy |
| **Risk** | **very high** / safety-critical |
| **Files** | `hyperdensity_parent_fabric_rollback.go` |
| **Next allowed action** | audit only |
| **Forbidden until** | resource_exchange policy complete; explicit safety review |

Rollback observed-state must not be wired in the same sprint as apply boundary closure or broad observation.

---

## 3. VM runtime

| Field | Value |
|-------|-------|
| **Status** | legacy |
| **Risk** | **high** |
| **Files** | `hyperdensity_parent_fabric_vm_linux_*.go` |
| **Next allowed action** | audit only |
| **Forbidden until** | dedicated VM observation proposal; no Windows runtime claim |

---

## 4. usage.go / other-review

| Field | Value |
|-------|-------|
| **Status** | review-needed |
| **Risk** | **medium** |
| **Files** | `hyperdensity_parent_fabric_usage.go`, adapter wiring files (delegation only) |
| **Next allowed action** | classification sprint |

Adapter files may appear in remaining-audit counts due to **delegation** to legacy helpers — not as wiring targets.

---

## Decision

| Question | Answer |
|----------|--------|
| Is apply-observation complete? | **Yes** (Sprint 65–76) |
| Is next sprint broad observation? | **No** |
| Preferred next track | **resource_exchange_callsite_wiring** (Sprint 82 readiness complete; full-helper only) |
| `ObservationWiredV1` | **`false`** (deliberate) |
| `ProductionWiredV1` | **`false`** |

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_MIGRATION_BOUNDARY.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_OBSERVATION_REAUDIT.md`


---

## Sprint 83 (call-site wiring)

Sprint 83 wires all 32 production `resource_exchange_*` observation call sites to full-helper staged wrappers (8 CPU + 12 ready + 12 restart). `ResourceExchangeObservationWiredV1` and `ResourceExchangeObservationCandidateRuntimeUsedV1` remain **false**; effective runtime path is **legacy**. Direct candidate calls in production remain forbidden. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING.md` (Hyperdensity) and `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING_M83.md` (Dashboard).


---

## Sprint 84 (candidate-runtime staging)

Sprint 84 sets `ResourceExchangeObservationCandidateRuntimeUsedV1=true` while `ResourceExchangeObservationWiredV1=false`. AND gate keeps effective runtime on legacy; candidate branch inactive. Production call-sites remain wrappers from Sprint 83. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CANDIDATE_RUNTIME_STAGING.md`.
