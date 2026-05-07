# Node Fluid Actuator Contract

Implemented `KARLNodeFluidActuatorContract` with required fields:

- Identity: `actuatorId`, `requestId`, `actuatorVersion`, `createdAt`
- Target mapping: `shellRef`, `vmRef`, `namespace`, `virtLauncherPodRef`, `podUid`, `qemuPid`, `qemuStartTime`, `cgroupPath`
- CPU mutation: `requestedCpuMax`, `previousCpuMax`, `appliedCpuMax`, `rollbackCpuMax`
- Policy/scope: `allowedControllers`, `mutationScope`, `allowlistDecision`, `policyDecision`
- Evidence: `beforeEvidence`, `afterEvidence`, `rollbackEvidence`, `returnToFloorEvidence`, `blockers`

Validation enforces allowlist, mutation scope, non-ambiguous target, evidence presence, and CPU entitlement field consistency.
