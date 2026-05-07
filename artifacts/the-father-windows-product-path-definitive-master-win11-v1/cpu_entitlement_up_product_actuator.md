# CPU Entitlement Up (Product Actuator)

Mechanism:

- CPU is entitlement liquidity via controlled node actuator (`cpu.max` only).

Cycle performed:

1. Establish floor (`600000 -> 300000`) using `actuator_request_set_floor.json`.
2. Apply CPU up (`300000 -> 600000`) using `actuator_request_cpu_up.json`.

Evidence:

- `raw_logs_sanitized/actuator_set_floor_output.json`
- `raw_logs_sanitized/actuator_cpu_up_output.json`

Workload metrics (guest load + host cpu.stat deltas):

- floor run1:
  - elapsed: `27533 ms`
  - usage delta: `80225139 usec`
  - throttled delta: `58796712 usec`
  - nr_throttled delta: `216`
- ceiling run:
  - elapsed: `28501 ms`
  - usage delta: `142286379 usec`
  - throttled delta: `278773 usec`
  - nr_throttled delta: `113`

Interpretation:

- ceiling shows materially higher CPU usage budget with strongly reduced throttling.

Result:

- `CPU_ENTITLEMENT_UP_CONFIRMED=true`
