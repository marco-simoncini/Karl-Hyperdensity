# Replayed Safety Checks

Safety checks replayed against contract gates include:
- contract boundary loaded
- node allowlist/manual approval/lease ttl/rollback/return-to-floor/audit/kill-switch gates
- no raw control exposure
- no autonomous apply
- no production apply
- no windows ga claim

Each check is marked as `passed` or `replayed`, and `requiredBeforeRuntimeMvp=true`.
