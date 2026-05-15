# CPU Entitlement Down / Return-to-Floor

Before down:

- `cpu.max=600000 100000` (ceiling)
- `cpu.stat` before down run:
  - `usage_usec=1081996626`
  - `nr_throttled=1013`
  - `throttled_usec=138579738`

Action:

- set `cpu.max` back to floor `300000 100000`.

After down run:

- compute-side readback: `cpu.max=300000 100000`
- `cpu.stat` after down run:
  - `usage_usec=1162226180`
  - `nr_throttled=1227`
  - `throttled_usec=183382937`

Workload evidence:

- down floor workload elapsed: `28642 ms`
- down floor host deltas:
  - CPU usage delta: `80229554 usec`
  - throttled delta: `44803199 usec` (`+214`)
- behavior is coherent with floor budget (high throttling vs ceiling run).

Continuity:

- guest ACK true
- same QEMU PID/start (`96`, `Thu May 7 18:58:03 2026`)
- same pod/node
- same Windows boot and machine hash
- no VMIM / no migration

Result:

- `CPU_ENTITLEMENT_DOWN_CONFIRMED=true`
