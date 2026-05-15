# Actuator Boundary Schema

## Root Contract
`WindowsFluidVirtNodeActuatorContract`

Required root fields:
- `actuatorContractId=windows_fluidvirt_node_actuator_contract_boundary_v1`
- `actuatorContractVersion=v1`
- `releaseTrack=technical_preview`
- `laneStatus=gated_preview`
- `contractMode=boundary_only`
- runtime/apply/autonomy/production flags all disabled

## CPU Boundary
- mechanism: `cgroup_v2_cpu_max_entitlement_liquidity`
- mechanism state: `contract_defined_not_runtime_enabled`
- host scope: `node_local`
- write/cgroup/cpu.max mutation: disabled
- requires allowlist/manual approval/guest witness/same-boot/same-qemu/rollback/return-to-floor/audit/kill-switch/lease-ttl

## RAM Boundary
- mechanism: `qmp_balloon_liquidity_model`
- mechanism state: `contract_reference_only`
- qmp command execution/raw qmp controls/memory apply disabled
- generic kubevirt vm ram template mutation disabled
