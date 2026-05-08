# Node Actuator Contract Boundary Design

## Milestone
`hyperdensity_windows_fluidvirt_node_actuator_contract_boundary_v1`

## Intent
Define contract-only boundaries, gates, blockers, and attestations required before any Node Fluid Actuator MVP runtime port.

## Design Summary
- new model: `WindowsFluidVirtNodeActuatorContract`
- mode: `boundary_only`
- release: `technical_preview`
- lane: `gated_preview`
- runtime actuator remains disabled
- no cgroup write path enabled
- no QMP/QGA command execution
- no executor/controlled apply integration in this milestone
