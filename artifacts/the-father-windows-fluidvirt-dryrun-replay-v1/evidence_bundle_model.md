# evidence_bundle_model

`WindowsFluidRuntimeEvidenceBundle` includes:

- shell declared/runtime targets and observed resource state;
- KubeVirt continuity evidence (`before` and `after`);
- optional read-only QMP evidence;
- optional guest evidence from `modules.fluidShell`;
- optional non-mutating lease intent;
- policy gates, observed blockers, timestamps, source metadata, sanitization status.

Key safety behavior:

- missing QMP evidence => blocker (`qmp_socket_unavailable`);
- missing guest evidence => blocker (`guest_ack_missing`);
- incomplete identity evidence => blocked with synthetic continuity placeholders and blocker;
- pool source (`win11-pool-*` or explicit context flag) => `BLOCKED_POOL_REPLICA_MODEL`.
