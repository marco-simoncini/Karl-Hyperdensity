# Fleet/multi-target guarded execution boundary (KHR-EE)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-EE |
| **Evidence** | Karl-Installer `docs/evidence/karl2-fleet-multitarget-guarded-execution/committed-khr-ee-v1/` |

---

## Guarded execution

| Field | Value |
|-------|-------|
| `fleetMultiTargetAuthorized` | `true` |
| `fleetApply` | `true` |
| `multiTargetApply` | `true` |
| `fleetMode` | `single-connected-reference-cluster` |
| `allowedScope` | `karl-metal-01@ovh` |
| `crossClusterMutation` | `false` |
| `targetMutation` | `false` |

State: `docs/khr/fleet-multitarget-state.json`. Enforcement deny-only from KHR-EC retained in `docs/khr/enforcement-model-state.json`.

---

## Scope boundary

Fleet/multi-target posture markers applied on dashboard deployment only. Hyperdensity apply paths for cross-cluster, workload, ResourceLease, ResourcePort loop, and autonomous orchestration remain denied. Dashboard projection and rdp-GW accessgraph/session remain read-only.
