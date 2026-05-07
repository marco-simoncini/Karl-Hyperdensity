# CPU Entitlement Up

Mechanism:

- authorized node-local actuator writing host cgroup v2 `cpu.max` for target QEMU container scope.

Before:

- `cpu.max=300000 100000` (floor)
- `cpu.stat` before floor run:
  - `usage_usec=876421870`
  - `nr_throttled=713`
  - `throttled_usec=92310572`

Action:

- set `cpu.max` to `600000 100000` (ceiling).

After:

- compute-side readback: `cpu.max=600000 100000`
- `cpu.stat` after ceiling run:
  - `usage_usec=1081996568`
  - `nr_throttled=1013`
  - `throttled_usec=138579738`

Workload evidence (same 20s guest load profile):

- floor run elapsed: `28604 ms`
- ceiling run elapsed: `28618 ms`
- host cgroup delta comparison:
  - floor run CPU usage delta: `80654739 usec`
  - ceiling run CPU usage delta: `124919901 usec`
  - floor run throttled delta: `46132675 usec` (`+215`)
  - ceiling run throttled delta: `136491 usec` (`+85`)

Continuity checks:

- guest ACK true
- same QEMU PID/start (`96`, `Thu May 7 18:58:03 2026`)
- same pod/node
- same Windows boot and machine hash

Result:

- `CPU_ENTITLEMENT_UP_CONFIRMED=true`
