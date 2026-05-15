# CPU Entitlement Up

Status: `CONFIRMED`

Entitlement model:

- floor: `cpu.max=300000 100000` (3 CPU-equivalent budget)
- ceiling: `cpu.max=600000 100000` (6 CPU-equivalent budget)

Before up:

- `cpu.max=300000 100000`
- `cpu.stat`:
  - `usage_usec=573400875`
  - `nr_throttled=408`
  - `throttled_usec=45678959`

Action:

- write `600000 100000` on target host cgroup `cpu.max`.

After up:

- `cpu.max=600000 100000` (host and compute readback aligned)
- `cpu.stat` after same 20s guest load:
  - `usage_usec=698059928`
  - `nr_throttled=499`
  - `throttled_usec=45806823`

Metric coherence (same workload shape):

- CPU usage delta under ceiling load:
  - `124659053 usec` (698059928 - 573400875)
- Throttle delta under ceiling load:
  - `+91 periods` and `+127864 usec`

Guest and continuity:

- guest ACK true
- same QEMU PID/start (`96`, `Thu May 7 18:58:03 2026`)
- same Windows boot (`/Date(1778180311500)/`)
- same pod/node

Result:

- `CPU_ENTITLEMENT_UP_CONFIRMED=true`
