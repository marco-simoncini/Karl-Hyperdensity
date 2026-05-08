# Replayed Events

Included deterministic replay events:
- compliance_replay_started
- product_model_loaded
- actuator_boundary_loaded
- readonly_replay_loaded
- fake_runtime_boundary_verified
- blocker_taxonomy_verified
- candidate_transition_verified
- return_to_floor_requirement_verified
- rollback_requirement_verified
- audit_hash_chain_built
- audit_hash_chain_verified
- compliance_replay_completed
- no_runtime_mutation_attested

Each event is replayed-only, ordered deterministically, and contains no secret material.
