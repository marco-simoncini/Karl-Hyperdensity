# Live Recovery Attempt

Attempt objective:

- recover mismatch (`QMP/libvirt=7`, `guest=6`) without reboot.

Before:

- QMP/libvirt CPU: `7`
- guest CPU: `6`
- extra hotplugged CPU device present: `vcpu6`

Mutation executed:

- QMP command:
  - `{"execute":"device_del","arguments":{"id":"vcpu6"}}`
- QMP accepted command (`return: {}`).

After stabilization:

- QMP/libvirt CPU: `7` (unchanged)
- `query-cpus-fast`: still includes `/machine/peripheral/vcpu6`
- guest CPU: `6` (unchanged)
- mismatch persists.

Continuity during live attempt:

- QEMU PID unchanged (`92`)
- QEMU start unchanged (`Thu May 7 18:16:51 2026`)
- pod UID unchanged (`6e0a597f-b508-450d-a4ac-d5e4262c8615`)
- node unchanged (`karl-lab-metal-01`)
- Windows boot unchanged (`/Date(1778177812500)/`)
- machineGuidHash unchanged
- no VMIM / no migration evidence

Result:

- live recovery path failed; state remained quarantined.
