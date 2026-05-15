# Live State Baseline

Kubernetes:

- VM: `master-win11`
- Namespace: `karl`
- VM UID: `c81b95dc-d955-4fb3-a1af-59d979f48bcb`
- VMI UID: `a565fc95-d79f-4cbd-a852-e0467e5f2110`
- VMI phase: `Running`
- virt-launcher pod: `virt-launcher-master-win11-w7mxv`
- pod UID: `6e0a597f-b508-450d-a4ac-d5e4262c8615`
- node: `karl-lab-metal-01`
- restart count: `0`
- VMIM objects: none

QEMU/libvirt/QMP:

- QEMU PID: `92`
- QEMU start time: `Thu May 7 18:16:51 2026`
- QMP `query-status`: `running`
- QMP `query-cpus-fast`: `7` vCPU objects (`cpu-index 0..6`)
- QMP hotplug path for extra CPU: `/machine/peripheral/vcpu6`
- libvirt `dominfo CPU(s)`: `7`
- libvirt `vcpucount --live`: `7`
- libvirt `vcpucount --config`: `6`
- max vCPU (`--maximum --live`/`--config`): `24`
- memory summary read-only: base only, plugged memory `0`

Guest Windows:

- guestAck (`guest-ping`): true
- guest processor count (`guest-get-vcpus`): `6`
- `Win32_ComputerSystem.NumberOfLogicalProcessors`: `6`
- `Win32_Processor.NumberOfLogicalProcessors`: `6`
- last boot: `/Date(1778177812500)/`
- machineGuidHash: `d1fdf2ce69932d7ac2f9e3497be7f845d4a3ef9140c45bd0d7b0ea47c195508e`
- pending reboot: false
- critical events 24h: `0`

Classification:

- QMP CPU count: `7`
- libvirt CPU count: `7`
- guest CPU count: `6`
- mismatch class: `QMP_LIBVIRT_GUEST_CPU_DIVERGENCE_PERSISTENT`
