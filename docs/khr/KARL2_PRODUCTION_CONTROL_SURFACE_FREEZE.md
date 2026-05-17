# Hyperdensity control-surface freeze (KHR-DZ)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-DZ |
| **Evidence** | Karl-Installer `docs/evidence/karl2-production-control-surface-audit/committed-khr-dz-v1/` |

---

## Frozen control surface

Live dashboard KHR backend projection audited read-only after KHR-DY rollout.

| Field | Posture |
|-------|---------|
| `enforcementEnabled` | false |
| `autonomousOrchestration` | false |
| `resourceLeaseApplyExposed` | false |
| `resourcePortPersistentLoopExposed` | false |
| `fleetApply` | false |
| `multiTargetApply` | false |

---

## Historical env note

Dashboard deployment retains historical Hyperdensity/Wave env vars from pre-KHR production path. Live projection confirms these remain **legacy-inert**; no enforcement, ResourceLease apply, ResourcePort persistent loop, or fleet/multi-target action surface is exposed.

---

## Blocker freeze

`enforcementModelFrozen=true` — enforcement model remains frozen; no enablement in KHR-DZ.
