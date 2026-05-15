# apply_governance_contract

Model: `WindowsFluidApplyGovernanceContract`

Core properties:

- links source admission decision and action slate references;
- encodes requested future action (`future-cpu-apply`, `future-memory-apply`, etc.);
- governance phases:
  - `CONTRACT_PREPARED`
  - `CONTRACT_BLOCKED`
  - `CONTRACT_QUARANTINED`
  - `NEEDS_REVALIDATION`
- hard non-execution locks:
  - `mutationAllowed=false`
  - `applyAllowed=false`
  - `runtimeMode=in-place-qmp`
- carries required evidence/revalidation/post-apply verification requirements;
- captures blockers/denial reasons and auditable references.
