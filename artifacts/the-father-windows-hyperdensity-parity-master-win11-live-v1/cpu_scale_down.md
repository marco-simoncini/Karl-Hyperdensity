# CPU Scale Down

`CPU_DOWN_CONFIRMED=false`

## Action attempted

- `virsh setvcpus karl_master-win11 6 --live`

## Result

- operation error: `vcpu unplug request timed out`
- QMP CPU remains `7`
- guest logical CPU remains `6`

Deterministic in-place return to floor was not proven.
