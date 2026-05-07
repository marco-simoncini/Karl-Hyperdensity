# Blocker Report

Final state verdict:

- `WINDOWS_CPU_STATE_RECOVERED_BY_LAB_RESET`

Root-cause classification:

- `WINDOWS_CPU_HOTPLUG_ROOT_CAUSE_IDENTIFIED`

Identified blockers in live path:

1. `guest_not_onlining_hotplugged_cpu`
   - QMP/libvirt sees CPU `vcpu6`; guest logical processors stay at 6.
2. `cpu_unplug_timeout_or_no_effect_live`
   - `setvcpus --live` previously timed out.
   - direct `device_del vcpu6` accepted but no effective removal.
3. `guest_cpu_offline_control_unavailable`
   - QGA `guest-set-vcpus` disabled.
   - `guest-get-vcpus` reports `can-offline=false`.
4. `qmp_guest_mismatch_quarantine`
   - persistent `7/6` mismatch required quarantine and non-live recovery.

Controlled reset outcome:

- restored coherent baseline `6/6` but changed runtime continuity anchors (new VMI/pod/QEMU/boot).
- therefore this recovery is operational-lab only, not Hyperdensity proof.
