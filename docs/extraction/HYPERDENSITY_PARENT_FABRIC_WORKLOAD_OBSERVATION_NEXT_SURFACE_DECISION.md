# Hyperdensity Parent Fabric — observation next-surface decision (Sprint 77)

## Summary

**Sprint 77** classifies **remaining** legacy observation surfaces after apply-observation track completion. **`ObservationWiredV1` remains `false`.** Next sprint must **not** be broad observation.

**Sprint 78** completed **resource_exchange observation audit** (proposal only). **No** wiring. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_AUDIT.md`.

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
| Preferred next track | **resource_exchange_observation_shadow_matrix** (Sprint 78 audit complete) or **usage.go classification** |
| `ObservationWiredV1` | **`false`** (deliberate) |
| `ProductionWiredV1` | **`false`** |

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_MIGRATION_BOUNDARY.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_OBSERVATION_REAUDIT.md`
