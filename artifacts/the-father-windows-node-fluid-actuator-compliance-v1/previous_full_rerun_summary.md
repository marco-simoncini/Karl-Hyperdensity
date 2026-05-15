# Previous Full Rerun Summary

Source artifact: `artifacts/the-father-windows-prearmed-fluid-envelope-full-rerun-master-win11-v1/`

- CPU floor/ceiling: `300000 100000` -> `600000 100000` (`cpu.max`)
- RAM floor/ceiling: `12884901888` -> `13958643712` bytes
- cgroup target path evidence: process-visible `/` with compute-side `cpu.max` readback; node actuator mechanism used for host mutation
- actuator lab mechanism: authorized privileged node-local pod writing cgroup v2 `cpu.max`
- QEMU process evidence: PID `96`, start `Thu May 7 18:58:03 2026`, unchanged
- pod/node evidence: pod `virt-launcher-master-win11-kmwgg`, pod UID `7b6a904a-1c9a-4a44-9b37-1dc737304773`, node `karl-lab-metal-01`
- guest ACK evidence: ACK true across CPU and RAM transitions
- rollback/return evidence: CPU down to floor and baseline restore; RAM down to floor confirmed by QMP/libvirt/guest
- residual blocker note: non-privileged compute context cannot write CPU entitlement directly
- architecture decision: Prearmed Fluid Envelope v2 confirmed (CPU via node actuator cgroup lease, RAM via QMP balloon, continuity invariants mandatory)
