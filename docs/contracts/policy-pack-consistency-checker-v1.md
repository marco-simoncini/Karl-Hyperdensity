# policy-pack-consistency-checker-v1

Contract ID: `hyperdensity_policy_pack_consistency_checker_v1`

Validation-only consistency checker for `hyperdensity_policy_pack_v1`.

## Purpose

`policy-pack-consistency-checker-v1` verifies that Policy Pack v1 remains aligned with required Hyperdensity governance surfaces and safety gates without changing runtime behavior.

It is drift detection and proof, not enforcement.

## Required assertions

- `consistencyMode=validation_only`
- enforcement remains disabled
- no mutation path is introduced
- no autonomous apply path is introduced
- missing required safety gates are blockers
- warnings may exist for optional/future coverage
- Windows lane remains out-of-scope

## Required check domains

- policy identity and mode invariants
- section coverage
- safety gate coverage
- source-surface representation coverage
- readiness semantics (`warming_up`, `partial`, `blocked` are not ready)
- mutation safety (`admission_guard` remains audit-only, `mutate_preview` remains preview-only)

## Blocking behavior

The checker must report `consistent=false` when required sections, required safety gates, or blocking invariants are missing or violated.
