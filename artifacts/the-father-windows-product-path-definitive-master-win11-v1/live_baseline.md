# Live Baseline

Cluster/context:

- context: `karl-metal-01@ovh`
- namespace: `karl`
- vm: `master-win11`

Identity baseline:

- VM UID: `c81b95dc-d955-4fb3-a1af-59d979f48bcb`
- VMI UID: `b31dcab6-5a99-432e-a6fe-21607a9e3403`
- VMI phase: `Running`
- virt-launcher pod: `virt-launcher-master-win11-kmwgg`
- pod UID: `7b6a904a-1c9a-4a44-9b37-1dc737304773`
- node: `karl-lab-metal-01`
- restart count: `0`
- VMIM objects: none

Runtime baseline:

- QEMU PID (guest namespace): `96`
- QEMU start: `Thu May 7 18:58:03 2026`
- host qemu pid: `2480927`
- host cgroup path: `/kubepods.slice/kubepods-burstable.slice/kubepods-burstable-pod7b6a904a_1c9a_4a44_9b37_1dc737304773.slice/cri-containerd-6b3736cf90a6e2f41147f27714926ac68e57ff1a2d2b22da2375a4a74fcf0c87.scope`
- cpu.max baseline: `600000 100000`
- QMP status: `running`
- QMP balloon baseline: `12884901888`

Guest baseline:

- guest ACK: true
- Windows boot: `/Date(1778180311500)/`
- machineGuidHash: `d1fdf2ce69932d7ac2f9e3497be7f845d4a3ef9140c45bd0d7b0ea47c195508e`
- logical processors: `6`
- visible memory bytes: `13938089984`
- pending reboot: false
- critical events 1h: `0`
