# CPU Control Surface Discovery

Target live identity:

- VM: `master-win11` (`karl`)
- VMI UID: `b31dcab6-5a99-432e-a6fe-21607a9e3403`
- pod: `virt-launcher-master-win11-kmwgg` (`7b6a904a-1c9a-4a44-9b37-1dc737304773`)
- node: `karl-lab-metal-01`
- QEMU PID/start: `96` / `Thu May 7 18:58:03 2026`

Pod/runtime surface:

- container runtime ID: `containerd://6b3736cf90a6e2f41147f27714926ac68e57ff1a2d2b22da2375a4a74fcf0c87`
- QoS class: `Burstable`
- runtimeClass: none
- compute container security context:
  - `allowPrivilegeEscalation=false`
  - capabilities dropped (`ALL`) with only `NET_BIND_SERVICE` added
  - non-root user/group `107`

cgroup surface:

- cgroup version: `v2`
- controllers: `cpuset cpu io memory ...`
- QEMU host cgroup (resolved from host PID):
  - `/kubepods.slice/kubepods-burstable.slice/kubepods-burstable-pod7b6a904a_1c9a_4a44_9b37_1dc737304773.slice/cri-containerd-6b3736cf90a6e2f41147f27714926ac68e57ff1a2d2b22da2375a4a74fcf0c87.scope`
- cpu controls present:
  - `cpu.max=600000 100000`
  - `cpu.weight=67`
  - `cpuset.cpus.effective=0-31`
- from compute container mount: cgroup2 mounted read-only, writes blocked.

Libvirt surface:

- active URI: `qemu:///session`
- system URI unavailable (`qemu:///system` rejected)
- `schedinfo` tune operation not supported in session mode
- pin info readable (`vcpupin`, `emulatorpin`)

QEMU/vCPU thread surface:

- QEMU thread names:
  - `qemu-kvm` (main)
  - `CPU 0/KVM`..`CPU 5/KVM`
- affinity readable (`0-31`)
- affinity writes from compute container blocked (`taskset` not permitted).

Guest surface:

- guest ACK true
- `guest-set-vcpus` disabled (expected; not used)
- processor count stable at `6`
- workload execution via `guest-exec` available (used for controlled CPU load generation).
