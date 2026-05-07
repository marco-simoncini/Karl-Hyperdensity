# Product Architecture Update

Decision:

- `KARL Windows Prearmed Fluid Envelope v2` CPU control is validated via authorized runtime entitlement actuator, without guest-visible logical processor scaling.

v2 shape:

- CPU entitlement lease:
  - authorized node-local actuator writes cgroup CPU controller for target QEMU container scope.
  - entitlement moves by quota (`cpu.max`) with deterministic rollback.
- RAM entitlement lease:
  - QMP balloon target (already validated in previous proof).
- Guest-visible processor count:
  - remains constant and is not used as success criterion.
- Success semantics:
  - effective host/runtime entitlement change
  - guest ACK continuity
  - same-QEMU/same-boot/no-migration invariants.

Current classification:

- `WINDOWS_PREARMED_FLUID_ENVELOPE_CONFIRMED_PENDING_FULL_RERUN`
  - CPU actuator confirmed in this run.
  - RAM proof reused from prior run.
  - full combined CPU+RAM rerun in one sprint still pending.
