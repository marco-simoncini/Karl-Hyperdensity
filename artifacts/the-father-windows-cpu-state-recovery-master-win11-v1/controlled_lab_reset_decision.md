# Controlled Lab Reset Decision

Decision: **execute controlled lab reset**.

Why:

- mismatch persisted after two live recovery paths:
  - `virsh setvcpus ... 6 --live` timeout (previous phase),
  - direct QMP `device_del id=vcpu6` accepted but no state change.
- state remains unsafe for further live CPU trials (`QMP/libvirt=7`, `guest=6`).

Policy framing:

- reset/cold start is **outside Hyperdensity criteria**;
- used only to restore coherent baseline for future diagnostics;
- not a proof of runtime deterministic CPU hotplug support.

Guardrails for reset action:

- no pool/replica usage;
- no LiveMigration/VMIM;
- no new VM creation;
- preserve same VM identity;
- collect before/after evidence.

Chosen operation sequence:

1. `virtctl stop master-win11 -n karl`
2. wait until VM/VMI reach stopped state
3. `virtctl start master-win11 -n karl`
4. wait running
5. collect post-reset baseline and continuity deltas
