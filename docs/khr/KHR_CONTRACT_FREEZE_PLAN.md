# KHR Contract Freeze Plan (KHR-AD)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AD |
| **Applies to** | TP boundary → beta-1 contract hardening |
| **Production enablement** | **Out of scope** |

---

## Freeze tiers

| Tier | Contracts | Policy through beta-1 |
|------|-----------|------------------------|
| **TP freeze** | `khr-projection-v1alpha1-readonly-y`, `khr-tp-readiness-summary-v1alpha1`, Host/Shell/Cell/ResourcePort/ResourceFuture CRDs + docs, Inventory posture schema (read-only fields) | Additive fields only; no removals; no semantic change to `readOnly: true` |
| **Experimental** | ResourceLease guarded apply, certification registry enforcement semantics, native-live cert claims, action approval apply path | May change with new evidence; not consumer-stable |
| **Beta target** | Inventory `TPObservationExport`, tp-readiness default-on in reference env, registry refresh export | Freeze at beta-1 RC after gap closure |
| **Deprecated** | Multus/NAD as target; unqualified KubeVirt-as-core wording | Documentation only |

---

## Contract version map

| Surface | Current version | Next (beta) | Breaking allowed? |
|---------|-----------------|-------------|-------------------|
| Dashboard projection | `khr-projection-v1alpha1-readonly-y` | `readonly-z` or `v1beta1-readonly` only with ADR | Additive only until ADR |
| TP readiness API | `khr-tp-readiness-summary-v1alpha1` | Same + optional fields | Additive only |
| runtime.karl.io CRDs | `v1alpha1` | `v1beta1` (future) | Requires conversion webhook plan |
| hyperdensity.karl.io CRDs | `v1alpha1` | `v1beta1` (future) | Same |
| Evidence summaries | sprint-tagged JSON | version field in summary | Additive metadata only |

---

## What to freeze for TP (now)

1. **Read-only projection fields** in `KHR_PROJECTION_V1.md` (Dashboard).
2. **TP package docs** and operator runbook procedures.
3. **Evidence anchor files** (`certification-summary.json`, registry/provenance summaries).
4. **CRD OpenAPI shapes** as installed by `install_khr_crds` (no field deletion).
5. **Guard scripts** (`guard_khr_docs_scope`, `guard_khr_iso_boundaries`, `khr_tp_package_check`).

---

## What stays experimental

| Item | Reason |
|------|--------|
| Guarded apply / cgroup markers | Sandbox-only; not production |
| `certified-preview` semantics | Not GA |
| Policy gate → apply coupling | Simulation only |
| Action approval → apply | Evidence only |
| Windows native-live lane | Observation gap |
| ResourcePort CR cluster apply | Disabled on ISO |

---

## What beta requires (non-runtime)

| Deliverable | Type |
|-------------|------|
| Signed contract inventory PASS | `khr_contract_inventory.sh` |
| Per-repo beta gap docs | ADR-AD series |
| Inventory live export spec frozen | schema + job contract |
| Breaking-change policy acknowledged by consumers | this document |
| Reference cluster evidence refresh SOP | operator runbook appendix |

---

## Breaking-change policy

| Change type | TP | Beta-1 | Post-beta |
|-------------|-----|--------|-----------|
| Add optional JSON field | Allowed | Allowed | Allowed |
| Add enum value (documented) | Allowed | Allowed | Allowed |
| Rename field | **Blocked** | ADR + minor bump | ADR + bump |
| Remove field | **Blocked** | Major bump only | Major bump |
| `readOnly: true` → false | **Forbidden** | **Forbidden** | **Forbidden** |
| New mutating API without sandbox gate | **Forbidden** | **Forbidden** | ADR + sandbox |

Consumers must pin `contractVersion` and treat unknown fields as opaque.

---

## Validation

```bash
./scripts/khr_contract_inventory.sh
./scripts/khr_tp_package_check.sh
./scripts/guard_khr_docs_scope.sh
```

---

## Cross-repo freeze ownership

| Repo | Owner doc |
|------|-----------|
| Karl-Hyperdensity | This plan + `BETA_READINESS_GAP_ANALYSIS.md` |
| Karl-Dashboard | `DASHBOARD_BETA_READINESS_GAP_ANALYSIS.md` |
| Karl-Inventory | `INVENTORY_BETA_READINESS_GAP_ANALYSIS.md` |
| Karl-OS-ISO | `ISO_BETA_READINESS_GAP_ANALYSIS.md` |

---

## Explicit non-goals

Production ready, autonomous orchestration, systemd auto-enable, Dashboard rewrite — remain **out of scope** through beta-1.
