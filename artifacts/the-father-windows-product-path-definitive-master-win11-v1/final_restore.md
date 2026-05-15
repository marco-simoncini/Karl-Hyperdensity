# Final Restore

CPU restore:

- rollback command applied with `actuator_request_set_floor.json`
- evidence: `raw_logs_sanitized/actuator_final_restore_output.json`
- cpu.max final: `600000 100000`

RAM restore:

- final query-balloon: `12884901888`

Guest final:

- guest ACK: true
- pending reboot: false
- critical events 1h: `0`

Continuity final:

- same QEMU PID: `96`
- same QEMU start: `Thu May 7 18:58:03 2026`
- same boot: `/Date(1778180311500)/`
- same machineGuidHash: `d1fdf2ce69932d7ac2f9e3497be7f845d4a3ef9140c45bd0d7b0ea47c195508e`
- same pod UID: `7b6a904a-1c9a-4a44-9b37-1dc737304773`
- same node: `karl-lab-metal-01`
