# Claim Boundary Attestation

This milestone attests the Windows FluidVirt core contract/model boundary:

## Allowed
- Windows FluidVirt Technical Preview candidate
- CPU liquidity via entitlement / cgroup v2 `cpu.max` model
- RAM liquidity via QMP balloon model
- guest witness dependency (`KARL Agent` + `fluidShell` + `QGA`) as future integration dependency
- same-QEMU / same-boot / guest ACK / rollback / return-to-floor / audit as model requirements

## Forbidden (explicitly enforced)
- Windows GA
- Windows production-ready
- Windows execution-ready by default
- vCPU hotplug support claim
- logical CPU scaling support claim
- pool scaling support claim
- LiveMigration / VMIM support claim
- reboot/recreate/rollout mechanism claim
- autonomous apply
- production AUTO
- raw runtime controls exposure
- generic KubeVirt VM RAM template mutation support claim
