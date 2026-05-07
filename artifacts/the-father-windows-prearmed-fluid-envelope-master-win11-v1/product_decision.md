# Product Decision

Decision in this run:

- `KARL Windows Prearmed Fluid Envelope` is **not fully confirmed** on `master-win11` because CPU entitlement control is blocked in current runtime environment.

What was validated:

- RAM entitlement can be controlled runtime via QMP balloon target and returned to floor safely without reboot/recreate/migration.

What is blocked:

- CPU entitlement up/down control via available host mechanisms in current pod/libvirt-session permissions.

Design direction:

1. Keep the Windows model as prearmed envelope + entitlement leases.
2. Separate CPU control plane from guest-visible vCPU topology changes.
3. Enable an authorized CPU runtime controller path (outside current restrictions), then rerun CPU entitlement up/down proof on same-boot same-QEMU window.
4. Keep guest ACK + host/QMP evidence mandatory for claim.
