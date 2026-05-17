# KARL 2.0 production readiness declaration (KHR-EF)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-EF |
| **Evidence** | Karl-Installer `docs/evidence/karl2-production-readiness-declaration/committed-khr-ef-v1/` |

---

## Declaration

| Field | Value |
|-------|-------|
| `productionReady` | `true` |
| `promotionAllowed` | `true` |
| `readinessScope` | `current-connected-reference-production` |
| `productionMode` | `karl2-baremetal-khr-native-single-connected-reference-cluster` |

State: `docs/khr/production-readiness-state.json`. Enforcement and fleet states unchanged.

---

## Karl-App exclusion

Karl-App formally excluded from KARL 2.0 production readiness. `karlAppRequiredForProductionDecision=false`. No Karl-App work required.

---

## Safety invariants preserved

Dashboard projection and rdp-GW accessgraph/session remain read-only. No cross-cluster mutation, workload mutation, ResourceLease apply, ResourcePort loop, or autonomous orchestration.
