# policy-apply-v1

Contract ID: `hyperdensity_policy_apply_v1`

Defines operator-controlled apply intent and guard result envelope.

Apply is valid only when guards pass and rollback is provable.

Blocking examples:
- restart-bound
- rollout-bound
- destructive migration-bound
- missing evidence
