# Mismatch Recovery Plan

Current state before recovery attempt:

- QMP/libvirt CPU: `7`
- Guest CPU: `6`
- Divergence: present
- Pending reboot: `false`
- Same QEMU process and same Windows boot: confirmed pre-action

Risk model:

- Prior `virsh setvcpus ... 6 --live` timed out during unplug.
- State suggests a hotplugged vCPU object exists (`/machine/peripheral/vcpu6`) but guest did not consume/online it.

Recovery objective:

- Return to floor `6/6` in place without reboot, recreate, migration, rollout.

Chosen command path (least-scope mutation first):

1. Attempt libvirt-managed live CPU down:
   - `virsh setvcpus karl_master-win11 6 --live`
2. If still divergent and safe to continue, evaluate direct QMP `device_del` against `vcpu6` in a follow-up phase (not mixed in same attempt block).

Success criteria for this phase:

- QMP CPU count returns to `6`
- Guest CPU count remains/returns to `6`
- QEMU PID unchanged
- QEMU start time unchanged
- Pod UID unchanged
- Node unchanged
- Windows last boot unchanged
- Machine hash unchanged
- No VMIM / migration / recreate evidence

Fail-fast criteria:

- Unplug timeout persists
- QMP remains `7`
- Guest remains `6`

If fail-fast criteria hit, stop and classify:

- `WINDOWS_CPU_RETURN_TO_FLOOR_BLOCKED_UNPLUG_TIMEOUT`
- and quarantine note for persistent mismatch risk.
