# runtime_invariant_set

Model: `WindowsFluidRuntimeInvariantSet`

Mandatory invariants include:

- no live migration / no VMIM / no recreate / no rollout / no reboot
- same node / same virt-launcher pod / same qemu process
- same windows boot / same machine identity
- qmp ack / guest ack
- rollback ready / return-to-floor ready
- kill-switch ready
- evidence freshness
- qmp read-only until apply phase

Rules:

- identity invariants failing -> quarantine semantics
- readiness invariants failing -> blocked semantics
- fundamental invariants are required (not warning-only)
