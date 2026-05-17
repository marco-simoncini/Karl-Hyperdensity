# KHR-native VM Intent contract (KHR-VM-B)

| Field | Value |
|-------|-------|
| Sprint | KHR-VM-B |
| Scope | Contract definition only |
| Evidence | Karl-Installer `docs/evidence/khr-vm-native-intent-contract/committed-khr-vm-b-v1/` |

## Contract coverage

- Canonical `VMIntent` structure (`apiVersion`, `kind`, `metadata`, `spec`, `status`).
- Provider compatibility modes: `khr-native`, `kubevirt-compatibility`, `hybrid-transition`.
- Lifecycle action contract map: create/start/stop/restart/delete/snapshot/migrate/recover.
- Dry-run semantics for this sprint (`dryRunOnly=true`, `applyExecuted=false`).
- Required status conditions and rollback planning constraints.

## Hyperdensity mapping boundary

- VM Intent plans Shell/Cell, ResourcePort, ResourceLease, and continuity relationships.
- Mapping is planning-only in KHR-VM-B; no runtime apply path and no mutation.
- VM lifecycle KHR-native readiness remains false in this sprint.
