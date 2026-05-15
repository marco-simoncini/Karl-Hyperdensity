# Standalone vs Pool-Child Model

- Standalone Windows VM: evaluated directly as shell target.
- Pool-child Windows VM: accepted only when treated as per-VM shell and pool remains provisioning context.
- Pool scaling as runtime mechanism: blocked with remediation `block_pool_scaling_as_mechanism`.
