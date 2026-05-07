# CPU Entitlement Up Test

Status: `NOT CONFIRMED`

Reason:

- No permitted runtime mechanism was available to increase/decrease CPU entitlement on running QEMU in this environment.

Evidence:

1. `virsh schedinfo karl_master-win11 --current`
   - `Operation not supported: CPU tuning is not available in session mode`.
2. cgroup control write
   - `/sys/fs/cgroup/cpu.max` is read-only from compute container.
3. thread affinity control
   - `taskset -cp ...` on QEMU/vCPU threads returns `Operation not permitted`.

Result:

- `CPU_ENTITLEMENT_UP_CONFIRMED=false`
