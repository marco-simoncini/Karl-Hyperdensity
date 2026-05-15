# Final Baseline Restore

Final restore target:

- CPU baseline restored to original `cpu.max=600000 100000`.
- RAM floor restored to `query-balloon=12884901888` (12 GiB).

Verification:

- compute reads `cpu.max=600000 100000`
- QMP balloon actual `12884901888`
- libvirt balloon `balloon.current=12582912 KiB`
- guest ACK true
- pending reboot false
- critical events 1h `0`

Continuity at final state:

- QEMU PID/start unchanged: `96`, `Thu May 7 18:58:03 2026`
- Windows boot unchanged: `/Date(1778180311500)/`
- machineGuidHash unchanged: `d1fdf2ce69932d7ac2f9e3497be7f845d4a3ef9140c45bd0d7b0ea47c195508e`
- same pod UID: `7b6a904a-1c9a-4a44-9b37-1dc737304773`
- same node: `karl-lab-metal-01`
- VMIM objects: none
- migration evidence: none
