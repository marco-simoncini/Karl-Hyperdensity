# Audit Bundle

Persisted evidence bundle includes:

- actuator requests/allowlist/kill-switch state
- actuator dry-run/apply/down/restore outputs
- compliance replay output with attestation and bundle index

Files:

- `raw_logs_sanitized/compliance_replay_before_apply.json`
- `raw_logs_sanitized/actuator_request_set_floor.json`
- `raw_logs_sanitized/actuator_request_cpu_up.json`
- `raw_logs_sanitized/actuator_allowlist_master-win11.json`
- `raw_logs_sanitized/actuator_dry_run_output.json`
- `raw_logs_sanitized/actuator_set_floor_output.json`
- `raw_logs_sanitized/actuator_cpu_up_output.json`
- `raw_logs_sanitized/actuator_cpu_down_output.json`
- `raw_logs_sanitized/actuator_final_restore_output.json`

Note:

- attestation mode: `future-signable`
- hash chain mode: `local-deterministic-hash-chain`
