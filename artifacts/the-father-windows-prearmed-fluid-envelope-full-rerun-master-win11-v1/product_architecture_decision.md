# Product Architecture Decision

`KARL Windows Prearmed Fluid Envelope v2` is confirmed in this full combined rerun.

Architecture statement:

- CPU lease:
  - runtime entitlement by authorized node-local actuator writing target QEMU cgroup `cpu.max`.
- RAM lease:
  - runtime entitlement by QMP balloon target (`query-balloon` + libvirt balloon evidence).
- Continuity contract:
  - same QEMU process and start time
  - same Windows boot and machine identity
  - same pod and node
  - no VMIM / no migration / no rollout / no recreate
- Guest semantics:
  - guest ACK maintained throughout
  - processor count allowed to remain constant
  - success criteria are effective entitlement + continuity + rollback/return proofs.

Scope note:

- This is **Prearmed Fluid Envelope parity** confirmation.
- No claim is made on logical CPU scaling via vCPU hotplug.
