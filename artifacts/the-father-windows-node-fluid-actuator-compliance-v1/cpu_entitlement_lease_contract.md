# CPU Entitlement Lease Contract

Implemented `WindowsCpuEntitlementLease`:

- Lease identity and target: `leaseId`, `shellRef`, `targetVm`
- Model controls: `mode=prearmed-fluid-envelope`, `mechanism=cgroup-v2-cpu-max`
- Envelope values: `floorCpuMax`, `ceilingCpuMax`, `currentCpuMax`, `requestedCpuMax`
- Safety gates: same-QEMU, same-boot, same-pod, same-node, guest ACK, actuator ACK
- Recovery controls: `rollbackTarget`, `returnToFloorTarget`, `ttlSeconds`
- Audit: `status`, `blockers`, `evidenceRefs`

Validation rejects:

- vCPU hotplug/unplug requests
- VM spec patch mutation intent
- missing rollback/return targets
- stale/invalid TTL
