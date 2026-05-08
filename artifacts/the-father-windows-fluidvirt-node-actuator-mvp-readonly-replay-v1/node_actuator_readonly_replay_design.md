# Node Actuator Readonly Replay Design

## Milestone
`hyperdensity_windows_fluidvirt_node_actuator_mvp_readonly_replay_v1`

## Goal
Represent and validate Node Fluid Actuator MVP behavior as readonly replay/fake-runtime simulation without runtime mutation.

## Key Outcome
- replay available (`nodeActuatorMvpReplayAvailable=true`)
- runtime not ported (`nodeActuatorMvpPortedAsRuntime=false`)
- runtime actuator disabled
- cgroup/QMP/QGA untouched
