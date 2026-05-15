# KARL Windows Prearmed Fluid Envelope Plan

CPU envelope plan (for this VM):

- Floor entitlement target: 3 host CPUs equivalent.
- Ceiling entitlement target: 6 host CPUs equivalent.
- Candidate mechanisms evaluated:
  1. `virsh schedinfo` / `cputune` (session mode): unavailable.
  2. cgroup `cpu.max` write: blocked (read-only filesystem).
  3. vCPU thread affinity (`taskset`): blocked (`Operation not permitted`).
- Measurement plan that would be used if control existed:
  - host: cgroup cpu quota/weight or affinity + `domstats --vcpu`.
  - guest: fixed CPU-bound synthetic workload duration.
  - rollback: restore floor setting.

CPU plan conclusion:

- runtime entitlement control is blocked in current execution model (libvirt session + container permissions).

RAM envelope plan (for this VM):

- Current ceiling: `13 GiB` (`13631488 KiB` / `13958643712` bytes via QMP balloon actual).
- Floor target: `12 GiB` (`12582912 KiB` / `12884901888` bytes).
- Control mechanism selected:
  - QMP `balloon` command with explicit byte target.
- Host verification:
  - QMP `query-balloon`.
  - libvirt `domstats --balloon` (`balloon.current`).
- Guest verification:
  - guest ACK true.
  - free memory telemetry and safety checks.
- Return-to-floor safety checks:
  - pending reboot false.
  - no migration/VMIM.
  - no critical error spikes.
  - adequate free memory before reclaim.
