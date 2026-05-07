# RAM Entitlement Up Test

Mechanism:

- QMP balloon target increase (`balloon value`), no memory hotplug.

Baseline for up test:

- floor target active: `12884901888` bytes (`12 GiB`)
- QEMU PID/start: `96` / `Thu May 7 18:58:03 2026`
- guest ACK: true

Action:

- `{"execute":"balloon","arguments":{"value":13958643712}}`

Host/QMP verification:

- `query-balloon` progression:
  - `13121880064`
  - `13354663936`
  - `13593739264`
  - `13828620288`
  - `13958643712`
- `domstats --balloon` returned to `balloon.current=13631488 KiB`.

Guest verification:

- guest ACK still true.
- guest telemetry remained available and stable (no reboot, no pending reboot).

Result:

- `RAM_ENTITLEMENT_UP_CONFIRMED=true`
