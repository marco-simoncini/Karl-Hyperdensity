# shell-factory-v1

Contract ID: `hyperdensity_shell_factory_v1`

Defines the canonical creation and profile-readiness model for Hyperdensity-ready shells.

## Product principle

- No raw resource creation.
- Only Hyperdensity-ready shell creation via `HyperdensityShellClaim` + `HyperdensityShellProfile`.

## Canonical contract objects

- `HyperdensityShellClaim`
  - binds one object reference to one shell profile and readiness score.
- `HyperdensityShellProfile`
  - declares shell kind, envelope policy, telemetry/rollback/runtime-overlay requirements, and exchange eligibility target.
- `HyperdensityShellEnvelope`
  - CPU + memory floor/baseline/burstStep/ceiling.
- `HyperdensityReadinessScore`
  - readiness score (0-100), state, missing requirements, remediation lane.
- `HyperdensityFactoryProfile`
  - factory mode and policy assertions (raw creation forbidden).

## Built-in profile set (v1)

- `liquidity_donor_v1`
- `liquidity_receiver_v1`
- `linux_container_reference_v1`
- `linux_vm_reference_v1`
- `daas_interactive_v1` (placeholder)
- `service_burstable_v1` (placeholder)
- `batch_elastic_v1` (placeholder)
- `single_node_edge_v1` (placeholder)

All unsupported or not-yet-converged cases must map to a remediation lane, not to a definitive unsupported product state.
