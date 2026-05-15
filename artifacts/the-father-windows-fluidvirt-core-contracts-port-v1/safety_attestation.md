# Safety Attestation

## Safety Preservation
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

## Runtime Safety
- action slate is planning/model-only
- no cgroup write execution
- no QMP command execution
- no actuator invocation
- no executor invocation
- no privileged operation path

## Delivery Safety
- no direct merge
- no dashboard/inventory/os-iso touch
- no artifact bulk import
