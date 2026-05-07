# Same QEMU / Same Boot Proof

Before actuator mutations:

- QEMU PID: `96`
- QEMU start: `Thu May 7 18:58:03 2026`
- Windows last boot: `/Date(1778180311500)/`
- pod: `virt-launcher-master-win11-kmwgg`
- pod UID: `7b6a904a-1c9a-4a44-9b37-1dc737304773`
- node: `karl-lab-metal-01`

After CPU entitlement up/down and rollback:

- QEMU PID: `96`
- QEMU start: `Thu May 7 18:58:03 2026`
- Windows last boot: `/Date(1778180311500)/`
- pod: unchanged
- pod UID: unchanged
- node: unchanged
- compute restart count: `0`

Runtime invariants:

- no VMIM objects in namespace
- no migration evidence during test window
- no rollout/recreate during proof window
