# master_win11_replay_summary

Source reference reused from previous milestone log:

- `artifacts/the-father-windows-fluidvirt-runtime-gate-qmp-v1/build_or_test_logs/cluster-probe.log`

Observed read-only metadata used for fixture shaping:

- namespace: `karl`
- VM name: `master-win11`
- status in probe snapshot: `Stopped`

Replay results:

- ready fixture => `READY_FOR_FLUID_SHELL_CERTIFICATION` / phase `READY`
- missing annotations fixture => `BLOCKED_GENERIC_WINDOWS_VM`
- missing QMP fixture => `BLOCKED_MISSING_QMP`
- missing guest fixture => `BLOCKED_MISSING_GUEST_ACK`
