# RAM Entitlement Down (Return-to-Floor)

Action:

- QMP balloon `13958643712 -> 12884901888`.

Observed evidence:

- query-balloon returned to `12884901888`
- domstats balloon returned to floor range (`balloon.current` around `12670976 KiB`)
- guest ACK: true
- pending reboot: false
- critical events 1h: `0`

Result:

- `RAM_ENTITLEMENT_DOWN_CONFIRMED=true`
