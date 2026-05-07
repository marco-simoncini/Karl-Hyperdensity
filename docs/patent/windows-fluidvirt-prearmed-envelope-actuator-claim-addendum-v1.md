# Windows FluidVirt Prearmed Envelope Actuator Claim Addendum v1

Draft tecnico per revisione brevettuale.

## Claimed technical direction

- Prearmed envelope for Windows VM shells.
- CPU entitlement modulation by node-local actuator writing cgroup v2 `cpu.max`.
- RAM entitlement modulation through QMP balloon control.
- Compliance engine that transforms heterogeneous Windows VM sources into Hyperdensity Ready shell targets.

## Core novelty anchors

- Standalone and pool-child VMs are treated as individual shell identities.
- Runtime continuity evidence is mandatory: same QEMU process, same Windows boot, same pod/node, no migration/recreate rollout path.
- Resource transitions require rollback and return-to-floor proof.
- Pool context is accepted for provisioning lineage only, not runtime scaling.

## Explicit non-claims

- No vCPU hotplug/unplug success path.
- No replica scaling mechanism.
- No LiveMigration/VMIM-based entitlement mechanism.
