# QMP Access Proof

QMP access status: `confirmed`

- Socket path: `/var/run/kubevirt-private/libvirt/qemu/lib/domain-1-karl_master-win11/monitor.sock`
- Domain via libvirt session: `karl_master-win11` (`running`)
- QEMU PID: `92`
- QEMU process start time: `Thu May 7 18:16:51 2026`
- QMP commands executed (read-only):
  - `qmp_capabilities`
  - `query-status`
  - `query-cpus-fast`
  - `query-hotpluggable-cpus`
  - `query-memory-devices`
  - `query-memory-size-summary`

No forbidden QMP mutating command was issued in read-only phase.
