# Windows Claim Boundary

## Allowed Claims
- Windows FluidVirt is a **Technical Preview candidate**.
- CPU liquidity path: entitlement via node-local cgroup v2 `cpu.max` actuator.
- RAM liquidity path: balloon via QMP.
- Continuity proof model: same QEMU, same boot, same pod, same node.
- Guest witness path: KARL Agent / `fluidShell` / guest ACK evidence.
- Safety lifecycle includes rollback, return-to-floor, and auditable replay evidence.

## Forbidden Claims
- Windows GA.
- Windows production-ready.
- VM lane execution-ready by default.
- vCPU hotplug support.
- logical CPU scaling support.
- pool scaling as runtime scaling mechanism.
- LiveMigration / VMIM as part of fluid scaling path.
- reboot/recreate/rollout as product mechanism.
- autonomous apply / production AUTO enablement.
- raw runtime controls exposure (QMP/libvirt/QGA/QOM/K8s patch knobs).
- generic KubeVirt VM RAM template mutation claim.
