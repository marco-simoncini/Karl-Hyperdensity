# Blocker Report

Replay status:

- Real `master-win11` replay: no blockers (`HYPERDENSITY_READY_WINDOWS_SHELL`).
- Pool-child replay: no blockers (`HYPERDENSITY_READY_WINDOWS_SHELL`).
- Pool-scaling-mechanism replay: blocked with `pool_scaling_as_mechanism`.

Capability gap replayed:

- missing CPU actuator -> `BLOCKED_WITH_REMEDIATION` (`node_fluid_actuator_unavailable`, `cgroup_path_mismatch`)
- missing RAM balloon -> `BLOCKED_WITH_REMEDIATION` (`ram_balloon_unavailable`, `return_to_floor_not_ready`)
