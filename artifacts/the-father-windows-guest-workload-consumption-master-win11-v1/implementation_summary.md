# Implementation Summary

Executed a live bounded guest workload proof on `master-win11` using KARL Windows Prearmed Fluid Envelope v2.

- CPU was consumed as entitlement liquidity via controlled node actuator (`cpu.max` on host cgroup v2).
- RAM was consumed as balloon liquidity via QMP (`balloon` command + guest allocation witness).
- Guest execution used QGA/guest-exec PowerShell scripts with JSON output.
- Continuity invariants remained pinned: same pod, node, QEMU PID/start, Windows boot, machineGuidHash.
- Evidence captured under `raw_logs_sanitized/` and summarized across CPU floor/ceiling/return and RAM floor/ceiling/return.

Result: **WINDOWS_GUEST_WORKLOAD_RESOURCE_CONSUMPTION_CONFIRMED**.
