# Windows FluidVirt Product Model v1

`windows-fluidvirt-product-model-v1` introduces core contracts/models only for Windows FluidVirt in KARL Hyperdensity.

## Release Boundary

- release track: `technical_preview`
- lane status: `technical_preview_candidate` or `gated_preview`
- this milestone is model/planning only
- no runtime mutation
- no autonomous apply
- no production apply
- no raw runtime control exposure

## Product Scope

This contract defines:

- product identity/version metadata
- support boundary and claim boundary
- CPU liquidity model through entitlement (`cgroup v2 cpu.max`)
- RAM liquidity model through QMP balloon **as model only**
- guest witness dependency (`KARL Agent`, `fluidShell`, `QGA`) as integration dependency
- blocker taxonomy and planning-only action slate
- readiness state model and safety defaults

## Forbidden Claims And Mechanisms

The following remain forbidden in this milestone:

- Windows GA claim
- Windows production-ready claim
- Windows execution-ready by default
- vCPU hotplug support claim
- logical CPU scaling support claim
- pool scaling support claim
- LiveMigration / VMIM support claim
- reboot/recreate/rollout as product mechanism
- autonomous apply
- production AUTO
- raw QMP/libvirt/QGA/QOM/K8s patch controls
- generic KubeVirt VM RAM template mutation support claim

## Safety Defaults

Core model defaults must remain:

- `productionMutationAllowed=false`
- `autonomousApplyAllowed=false`
- `enforcementMode=disabled`
- `windowsGaClaimAllowed=false`
- `windowsProductionReadyClaimAllowed=false`
- `windowsExecutionReadyByDefault=false`
- `vcpuHotplugClaimAllowed=false`
- `logicalCpuScalingClaimAllowed=false`
- `poolScalingClaimAllowed=false`
- `liveMigrationClaimAllowed=false`
- `rebootRecreateRolloutMechanismAllowed=false`
- `rawRuntimeControlsExposed=false`

## Action Slate Boundary

`WindowsFluidVirtActionSlate` is planning/model only.

It does not:

- execute cgroup writes
- execute QMP commands
- invoke actuator code paths
- invoke executor code paths
- execute privileged runtime operations

## Readiness Gates

Before any future apply progression, the model requires:

- guest witness and guest ACK continuity
- same-QEMU and same-boot continuity evidence
- rollback plan
- return-to-floor plan
- audit chain reference
- explicit manual approval in gated preview flows

This contract alone does not make Windows FluidVirt executable.
