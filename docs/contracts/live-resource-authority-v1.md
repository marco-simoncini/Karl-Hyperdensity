# live-resource-authority-v1

Contract ID: `hyperdensity_live_resource_authority_v1`

Defines the KARL Live Resource Authority v1 operational contract for one governed and auditable live CPU/RAM control surface across Linux containers and evidence-backed Linux VMs.

## Product principle

- The Authority is one KARL-owned product surface.
- Runtime drivers are implementation adapters under the surface.
- Uniformity is at the operational contract layer, not the low-level mechanism layer.
- No raw runtime controls are exposed.

## Scope and safety invariants

- `authorityMode=unified_runtime_control_surface`
- `releaseTrack=technical_preview`
- `operationMode=operator_controlled_projection`
- `autonomousApplyAllowed=false`
- `enforcementMode=disabled`
- `productionMutationAllowed=false`
- `evidenceScope=evidence_namespace_only`
- `supportMatrixId=hyperdensity_release_support_matrix_v1`
- `evidenceBundleId=hyperdensity_evidence_bundle_demo_scenario_pack_v1`
- `policyPackId=hyperdensity_policy_pack_v1`
- `profilePackId=hyperdensity_shell_claim_templates_profile_pack_v1`

## Runtime driver semantics

- `container_linux_pod_resize_driver`
  - Linux container CPU/RAM via Kubernetes pod resize/runtime request/cgroup verification paths where proven.
- `vm_linux_cpu_libvirt_qga_driver`
  - Linux VM CPU via libvirt/QGA/runtime CPU path where proven.
- `vm_linux_memory_virtiomem_qmp_driver`
  - Linux VM RAM via runtime overlay through virtio-mem/QMP/QOM requested-size where proven.

Drivers are metadata only for governance and audit projection; they are not raw user-executable control APIs.

## Uniform contract phases

The Authority must project all required phases:

1. `live_resource_intent`
2. `capability_check`
3. `preflight`
4. `dry_run_or_dry_run_like_validation`
5. `runtime_lease_or_overlay`
6. `apply`
7. `verify`
8. `audit`
9. `rollback`
10. `reconcile_or_expire`

## Support boundary wording

Approved wording must include:

- General authority wording for one governed/auditable surface across Linux containers and evidence-backed Linux VMs.
- Container wording bounded to Technical Preview proven paths with verification and rollback.
- VM wording bounded to evidence-backed object-specific paths, with VM RAM runtime overlay wording.
- Uniformity wording bounded to the operational contract phases.

Rejected wording must include:

- claims of identical low-level mechanisms across containers and VMs.
- raw QMP/libvirt/QGA/cgroup control exposure.
- generic KubeVirt template mutation claim for VM RAM.
- universal VM/Kubernetes support claims.
- production autonomous movement claims.
- Windows support claims.
- dry-run success as production readiness.

## Required state handling

- Missing required source surfaces must degrade/block authority state.
- `warming_up`, `partial`, and `blocked` are not ready states.
- Windows lane remains `out_of_scope`/`frozen` and must not count as supported.
