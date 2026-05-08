# Replayed CPU Entitlement Transitions

Model-only transitions included:
- floor_to_ceiling_candidate
- ceiling_to_floor_return
- blocked_without_approval
- blocked_without_guest_witness
- blocked_without_return_to_floor
- blocked_without_audit_chain

For each transition:
- `transitionState=replayed_model_only`
- `appliedToRealCgroup=false`
- `fakeRuntimeOnly=true`
