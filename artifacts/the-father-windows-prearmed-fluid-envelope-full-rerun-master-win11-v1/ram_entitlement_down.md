# RAM Entitlement Down / Return-to-Floor

Before down:

- QMP balloon actual: `13958643712` (13 GiB)
- libvirt balloon: `balloon.current=13631488 KiB`
- guest ACK true

Action:

- `{"execute":"balloon","arguments":{"value":12884901888}}`

After down:

- QMP convergence/steady state: `12884901888` (12 GiB floor)
- libvirt balloon: `balloon.current=12582912 KiB`
- guest free memory KB sample: `11477284`
- guest ACK true
- pending reboot false
- critical events 1h: `0`

Continuity:

- same QEMU PID/start (`96`, `Thu May 7 18:58:03 2026`)
- same pod/node
- same boot/hash
- no VMIM / no migration

Result:

- `RAM_ENTITLEMENT_DOWN_CONFIRMED=true`
