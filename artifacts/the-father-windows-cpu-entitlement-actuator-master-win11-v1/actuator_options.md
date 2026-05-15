# Actuator Options

## Option A — Privileged Node-Local CPU Entitlement Actuator

- Mechanism: privileged pod/daemon on target node writes CPU controller (`cpu.max`, `cpu.weight`) in host cgroup path for target QEMU container scope.
- Privileges: hostPID + host `/proc` + host `/sys/fs/cgroup` writable mount.
- Blast radius: limited by explicit VM/pod allowlist and node pinning.
- Rollback: immediate restore of original `cpu.max`.
- Evidence: host cgroup before/after + compute-side readback + continuity + guest ACK.
- Compatibility: high with same-QEMU/no-reboot model.
- Suitability: very good for single-node lab.
- Result in this phase: **selected and validated**.

## Option B — Libvirt System-Mode cputune Bridge

- Mechanism: runtime cputune via libvirt system URI.
- Blocker: system URI not available in virt-launcher context (`qemu:///system` rejected).
- Suitability now: blocked by libvirt model.

## Option C — Kubelet/Pod Resource Envelope Runtime Control

- Mechanism: runtime change through pod resource envelope/cgroup parent.
- Blocker: no direct runtime mutation path available from current non-privileged compute context; would require external privileged controller anyway.
- Suitability now: partial, indirect.

## Option D — QEMU Thread Scheduler Adapter

- Mechanism: set affinity/scheduler policy on QEMU/vCPU threads.
- Blocker in current context: `taskset` writes not permitted from compute container.
- Suitability now: viable only through privileged host actuator.
