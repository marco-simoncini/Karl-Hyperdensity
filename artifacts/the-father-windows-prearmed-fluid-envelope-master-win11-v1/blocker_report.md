# Blocker Report

Final verdict:

- `WINDOWS_PREARMED_FLUID_ENVELOPE_BLOCKED_BY_CPU_CONTROL`

Blockers:

1. `cpu_tuning_unavailable_libvirt_session_mode`
   - `virsh schedinfo` reports CPU tuning not available in session mode.
2. `cgroup_cpu_controls_readonly`
   - `/sys/fs/cgroup/cpu.max` cannot be modified (read-only filesystem).
3. `thread_affinity_not_permitted`
   - `taskset -cp` on QEMU/vCPU threads returns operation not permitted.

Non-blocked findings:

- RAM entitlement path via QMP balloon is functional and reversible.
- Runtime continuity guardrails remained true during proof mutations:
  - same QEMU, same boot, same pod/node, no VMIM, no migration.
