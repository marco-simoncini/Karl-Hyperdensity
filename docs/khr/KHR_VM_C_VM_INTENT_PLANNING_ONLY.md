# Hyperdensity VM Intent planning-only validation (KHR-VM-C)

| Field | Value |
|-------|-------|
| Sprint | KHR-VM-C |
| Scope | Planning-only validation |
| Evidence | Karl-Installer `docs/evidence/khr-vm-dashboard-intent-dryrun/committed-khr-vm-c-v1/` |

## Mapping status

- VMIntent is consumed in planning-only mode for Shell/Cell mapping.
- ResourcePort mapping remains planning-only (`persistentLoop=false`).
- ResourceLease mapping remains planning-only (`applyExecuted=false`).
- Continuity mapping remains read-only with no runtime mutation.

## Guardrails

- `hyperdensityPlanningOnly=true`
- `resourceLeaseApply=false`
- `resourcePortPersistentLoop=false`
- `vmLifecycleKHRNativeReady=false`
