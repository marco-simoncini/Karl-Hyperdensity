# RAM Return-to-Floor

After guest release, QMP balloon down executed.

- query-balloon after down: `12884901888`
- expected floor: `12884901888`
- return matched floor: `True`
- guestAck after down: `True`
- pending reboot after down: `False`
- critical events 1h after down: `0`

RAM_WORKLOAD_DOWN_CONFIRMED = true.
