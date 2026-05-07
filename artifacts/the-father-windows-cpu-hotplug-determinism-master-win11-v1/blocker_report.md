# Blocker Report

Final verdict:

- `WINDOWS_CPU_RETURN_TO_FLOOR_BLOCKED_UNPLUG_TIMEOUT`

Observed blockers:

1. `cpu_unplug_timeout_live_7_to_6`
   - `virsh setvcpus ... 6 --live` times out.
2. `qmp_guest_cpu_mismatch_persistent`
   - QMP/libvirt remains at `7`.
   - Guest remains at `6`.
3. `guest_cpu_online_offline_control_missing`
   - Guest reports only 6 CPUs and all with `can-offline=false`.
   - QGA does not expose enabled guest-side `set-vcpus`.

Why determinism is not confirmed:

- No guest-confirmed `+CPU` event in this phase.
- No successful return-to-floor `7 -> 6` across both QMP and guest.
- State remains divergent; proceeding would violate deterministic safety rules.
