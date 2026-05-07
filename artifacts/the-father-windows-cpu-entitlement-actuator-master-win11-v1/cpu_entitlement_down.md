# CPU Entitlement Down / Return to Floor

Status: `CONFIRMED`

Down action:

- from ceiling `cpu.max=600000 100000` to floor `cpu.max=300000 100000`.

Before down workload:

- `usage_usec=698480304`
- `nr_throttled=499`
- `throttled_usec=45806823`

After down workload (same 20s guest load shape):

- `usage_usec=778361552`
- `nr_throttled=713`
- `throttled_usec=92310572`

Metric coherence:

- CPU usage delta under floor:
  - `79881248 usec` (778361552 - 698480304)
- Throttle delta under floor:
  - `+214 periods` and `+46503749 usec`
- Compared to ceiling run, floor run shows materially higher throttling and lower CPU usage gain.

Return-to-floor and rollback:

- floor set confirmed (`300000 100000`).
- final rollback to original baseline performed:
  - restored `cpu.max=600000 100000`.

Continuity:

- guest ACK true
- same QEMU PID/start (`96`, `Thu May 7 18:58:03 2026`)
- same Windows boot (`/Date(1778180311500)/`)
- same pod/node
- no VMIM / no migration observed

Result:

- `CPU_ENTITLEMENT_DOWN_CONFIRMED=true`
