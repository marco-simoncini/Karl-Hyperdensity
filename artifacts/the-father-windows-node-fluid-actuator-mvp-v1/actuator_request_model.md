# Actuator Request Model

Model: `KARLNodeFluidActuatorRequest`

Core fields:

- request metadata: `requestId`, `requestVersion`, `action`
- identity pinning: `namespace`, `vmName`, `vmUid`, `vmiUid`, `podUid`, `nodeName`, `qemuPid`, `qemuStartTime`
- target pinning: `cgroupPath`, `controller`
- cpu bounds/targets: `previousCpuMax`, `requestedCpuMax`, `rollbackCpuMax`, `minCpuMax`, `maxCpuMax`
- timing: `ttlSeconds`, `createdAt`, `expiresAt`
- policy/audit: `reason`, `risk`, `evidenceRefs`, `policyVersion`, optional `attestationRef`

Enforced:

- stale request reject
- missing rollback reject
- missing pod/qemu identity reject
- controller different from `cpu.max` reject
