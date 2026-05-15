# Cluster Baseline

- Cluster context: `karl-metal-01@ovh`
- Namespace: `karl`
- Target VM: `master-win11`
- VM UID: `c81b95dc-d955-4fb3-a1af-59d979f48bcb`
- VM runStrategy: `Halted`
- VM printableStatus: `Stopped`
- VMI: not found (`VMINotExists`)
- virt-launcher pod: not found
- VMIM objects in namespace: none
- LiveMigration objects in namespace: none observed for target
- Pool context-only VMs:
  - `win11-pool-0` (Stopped)
  - `win11-pool-1` (Stopped)

## Safety decision

Runtime proof sprint is blocked before mutation because the target is not running and does not expose a live VMI/virt-launcher runtime to prove in-place continuity.
