# Selected Actuator Plan

Selected path:

- Option A (lab-scoped privileged node-local actuator).

Why selected:

- Only path that can legally write target cgroup CPU controller while preserving same VM/QEMU/boot.
- Fully reversible by restoring original `cpu.max`.

Plan and guardrails:

1. Create temporary privileged pod pinned to `karl-lab-metal-01` with host `/proc` and host `/sys/fs/cgroup`.
2. Resolve host QEMU PID by `guest=karl_master-win11`.
3. Resolve exact host cgroup path from `/host-proc/<pid>/cgroup`.
4. Baseline:
   - `cpu.max`, `cpu.stat`, QEMU PID/start, guest ACK, Windows boot.
5. CPU entitlement mutations:
   - floor: `300000 100000`
   - ceiling: `600000 100000`
6. Workload probe:
   - guest 20s multi-worker CPU load via QGA `guest-exec`.
7. Verify each phase:
   - cgroup setting changed/readback from compute container.
   - `cpu.stat` deltas coherent with entitlement.
   - same QEMU, same boot, same pod/node, guest ACK.
8. Rollback:
   - restore original `cpu.max=600000 100000`.
9. Teardown:
   - delete temporary actuator pod.

Stop conditions:

- inability to resolve cgroup path
- write failure on target controller
- continuity break
- guest ACK loss
