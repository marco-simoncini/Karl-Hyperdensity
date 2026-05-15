# CPU Return-to-Floor Recovery Attempt (7 -> 6)

Before attempt:

- libvirt/QMP CPU count: `7`
- guest CPU count: `6`
- mismatch present: yes

Command executed:

- `virsh -c qemu:///session setvcpus karl_master-win11 6 --live`

Command result:

- failed with timeout:
  - `Timed out during operation: vcpu unplug request timed out. Unplug result must be manually inspected in the domain`

After attempt:

- libvirt/QMP CPU count: `7` (unchanged)
- guest CPU count: `6` (unchanged)
- mismatch present: yes

Continuity checks after failure:

- same QEMU PID/start: yes (`92`, `Thu May 7 18:16:51 2026`)
- same pod UID/node: yes
- same Windows boot: yes (`/Date(1778177812500)/`)
- same machine hash: yes
- VMIM/migration evidence: none

Result:

- recovery attempt did not return VM to `6/6`.
