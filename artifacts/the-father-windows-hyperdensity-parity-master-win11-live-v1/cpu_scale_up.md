# CPU Scale Up

`CPU_UP_CONFIRMED=false`

## Action attempted

- `virsh setvcpus karl_master-win11 7 --live` (inside target virt-launcher)

## Verification

- QMP CPU changed from `6` to `7` (QMP positive)
- Guest logical CPU remained `6` (guest confirmation missing)
- Same QEMU PID: `92` (unchanged)
- Same node: `karl-lab-metal-01` (unchanged)
- Same pod UID: `6e0a597f-b508-450d-a4ac-d5e4262c8615` (unchanged)
- Pending reboot: false
- Critical events: 0

Because guest actual state did not confirm the CPU increase, `CPU_UP_CONFIRMED` stays false.
