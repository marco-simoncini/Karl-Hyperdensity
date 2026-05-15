# RAM Entitlement Down / Return-to-Floor Test

Mechanism:

- QMP balloon target decrease (`balloon value`), no memory hotplug.

Safety baseline:

- before-down `query-balloon actual`: `13958643712` bytes (`13 GiB`)
- guest free physical memory baseline sample: `11674088 KB`
- pending reboot: false
- guest ACK: true

Action:

- `{"execute":"balloon","arguments":{"value":12884901888}}`

Host/QMP verification:

- `query-balloon actual` converged to `12884901888`.
- `domstats --balloon` moved to `balloon.current=12582912 KiB` (`12 GiB` floor).

Guest and continuity verification:

- guest ACK true after mutation.
- same QEMU PID/start (`96`, `Thu May 7 18:58:03 2026`).
- same pod/node (`virt-launcher-master-win11-kmwgg`, `karl-lab-metal-01`).
- same Windows boot (`/Date(1778180311500)/`).
- no VMIM; no migration evidence.

Result:

- `RAM_ENTITLEMENT_DOWN_CONFIRMED=true`
- return-to-floor for RAM: verified (`12 GiB`).
