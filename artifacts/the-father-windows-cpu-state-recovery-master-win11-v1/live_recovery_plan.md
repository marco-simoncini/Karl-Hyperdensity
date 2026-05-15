# Live Recovery Plan

Current state before recovery:

- QMP/libvirt CPU: `7`
- guest CPU: `6`
- mismatch: yes
- continuity baseline: same VM/pod/node/QEMU/boot confirmed

Chosen strategy:

- **B. QMP `device_del` of hotplugged CPU device-id**.
- Target device-id: `vcpu6`
- Device-id confidence:
  - QMP/QOM path `/machine/peripheral/vcpu6`
  - `qom-list /machine/peripheral` exposes child named `vcpu6`
  - `qom-get` shows `realized=true`, `hotplugged=true`

Why this is safer than repeating generic unplug:

- previous `virsh setvcpus ... 6 --live` timed out and did not change state;
- direct device-id operation targets only the known extra hotplugged CPU object.

Command plan:

1. before evidence snapshot (QMP/libvirt/guest + continuity)
2. execute:
   - `{"execute":"device_del","arguments":{"id":"vcpu6"}}`
3. wait short stabilization window
4. after evidence snapshot

Expected result:

- QMP/libvirt CPU returns to `6`
- guest remains or returns to `6`
- no reboot/migration/recreate
- same QEMU PID and start time

Stop conditions:

- QMP remains `7` after stabilization
- command errors
- continuity break

Fallback / rollback:

- no further random unplug mutations
- if live recovery fails, quarantine state and evaluate controlled lab reset as out-of-criteria recovery only.

Explicit scope note:

- This is a state recovery action, not Hyperdensity success proof.
