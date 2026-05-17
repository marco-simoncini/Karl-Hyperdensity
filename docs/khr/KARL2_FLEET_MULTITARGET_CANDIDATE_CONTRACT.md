# Fleet/multi-target candidate contract (KHR-ED)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-ED |
| **Evidence** | Karl-Installer `docs/evidence/karl2-fleet-multitarget-authorization-gate/committed-khr-ed-v1/` |

---

## Candidate fleet/multi-target model

| Field | Value |
|-------|-------|
| `mode` | `candidate-dry-run` |
| `allowedScope` | `read-only-target-discovery` |
| `singleReferenceCluster` | `karl-metal-01@ovh` |
| `candidateClusterCount` | `1` |
| `fleetApply` | `false` |
| `multiTargetApply` | `false` |

---

## Scope boundary

Hyperdensity fleet/multi-target apply paths remain disabled. Dashboard KHR backend projection and rdp-GW accessgraph/session remain read-only. Target inventory dry-run classifies targets without cluster mutation.

---

## Authorization

`KARL2_FLEET_MULTITARGET_I_UNDERSTAND` guard required for any future fleet/multi-target execution sprint. `guardSatisfied=false` and `fleetMultiTargetAuthorized=false` in KHR-ED.

---

## Enforcement posture (retained from KHR-EC)

| Field | Value |
|-------|-------|
| `enforcementEnabled` | `true` |
| `enforcementMode` | `deny-only` |
| `allowedScope` | `single-reference-cluster` |

State: `docs/khr/enforcement-model-state.json`. No enforcement rollback in KHR-ED.

---

## Guarded execution (KHR-EE)

| Field | Value |
|-------|-------|
| `fleetMultiTargetAuthorized` | true |
| `fleetApply` | true |
| `multiTargetApply` | true |
| `fleetMode` | single-connected-reference-cluster |
| `allowedScope` | karl-metal-01@ovh |

Evidence: Karl-Installer `docs/evidence/karl2-fleet-multitarget-guarded-execution/committed-khr-ee-v1/`. State: `docs/khr/fleet-multitarget-state.json`.

---

## Production readiness (KHR-EF)

| Field | Value |
|-------|-------|
| `productionReady` | true |
| `promotionAllowed` | true |
| `karlAppExcludedFromProductionReadiness` | true |

Evidence: Karl-Installer `docs/evidence/karl2-production-readiness-declaration/committed-khr-ef-v1/`. State: `docs/khr/production-readiness-state.json`.

---

## Production release (KHR-EG)

| Field | Value |
|-------|-------|
| `karl2ProductionRelease` | true |
| `productionReady` | true |
| `promotionAllowed` | true |

Evidence: Karl-Installer `docs/evidence/karl2-production-release-declaration/committed-khr-eg-v1/`. State: `docs/khr/production-release-state.json`.

---

## Production post-release operating window (KHR-EH)

| Field | Value |
|-------|-------|
| `postReleaseWindowPassed` | true |
| `stableAcrossWindow` | true |
| `productionReady` | true |

Evidence: Karl-Installer `docs/evidence/karl2-production-postrelease-operating-window/committed-khr-eh-v1/`. State: `docs/khr/production-postrelease-stability-state.json`.
