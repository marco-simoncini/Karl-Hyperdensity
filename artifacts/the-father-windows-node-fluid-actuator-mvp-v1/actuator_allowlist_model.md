# Actuator Allowlist Model

Model: `KARLNodeFluidActuatorAllowlist`

Fields:

- `allowlistId`
- identity: `nodeName`, `namespace`, `vmName`, `vmUid`, `podUid`, `qemuPid`, `qemuStartTime`
- scope: `allowedCgroupPath`, `allowedControllers`, `allowedActions`
- bounds: `minCpuMax`, `maxCpuMax`
- safety controls:
  - `allowParentCgroupWrite=false`
  - `allowArbitraryWrite=false`
  - `allowSymlinkTraversal=false`
- validity window: `createdAt`, `expiresAt`

Design goal: exact identity+path binding to one shell target.
