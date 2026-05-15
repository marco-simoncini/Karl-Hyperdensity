# Runtime Baseline

Kubernetes identity:

- VM: `master-win11`
- Namespace: `karl`
- VM UID: `c81b95dc-d955-4fb3-a1af-59d979f48bcb`
- VMI UID: `a565fc95-d79f-4cbd-a852-e0467e5f2110`
- VMI phase: `Running`
- virt-launcher pod: `virt-launcher-master-win11-w7mxv`
- Pod UID: `6e0a597f-b508-450d-a4ac-d5e4262c8615`
- Node: `karl-lab-metal-01`
- Pod restarts: `0`
- VMIM objects in `karl`: none
- LiveMigration objects for target: none observed

CPU mismatch baseline:

- QMP/libvirt CPU count: `7`
- Guest CPU count (QGA `guest-get-vcpus`): `6 online`
- Mismatch classification: `QMP_GUEST_CPU_DIVERGENCE`

Continuity baseline:

- QEMU PID: `92`
- QEMU start time: `Thu May 7 18:16:51 2026`
- Windows last boot: `/Date(1778177812500)/`
- Machine identity hash: `d1fdf2ce69932d7ac2f9e3497be7f845d4a3ef9140c45bd0d7b0ea47c195508e`
- Pending reboot: `false`
- Critical events 24h: `0`
