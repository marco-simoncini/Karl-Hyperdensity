# Windows Fluid Node Actuator Readonly Replay v1

This contract introduces a **readonly replay / fake-runtime simulation** for the Windows Node Fluid Actuator MVP behavior.

## Boundary

- replay-only mode (`readonly_replay_only`)
- fake-runtime only
- no runtime actuator enabled
- no real cgroup write
- no QMP execution
- no QGA execution
- no production apply
- no autonomous apply
- no Windows GA claim
- no Windows production-ready claim

## Intent

The replay demonstrates and validates model transitions and safety checks against:
- `windows_fluidvirt_node_actuator_contract_boundary_v1`

without performing runtime mutations.

## What Is Included

- replay scenario model
- fake-runtime boundary
- replayed CPU entitlement transitions (model-only)
- replayed safety checks
- replayed blockers
- replayed audit events
- replay validation result

## What Is Excluded

- runtime actuator executable
- cgroup runtime write path
- QMP/QGA command execution
- controlled apply / executor integration
- Dashboard/Inventory/OS-ISO integration

## Next Step

A future milestone may port Node Actuator MVP runtime only behind additional explicit gates, validations, and safety controls.
