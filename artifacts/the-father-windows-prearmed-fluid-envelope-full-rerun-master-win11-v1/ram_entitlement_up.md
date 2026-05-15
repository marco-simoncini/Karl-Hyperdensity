# RAM Entitlement Up

Mechanism:

- QMP balloon target increase (`query-balloon` / `balloon`).

Before up:

- QMP balloon actual: `12884901888` (12 GiB)
- libvirt balloon: `balloon.current=12582912 KiB`
- guest free memory KB: `10755268`
- guest ACK true

Action:

- `{"execute":"balloon","arguments":{"value":13958643712}}`

After up:

- QMP convergence:
  - `13115588608`
  - `13350469632`
  - `13583253504`
  - `13811843072`
  - `13958643712` (13 GiB)
- libvirt balloon returned to `balloon.current=13631488 KiB`
- guest free memory KB sample: `10750084`
- guest ACK true

Continuity:

- same QEMU PID/start (`96`, `Thu May 7 18:58:03 2026`)
- same pod/node
- same boot/hash

Result:

- `RAM_ENTITLEMENT_UP_CONFIRMED=true`
