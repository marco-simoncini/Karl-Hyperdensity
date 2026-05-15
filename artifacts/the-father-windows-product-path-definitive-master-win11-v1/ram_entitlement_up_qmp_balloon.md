# RAM Entitlement Up (QMP Balloon)

Mechanism:

- RAM is balloon liquidity via QMP.

Action:

- QMP balloon `12884901888 -> 13958643712`.

Observed evidence:

- query-balloon reached `13958643712`
- domstats balloon increased (`balloon.current` to `13606912 KiB` range)
- guest ACK: true
- pending reboot: false
- critical events 1h: `0`

Result:

- `RAM_ENTITLEMENT_UP_CONFIRMED=true`
