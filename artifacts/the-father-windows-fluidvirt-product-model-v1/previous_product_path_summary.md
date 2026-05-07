# Previous Product Path Summary

Inputs consumed:

- definitive product proof: `ce05534`, `WINDOWS_FLUIDVIRT_PRODUCT_PATH_CONFIRMED`
- actuator MVP hardening: `d7c0226`, `windows_node_fluid_actuator_mvp_ready`

Extracted constants:

- CPU floor/ceiling: `300000 100000` / `600000 100000`
- RAM floor/ceiling: `12884901888` / `13958643712`
- continuity requirement: same QEMU + same boot + no migration/recreate/rollout
