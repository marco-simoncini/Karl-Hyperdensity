# Blocker Taxonomy

## Blocking Categories

Hard blockers:
- `restart_bound`
- `rollout_bound`
- `migration_destructive`
- `missing_evidence`
- `recent_correlated_warning`

## Warning Policy

Historical warning debt can be treated as non-blocking only under bounded warning policy.

Recent correlated warnings remain blocking until explicitly cleared by policy conditions.

## Compliance Integrity Rule

No fake compliance claim is allowed when a blocker remains active.
