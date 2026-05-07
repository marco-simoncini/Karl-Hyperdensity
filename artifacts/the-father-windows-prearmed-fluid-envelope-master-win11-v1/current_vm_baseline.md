# Current VM Baseline

Kubernetes:

- VM: `master-win11`
- namespace: `karl`
- VM UID: `c81b95dc-d955-4fb3-a1af-59d979f48bcb`
- VMI UID: `b31dcab6-5a99-432e-a6fe-21607a9e3403`
- pod: `virt-launcher-master-win11-kmwgg`
- pod UID: `7b6a904a-1c9a-4a44-9b37-1dc737304773`
- node: `karl-lab-metal-01`
- pod restart count: `0`
- VMIM objects: none

QEMU/libvirt/QMP:

- QEMU PID: `96`
- QEMU start: `Thu May 7 18:58:03 2026`
- vCPU runtime count: `6`
- max vCPU: `24`
- memory current/max: `13631488 KiB` / `13631488 KiB`
- balloon support: present (`virtio-balloon`, QMP `query-balloon`)
- QMP memory-devices: empty (no memory hotplug devices)
- vCPU thread IDs: `101,102,103,104,105,106`
- scheduler tuning via `virsh schedinfo`: unavailable in session mode

cgroup and CPU controls:

- QEMU cgroup path: `/`
- `cpu.max`: `600000 100000`
- `cpu.weight`: `67`
- cgroup filesystem for write controls: read-only from compute container
- direct `taskset` affinity change on QEMU/vCPU threads: not permitted

Guest:

- guest ACK: true
- processor count: `6`
- total visible memory bytes: `13938089984`
- free physical memory KB (baseline sample): `11868112`
- balloon/virtio indicators:
  - `VirtIO Balloon Driver` detected
  - `VirtIO Viomem Driver` detected (Unknown status)
- last boot: `/Date(1778180311500)/`
- machineGuidHash: `d1fdf2ce69932d7ac2f9e3497be7f845d4a3ef9140c45bd0d7b0ea47c195508e`
- pending reboot: false
- critical events (last hour): `1`
